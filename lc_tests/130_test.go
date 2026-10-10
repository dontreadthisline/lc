package lc_tests

import (
	"slices"
	"testing"

	"demo/lc"
)

// 130. 被围绕的区域
func TestSolve(t *testing.T) {
	cases := []struct {
		name string
		in   [][]byte
		want [][]byte
	}{
		{
			"官方用例1",
			[][]byte{
				{'X', 'X', 'X', 'X'},
				{'X', 'O', 'O', 'X'},
				{'X', 'X', 'O', 'X'},
				{'X', 'O', 'X', 'X'},
			},
			[][]byte{
				{'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X'},
				{'X', 'O', 'X', 'X'},
			},
		},
		{
			"官方用例2_单格",
			[][]byte{{'X'}},
			[][]byte{{'X'}},
		},
		{
			"边缘O通过链保护内陆O",
			[][]byte{
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'O'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
			},
			[][]byte{
				{'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'X'},
				{'X', 'X', 'X', 'X', 'O', 'O'},
				{'X', 'X', 'X', 'X', 'X', 'X'},
			},
		},
		{
			"全O",
			[][]byte{
				{'O', 'O', 'O'},
				{'O', 'O', 'O'},
				{'O', 'O', 'O'},
			},
			[][]byte{
				{'O', 'O', 'O'},
				{'O', 'O', 'O'},
				{'O', 'O', 'O'},
			},
		},
		{
			"全X",
			[][]byte{
				{'X', 'X'},
				{'X', 'X'},
			},
			[][]byte{
				{'X', 'X'},
				{'X', 'X'},
			},
		},
		{
			"中央孤O被捕获",
			[][]byte{
				{'X', 'X', 'X'},
				{'X', 'O', 'X'},
				{'X', 'X', 'X'},
			},
			[][]byte{
				{'X', 'X', 'X'},
				{'X', 'X', 'X'},
				{'X', 'X', 'X'},
			},
		},
		{
			"单行",
			[][]byte{{'O', 'X', 'O'}},
			[][]byte{{'O', 'X', 'O'}},
		},
		{
			"单列",
			[][]byte{{'O'}, {'X'}, {'O'}},
			[][]byte{{'O'}, {'X'}, {'O'}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			board := make([][]byte, len(c.in))
			for i := range c.in {
				board[i] = slices.Clone(c.in[i])
			}
			lc.Solve(board)
			for i := range c.want {
				if !slices.Equal(board[i], c.want[i]) {
					t.Errorf("Solve()\n got  %v\n want %v", board, c.want)
					return
				}
			}
		})
	}
}
