package rv

import (
	"fmt"
	"reflect"
	"strconv"
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
	Default string
}

func NewDefaultHandler(args []string) (FieldHandler, error) {
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
