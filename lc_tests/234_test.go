package lc_tests

import (
	"demo/lc"
	"fmt"
	"testing"
)


func TestMiddleOfLinkList(t *testing.T) {
	arrs := [][]int{
		{},
		{1},
		{1,2},
		{1,2,3},
		{1,2,3,4},
	}
	for _,arr := range arrs {
		head := lc.NewLinkList(arr)
		mid := lc.MiddleOfLinkList(head)
		fmt.Printf("list:%s\n",head)
		fmt.Printf("mid:%s\n",mid)
	}
}
