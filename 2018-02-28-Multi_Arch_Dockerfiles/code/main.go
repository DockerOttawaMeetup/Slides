package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Hello, %s/%s!\n", runtime.GOOS, runtime.GOARCH)
}
