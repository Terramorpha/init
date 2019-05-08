package main

import (
	"fmt"
	"runtime"
)

func getLine() string {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", f, l)
}
