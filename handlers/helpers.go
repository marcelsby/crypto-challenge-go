package handlers

import "reflect"

// Just declare your handler functions as h struct public member functions,
// and when called this function will mount the routes.
func mountHandlers(h interface{}) {
	v := reflect.ValueOf(h)
	t := v.Type()

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)

		if method.PkgPath == "" {
			methodValue := v.MethodByName(method.Name)
			if methodValue.Kind() == reflect.Func {
				methodValue.Call(nil)
			}
		}
	}
}
