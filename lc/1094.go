package lc

/*

车辆容量为 capacity 个座位。给定行程数组 trips，其中 trips[i] = [numPassengers_i, from_i, to_i] 表示第 i 组乘客从 from_i 上车、to_i 下车。判断能否一次性接送所有乘客（任意时刻车上人数不能超过 capacity）。

**示例 1：**
```
输入：trips = [[2,1,5],[3,3,7]], capacity = 4
输出：false
```
**示例 2：**
```
输入：trips = [[2,1,5],[3,3,7]], capacity = 5
输出：true
```
**提示：**

- 1 <= trips.length <= 1000
- 1 <= numPassengers <= 100
- 0 <= from < to <= 1000
- 1 <= capacity <= 10^5
*/

func CarPooling(trips [][]int, capacity int) bool {
	n := 0
	for _, trip := range trips {
		n = max(n, trip[2])
	}
	diff := make([]int, n+2)

	for _, trip := range trips {
		diff[trip[1]] += trip[0]
		diff[trip[2]] -= trip[0]
	}

	var val int
	for i := 0; i <= n; i++ {
		val += diff[i]
		if val > capacity {
			return false
		}
	}
	return true
}
