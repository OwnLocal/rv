package rv

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

type TypeHandler struct {
	Type string
}

var y = struct{}{}
var acceptedTypes = map[string]struct{}{
	"bool": y,
	"int":  y, "int8": y, "int16": y, "int32": y, "int64": y,
	"uint": y, "uint8": y, "uint16": y, "uint32": y, "uint64": y,
	"float32": y, "float64": y,
	"string": y,
	"time":   y,
}

func NewTypeHandler(args []string) (FieldHandler, error) {
	// TODO: add support for array types of basic types (right now the type of array elements is not passed along)
	if _, accepted := acceptedTypes[args[0]]; !accepted {
		return nil, fmt.Errorf("'%s' is not a supported type", args[0])
	}
	return TypeHandler{args[0]}, nil
}

func (h TypeHandler) Precidence() int { return 800 }
func (h TypeHandler) Run(r Request, f *Field) {
	if f.Value == nil {
		return
	}

	var err error

	switch h.Type {
	case "bool":
		err = toBool(&f.Value)
	case "int", "int8", "int16", "int32", "int64":
		err = toInt(&f.Value, h.Type)
	case "uint", "uint8", "uint16", "uint32", "uint64":
		err = toUint(&f.Value, h.Type)
	case "float32":
		err = toFloat(&f.Value, 32)
	case "float64":
		err = toFloat(&f.Value, 64)
	case "string":
		err = toString(&f.Value)
	case "time":
		err = toTime(&f.Value)
	default:
		err = fmt.Errorf("don't know how to convert to %s", h.Type)
	}

	if err != nil {
		f.Errors = append(f.Errors, err)
	}
}

// Follows the bool string options in strconv.ParseBool http://golang.org/pkg/strconv/#ParseBool
func toBool(val *interface{}) (err error) {
	switch v := (*val).(type) {
	case bool:
		// already ok
	case string:
		if b, e := strconv.ParseBool(v); e == nil {
			*val = b
		} else {
			err = e
		}
	case int, int8, int16, int32, int64:
		switch reflect.ValueOf(v).Int() {
		case 0:
			*val = false
		case 1:
			*val = true
		default:
			err = fmt.Errorf("bool int expected 1 or 0, got %d", v)
		}
	case uint, uint8, uint16, uint32, uint64:
		switch reflect.ValueOf(v).Uint() {
		case 0:
			*val = false
		case 1:
			*val = true
		default:
			err = fmt.Errorf("bool uint expected 1 or 0, got %d", v)
		}
	}

	return err
}

const maxInt = int64(int(^uint(0) >> 1))
const minInt = int64(-maxInt - 1)
const maxUint = uint64(^uint(0))

func toInt(val *interface{}, intType string) (err error) {
	var i int64

	switch v := (*val).(type) {

	case int, int8, int16, int32, int64:
		ival := reflect.ValueOf(v)
		i = ival.Int()

	case string:
		i, err = strconv.ParseInt(v, 0, 64)

	default:
		err = fmt.Errorf("don't know how to convert %T to %s", *val, intType)

	}

	if err == nil {
		switch {
		case intType == "int" && i <= maxInt && i >= minInt:
			*val = int(i)
		case intType == "int8" && i <= math.MaxInt8 && i >= math.MinInt8:
			*val = int8(i)
		case intType == "int16" && i <= math.MaxInt16 && i >= math.MinInt16:
			*val = int16(i)
		case intType == "int32" && i <= math.MaxInt32 && i >= math.MinInt32:
			*val = int32(i)
		case intType == "int64" && i <= math.MaxInt64 && i >= math.MinInt64:
			*val = int64(i)
		default:
			err = fmt.Errorf("int %d can't be represented as %s", i, intType)
		}
	}
	return err
}

func toUint(val *interface{}, intType string) (err error) {
	var ui uint64

	switch v := (*val).(type) {

	case int, int8, int16, int32, int64:
		ival := reflect.ValueOf(v)
		i := ival.Int()
		if i < 0 {
			err = fmt.Errorf("int %d can't be represented as %s", i, intType)
		} else {
			ui = uint64(i)
		}

	case uint, uint8, uint16, uint32, uint64:
		ival := reflect.ValueOf(v)
		ui = ival.Uint()

	case string:
		ui, err = strconv.ParseUint(v, 0, 64)

	default:
		err = fmt.Errorf("don't know how to convert %T to %s", *val, intType)

	}

	if err == nil {
		switch {
		case intType == "uint" && ui <= maxUint:
			*val = uint(ui)
		case intType == "uint8" && ui <= math.MaxUint8:
			*val = uint8(ui)
		case intType == "uint16" && ui <= math.MaxUint16:
			*val = uint16(ui)
		case intType == "uint32" && ui <= math.MaxUint32:
			*val = uint32(ui)
		case intType == "uint64" && ui <= math.MaxUint64:
			*val = uint64(ui)
		default:
			err = fmt.Errorf("uint %d can't be represented as %s", ui, intType)
		}
	}
	return err
}

func toFloat(val *interface{}, floatSize int) (err error) {
	var f float64

	switch v := (*val).(type) {
	case int, int8, int16, int32, int64:
		ival := reflect.ValueOf(v)
		f = float64(ival.Int())

	case uint, uint8, uint16, uint32, uint64:
		ival := reflect.ValueOf(v)
		f = float64(ival.Uint())

	case float32:
		if floatSize == 32 {
			// Leave it alone if it is already float32 so we don't lose any precision
			return nil
		}
		f = float64(v)

	case float64:
		f = v

	case string:
		f, err = strconv.ParseFloat(v, floatSize)

	default:
		err = fmt.Errorf("don't know how to convert %T to float%d", *val, floatSize)
	}

	if err == nil {
		switch {
		case floatSize == 32 && f <= math.MaxFloat32 && f >= -math.MaxFloat32:
			*val = float32(f)
		case floatSize == 64:
			*val = f
		}
	}
	return err
}

func toString(val *interface{}) (err error) {
	*val = fmt.Sprintf("%v", *val)
	return err
}

func toTime(val *interface{}) (err error) {
	switch v := (*val).(type) {
	case string:
		for _, time_fmt := range []string{
			"2006-01-02",
			"2006-01-02T15:04",
			"2006-01-02T15:04:05",
			"2006-01-02T15:04:05Z07:00",
		} {
			if *val, err = time.Parse(time_fmt, v); err == nil {
				break
			}
		}
	default:
		err = fmt.Errorf("don't know how to convert %T to time", *val)
	}
	return err
}
