package algo

type User struct {
	Id  int
	Age int
}

const (
	MaxAge = 120
)

// 桶排序,二维数组,数据范围是小范围(或者能hash到小范围),利用分组去排序
func BucketSort(users []*User) {
	buckets := make([][]*User, MaxAge+1)
	for _, user := range users {
		buckets[user.Age] = append(buckets[user.Age], user)
	}

	i := 0
	for _, bucket := range buckets {
		if len(bucket) > 0 {
			for _, user := range bucket {
				users[i] = user
				i += 1
			}
		}
	}
}

// 计数排序,统计每个元素出现的次数
func CountSort(users []*User) []*User {
	count := make([]int, MaxAge+1)
	for _, user := range users {
		count[user.Age] += 1
	}

	prefixSum := make([]int, MaxAge+1)
	for i := 0; i <= MaxAge; i++ {
		if i != 0 {
			prefixSum[i] = prefixSum[i-1] + count[i]
		} else {
			prefixSum[i] = count[i]
		}
	}

	oUsers := make([]*User, len(users))
	for i := len(users) - 1; i >= 0; i-- {
		prefixSum[users[i].Age] -= 1
		oUsers[prefixSum[users[i].Age]] = users[i]
	}
	return oUsers
}

/*
*  给整数数组 nums 和整数 k，返回出现频率前 k 高的元素，任意顺序。约束 n ≤ 10^5，答案集合唯一。进阶要求：必须优于 O(n log n)。

 ```
   示例: nums = [1,1,1,2,2,3], k = 2 → [1,2]
         (1 出现 3 次、2 出现 2 次、3 出现 1 次，前 2 高频)
 ```
*/

func MaxKFreq(nums []int, k int) []int {
	offset := 10000 // 10^4
	freqs := make([]int, 20001)
	for _, num := range nums {
		freqs[num+offset] += 1
	}
	buckets := make([][]int, len(nums)+1)
	for idx, f := range freqs {
		if f > 0 {
			buckets[f] = append(buckets[f], idx-offset)
		}
	}
	res := make([]int, 0, k)
	for f := len(buckets) - 1; f >= 0 && len(res) < k; f-- {
		res = append(res, buckets[f]...)
	}
	return res[:k]
}
