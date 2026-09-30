package lc

/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

有效字符串需满足：

- 左括号必须用相同类型的右括号闭合。

- 左括号必须以正确的顺序闭合。

- 每个右括号都有一个对应的相同类型的左括号。

示例 1：

输入：s = "()"

输出：true

示例 2：

输入：s = "()[]{}"

输出：true

示例 3：

输入：s = "(]"

输出：false

示例 4：

输入：s = "([])"

输出：true

示例 5：

输入：s = "([)]"

输出：false

提示：

- 1 <= s.length <= 10^4

- s 仅由括号 '()[]{}' 组成
*/

func IssValid(s string) bool {
	n := len(s)
	stack := make([]byte,n)
	pair := map[byte]byte{
		']':'[',
		')':'(',
		'}':'{',
	}
	for i := range n {
		c := s[i]
		switch c {
		case ']','}',')':
			l := len(stack)
			if l <= 0 {
				return false
			}
			top := stack[l-1]
			if top != pair[c] {
				return false
			}
			stack = stack[:l-1]
		default:
			stack = append(stack,c)
		}
	}
	return len(stack) == 0
}
