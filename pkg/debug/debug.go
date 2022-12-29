package debug

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func DumpStruct(obj any) {
	if b, err := json.MarshalIndent(obj, "", "  "); err != nil {
		fmt.Println("Cannot print", reflect.TypeOf(obj).String())
	} else {
		fmt.Println(reflect.TypeOf(obj).String(), string(b))
	}
}
