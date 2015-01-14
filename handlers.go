package rv

import (
	"fmt"
	"net/url"
	"strings"
)

type FieldHandlerCreator func(args []string) (FieldHandler, error)

type FieldHandler interface {
	Run(Request, *Field)
}

type PrecidenceFieldHandler interface {
	FieldHandler
	Precidence() int
}

type Field struct {
	Value  interface{}
	Errors []error
}

// SourceFieldHandler takes source and field names from the first
// argument in the rv struct tag, pulls the value from the request and
// puts it into the field.
type SourceFieldHandler struct {
	Source Source
	Field  string
}

func NewSourceFieldHandler(args []string) (FieldHandler, error) {
	source_field := strings.Split(args[0], ".")
	if len(source_field) != 2 {
		return nil, fmt.Errorf("Expected 'source.field', got '%s'", args[0])
	}
	source, field := sourceMap[source_field[0]], source_field[1]
	if source == UNDEFINED {
		return nil, fmt.Errorf("Expected one of %v, got '%s'", sources, source_field[0])
	}
	return SourceFieldHandler{source, field}, nil
}

func (h SourceFieldHandler) Precidence() int { return 1000 }

func (h SourceFieldHandler) Run(r Request, f *Field) {
	var (
		err error
		val interface{}
		ok  bool
	)

	switch h.Source {
	case PATH:
		var pathArgs map[string]string
		pathArgs, err = r.PathArgs()
		val, ok = pathArgs[h.Field]

	case QUERY:
		var queryArgs url.Values
		queryArgs, err = r.QueryArgs()
		val = queryArgs.Get(h.Field)
		_, ok = queryArgs[h.Field]

	case JSON:
		var json map[string]interface{}
		json, err = r.BodyJson()
		val, ok = json[h.Field]

	case FORM:
		var form url.Values
		form, err = r.BodyForm()
		val = form.Get(h.Field)
		_, ok = form[h.Field]

	}

	if err != nil {
		f.Errors = append(f.Errors, err)
	} else if ok {
		f.Value = val
	}
}

type TypeHandler struct {
	Type string
}

func NewTypeHandler(args []string) (FieldHandler, error) {
	return TypeHandler{args[0]}, nil
}

func (h TypeHandler) Precidence() int { return 999 }
func (h TypeHandler) Run(r Request, f *Field) {
	// TODO: Validate type
}
