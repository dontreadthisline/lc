package lc

/*
给定一个二维矩阵 matrix，以下类型的多个请求：

- 计算其子矩形范围内元素的总和，该子矩阵的 左上角 为 (row1, col1) ，右下角 为 (row2, col2) 。

实现 NumMatrix 类：

- NumMatrix(int[][] matrix) 给定整数矩阵 matrix 进行初始化

- int sumRegion(int row1, int col1, int row2, int col2) 返回 左上角 (row1, col1) 、右下角 (row2, col2) 所描述的子矩阵的元素 总和 。

示例 1：

输入:
["NumMatrix","sumRegion","sumRegion","sumRegion"]
[[[
[3,0,1,4,2],
[5,6,3,2,1],
[1,2,0,1,5],
[4,1,0,1,7],
[1,0,3,0,5]]],

[2,1,4,3],[1,1,2,2],[1,2,2,4]]
输出:
[null, 8, 11, 12]

解释:
NumMatrix numMatrix = new NumMatrix([[3,0,1,4,2],[5,6,3,2,1],[1,2,0,1,5],[4,1,0,1,7],[1,0,3,0,5]]);
numMatrix.sumRegion(2, 1, 4, 3); // return 8 (红色矩形框的元素总和)
numMatrix.sumRegion(1, 1, 2, 2); // return 11 (绿色矩形框的元素总和)
numMatrix.sumRegion(1, 2, 2, 4); // return 12 (蓝色矩形框的元素总和)

提示：

- m == matrix.length

- n == matrix[i].length

- 1 <= m, n <= 200

- -10^5 <= matrix[i][j] <= 10^5

- 0 <= row1 <= row2 < m

- 0 <= col1 <= col2 < n

- 最多调用 10^4 次 sumRegion 方法
*/

// 注意: 构造器按仓库惯例改名 NewXxx（同包内多个 Constructor 会重名冲突）

type NumMatrix struct {
	prefixSum map[int]int
	m         int
	n         int
	matrix    [][]int
}

func NewNumMatrix(matrix [][]int) *NumMatrix {
	nm := NumMatrix{
		m:      len(matrix),
		n:      len(matrix[0]),
		matrix: matrix,
	}
	nm.prefixSum = make(map[int]int, nm.m*nm.n)
	//去填充prefixSum
	return &nm
}

func (m *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	sum := 0
	for i := row1; i <= row2; i++ {
		for j := col1; j <= col2; j++ {
			sum += m.matrix[i][j]
		}
	}
	return sum
}
