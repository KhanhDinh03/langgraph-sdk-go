package client

import (
	"reflect"
	"strings"
)

func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		emptyStruct := reflect.New(v.Type()).Elem().Interface()
		return reflect.DeepEqual(value, emptyStruct)
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0.0
	}

	return false
}

func removeEmptyFields(data any) any {
	if data == nil {
		return nil
	}

	val := reflect.ValueOf(data)

	// Dereference pointer
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Map:
		cleanedMap := make(map[string]any)
		for _, key := range val.MapKeys() {
			v := val.MapIndex(key)
			cleanedVal := removeEmptyFields(v.Interface())
			if !isEmpty(cleanedVal) {
				cleanedMap[key.String()] = cleanedVal
			}
		}
		return cleanedMap

	case reflect.Slice, reflect.Array:
		cleanedSlice := []any{}
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			cleanedVal := removeEmptyFields(item)
			if !isEmpty(cleanedVal) {
				cleanedSlice = append(cleanedSlice, cleanedVal)
			}
		}
		return cleanedSlice

	default:
		return val.Interface()
	}
}

func containsDot(s string) bool {
	return strings.Contains(s, ".")
}
