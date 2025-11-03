package main

import (
	"fmt"
	"time"
)

type function func()

func main() {

	//timedRun(parseLogsInJson)
	//timedRun(analyseLogs)
	//timedRun(analyseExtractedData)
	//timedRun(starter)
	starter()
}

func timedRun(fn function) {

	now := time.Now().UnixMilli()
	fn()
	after := time.Now().UnixMilli()

	runtime := after - now
	fmt.Printf("\nTotal runtime: %v", runtime)
}
