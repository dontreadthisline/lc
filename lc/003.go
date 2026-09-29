package lc

func NoRepeateLongestSubStr(s string) int {
	return noRepeateLongestSubStr(s)
}

/*
思考,为什么能想到滑动窗口？说实话,只是凭着直觉
因为是这种字串性质的，天然适合滑动窗口
而且数量又是10^5,只能用o(n)/o(nlog(n))的方法
*/

func noRepeateLongestSubStr(s string) int {
	i, j, n := 0, 0, len(s)
	m := make(map[byte]int, n)
	res := 0
	for ; i < n && j < n; j++ {
		//update the window res
		c := s[j]
		m[c] += 1
		for i < j && m[c] > 1 {
			m[s[i]] -= 1
			i += 1
		}
		res = max(res, j-i+1)
	}
	return res
}

// 209. 长度最小的子数组
// 题面: ../mother-problems/0209-长度最小的子数组.md
// 滑窗标准入门：和 ≥ target 的最短连续子数组

func MinSubArrayLen(target int, nums []int) int {
	i, j, n := 0, 0, len(nums)
	res, s := n+1, 0
	for ; j < n; j++ {
		s += nums[j]                       //无条件的扩充
		for ; i <= j && s >= target; i++ { //一旦符合了,就尝试收缩
			res = min(res, j-i+1)
			s -= nums[i]
		}
	}

	if res == n+1 {
		return 0
	}

	return res
}

func atMostSum(nums []int, goal int) int {
	if goal < 0 {
		return 0
	}
	i, j, n := 0, 0, len(nums)
	res, sum := 0, 0
	for ; i < n && j < n; j++ {
		sum += nums[j]
		for ; i <= j && sum > goal; i++ {
			sum -= nums[i]
		}
		//sum <= goal where i <= k <= j
		res += (j - i + 1)
	}
	return res
}

// 713. 乘积小于 K 的子数组
// 题面: ../mother-problems/0713-乘积小于K的子数组.md
// 链接: https://leetcode.cn/problems/subarray-product-less-than-k/

func NumSubarrayProductLessThanK(nums []int, k int) int {
	res, product := 0, 1
	i, j, n := 0, 0, len(nums)
	for ; i < n && j < n; j++ {
		product *= nums[j]
		for ; i <= j && product >= k; i++ {
			product /= nums[i]
		}
		res += (j - i + 1) //{f(n)=f(n-1) + 1,f(1) =1}=> f(n) = n => j - i + 1
	}
	return res
}

// 930. 和相同的二元子数组
// 题面: ../mother-problems/0930-和相同的二元子数组.md
// 链接: https://leetcode.cn/problems/binary-subarrays-with-sum/
// 备注,这个题目,不需要强约束为nums[i] in set {0,1} 只要规定nums[i] >= 0 题目就完全能做,和713完全一样的约束

func NumSubarraysWithSum(nums []int, goal int) int {
	return atMostSum(nums, goal) - atMostSum(nums, goal-1) //eqaul goal
}

// 438. 找到字符串中所有字母异位词
// 题面: ../mother-problems/0438-找到字符串中所有字母异位词.md
// 固定窗口 + 计数比较（p 的签名 vs 窗口签名）

func FindAnagrams(s string, p string) []int {
	table1, table2 := [26]int{}, [26]int{}
	for _, c := range p {
		table2[c-'a'] += 1
	}
	anagrams := func(arr1, arr2 [26]int) bool {
		for i, num := range arr1 {
			if num != arr2[i] {
				return false
			}
		}
		return true
	}

	res := make([]int, 0)
	for i, j := 0, 0; j < len(s); j++ {
		table1[s[j]-'a'] += 1
		if j < len(p)-1 {
			continue
		}
		if anagrams(table1, table2) {
			res = append(res, i)
		}
		table1[s[i]-'a'] -= 1 //移除i
		i += 1
	}
	return res
}

// 904. 水果成篮
// 题面: ../mother-problems/0904-水果成篮.md
// 最多含两种元素的最长连续段：滑窗 + 计数 map（k=2 的泛化形态）

func TotalFruit(fruits []int) int {
	i, j, n := 0, 0, len(fruits)
	res := 0
	m := make(map[int]int, n)
	for ; i < n && j < n; j++ {
		m[fruits[j]] += 1
		for ; i <= j && len(m) > 2; i++ {
			m[fruits[i]] -= 1
			if m[fruits[i]] <= 0 {
				delete(m, fruits[i])
			}
		}
		res = max(res, j-i+1)
	}
	return res
}

// 424. 替换后的最长重复字符
// 题面: ../mother-problems/0424-替换后的最长重复字符.md
// 滑窗（流派 A：求最长合法窗口）：不合法 → 左缩；合法 → 更新答案
// 判定：窗口内替换次数 = 窗口长度 - 最大出现次数 ≤ k
// maxCnt 偏大不影响正确性：res 只在扩张时可能创新高，maxCnt 在扩张时总是新鲜

func CharacterReplacement(s string, k int) int {
	m := make(map[byte]int, 26)
	i, j, n := 0, 0, len(s)
	res, maxCnt := 0, 0
	for ; i < n && j < n; j++ {
		m[s[j]] += 1
		maxCnt = max(maxCnt, m[s[j]])
		for ; i <= j && j-i+1-maxCnt > k; i++ {
			m[s[i]] -= 1
		}
		res = max(res, j-i+1)
	}
	return res
}

// 1004. 最大连续1的个数 III
// 题面: ../mother-problems/1004-最大连续1的个数III.md
// 链接: https://leetcode.cn/problems/max-consecutive-ones-iii/

func LongestOnes(nums []int, k int) int {
	zero, res := 0, 0
	i, j, n := 0, 0, len(nums)
	for ; i < n && j < n; j++ {
		if nums[j] == 0 {
			zero++
		}
		for ; i <= j && zero > k; i++ {
			if nums[i] == 0 {
				zero--
			}
		}
		res = max(res, j-i+1)
	}
	return res
}

// 1493. 删掉一个元素以后全为 1 的最长子数组
// 题面: ../mother-problems/1493-删掉一个元素以后全为1的最长子数组.md
// 链接: https://leetcode.cn/problems/longest-subarray-of-1s-after-deleting-one-element/

func LongestSubarray(nums []int) int {
	zero, res := 0, 0
	i, j, n := 0, 0, len(nums)
	for ; i < n && j < n; j++ {
		if nums[j] == 0 {
			zero++
		}
		for ; i <= j && zero > 1; i++ {
			if nums[i] == 0 {
				zero--
			}
		}
		res = max(res, j-i)
	}
	return res
}

// 395. 至少有 K 个重复字符的最长子串
// 题面: ../mother-problems/0395-至少有K个重复字符的最长子串.md
// 链接: https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/
// 流派 A（最长合法窗口）的进阶形态，骨架与 sliding_window_template.go 一致：
// 扩 → 不合法才缩 → 合法才记。
// 难点：题目原条件"每种字符部 >= k"对收缩不单调（次数不够只能等右指针，
// 左缩只会更少），不能直接当收缩条件。改为外层枚举窗口允许的字符种数 t，
// 用种数预算驱动收缩（904 的 len(m)>2 同一性质，t 换成枚举 1..26），
// 原条件降级为记录条件（less==0）。最优解的种数必在 1..26 内，不漏解。O(26n)。

func LongestSubstring(s string, k int) int {
	res := 0
	for t := 1; t <= 26; t++ {
		m := make(map[byte]int)
		less := 0 // 出现次数 < k 的字符种数
		i := 0
		for j := 0; j < len(s); j++ {
			// 扩：右端进窗，维护计数
			c := s[j]
			if m[c] == 0 {
				less++ // 新字符 0->1，必 < k（k==1 时下一行立刻抵消）
			}
			m[c]++
			if m[c] == k {
				less-- // 攒够 k 次，移出 less
			}
			// 缩：不合法（种数超预算）才缩，预算型条件收缩可恢复
			for len(m) > t {
				d := s[i]
				if m[d] == k {
					less++ // 从 k 线上掉下来，进入 less
				}
				m[d]--
				if m[d] == 0 {
					delete(m, d)
					less-- // 彻底离窗，移出 less
				}
				i++
			}
			// 合法才记：种数 <= t 已由上面的缩保证，这里只查原条件
			if less == 0 {
				res = max(res, j-i+1)
			}
		}
	}
	return res
}
