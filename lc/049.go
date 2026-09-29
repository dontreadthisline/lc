package lc

import (
	"slices"
	"strings"
)

/*
给你一个字符串数组，请你将 字母异位词 组合在一起。可以按任意顺序返回结果列表。

示例 1:

输入: strs = ["eat", "tea", "tan", "ate", "nat", "bat"]

输出: [["bat"],["nat","tan"],["ate","eat","tea"]]

解释：

- 在 strs 中没有字符串可以通过重新排列来形成 "bat"。

- 字符串 "nat" 和 "tan" 是字母异位词，因为它们可以重新排列以形成彼此。

- 字符串 "ate" ，"eat" 和 "tea" 是字母异位词，因为它们可以重新排列以形成彼此。

示例 2:

输入: strs = [""]

输出: [[""]]

示例 3:

输入: strs = ["a"]

输出: [["a"]]

提示：

- 1 <= strs.length <= 10^4

- 0 <= strs[i].length <= 100

- strs[i] 仅包含小写字母
*/

func groupAnagramsSortWay(strs []string) [][]string {
	anagrams := func(s1, s2 string) int {
		if len(s1) != len(s2) {
			return strings.Compare(s1, s2)
		}
		cnt1, cnt2 := [26]int{}, [26]int{}
		for i, c := range s1 {
			cnt1[c-'a'] += 1
			cnt2[s2[i]-'a'] += 1
		}
		for i := range 26 {
			if cnt1[i] != cnt2[i] {
				return strings.Compare(s1, s2)
			}
		}
		return 0
	}

	slices.SortFunc(strs, anagrams)
	//排好序之后分组 又用到了双指针
	groups := make([][]string, 0, len(strs))
	i, j, n := 0, 0, len(strs)
	for i < n {
		group := make([]string, 0)
		for j = i + 1; j < n && anagrams(strs[j], strs[i]) == 0; j += 1 {
			group = append(group, strs[j])
		}
		groups = append(groups, group)
		i = j
	}
	return groups
}

func groupAnagramsHashWay(strs []string) [][]string {
	anagramsKey := func(s string) string {
		chars := []byte(s)
		slices.Sort(chars)
		return string(chars)
	}
	m := make(map[string][]int, len(strs))
	res := make([][]string, 0, len(strs))
	for i, str := range strs {
		key := anagramsKey(str)
		m[key] = append(m[key], i)
	}
	for _, idxs := range m {
		anagrams := make([]string, 0, len(idxs))
		for _, idx := range idxs {
			anagrams = append(anagrams, strs[idx])
		}
		res = append(res, anagrams)
	}
	return res
}

func GroupAnagrams(strs []string) [][]string {
	return groupAnagramsHashWay(strs)
}
