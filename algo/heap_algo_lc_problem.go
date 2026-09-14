package algo

import (
	"math"
)

/*
* LC 692 前K个高频单词：给单词列表 words 和整数 k，返回前 k 个出现次数最多的单词。
* 结果按频率从高到低排序；频率相同的单词按字典序排列（字典序小者在前）。
* 与 347 的关键区别：无“答案唯一”保证，同频必须 tiebreak（字典序二级排序）。
*
* 约束：1 <= words.length <= 10^5；1 <= k <= 不同单词数；单词仅含小写字母。
* 提示：小顶堆门槛法（同 countSort 的多级比较器思维）：堆里留“最好的 k 个”，
* “最差”= 频次低，或同频时字典序大（字典序大者更靠近堆顶等着被踢）。
* 注意：堆弹出的是“最差优先”，最终答案需逆序收集。

 ```
 示例: words = ["i","love","leetcode","i","love","coding"], k = 2 → ["i","love"]
       (i=2, love=2, leetcode=1, coding=1；同频的 i < love 字典序，i 在前)
       words = ["the","day","is","sunny","the","the","the","sunny","is","is"], k = 4
       → ["the","is","sunny","day"]  (the=4, is=3, sunny=2, day=1)
 ```
*/

type Pair struct {
	word string
	cnt  int
}

func MaxKFreqWords(words []string, k int) []string {
	res := make([]string, 0, k)
	m := map[string]int{}
	less := func(a, b Pair) bool {
		if a.cnt != b.cnt {
			return a.cnt < b.cnt
		}
		return a.word > b.word
	}
	for _, word := range words {
		m[word] += 1
	}
	h := NewHeap(less)
	for word, cnt := range m {
		if h.Len() == k {
			min := h.Peek()
			if cnt < min.cnt || (cnt == min.cnt && word >= min.word) {
				continue
			}
			h.Update(0, Pair{
				word: word,
				cnt:  cnt,
			})
		} else {
			h.Push(Pair{
				word: word,
				cnt:  cnt,
			})
		}
	}
	for h.Len() > 0 {
		res = append(res, h.Pop().word)
	}
	for i, j := 0, k-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

/*
* LC 703 数据流中的第 K 大元素：设计一个类，初始化传整数数组 nums 和整数 k，
* add(val) 每次向流中加入一个数并返回当前第 k 大的元素。

 ```
   示例: nums = [4,5,8,2], k = 3
         add(3) → 4    add(5) → 5    add(10) → 5    add(9) → 8
         add(4) → 8    （堆内始终保留最大的 3 个，堆顶是其中最小 = 第 3 大）
 ```
*/

type KthLargest struct {
	h *Heap[int]
	k int
}

func NewKthLargest(k int, nums []int) *KthLargest {
	kl := KthLargest{
		h: NewHeapFrom(nums[:min(k, len(nums))], func(a int, b int) bool {
			return a < b
		}),
		k: k,
	}

	for i := k; i < len(nums); i++ {
		kl.add(nums[i])
	}
	return &kl
}

func (kl *KthLargest) Add(val int) int {
	return kl.add(val)
}

func (kl *KthLargest) add(val int) int {
	if kl.h.Len() < kl.k {
		kl.h.Push(val)
	} else if val > kl.h.Peek() {
		kl.h.Update(0, val)
	}
	return kl.h.Peek()
}

/*
* LC 1046 最后一块石头的重量：每轮选最重的两块石头 x >= y 粉碎：
* x == y → 两块都碎；x != y → 剩下重量 x - y。返回最后剩下石头的重量（没有则 0）。

 ```
   示例: stones = [2,7,4,1,8,1] → 1
         (7,8 碰 → 剩 1；1,1 碎光；2,4 碰 → 剩 2；1,2 碰 → 剩 1)
 ```
*/
//吐槽一下,你给的这个模拟实例,就是一坨屎,根本不对

func LastStoneWeight(stones []int) int {
	h := NewHeapFrom(stones, func(a, b int) bool {
		return a > b
	})

	for h.Len() > 1 {
		a, b := h.Pop(), h.Pop()
		if a != b {
			h.Push(max(a, b) - min(a, b))
		}
	}

	if h.Len() == 0 {
		return 0
	}
	return h.Pop()
}

/*
* LC 215 数组中的第 K 个最大元素：返回排序后第 k 个最大元素（注意不是第 k 个不同的）。
* 要求 O(n)。思路：门槛堆 O(n log k) 可过；理论最优是快速选择（algo/partation.go 已有实现可对照）。

 ```
   示例: nums = [3,2,1,5,6,4], k = 2 → 5
         nums = [3,2,3,1,2,4,5,5,6], k = 4 → 4
 ```
*/

func FindKthLargestHeap(nums []int, k int) int {
	//我会用partation的方法,用堆来解决吧
	k = min(k, len(nums))
	h := NewHeapFrom(nums[:k], func(a, b int) bool {
		return a < b
	})

	for i := k; i < len(nums); i++ {
		if nums[i] > h.Peek() {
			h.Update(0, nums[i])
		}
	}
	return h.Pop()
}

/*
* LC 973 最接近原点的 K 个点：平面上有 points 点，返回离原点 (0,0) 最近的 k 个点。
* 答案任意顺序。距离按欧几里得距离（不开根号直接比平方，省精度省计算）。
 ```
   示例: points = [[1,3],[-2,2]], k = 1 → [[-2,2]]  (距离平方 4 < 10)
         points = [[3,3],[5,-1],[-2,4]], k = 2 → [[3,3],[-2,4]]
 ```
*/

func KClosest(points [][]int, k int) [][]int {
	dist := func(point []int) int {
		return point[0]*point[0] + point[1]*point[1]
	}
	k = min(k, len(points))
	h := NewHeapFrom(points[:k], func(a, b []int) bool {
		return dist(a) > dist(b)
	})
	for i := k; i < len(points); i++ {
		if dist(points[i]) < dist(h.Peek()) {
			h.Update(0, points[i])
		}
	}
	res := make([][]int, 0, k)
	for h.Len() > 0 {
		res = append(res, h.Pop())
	}
	return res
}

/*
* LC 373 查找和最小的 K 对数字：两个升序数组 nums1、nums2，//如何利用这个有序的性质?
* 返回和最小的前 k 对 (u, v)，按 (u1+v1) <= (u2+v2) 任意顺序。
 ```
   示例: nums1 = [1,7,11],
         nums2 = [2,4,6] , k = 3
         → [[1,2],[1,4],[1,6]]  (和为 3,5,7)
 ```
*/

type pair struct {
	i int
	j int
}

func KSmallestPairs(nums1, nums2 []int, k int) [][]int {
	less := func(p1, p2 pair) bool {
		return nums1[p1.i]+nums2[p1.j] < nums1[p2.i]+nums2[p2.j]
	}
	m, n := len(nums1), len(nums2)
	pairs := make([]pair, 0, n)
	for i := range m {
		pairs = append(pairs, pair{
			i: i,
			j: 0,
		})
	}
	h := NewHeapFrom(pairs, less)
	res := make([][]int, 0, k)
	for h.Len() > 0 && len(res) < k {
		top := h.Peek()
		res = append(res, []int{nums1[top.i], nums2[top.j]})
		if top.j+1 < n {
			h.Update(0, pair{
				i: top.i,
				j: top.j + 1,
			}) // Update 内部已 Fix，不需要再调
		} else {
			h.Pop() // 列走完，退场；不退场会重复收同一个对
		}
	}
	return res
}

/*
* LC 1439 有序矩阵中的第 k 个最小数组和：m×n 矩阵只有每行升序（列无约束）。
* 从每一行各选 1 个元素组成数组，返回所有选法中和的第 k 小。
* 与 378 的区别：378 是矩阵里现成元素的第 k 小；本题是“每行选一”组合出的和的第 k 小。
* 与 373 的关系：本题的 m=2 退化版就是 373（两个升序数组各选一个）。

 ```
   示例: mat = [
					[1,3,11],
					[2,4,6]], k = 5 → 7
         前 5 小的和: [1,2]=3  [1,4]=5  [3,2]=5  [3,4]=7  [1,6]=7，第 5 名 = 7
         示例3: mat = [
					[1,10,10],
					[1,4,5],
					[2,3,6]], k = 7 → 9
 ```
*/

func KthSmallestPathSum(matrix [][]int, k int) int {

	m, n := len(matrix), len(matrix[0])
	less := func(idx1, idx2 []int) bool {
		res1, res2 := 0, 0
		for i := range m {
			res1 += matrix[i][idx1[i]]
			res2 += matrix[i][idx2[i]]
		}
		return res1 < res2
	}

	h := NewHeap(less)
	for j := range n {
		idx := make([]int, m)
		idx[0] = j
		h.Push(idx)
	}
	res := make([]int, 0, k)
	for h.Len() > 0 && len(res) < k {
		top := h.Peek()
		s := 0
		for i := range m {
			s += matrix[i][top[i]]
		}
		res = append(res, s)
		m := math.MaxInt
		var row int
		for i := 1; i < m; i++ {
			if matrix[i][top[i]] < m {
				m = matrix[i][top[i]]
				row = i
			}
		}
		if top[row]+1 < n {
			top[row] = +1
			h.Update(0, top)
		} else {
			h.Pop()
		}
	}
	return res[k-1]
}

/*
* LC 378 有序矩阵中第 K 小的元素：n x n 矩阵每行每列均升序，返回第 k 小元素。

 ```
   示例: matrix = [[1,5,9],[10,11,13],[12,13,15]], k = 8 → 13
 ```
*/

func KthSmallest(matrix [][]int, k int) int {
	m, n := len(matrix), len(matrix[0])
	pairs := make([]pair, 0, n)
	for i := range m {
		pairs = append(pairs, pair{
			i: i,
			j: 0,
		})
	}
	less := func(p1, p2 pair) bool {
		return matrix[p1.i][p1.j] < matrix[p2.i][p2.j]
	}
	h := NewHeapFrom(pairs, less)
	res := make([]int, 0, k)
	for h.Len() > 0 && len(res) < k {
		top := h.Peek()
		res = append(res, matrix[top.i][top.j])
		if top.j+1 < n {
			h.Update(0, pair{
				i: top.i,
				j: top.j + 1,
			})
		} else {
			h.Pop()
		}
	}
	return res[k-1]
}

/*
* LC 23 合并 K 个升序链表：给定链表数组 lists，合并为一条升序链表。
* 模板和多路归并一模一样，唯一的差别：路的"前进"从 j+1 变成 node.Next。
* 注意：Go 版需要自己定义 ListNode（LeetCode 平台已有）。

 ```
   示例: lists = [[1,4,5],[1,3,4],[2,6]] → 1→1→2→3→4→4→5→6
 ```
*/

type ListNode struct {
	Val  int
	Next *ListNode
}

func MergeKLists(lists []*ListNode) *ListNode {
	less := func(n1, n2 *ListNode) bool {
		return n1.Val < n2.Val
	}
	//k比较小,就用逐元素push的方法吧
	h := NewHeap(less)
	for _, node := range lists {
		if node != nil {
			h.Push(node)
		}
	}
	head := &ListNode{}
	cur := head
	for h.Len() > 0 {
		top := h.Peek()
		cur.Next = top
		cur = cur.Next
		if top.Next != nil {
			h.Update(0, top.Next)
		} else {
			h.Pop() //如果到了队尾,就pop出来！
		}
	}
	return head.Next
}

/*
* 练习：合并 K 个有序数组的前 m 个元素（多路归并教学题，非 LC 原题，LC 23 的数组版）。
* 给 k 个各自升序的数组，返回全局升序的前 m 个元素。不足 m 个则全部返回。
*
* 模型：候选池 = 每路"当前出头元素"。弹一个，该路游标 +1，补后继。

 ```
   示例: arrs = [[1,4,7,10],[2,5,8],[0,3,6,9]], m = 7
         → [0,1,2,3,4,5,6]
         (三路的出头元素 1/2/0 起步，堆顶依次弹出 0,1,2,3,4,5,6)
 ```
*/

type item struct {
	val int
	idx int
	row int
}

func MergeKSorted(arrs [][]int, k int) []int {
	h := NewHeap(func(a, b item) bool {
		if a.val != b.val {
			return a.val < b.val
		}
		return a.row < b.row
	})
	m := len(arrs)
	for i := range m {
		if len(arrs[i]) > 0 {
			h.Push(item{
				val: arrs[i][0],
				idx: 0,
				row: i,
			})
		}
	}
	var res []int
	for len(res) < k && h.Len() > 0 {
		min := h.Peek()
		res = append(res, min.val)
		if min.idx+1 < len(arrs[min.row]) {
			h.Update(0, item{
				val: arrs[min.row][min.idx+1],
				idx: min.idx + 1,
				row: min.row,
			})
		} else {
			h.Pop()
		}
	}
	return res
}

/*
* 最小区间：k 个升序数组，从每个数组里各取一个数，
* 使得取出的数中最大值 - 最小值最小。
* 返回的是区间 [最小值, 最大值]，不是选出的那 k 个数。

 ```
   示例: nums = [[4,10,15,24,26],[0,9,12,20],[5,18,22,30]]
         → [20,24]
         最优选法：三路分别取 24、20、22（构成区间 [20,24]，差值 4；
         注意返回的是区间，不是这三个数）
 ```
*/

func SmallestRange(nums [][]int) []int {
	k := len(nums)
	less := func(p1, p2 pair) bool {
		return nums[p1.i][p1.j] < nums[p2.i][p2.j]
	}
	h := NewHeap(less)
	maxHead := math.MinInt
	for i := range k {
		if len(nums[i]) > 0 {
			h.Push(pair{
				i: i,
				j: 0,
			})
			maxHead = max(maxHead, nums[i][0])
		}
	}
	best := math.MaxInt
	res := []int{}
	for h.Len() == k {
		top := h.Peek()
		diff := maxHead - nums[top.i][top.j]
		if diff < best {
			res = []int{nums[top.i][top.j], maxHead}
			best = diff
		}
		if top.j+1 < len(nums[top.i]) {
			maxHead = max(maxHead, nums[top.i][top.j+1])
			h.Update(0, pair{
				i: top.i,
				j: top.j + 1,
			})
		} else {
			h.Pop()
		}
	}
	return res
}

/*
* LC 313 超级丑数：质数数组 primes，丑数 = 只含这些质因子的正整数，
* 返回第 n 个（升序）。第 1 个是 1。

 ```
   示例: n = 12, primes = [2,7,13,19] → 32
         (序列前12个: 1,2,4,7,8,13,14,16,19,26,28,32)
 ```
*/

func NthSuperUglyNumber(n int, primes []int) int {
	return 0
}

/*
* LC 295 数据流的中位数
* AddNum(num) 加入一个数；FindMedian() 返回当前全部数的中位数。

 ```
   示例: addNum(1) addNum(2) → median = 1.5
         addNum(3)            → median = 2
 ```
*/

type MedianFinder struct {
	left  *Heap[int]
	right *Heap[int]
	n     int
}

func NewMedianFinder() *MedianFinder {
	m := MedianFinder{
		left: NewHeap(func(a, b int) bool {
			return a > b
		}),
		right: NewHeap(func(a, b int) bool {
			return a < b
		}),
	}
	return &m
}
func (mf *MedianFinder) AddNum(num int) {
	mf.left.Push(num)
	mf.n += 1
	if mf.n == 1 {
		return
	}
	//先修数量关系
	for mf.left.Len() > mf.right.Len()+1 {
		mf.right.Push(mf.left.Pop())
	}
	//数量关系满足了,大小关系还有可能不满足
	if l, r := mf.left.Peek(), mf.right.Peek(); l > r {
		mf.left.Update(0, r)
		mf.right.Update(0, l) //思考,为什么数量关系满足了,只需要修一次就够了？
		//关键点 因为我们根据一开始的设计,left的最大值都要小于 right的最小值。初始时,满足这种情况,假设当前破坏了这种关系,那至多也左右各破坏一个元素！
	}
}

func (mf *MedianFinder) RemoveNum(num int) {
	idx := -1
	var left bool
	if mf.left != nil {
		for i := 0; i < mf.left.Len(); i++ {
			if mf.left.At(i) == num {
				idx = i
				left = true
				break
			}
		}
	}

	if !left && mf.right != nil {
		for i := 0; i < mf.right.Len(); i++ {
			if mf.right.At(i) == num {
				idx = i
				left = false
				break
			}
		}
	}

	if left && idx != -1 {
		mf.left.Remove(idx)
	}
	if !left && idx != -1 {
		mf.right.Remove(idx)
	}
	mf.n -= 1

	//先修数量关系
	for mf.left.Len() > mf.right.Len()+1 {
		mf.right.Push(mf.left.Pop())
	}

	for mf.left.Len() < mf.right.Len() {
		mf.left.Push(mf.right.Pop())
	}
	//删除以后,左右都可能为空(当该死的k==1时候)
	//数量关系满足了,大小关系还有可能不满足
	if mf.left.Len() <= 0 || mf.right.Len() <= 0 {
		return
	}
	if l, r := mf.left.Peek(), mf.right.Peek(); l > r {
		mf.left.Update(0, r)
		mf.right.Update(0, l) //思考,为什么数量关系满足了,只需要修一次就够了？
		//关键点 因为我们根据一开始的设计,left的最大值都要小于 right的最小值。初始时,满足这种情况,假设当前破坏了这种关系,那至多也左右各破坏一个元素！
	}
	//删除以后balance
}
func (mf *MedianFinder) FindMedian() float64 {
	if mf.n == 0 {
		panic("here")
	}
	if mf.n%2 == 1 {
		return float64(mf.left.Peek())
	} else {
		return float64(mf.left.Peek()+mf.right.Peek()) / 2
	}
}

/*
* LC 480 滑动窗口中位数（阶段 3）：窗口从最左滑到最右，每次右移一位，
* 返回每次窗口内元素的中位数数组。

 ```
   示例: nums = [1,3,-1,-3,5,3,6,7], k = 3
         → [1.00000,-1.00000,-1.00000,3.00000,5.00000,6.00000]
 ```
*/
//1:需要堆能移除某个元素的能力.
//2:如何知道某个元素在堆内的位置,是否要在上一个题目的基础上改造,支持删除？
//3:难点,删除元素,打破左右关系,之后调整,比较难
func MedianSlidingWindow(nums []int, k int) []float64 {
	n := len(nums)
	res := make([]float64, 0, n-k+1)
	mf := NewMedianFinder()
	for i, num := range nums {
		mf.AddNum(num)
		if i >= k-1 {
			//n-1-k+1+1 = n - k + 1
			res = append(res, mf.FindMedian())
			//remove nums[i-k+1]
			mf.RemoveNum(nums[i-k+1])
		}
	}
	return res
}
