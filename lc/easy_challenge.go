package lc

import "math"

/*
挑战背景：这 5 道题是从 leetcode.cn 算法分类 3396 道非会员题中，
按种子 20260915 随机抽出的"未被 18 家族课程标签覆盖"的简单题。
标签上它们不落入任何算法家族（只挂数组/字符串/数学等载体标签）——
预测：纯操作题，应当全部秒杀。函数体留空，写完去 lc_tests 跑对拍。

═══════════════════════════════════════════════════════════════

2908. 元素和最小的山形三元组 I
难度：简单  标签：数组
链接：https://leetcode.cn/problems/minimum-sum-of-mountain-triplets-i/

给你一个下标从 0 开始的整数数组 nums 。如果三元组 (i, j, k) 满足下述全部条件，则认为它是一个山形三元组：
  - i < j < k
  - nums[i] < nums[j] 且 nums[k] < nums[j]
请你返回 nums 中 山形三元组 的 最小元素和 。如果不存在山形三元组，返回 -1 。

示例 1：
输入：nums = [8,6,1,5,3]
输出：9
解释：三元组 (2,3,4) 是一个山形三元组，元素和等于 9 。

示例 2：
输入：nums = [5,4,8,7,10,2]
输出：13
解释：三元组 (1,3,5) 是一个山形三元组，元素和等于 13 。

示例 3：
输入：nums = [6,5,4,3,4,5]
输出：-1
解释：不存在山形三元组。

提示：
- 3 <= nums.length <= 50
- 1 <= nums[i] <= 50

═══════════════════════════════════════════════════════════════

2644. 找出可整除性得分最大的整数
难度：简单  标签：数组
链接：https://leetcode.cn/problems/find-the-maximum-divisibility-score/

divisors[i] 的 可整除性得分 等于满足 nums[j] 能被 divisors[i] 整除的下标 j 的数量。

返回 可整除性得分 最大的整数 divisors[i] 。如果有多个整数具有最大得分，则返回数值最小的一个。

给你两个正整数数组 nums 和 numsDivide 。你可以从 nums 中删除任意个元素。
你的得分定义为一个整数 d，d 满足：
  - nums 中剩余元素都能被 d 整除（你可以选择删除若干元素使得剩余元素都满足）
本题问的是：numsDivide 中可以被 d 整除的元素个数最多时的 d。
（准确题意见原题）返回可整除性得分最大的整数 d 。如果有多个整数得分相同，返回最小的 d 。
（注：以下为实现本题在测试中的实际语义版本，见函数签名）

示例 1：
输入：nums = [4,7,8,15,16,13,6,5], numsDivide = [8]
输出：8（语义详见原题与测试）

提示：
- 1 <= nums.length, numsDivide.length <= 10^5
- 1 <= nums[i], numsDivide[i] <= 10^5

═══════════════════════════════════════════════════════════════

3127. 构造相同颜色的正方形
难度：简单  标签：数组、枚举、矩阵
链接：https://leetcode.cn/problems/make-a-square-with-the-same-color/

给你一个二维 3 x 3 矩阵 grid ，其中每个元素是 'B' 或 'W'。
你最多可以修改一次格子的颜色，问是否存在一个 2 x 2 的同色正方形。

如果 grid 中存在一个 2 x 2 的正方形，其四个格子颜色全部相同，返回 true ；
否则，如果可以通过修改恰好一个格子的颜色使得某个 2 x 2 正方形同色，也返回 true ；
否则返回 false 。

示例 1：
输入：grid = [["B","W","B"],["B","W","W"],["B","W","B"]]
输出：true
解释：修改左上角 (0,0) 为 'W' 后，左上 2x2 全白。（已有 3 个 W，改 1 个即可）

提示：
- grid.length == 3
- grid[i].length == 3
- grid[i][j] 是 'W' 或 'B'

═══════════════════════════════════════════════════════════════

806. 写字符串需要的行数
难度：简单  标签：数组、字符串
链接：https://leetcode.cn/problems/number-of-lines-to-write-string/

我们要把给定的字符串 S 从左到右写到每一行上，每一行的最大宽度为 100 个单位，
同一个字体大小的每个字符宽度是固定的：widths[0] 是 'a' 的宽度，widths[1] 是 'b' 的宽度，……，widths[25] 是 'z' 的宽度。

写字符串时需要遵守如下规则：
  - 每行最多 100 个单位宽度，若某个字符导致当前行超过 100，则换行写这个字符
返回一个长度为 2 的列表：[总行数, 最后一行的宽度（单位）]。

示例 1：
输入：widths = [10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10], s = "abcdefghijklmnopqrstuvwxyz"
输出：[3, 60]
解释：前两行各写 10 个字符（100 单位），第三行写 6 个字符（60 单位）。

示例 2：
输入：widths = [4,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10], s = "bbbcccdddaaa"
输出：[2, 4]
解释：前 10 个字符（bbbcccdddaa）占 98 单位；再加第三个 a 超过 100，换行，末行 1 个 a 占 4 单位。

提示：
- lengths of widths == 26
- 2 <= widths[i] <= 10
- 1 <= s.length <= 1000
- s 只包含小写字母

═══════════════════════════════════════════════════════════════

3232. 判断是否可以赢得数字游戏
难度：简单  标签：数组、数学
链接：https://leetcode.cn/problems/find-if-digit-game-can-be-won/

给你一个正整数数组 nums 。
Alice 和 Bob 轮流进行如下选择游戏，Alice 先手：
  - 每一轮，当前玩家从 nums 中选出一个数字（数字被移除）
Alice 只能选 单个数字位 的数（1 到 9），Bob 只能选 至少两个数字位 的数（10 及以上）。
如果没有玩家能选（轮到某人但无符合条件的数字），该玩家输。
若 Alice 获胜返回 true，Bob 获胜返回 false。

示例 1：
输入：nums = [1,2,3,4,10]
输出：true
解释：Alice 先手拿走 1..4 中的任意一个；最终 Alice 拿的单 digit 数多于……（策略详见原题）

示例 2：
输入：nums = [1,2,3,4,5,14]
输出：true

示例 3：
输入：nums = [5,5,5,25]
输出：true

提示：
- 1 <= nums.length <= 100
- 1 <= nums[i] <= 10^5
*/

// ============================ 挑战函数 ============================

// MinMountainSum 2908：山形三元组最小元素和，不存在返回 -1
func MinMountainSum(nums []int) int {
	res := math.MaxInt
	n := len(nums)
	for i := range n {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; j++ {
				if nums[j] > nums[i] && nums[j] > nums[k] {
					res = min(res,nums[i]+nums[j]+nums[k])
				}
			}
		}
	}
	if res == math.MaxInt {
		return -1
	}
	return res
}

// MaxDivisibilityScoreDivisor 2644：语义按测试用例（见 lc_tests/easy_challenge_test.go 顶注）
func MaxDivisibilityScoreDivisor(nums []int, divisors []int) int {
	// TODO: 挑战 2（语义：返回能整除 nums 中最多元素的 divisors 值，并列取最小）
	res,max := 0,-1
	for _,divisor := range divisors {
		cnt := 0
		for _,num := range nums {
			if num % divisor == 0 {
				cnt++
			}
		}
		if max == -1 && cnt > 0{
			res,max = divisor,cnt
		} else if cnt >= max && divisor < res {
			res,max = divisor,cnt
		}
	}
	return res
}

// MakeSameSquare 3127：能否（至多改一格）得到同色 2x2 正方形
func MakeSameSquare(grid [][]byte) bool {
	//扫描,然后数四个框,然后数b,w的个数,任意一个<=1就可以
	return false
}

// NumberOfLines806 806：写字符串的 [行数, 末行宽度]
func NumberOfLines806(widths []int, s string) []int {
	// TODO: 挑战 4
	panic("TODO")
}

// DigitGameWin 3232：Alice 是否获胜
func DigitGameWin(nums []int) bool {
	return true
}
