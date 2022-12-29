package debug

import (
	"fmt"
	"reflect"
)

func DumpStruct(obj any) {
	fmt.Println("=================================================")
	fmt.Println(reflect.TypeOf(obj).String())
	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf(">> %s %s = %#v\n",
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	fmt.Println("=================================================")
}
