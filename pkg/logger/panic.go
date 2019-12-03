package logger

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

func AnalyseFunc(fn interface{}) {
	//Reflection type of the underlying data of the interface
	x := reflect.TypeOf(fn)
	xx := reflect.ValueOf(fn)
	numIn := x.NumIn()   //Count inbound parameters
	numOut := x.NumOut() //Count outbounding parameters
	fmt.Println("Method:", runtime.FuncForPC(xx.Pointer()).Name())
	fmt.Println("Variadic:", x.IsVariadic()) // Used (<type> ...) ?
	fmt.Println("Package:", x.PkgPath())

	for i := 0; i < numIn; i++ {
		inV := x.In(i)
		vv := reflect.ValueOf(inV)
		in_Kind := inV.Kind() //func
		fmt.Printf("\nParameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\nValue: %v\n-----------", in_Kind, inV.Name(), vv.Interface())

	}
	for o := 0; o < numOut; o++ {

		returnV := x.Out(0)
		return_Kind := returnV.Kind()
		fmt.Printf("\nParameter OUT: "+strconv.Itoa(o)+"\nKind: %v\nName: %v\n", return_Kind, returnV.Name())
	}
}
