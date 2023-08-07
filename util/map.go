package util

import "reflect"

func ReverseMapO(m map[interface{}]interface{}) (r map[interface{}]interface{}) {
	r = make(map[interface{}]interface{}, len(m))
	for k, v := range m {
		r[v] = k
	}
	return
}

func ReverseMap(m interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		panic("m must be a map")
	}
	t := reflect.MapOf(v.Type().Elem(), v.Type().Key())
	r := reflect.MakeMap(t)
	for _, k := range v.MapKeys() {
		r.SetMapIndex(v.MapIndex(k), k)
	}
	return r.Interface()
}
