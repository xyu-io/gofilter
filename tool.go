package gofilter

import "reflect"

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
