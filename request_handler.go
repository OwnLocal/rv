package rv

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

var handlerMap = map[string]FieldHandlerCreator{
	"source":   NewSourceFieldHandler,
	"type":     NewTypeHandler,
	"default":  NewDefaultHandler,
	"range":    NewRangeHandler,
	"options":  NewOptionsHandler,
	"required": NewRequiredHandler,
}

type RequestHandler struct {
	Fields      map[string]FieldHandlers
	requestType reflect.Type

	indexCache     map[reflect.Type]int
	indexCacheLock sync.Mutex
}

// NewRequestHandler builds a RequestHandler which will extract and
// validate values from a request based on the "rv" tags on the struct
// fields.
func NewRequestHandler(requestStruct interface{}) (*RequestHandler, error) {
	tags, err := extractTags(requestStruct)
	if err != nil {
		return nil, err
	}

	requestHandler := RequestHandler{
		requestType: reflect.TypeOf(requestStruct),
		indexCache:  make(map[reflect.Type]int)}

	handlers := map[string]FieldHandlers{}
	for field, opts := range tags {
		fieldHandlers := FieldHandlers{}
		isList := false
		var listHandler ListHandler

		for opt, args := range opts {
			var err error
			if opt == "type" && args[0] == "slice" {
				isList = true
				listHandler = ListHandler{}
				listHandler.SubHandlers, err = addRegularHandler(FieldHandlers{}, "type", args[1:2])
			} else {
				fieldHandlers, err = addRegularHandler(fieldHandlers, opt, args)
			}
			if err != nil {
				return nil, err
			}

		}
		sort.Stable(fieldHandlers)
		if isList {
			fieldHandlers = addListHandler(fieldHandlers, listHandler)
		}
		handlers[field] = fieldHandlers
	}
	requestHandler.Fields = handlers
	return &requestHandler, nil
}

// Run fills the provided struct with data from the request, as
// specified in the "rv" tags on the struct fields.
func (h *RequestHandler) Run(req Request, requestStruct interface{}) (argErr error, fieldErrors map[string]Field) {
	val := reflect.ValueOf(requestStruct)
	if val.Type().Kind() != reflect.Ptr || val.Type().Elem() != h.requestType {
		return fmt.Errorf("Expected *%v, got %v", h.requestType, val.Type()), nil
	}
	val = val.Elem()

	fieldErrors = make(map[string]Field)
	for name, handlers := range h.Fields {
		field := Field{}
		for _, handler := range handlers {
			handler.Run(req, &field)
		}
		if len(field.Errors) > 0 {
			fieldErrors[name] = field
		} else if field.Value != nil {
			val.FieldByName(name).Set(reflect.ValueOf(field.Value))
		}
	}

	return nil, nil
}

// Bind searches the container for a field matching the
// RequestHandler's field type, then fills it by calling
// RequestHandler.Run with the specified Request and the matching
// field.
func (h *RequestHandler) Bind(req Request, container interface{}) (argErr error, fieldErrors map[string]Field) {
	val := reflect.ValueOf(container)
	if val.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("Expected pointer to struct, got %T", container), nil
	}

	val = val.Elem()
	if val.Type().Kind() != reflect.Struct {
		return fmt.Errorf("Expected pointer to struct, got %T", container), nil
	}

	i, err := h.fieldIndex(val.Type())
	if err != nil {
		return err, nil
	}

	return h.Run(req, val.Field(i).Addr().Interface())
}

func (h *RequestHandler) fieldIndex(container reflect.Type) (int, error) {
	h.indexCacheLock.Lock()
	i, ok := h.indexCache[container]
	h.indexCacheLock.Unlock()
	if ok {
		return i, nil
	}

	for i := 0; i < container.NumField(); i++ {
		field := container.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType == h.requestType {
			h.indexCacheLock.Lock()
			h.indexCache[container] = i
			h.indexCacheLock.Unlock()

			return i, nil
		}
	}
	return 0, fmt.Errorf("No %v field found in provided %v", h.requestType, container)
}

func addRegularHandler(fieldHandlers FieldHandlers, opt string, args []string) (FieldHandlers, error) {
	handlerCreator, ok := handlerMap[opt]
	if !ok {
		return fieldHandlers, fmt.Errorf("Invalid handler: %s", opt)
	}
	handler, err := handlerCreator(args)
	if err != nil {
		return fieldHandlers, err
	}
	fieldHandlers = append(fieldHandlers, handler)
	return fieldHandlers, nil
}

func addListHandler(fieldHandlers FieldHandlers, listHandler ListHandler) (fh FieldHandlers) {
	for _, handler := range fieldHandlers {
		switch handler.(type) {
		case SourceFieldHandler, TypeHandler, DefaultHandler:
			fh = append(fh, handler)
		default:
			listHandler.SubHandlers = append(listHandler.SubHandlers, handler)
		}
	}
	fh = append(fh, listHandler)
	return fh
}

type FieldHandlers []FieldHandler

func (f FieldHandlers) precidence(i int) int {
	if fv, ok := f[i].(PrecidenceFieldHandler); ok {
		return fv.Precidence()
	}
	return 0
}

// Implement sort.Interface
func (f FieldHandlers) Len() int           { return len(f) }
func (f FieldHandlers) Less(i, j int) bool { return f.precidence(i) > f.precidence(j) }
func (f FieldHandlers) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
