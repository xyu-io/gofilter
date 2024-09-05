package gofilter

import (
	"errors"
	"reflect"
)

func AnyFind[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
		//if reflect.DeepEqual(v, item) {
		//	return true
		//}
	}

	return false
}

func AnyEqual[T comparable](arr []T, item T) (bool, error) {
	if len(arr) != 1 {
		return true, errors.New("params length must be 1")
	}

	element := arr[0]
	if element == item {
		return true, nil
	}

	return false, nil
}

func AnyLessThan[T comparable](arr []T, item T) (bool, error) {
	if len(arr) != 1 {
		return true, errors.New("params length must be 1")
	}
	element := arr[0]
	elementValue := reflect.ValueOf(element)
	itemValue := reflect.ValueOf(item)

	switch elementValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if itemValue.Int() < elementValue.Int() {
			return true, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uint64(itemValue.Uint()) < elementValue.Uint() {
			return true, nil
		}
	case reflect.Float32, reflect.Float64:
		if float64(itemValue.Float()) < elementValue.Float() {
			return true, nil
		}
	case reflect.Bool:
		return reflect.ValueOf(item).Bool(), nil
	default: // string/map/slice ...
		return false, nil
	}

	return false, nil
}

func AnyGreaterThan[T comparable](arr []T, item T) (bool, error) {
	if len(arr) != 1 {
		return true, errors.New("params length must be 1")
	}
	element := arr[0]
	elementValue := reflect.ValueOf(element)
	itemValue := reflect.ValueOf(item)
	switch elementValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if itemValue.Int() > elementValue.Int() {
			return true, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uint64(itemValue.Uint()) > elementValue.Uint() {
			return true, nil
		}
	case reflect.Float32, reflect.Float64:
		if float64(itemValue.Float()) > elementValue.Float() {
			return true, nil
		}
	case reflect.Bool:
		return reflect.ValueOf(item).Bool(), nil
	default: // string/map/slice ...
		return false, nil
	}

	return false, nil
}

func dealStructPtr(i any, ctype string) (string, any, bool) {
	var flag = false
	value := reflect.ValueOf(i)
	if value.Kind() != reflect.Ptr {
		return "", nil, false
	}

	flag = true
	value = value.Elem()

	if value.Kind() == reflect.Struct {
		typeOfValue := value.Type()
		for i := 0; i < value.NumField(); i++ {
			fieldType := typeOfValue.Field(i)
			if fieldType.Name != ctype {
				continue
			}
			return ctype, value.Field(i).Interface(), flag
		}
	}

	return "", nil, flag
}

func dealStruct(i any, ctype string) (string, any) {
	value := reflect.ValueOf(i)

	if value.Kind() == reflect.Struct {
		typeOfValue := value.Type()
		for i := 0; i < value.NumField(); i++ {
			fieldType := typeOfValue.Field(i)
			if fieldType.Name != ctype {
				continue
			}
			return ctype, value.Field(i).Interface()
		}
	}

	return "", nil
}
