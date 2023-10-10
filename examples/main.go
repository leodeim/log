package main

import (
	"fmt"
	"reflect"
	"runtime"
)

var examples = []func(){
	defaultLogger,
	customLevel,
	childLogger,
	jsonLogger,
	simpleTextLogger,
	nonBlocking,
}

func main() {
	for _, example := range examples {
		fmt.Println(runtime.FuncForPC(reflect.ValueOf(example).Pointer()).Name())
		example()
		fmt.Println("")
	}
}
