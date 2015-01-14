package rv

import (
	"fmt"
	"reflect"
	"strings"
)

func extractTags(reqStruct interface{}) (map[string]map[string][]string, error) {
	reqType := reflect.TypeOf(reqStruct)
	if kind := reqType.Kind(); kind != reflect.Struct {
		return nil, fmt.Errorf("Expected struct, got %s", kind)
	}

	tagMap := map[string]map[string][]string{}
	for i := 0; i < reqType.NumField(); i++ {
		opts := map[string][]string{}

		field := reqType.Field(i)
		if field.PkgPath != "" {
			continue // if PkgPath is set, the field is unexported: http://golang.org/pkg/reflect/#StructField
		}

		// Skip fields with no rv tag
		tag := field.Tag.Get("rv")
		if tag == "" {
			continue
		}

		kind := field.Type.Kind()
		opts["type"] = []string{kind.String()}
		for _, opt := range strings.Split(tag, " ") {
			keyVal := strings.SplitN(opt, "=", 2)
			if len(keyVal) == 1 {
				keyVal = []string{"source", keyVal[0]}
			}
			opts[keyVal[0]] = strings.Split(keyVal[1], ",")
		}

		tagMap[field.Name] = opts
	}
	return tagMap, nil
}
