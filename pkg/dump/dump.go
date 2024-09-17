package dump

import (
	"reflect"

	"github.com/dathuynh1108/clean-arch-base/pkg/comjson"
)

func DumpDataByJSON(fromObject any, toObject any) error {
	fromJSON, err := comjson.Marshal(fromObject)
	if err != nil {
		return err
	}
	return comjson.Unmarshal(fromJSON, toObject)
}

type KeyValue struct {
	Key string
	Val any
}

func KeyValueArrayToJSONArray(data []KeyValue) (map[string]any, error) {
	m := make(map[string]any)
	for _, kv := range data {
		if existingVal, exists := m[kv.Key]; exists {
			// If the key already exists, convert the existing value to a slice if it isn't already
			switch v := existingVal.(type) {
			case []any:
				m[kv.Key] = append(v, kv.Val)
			default:
				// Stackup
				newVal := make([]any, 0, 2)
				newVal = append(newVal, existingVal)
				newVal = append(newVal, kv.Val)
				m[kv.Key] = newVal
			}
		} else {
			m[kv.Key] = []any{kv.Val}
		}
	}
	return m, nil
}

func DumpKeyValueArrayToObject(data []KeyValue, dest any) error {
	m, err := KeyValueArrayToJSONArray(data)
	if err != nil {
		return err
	}

	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		fieldValue := destValue.Field(i)

		if fieldValue.Kind() == reflect.Slice {
			// Convert non-array source fields to arrays if the destination field is an array
			fieldName := field.Tag.Get("json")
			if fieldName == "" {
				fieldName = field.Name
			}
			if srcField, ok := m[fieldName]; ok {
				if reflect.TypeOf(srcField).Kind() != reflect.Slice {
					m[fieldName] = []any{srcField}
				}
			}
		} else {
			// Convert array source fields to non-array if the destination field is not an array
			fieldName := field.Tag.Get("json")
			if fieldName == "" {
				fieldName = field.Name
			}
			if srcField, ok := m[fieldName]; ok {
				if reflect.TypeOf(srcField).Kind() == reflect.Slice {
					srcSlice := reflect.ValueOf(srcField)
					if srcSlice.Len() > 0 {
						m[fieldName] = srcSlice.Index(0).Interface()
					}
				}
			}
		}
	}

	return DumpDataByJSON(m, dest)
}
