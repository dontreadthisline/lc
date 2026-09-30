package lc_tests

import (
	"demo/lc"
	"testing"
)

// 160. 相交链表
func TestGetIntersectionNode(t *testing.T) {
	t.Run("官方用例1_中途相交", func(t *testing.T) {
		a, b, want := buildIntersect([]int{4, 1}, []int{5, 6, 1}, []int{8, 4, 5})
		if got := lc.GetIntersectionNodeHashWay(a, b); got != want {
			t.Errorf("相交节点 = %v, want %v", got, want)
		}
	})
	t.Run("官方用例2_不相交", func(t *testing.T) {
		a, b, _ := buildIntersect([]int{2, 6, 4}, []int{1, 5}, nil)
		if got := lc.GetIntersectionNodeHashWay(a, b); got != nil {
			t.Errorf("不相交应返回 nil, got %v", got)
		}
	})
	t.Run("官方用例3_前段长度不等", func(t *testing.T) {
		a, b, want := buildIntersect([]int{1, 9, 1}, []int{3}, []int{2, 4})
		if got := lc.GetIntersectionNodeHashWay(a, b); got != want {
			t.Errorf("相交节点 = %v, want %v", got, want)
		}
	})
	t.Run("公共段起点即交点", func(t *testing.T) {
		a, b, want := buildIntersect([]int{1}, []int{4, 5}, []int{2, 3})
		if got := lc.GetIntersectionNodeHashWay(a, b); got != want {
			t.Errorf("相交节点 = %v, want %v", got, want)
		}
	})
	t.Run("一个为空", func(t *testing.T) {
		a, b, _ := buildIntersect(nil, []int{1}, nil)
		if got := lc.GetIntersectionNodeHashWay(a, b); got != nil {
			t.Errorf("空表参与应返回 nil, got %v", got)
		}
	})
	t.Run("两表均空", func(t *testing.T) {
		if got := lc.GetIntersectionNodeHashWay(nil, nil); got != nil {
			t.Errorf("均空应返回 nil, got %v", got)
		}
	})
}
