package lc_tests

import (
	"fmt"
	"sort"
	"testing"
)

// 测试先行的公共辅助: 被测函数是 panic 桩时不中止整个测试进程
func runSafely(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("被测函数 panic: %v (未实现?)", r)
		}
	}()
	fn()
}

// setKey: 把二维结果规范化成可比较的字符串(行内有序, 行间排序)
func setKey[T any](res []T) string {
	keys := make([]string, len(res))
	for i, row := range res {
		keys[i] = fmt.Sprint(row)
	}
	sort.Strings(keys)
	return fmt.Sprint(keys)
}
