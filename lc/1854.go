package lc

/*

给你二维整数数组 logs，logs[i] = [birth_i, death_i] 表示第 i 个人的出生和死亡年份。年份 x 的人口等于满足 birth <= x < death 的人数（死亡当年不算存活）。返回人口最多且最早的年份。

**示例 1：**
```
输入：logs = [[1993,1999],[2000,2010]]
输出：1993
```
**示例 2：**
```
输入：logs = [[1950,1961],[1960,1971],[1970,1981]]
输出：1960
解释：1960 年两人同时存活。
```
**提示：**

- 1 <= logs.length <= 100
- 1950 <= birth_i < death_i <= 2050
*/

func MaximumPopulation(logs [][]int) int {
	diff, base := make([]int, 101), 1950

	for _, log := range logs {
		diff[log[0]-base] += 1
		diff[log[1]-base] -= 1
	}

	max, val, year := 0, 0, 0
	for i := range 100 {
		val += diff[i]
		if val > max {
			max = val
			year = i + base
		}
	}

	return year
}

func MaximumPopulationBruteForce(logs [][]int) int {
	max, year := 0, 0
	var sum int
	for i := 1950; i <= 2050; i++ {
		sum = 0
		//一样的,都是把内层循环转成差分数组
		for _, log := range logs {
			if log[0] <= i && i < log[1] {
				sum += 1
			}
		}
		if sum > max {
			max = sum
			year = i
		}
	}
	return year
}
