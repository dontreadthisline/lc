package main

import "math"


func main() {
	/*
*     binary_answer1_test.go:83: MinSpeedOnTime([1 3 2], 1.9) = 9223372036854775806, want -1
    binary_answer1_test.go:83: MinSpeedOnTime([1 1 100000], 2.01) = 10000000, want -1
	*/
	nums := []float64{1.1,1.4,1.5,1.9}
	println(math.MaxInt)
	for _,num := range nums {
		println(int(num))
	}
}
