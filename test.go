package main

import "fmt"

type Test struct {
}

var test Test

func (t Test) print() {
	fmt.Println("test")
}

func destroy() {
	defer test.print()
}
