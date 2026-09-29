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
	var nums []int
	nums = nil
	fmt.Println(nums)
	nums = append(nums,1)
	fmt.Println(nums)
}
