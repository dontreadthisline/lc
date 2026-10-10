package lc

/*
给定一个整数数组 temperatures ，表示每天的温度，返回一个数组 answer ，其中 answer[i] 是指对于第 i 天，下一个更高温度出现在几天后。如果气温在这之后都不会升高，请在该位置用 0 来代替。

示例 1:

输入: temperatures = [73,74,75,71,69,72,76,73]
输出: [1,1,4,2,1,1,0,0]

示例 2:

输入: temperatures = [30,40,50,60]
输出: [1,1,1,0]

示例 3:

输入: temperatures = [30,60,90]
输出: [1,1,0]

提示：

- 1 <= temperatures.length <= 10^5

- 30 <= temperatures[i] <= 100
*/

func DailyTemperatures(temperatures []int) []int {
	n := len(temperatures)
	stack := make([]int,0,n)
	res := make([]int,n)

	for i,t := range temperatures {
		for n := len(stack); n  > 0  && t > temperatures[stack[n-1]]; n = len(stack){
			res[stack[n-1]] = i - stack[n-1]
			stack = stack[0:n-1]
		}
		stack = append(stack,i)
	}

	for n := len(stack); n > 0; n = len(stack) {
		res[stack[n-1]] = 0
		stack=stack[0:n-1]
	}

	return res
}


func DailyTemperatures1(temperatures []int) []int {
	n := len(temperatures)
	stack := NewStack[int]()
	res := make([]int,n)

	for i,t := range temperatures {
		for stack.Len() > 0 && t > temperatures[stack.Peek()] {
			res[stack.Peek()] = i - stack.Peek()
			stack.Pop()
		}
		stack.Push(i)
	}

	for stack.Len() > 0 {
		res[stack.Peek()] = 0
		stack.Pop()
	}

	return res
}


