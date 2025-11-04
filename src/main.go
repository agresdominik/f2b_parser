package main

import (
	"fmt"
	"time"
)

type function func()

func main() {
	starter()
}

func timedRun(fn function) {

	now := time.Now().UnixMilli()
	fn()
	after := time.Now().UnixMilli()

	runtime := after - now
	fmt.Printf("\nTotal runtime: %v", runtime)
}
