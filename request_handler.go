package rv

import (
	"fmt"
	"reflect"
	"sort"
)

var handlerMap = map[string]FieldHandlerCreator{
	"source":  NewSourceFieldHandler,
	"type":    NewTypeHandler,
	"default": NewDefaultHandler,
	"range":   NewRangeHandler,
	"options": NewOptionsHandler,
}

type RequestHandler struct {
	Fields      map[string]FieldHandlers
	requestType reflect.Type
}

func NewRequestHandler(requestStruct interface{}) (*RequestHandler, error) {
	tags, err := extractTags(requestStruct)
	if err != nil {
		return nil, err
	}

	requestHandler := RequestHandler{requestType: reflect.TypeOf(requestStruct)}

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
