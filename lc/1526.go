package lc

/*

给你整数数组 target 和初始全 0 的同长数组 arr。一次操作可选 arr 的任意连续子数组，将其每个元素加 1。返回使 arr 变为 target 的最少操作次数。

**示例 1：**
```
输入：target = [1,2,3,2,1] -> [0,1,2,1,0] -> [0,0,1,0,0] -> [0,0,0,0,0]
输出：3
解释：对 [0,4] 加 1、对 [1,3] 加 1、对 [2,2] 加 1。
```
**示例 2：**
```
输入：target = [3,1,1,2] -> [2,0,0,1] -> [1,0,0,1] -> [0,0,0,1] -> [0,0,0,0]
输出：4
```
**提示：**

- 1 <= target.length <= 10^5
- 1 <= target[i] <= 10^5
*/
/*
* 每次都模拟可以么？ 复杂度o(n) * max(nums[i]) ~ o(n^2)
 */
func MinNumberOperations(target []int) int {
	n := len(target)
	diff := make([]int,n)
	for i := 0; i < n; i++ {
		if i == 0 {
			diff[i] = target[i]
		} else {
			diff[i] = target[i] - target[i-1]
		}
	}
	res := 0
	for _,delta := range(diff) {
		if delta > 0 {
			res += delta
		}
	}
	return res
}
