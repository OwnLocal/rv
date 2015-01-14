package rv

import (
	"fmt"
	"sort"
)

var handlerMap = map[string]FieldHandlerCreator{
	"source":  NewSourceFieldHandler,
	"type":    NewTypeHandler,
	"default": NewDefaultHandler,
}

func NewRequestHandler(requestStruct interface{}) (*RequestHandler, error) {
	tags, err := extractTags(requestStruct)
	if err != nil {
		return nil, err
	}

	handlers := map[string]FieldHandlers{}
	for field, opts := range tags {
		fieldHandlers := FieldHandlers{}
		for opt, args := range opts {
			handlerCreator, ok := handlerMap[opt]
			if !ok {
				return nil, fmt.Errorf("Invalid handler: %s", opt)
			}
			handler, err := handlerCreator(args)
			if err != nil {
				return nil, err
			}
			fieldHandlers = append(fieldHandlers, handler)
		}
		sort.Stable(fieldHandlers)
		handlers[field] = fieldHandlers
	}
	return &RequestHandler{Fields: handlers}, nil
}

type RequestHandler struct {
	Fields map[string]FieldHandlers
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
