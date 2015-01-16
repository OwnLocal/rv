package rv

import (
	"fmt"
	"reflect"
	"strconv"
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

type DefaultHandler struct {
	Default interface{}
}

func NewDefaultHandler(args []string) (FieldHandler, error) {
	if len(args) > 1 {
		return DefaultHandler{args}, nil
	}
	return DefaultHandler{args[0]}, nil
}

func (h DefaultHandler) Run(req Request, field *Field) {
	if field.Value == nil {
		field.Value = h.Default
	}
}

// goes before TypeHandler so the default string will be transformed into the right type
func (h DefaultHandler) Precidence() int { return 900 }

type RangeHandler struct {
	Start string
	End   string
}

func NewRangeHandler(args []string) (FieldHandler, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("need two comma-separated arguments for range, got %#v", args)
	}
	return RangeHandler{args[0], args[1]}, nil
}

func (h RangeHandler) Run(req Request, field *Field) {
	var err error

	switch v := field.Value.(type) {
	case nil:
		err = fmt.Errorf("need value in range %s, %s, got no value", h.Start, h.End)
	case int, int8, int16, int32, int64:
		var min, max int64
		i := reflect.ValueOf(v).Int()
		min, max, err = h.intRange()
		if err == nil && i < min || i > max {
			err = fmt.Errorf("%d not in range %d, %d", i, min, max)
		}
	case uint, uint8, uint16, uint32, uint64:
		var min, max uint64
		i := reflect.ValueOf(v).Uint()
		min, max, err = h.uintRange()
		if err == nil && i < min || i > max {
			err = fmt.Errorf("%d not in range %d, %d", i, min, max)
		}
	case float32, float64:
		var min, max float64
		f := reflect.ValueOf(v).Float()
		min, max, err = h.floatRange()
		if err == nil && f < min || f > max {
			err = fmt.Errorf("%f not in range %f, %f", f, min, max)
		}
	case string:
		if v < h.Start || v > h.End {
			err = fmt.Errorf("%#v not in range %#v, %#v", v, h.Start, h.End)
		}
	default:
		err = fmt.Errorf("don't know how to determine range for %T(%v)", v, v)
	}

	if err != nil {
		field.Errors = append(field.Errors, err)
	}
}

func (h RangeHandler) intRange() (min, max int64, err error) {
	min, err = strconv.ParseInt(h.Start, 0, 64)
	if err == nil {
		max, err = strconv.ParseInt(h.End, 0, 64)
	}
	return min, max, err
}

func (h RangeHandler) uintRange() (min, max uint64, err error) {
	min, err = strconv.ParseUint(h.Start, 0, 64)
	if err == nil {
		max, err = strconv.ParseUint(h.End, 0, 64)
	}
	return min, max, err
}

func (h RangeHandler) floatRange() (min, max float64, err error) {
	min, err = strconv.ParseFloat(h.Start, 64)
	if err == nil {
		max, err = strconv.ParseFloat(h.End, 64)
	}
	return min, max, err
}

type OptionsHandler struct {
	Options map[string]struct{}
}

func NewOptionsHandler(args []string) (FieldHandler, error) {
	argSet := map[string]struct{}{}
	for _, arg := range args {
		argSet[arg] = struct{}{}
	}
	return OptionsHandler{Options: argSet}, nil
}

func (h OptionsHandler) Run(req Request, field *Field) {
	val := fmt.Sprintf("%v", field.Value)
	if _, valid := h.Options[val]; !valid {
		var options []string
		for opt, _ := range h.Options {
			options = append(options, opt)
		}
		field.Errors = append(field.Errors, fmt.Errorf("Expected one of %#v, got %#v", options, val))
	}
}

type ListHandler struct {
	SubHandlers FieldHandlers
}

func (h ListHandler) Run(req Request, field *Field) {
	var fields []*Field
	if field.Value == nil {
		return
	} else if val, ok := field.Value.(string); ok {
		for _, part := range strings.Split(val, ",") {
			fields = append(fields, &Field{Value: part})
		}
	} else if reflect.TypeOf(field.Value).Kind() == reflect.Slice {
		slice := reflect.ValueOf(field.Value)
		for i := 0; i < slice.Len(); i++ {
			fields = append(fields, &Field{Value: slice.Index(i).Interface()})
		}
	} else {
		fields = append(fields, &Field{Value: field.Value})
	}

	for _, subField := range fields {
		for _, handler := range h.SubHandlers {
			handler.Run(req, subField)
		}
	}

	if len(fields) == 0 {
		return
	}

	valSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(fields[0].Value)), 0, len(fields))
	for _, subField := range fields {
		valSlice = reflect.Append(valSlice, reflect.ValueOf(subField.Value))
		field.Errors = append(field.Errors, subField.Errors...)
	}
	field.Value = valSlice.Interface()
}

type RequiredHandler struct {
	Required bool
}

func NewRequiredHandler(args []string) (FieldHandler, error) {
	if required, err := strconv.ParseBool(args[0]); err == nil {
		return RequiredHandler{Required: required}, nil
	} else {
		return nil, err
	}

}

func (h RequiredHandler) Run(req Request, field *Field) {
	if h.Required && field.Value == nil {
		field.Errors = append(field.Errors, fmt.Errorf("required field missing"))
	}
}

func (h RequiredHandler) Precidence() int {
	return -100
}
