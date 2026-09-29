package lc

import (
	"math"
	"slices"
)

func FindMedian(nums1, nums2 []int) float64 {
	return FindMedianMergeSortAlgs(nums1, nums2)
}

func FindMedianMergeSortAlgs(nums1, nums2 []int) float64 {
	m, n := len(nums1), len(nums2)
	nums := make([]int, 0, m+n)
	i, j := 0, 0
	for i < m || j < n {
		if i < m && j < n {
			if nums1[i] < nums2[j] {
				nums = append(nums, nums1[i])
				i += 1
			} else {
				nums = append(nums, nums2[j])
				j += 1
			}
		} else if i < m {
			nums = append(nums, nums1[i])
			i += 1
		} else if j < n {
			nums = append(nums, nums2[j])
			j += 1
		}
	}
	// 1,3=> 0,1
	// 2,4=> 1,2
	mid := (m + n) / 2
	if (m+n)%2 == 1 {
		return float64(nums[mid])
	} else {
		return (float64(nums[mid]) + float64(nums[mid-1])) / 2.0
	}
}

// 如何想到二分,朴素方式,归并排序后,找出中点位置就可以. o(m+n)
// 两个数组已经有序了,能够充分利用呢?
// 已就位：swap 短数组 ✓  k=(m+n+1)/2 ✓
// 待改：r 应从 m 开始（i 可取 m，即 nums1 全部进左半）；空循环里填三分支：
//
//	nums1[i-1] > nums2[j] → r = i-1（i 太大）
//	nums2[j-1] > nums1[i] → l = i+1（i 太小）
//	两条都不违反 → 合法切分，越界处用 math.MinInt/MaxInt 当哨兵
//
// 锚点问题（先想这个，再看下面的分支）：
//
//	"我猜从 nums1 左边拿 i 个，这个猜测是拿多了还是拿少了？"
//	每轮回答这一个二选一问题，就能把 i 的候选区间砍掉一半。
//
// 判定所用的信息（由单调性保证，见 docs/leetcode-8020-mother-problems.md 思维3）：
//
//	nums1[i-1] > nums2[j]  → 拿多了 → 砍掉右半候选 (r = i-1)
//	nums2[j-1] > nums1[i]  → 拿少了 → 砍掉左半候选 (l = i+1)
//	两条都不违反 → 切分合法，中位数就在切分线两侧的边界上
//
// 锚点问题：我猜左半 = nums1 的前 i 个 + nums2 的前 k-i 个——
// 这个 i 是"拿多了"还是"拿少了"？
// 判定只有两个交叉检查（组内天然有序，不用查）：
//
//	nums1[i-1] > nums2[j] → nums1 拿多了（它比右半最小的还大）→ hi = i-1
//	nums2[j-1] > nums1[i] → nums1 拿少了（nums2 的第 j 个替它顶在前头）→ lo = i+1
//
// 两个都不违反 → 恰好：左半恰是全局最小的 k 个，中位数从切分线两侧读出
func findMedianBinarySearchAlgs(nums1, nums2 []int) float64 {
	m, n := len(nums1), len(nums2)
	if m > n {
		m, n = n, m
		nums1, nums2 = nums2, nums1
	}
	k := (m + n + 1) / 2 // 左半的大小（奇数时左半多一个）

	lo, hi := 0, m
	for lo <= hi {
		i := lo + (hi-lo)/2
		j := k - i
		// 哨兵：贡献 0 个的一侧，左边界取 -∞（不挡路），右边界取 +∞（不通过）
		L1, R1 := math.MinInt, math.MaxInt
		if i > 0 {
			L1 = nums1[i-1]
		}
		if i < m {
			R1 = nums1[i]
		}
		L2, R2 := math.MinInt, math.MaxInt
		if j > 0 {
			L2 = nums2[j-1]
		}
		if j < n {
			R2 = nums2[j]
		}
		if L1 > R2 {
			hi = i - 1 // nums1 拿多了
		} else if L2 > R1 {
			lo = i + 1 // nums1 拿少了
		} else {
			if (m+n)%2 == 0 {
				return (float64(max(L1, L2)) + float64(min(R1, R2))) / 2
			}
			return float64(max(L1, L2))
		}
	}
	return 0 // 理论不可达：合法切分必存在于 [0, m]
}

/*
* 二分查找 —— 裸模板，一切二分的祖先

识别信号：有序数组 + 找一个值。凭什么是二分？因为每比较一次，都能安全地扔掉一半。

模板：闭区间 [l, r]，不变量只有一句话——"答案若存在，必在 [l, r] 里"。
三个分支各自砍掉一半，区间缩到空（l > r）还没命中，就是不存在。

给定一个 n 个元素有序的（升序）整型数组 nums 和一个目标值 target  ，写一个函数搜索 nums 中的 target，如果 target 存在返回下标，否则返回 -1。

你必须编写一个具有 O(log n) 时间复杂度的算法。

示例 1:
输入: nums = [-1,0,3,5,9,12], target = 9
输出: 4
解释: 9 出现在 nums 中并且下标为 4
示例 2:
输入: nums = [-1,0,3,5,9,12], target = 2
输出: -1
解释: 2 不存在 nums 中因此返回 -1
提示：

- 你可以假设 nums 中的所有元素是不重复的。

- n 将在 [1, 10000]之间。

- nums 的每个元素都将在 [-9999, 9999]之间。
*/

func BinarySearch(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := (l + r) / 2
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}

/*
34. 排序数组中查找元素的第一个和最后一个位置

704 教你"找到一个就停"，这题偏不让你停：命中后可能左右还站着一样的值。
新知识点：把"找第一个 == target"翻译成"找第一个 >= target"（左边界 lowerBound）。
这是二分的第二层——查的不是"某个值在不在"，是"边界在哪"。

妙手：右边界 = lowerBound(nums, target+1) - 1。
"最后一个 8"就是"第一个 9 的前一个"，一个原版函数跑两遍，右边界免费拿到。

给你一个按照非递减顺序排列的整数数组 nums，和一个目标值 target。请你找出给定目标值在数组中的开始位置和结束位置。

如果数组中不存在目标值 target，返回 [-1, -1]。

你必须设计并实现时间复杂度为 O(log n) 的算法解决此问题。

示例 1：
输入：nums = [5,7,7,8,8,10], target = 8
输出：[3,4]
示例 2：
输入：nums = [5,7,7,8,8,10], target = 6
输出：[-1,-1]
示例 3：
输入：nums = [], target = 0
输出：[-1,-1]
提示：

- 0 <= nums.length <= 10^5

- -10^9 <= nums[i] <= 10^9

- nums 是一个非递减数组

- -10^9 <= target <= 10^9
*/

func SearchRange(nums []int, target int) []int {
	first := lowerBound(nums, target)
	if first >= len(nums) || nums[first] != target {
		return []int{-1, -1}
	}
	last := lowerBound(nums, target+1) - 1
	return []int{first, last}
}

// lowerBound：第一个 >= x 的下标（全都 < x 则返回 len(nums)）
// 与 704 的三点分支不同，这里是【二选一 + 收敛到一点】的固定搭配：
//
//	for l < r（区间缩到恰好 1 个元素时循环结束）
//	mid 取左中点（l+(r-l)/2 向下取整，保证 r=mid 时区间一定在缩）
//	nums[mid] >= x → mid 可能就是答案，不能丢 → r = mid
//	nums[mid] <  x → mid 绝不是答案，放心丢 → l = mid + 1
//
// 对照 004：最后那个 TODO 循环里 r=i-1 / l=i+1 的三分支，本质是
// "704 的命中分支"换成了"两个交叉检查"，骨架完全同源。
func lowerBound(nums []int, x int) int {
	l, r := 0, len(nums)
	// 为什么r是len(nums),而不是len(nums -1) 这里是为了能够表达找不到的情况,实际上也可以用len(nums-1)
	// 那得需要额外用-1来表达，还得加判断。所以这里为了统一,用len(nums)来表达不满足条件的情况
	for l < r { // r = mid收缩的情况下,退出条件只能是l < r 不然会陷入死循环
		mid := l + (r-l)/2
		if nums[mid] >= x { // 俩分支合并
			r = mid
		} else {
			l = mid + 1
		}
	} // break的时候,l == r 所以返回l和r都行,只是惯例返回l
	return l
}

/*
33. 搜索旋转排序数组

数组本来有序，被从某处拧了一圈：[0,1,2,4,5,6,7] → [4,5,6,7,0,1,2]
新知识点：整体无序，但 mid 一刀切下去，两半里必有一半是有序的。
（左半无序 ⇔ 断点在左半 ⇔ 右半有序，断点只有一个，跑不了。）

策略：每轮先认出有序的那半，再看 target 的值落不落在它的值域里——
落在 → 进那半（正常 704 逻辑）；不落在 → 只可能在另一半。
每次仍然稳定扔掉一半，log 轮收工。

整数数组 nums 按升序排列，数组中的值 互不相同 。

在传递给函数之前，nums 在预先未知的某个下标 k（0 <= k < nums.length）上进行了 向左旋转，使数组变为 [nums[k], nums[k+1], ..., nums[n-1], nums[0], nums[1], ..., nums[k-1]]（下标 从 0 开始 计数）。例如， [0,1,2,4,5,6,7] 下标 3 上向左旋转后可能变为 [4,5,6,7,0,1,2] 。

给你 旋转后 的数组 nums 和一个整数 target ，如果 nums 中存在这个目标值 target ，则返回它的下标，否则返回 -1 。

你必须设计一个时间复杂度为 O(log n) 的算法解决此问题。

示例 1：
输入：nums = [4,5,6,7,0,1,2], target = 0
输出：4
示例 2：
输入：nums = [4,5,6,7,0,1,2], target = 3
输出：-1
示例 3：
输入：nums = [1], target = 0
输出：-1
提示：

- 1 <= nums.length <= 5000

- -10^4 <= nums[i] <= 10^4

- nums 中的每个值都 独一无二

- 题目数据保证 nums 在预先未知的某个下标上进行了旋转

- -10^4 <= target <= 10^4
*/

func SearchRotated(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] == target {
			return mid
		}
		// 左半有序
		if nums[mid] >= nums[l] {
			// 看看是不是在左半
			if nums[l] <= target && target < nums[mid] {
				r = mid - 1
			} else {
				l = mid + 1
			}
		} else {
			// 真在右半
			if nums[mid] < target && target <= nums[r] {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}
	}
	return -1
}

// ============================================================
// 二分答案台阶：004 的前置练习（从这三道爬上去，再回头看切分法）
// 和 704/34/33 的区别：搜索空间不再由数组下标现成提供，
// 而是你自己圈出来的"答案值域" + 自己写的单调判定函数。
// ============================================================
//
// ---------- 台阶 1：答案值域 [0, x]，判定最直白 ----------
//
// 我的写法（对照你的三分支版）：
// 锚点问题：给定 mid，能否确认答案在它左边还是右边？
// 把答案定义为"最大的 k 使 k*k <= x"，等价于"第一个 k*k > x 的 k 的前驱"——
// 于是谓词只剩一个：mid*mid > x。窗口被它切成 [可行 | 不可行]，收敛式砍半。
//
// 两个细节：
//
//	r = x+1 是哨兵——保证区间里必然存在不可行的 k（x+1 一定超），呼应内置
//	max 签名里那个裸 x：都是"编译期/构造期保证区间非空、答案良定义"的思路
//	mid 取左中点即可——因为收缩形态是 r=mid / l=mid+1，和 34、sort.Search 同款，
//	不需要右中点（那是 l=mid 收缩形态才需要的防死循环搭配）

func SqrtSearchAlgs(x int) int {
	l, r := 0, x+1
	for l < r {
		mid := l + (r-l)/2
		if mid*mid > x { // 谓词：mid 不可行。mid ≤ 2^31，mid² ≤ 4.6e18，int64 下无溢出
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l - 1 // l 是第一个不可行的 k，前驱即最大可行 k
}

/*
给你一个非负整数 x ，计算并返回 x 的 算术平方根 。

由于返回类型是整数，结果只保留 整数部分 ，小数部分将被 舍去 。

注意：不允许使用任何内置指数函数和算符，例如 pow(x, 0.5) 或者 x  0.5 。

示例 1：
输入：x = 4
输出：2
示例 2：
输入：x = 8
输出：2
解释：8 的算术平方根是 2.82842..., 由于返回类型是整数，小数部分将被舍去。
提示：

- 0 <= x <= 2^31 - 1
r * r <= x && (r + 1) ^2 > x
*/

func Sqrt(x int64) int64 {
	l, r := int64(0), x
	for l <= r {
		mid := l + (r-l)/2
		a := mid * mid
		if a <= x && (mid+1)*(mid+1) > x {
			return mid
		} else if a > x {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return l
}

func SqrtNetonAlgs(x float64) float64 { // 牛顿法
	return 0.0
}

func SqrtFastAlgs(x float64) float64 { // 雷神之锤的那种快速解法(含近似)
	return 0.0
}

func SqetIterAlgs(x float64) float64 { // 迭代法

	return 0.0
}

/*
	 ---------- 台阶 2：值域 [1, max]，判定函数要自己写 ----------

	 珂珂喜欢吃香蕉。这里有 n 堆香蕉，第 i 堆中有 piles[i] 根香蕉。警卫已经离开了，将在 h 小时后回来。

	 珂珂可以决定她吃香蕉的速度 k （单位：根/小时）。每个小时，她将会选择一堆香蕉，从中吃掉 k 根。如果这堆香蕉少于 k 根，她将吃掉这堆的所有香蕉，然后这一小时内不会再吃更多的香蕉。

	 珂珂喜欢慢慢吃，但仍然想在警卫回来前吃掉所有的香蕉。

	 返回她可以在 h 小时内吃掉所有香蕉的最小速度 k（k 为整数）。

	 示例 1：
	 输入：piles = [3,6,7,11], h = 8

		[1,2,2,3]

	 输出：4
	 示例 2：
	 输入：piles = [30,11,23,4,20], h = 5
	 [1,1,1,1,1]
	 输出：30
	 示例 3：
	 输入：piles = [30,11,23,4,20], h = 6
	 输出：23
	 提示：

	 - 1 <= piles.length <= 10^4

	 - piles.length <= h <= 10^9

	 - 1 <= piles[i] <= 10^9
	 这个该死的题目我想起来了,当时就被搞的糊里糊涂的
	 思考,值域一定[1,max(piles)]
*/

func EtaFuckBananas(piles []int, h int) int {
	l, r := 1, slices.Max(piles)
	meet := func(k int) bool {
		sum := 0
		for _, pile := range piles {
			sum += (pile + k - 1) / k
		}
		return sum <= h
	}
	for l < r {
		mid := l + (r-l)/2
		if meet(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

/*
 ---------- 台阶 3：同款骨架换判定，练到肌肉记忆 ----------

 传送带上的包裹必须在 days 天内从一个港口运送到另一个港口。

 传送带上的第 i 个包裹的重量为 weights[i]。每一天，我们都会按给出重量（weights）的顺序往传送带上装载包裹。我们装载的重量不会超过船的最大运载重量。

 返回能在 days 天内将传送带上的所有包裹送达的船的最低运载能力。

 示例 1：
 输入：weights = [1,2,3,4,5,6,7,8,9,10], days = 5
 输出：15
 解释：
 船舶最低载重 15 就能够在 5 天内送达所有包裹，如下所示：
 第 1 天：1, 2, 3, 4, 5
 第 2 天：6, 7
 第 3 天：8
 第 4 天：9
 第 5 天：10

 请注意，货物必须按照给定的顺序装运，因此使用载重能力为 14 的船舶并将包装分成 (2, 3, 4, 5), (1, 6, 7), (8), (9), (10) 是不允许的。
 示例 2：
 输入：weights = [3,2,2,4,1,4], days = 3
 输出：6
 解释：
 船舶最低载重 6 就能够在 3 天内送达所有包裹，如下所示：
 第 1 天：3, 2
 第 2 天：2, 4
 第 3 天：1, 4
 示例 3：
 输入：weights = [1,2,3,1,1], days = 4
 输出：3
 解释：
 第 1 天：1
 第 2 天：2
 第 3 天：3
 第 4 天：1, 1
 提示：

 - 1 <= days <= weights.length <= 5 * 10^4

 - 1 <= weights[i] <= 500
*/

func Transmite(weights []int, days int) int {
	// 只能按照顺序传送
	s := 0
	for _, w := range weights {
		s += w
	}
	l, r := slices.Max(weights), s
	n := len(weights)
	meet := func(load int) bool {
		i, j := 0, 0
		cnt := 0
		for i < n && j < n {
			s := 0
			for j = i; j < n && s+weights[j] <= load; j += 1 {
				s += weights[j]
			}
			cnt += 1
			i = j
		}
		return cnt <= days
	}
	for l < r {
		mid := l + (r-l)/2
		if meet(mid) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

/*
形态三台阶：结构判定谓词（004 切分法的前置，此前缺失的两级）
这类题的判定不是"值比较"而是"结构性质"：

	162: 邻居比较（mid 和 mid+1 谁高）
	378: 计数谓词（矩阵里 <= mid 的元素个数是否 >= k）
	410: 贪心装箱谓词（1011 同款，值域更紧）

每道先问锚点：什么比较能确认答案在哪一半？
============================================================

---------- 台阶 A：162 寻找峰值（邻居比较定方向） ----------

峰值元素是指其值严格大于左右相邻值的元素。

给你一个整数数组 nums，找到峰值元素并返回其索引。数组可能包含多个峰值，在这种情况下，返回 任何一个峰值 所在位置即可。

你可以假设 nums[-1] = nums[n] = -∞ 。

你必须实现时间复杂度为 O(log n) 的算法来解决此问题。

示例 1：
输入：nums = [1,2,3,1]
输出：2
解释：3 是峰值元素，你的函数应该返回其索引 2。
示例 2：
输入：nums = [1,2,1,3,5,6,4]
输出：1 或 5
解释：你的函数可以返回索引 1，其峰值元素为 2；

	或者返回索引 5， 其峰值元素为 6。

提示：

- 1 <= nums.length <= 1000

- -2^31 <= nums[i] <= 2^31 - 1

- 对于所有有效的 i 都有 nums[i] != nums[i + 1]
*/

func SearchPeakBruteForceAlgs(nums []int) int {
	//不二分，那就暴力搜索,因为峰值一定存在,就遍历找符合条件的就是了
	n := len(nums)
	meet := func(i int) bool {
		if i == 0 {
			return nums[i] > nums[i+1]
		} else if i == n-1 {
			return nums[i] > nums[i-1]
		} else {
			return nums[i-1] < nums[i] && nums[i] > nums[i+1]
		}
	}
	for i := range n {
		if meet(i) {
			return i
		}
	}
	return -1
}

func SearchPeakBinarySearchAlgs(nums []int) int {
	n := len(nums)
	l, r := 0, n-1
	meet := func(i int) bool {
		if i == 0 {
			return nums[i] > nums[i+1]
		} else if i == n-1 {
			return nums[i] > nums[i-1]
		} else {
			return nums[i-1] < nums[i] && nums[i] > nums[i+1]
		}
	}
	//可以二分的前提:峰值一定存在,（对么）,所有的相邻元素都不相同，这给了我们什么启示？
	for l < r {
		mid := l + (r-l)/2
		if meet(mid) {
			return mid
		} else {
			if nums[mid] < nums[mid+1] { //上坡
				l = mid + 1 //因为mid 不可能是答案了,所以就用mid + 1
			} else {
				r = mid - 1 //同理,mid不可能是答案,所以mid - 1
			}
		}
	}
	return l
}

/*
* ---------- 台阶 B：378 计数谓词（004 的直系前置） ----------

	给你一个 n x n 矩阵 matrix ，其中每行和每列元素均按升序排序，找到矩阵中第 k 小的元素。

	请注意，它是 排序后 的第 k 小元素，而不是第 k 个 不同 的元素。

	你必须找到一个内存复杂度优于 O(n^2) 的解决方案。

	示例 1：
	输入：matrix = [
									[1,5,9],
									[10,11,13],
									[12,13,15]], k = 8
	输出：13
	解释：矩阵中的元素为 [1,5,9,10,11,12,13,13,15]，第 8 小元素是 13
	示例 2：
	输入：matrix = [[-5]], k = 1
	输出：-5
	提示：

	- n == matrix.length

	- n == matrix[i].length

	- 1 <= n <= 300

	- -10^9 <= matrix[i][j] <= 10^9

	- 题目数据 保证 matrix 中的所有行和列都按 非递减顺序 排列

	- 1 <= k <= n^2

	进阶：

	- 你能否用一个恒定的内存(即 O(1) 内存复杂度)来解决这个问题?

	- 你能在 O(n) 的时间复杂度下解决这个问题吗?这个方法对于面试来说可能太超前了，但是你会发现阅读这篇文章（ this paper ）很有趣。
	- 先说好这个O(n)的n是那个n? 是数据规模的n还是k? 我猜绝对是k吧
*/

func SearchMatrixSearchWay(matrix [][]int, k int) int {
	//自造值域+自定义谓词
	m, n := len(matrix), len(matrix[0])
	cnt := func(v int) int {
		num := 0
		r, c := m-1, 0
		for r >= 0 && c < n {
			if matrix[r][c] <= v {
				c += 1
				num += (r + 1)
			} else {
				r -= 1
			}
		}
		return num
	}

	lo, hi := math.MaxInt, math.MinInt
	for _, row := range matrix {
		for _, num := range row {
			if num < lo {
				lo = num
			}
			if num > hi {
				hi = num
			}
		}
	}

	l, r := lo, hi
	for l < r {
		mid := l + (r-l)/2
		if cnt(mid) >= k {
			r = mid
		} else {
			l = mid + 1
		}
	}

	return l
}

/*
 ---------- 台阶 C：410 分割数组的最大值（1011 加强，可选） ----------

 给定一个非负整数数组 nums 和一个整数 k ，你需要将这个数组分成 k 个非空的连续子数组，使得这 k 个子数组各自和的最大值 最小。

 返回分割后最小的和的最大值。

 子数组 是数组中连续的部分。

 示例 1：
 输入：nums = [7,2,5,10,8], k = 2
 输出：18
 解释：
 一共有四种方法将 nums 分割为 2 个子数组。
 其中最好的方式是将其分为 [7,2,5] 和 [10,8] 。
 因为此时这两个子数组各自的和的最大值为18，在所有情况中最小。
 示例 2：
 输入：nums = [1,2,3,4,5], k = 2
 输出：9
 示例 3：
 输入：nums = [1,4,4], k = 3
 输出：4
 提示：

 - 1 <= nums.length <= 1000

 - 0 <= nums[i] <= 10^6

 - 1 <= k <= min(50, nums.length)
*/
