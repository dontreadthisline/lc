package lc

/*
设计一个支持 push ，pop ，top 操作，并能在常数时间内检索到最小元素的栈。

实现 MinStack 类:

- MinStack() 初始化堆栈对象。

- void push(int value) 将元素 value 推入堆栈。

- void pop() 删除堆栈顶部的元素。

- int top() 获取堆栈顶部的元素。

- int getMin() 获取堆栈中的最小元素。

示例 1:

输入：
["MinStack","push","push","push","getMin","pop","top","getMin"]
[[],[-2],[0],[-3],[],[],[],[]]

输出：
[null,null,null,null,-3,null,0,-2]

解释：
MinStack minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
minStack.getMin();   --> 返回 -3.
minStack.pop();
minStack.top();      --> 返回 0.
minStack.getMin();   --> 返回 -2.

提示：

- -2^31 <= val <= 2^31 - 1

- pop、top 和 getMin 操作总是在 非空栈 上调用

- push, pop, top, and getMin最多被调用 3 * 10^4 次
*/

//MinStack 最小栈的定义
type MinStack struct {
	stack []int
	mins []int
	n int
}

func NewMinStack() *MinStack {
	return &MinStack{
		stack:make([]int,0),
		mins:make([]int,0),
		n:0,
	}
}

func (s *MinStack) Push(val int) {
	s.stack = append(s.stack,val)
	min := s.GetMin()
	if val < min {
		s.mins = append(s.mins,val)
	} else {
		s.mins = append(s.mins,min)
	}
	s.n += 1
}

func (s *MinStack) Pop() {
	if s.n <= 0 {
		return
	}
	s.stack = s.stack[:s.n-1]
	s.mins = s.mins[:s.n-1]
	s.n -=1
}

func (s *MinStack) Top() int {
	return s.stack[s.n-1]
}

func (s *MinStack) GetMin() int {
	return s.mins[s.n-1]
}
