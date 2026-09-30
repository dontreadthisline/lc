package main

import "fmt"

type Foo struct {
	a int
	b string
}
type Bar struct {
	a int
	b int
}

func main() {
	a := 123
	b := 2
	c := 123
	fmt.Println(a ^ c)
	fmt.Println(a  ^ b)
}
