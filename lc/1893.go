package lc

/*

给你区间列表 ranges，ranges[i] = [start_i, end_i] 覆盖 [start_i, end_i] 内所有整数。
再给整数 left 和 right，判断 [left, right] 内每个整数是否都被至少一个区间覆盖。

**示例 1：**
```
输入：ranges = [[1,3],[2,6]], left = 2, right = 5
输出：true
```
**示例 2：**
```
输入：ranges = [[1,2],[5,6]], left = 2, right = 5
输出：false
解释：3 和 4 未被覆盖。
```
**提示：**

- 1 <= ranges.length <= 50
- 0 <= start_i <= end_i <= 50
- 1 <= left <= right <= 50
*/
//这道题目翻译成数学表达式就是给定压缩后的差分数组,求是否存在 f(i) <= 0
//

func IsCovered(ranges [][]int, left int, right int) bool {
	diff ,val:= make([]int, 52),0

	for _, rg := range ranges {
		diff[rg[0]] += 1
		diff[rg[1]+1] -= 1
	}

	for i := 0; i <= right; i++ {
		val += diff[i]
		if left <= i && val <= 0 {
			return false
		}
	}

	return true
}

func IsCoveredBruteForce(ranges [][]int, left int, right int) bool {
	for i := left; i <= right; i++ {
		exists := false
		for _, rg := range ranges {
			if rg[0] <= i && i <= rg[1] {
				exists = true
				break
			}
		}
		if !exists {
			return false
		}
	}
	return true
}
