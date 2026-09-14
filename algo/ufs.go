package algo

type UFS struct {
	parent []int
	rank   []int
	n      int
}

func NewUFS(n int) *UFS {
	ufs := UFS{
		parent: make([]int, n),
		rank:   make([]int, n),
		n:      n,
	}

	for i := range n {
		ufs.parent[i] = i
	}

	return &ufs
}

func (u *UFS) Find(x int) int {
	if x >= u.n || x < 0 {
		return -1
	}

	if x != u.parent[x] {
		//顺便做个路径压缩
		u.parent[x] = u.Find(u.parent[x])
	}
	//为什么是返回u.parent[x] 因为u.parent[x] 才是最终的根,而不是x.
	return u.parent[x]
}

func (u *UFS) Union(x, y int) bool {
	rx, ry := u.Find(x), u.Find(y)
	if rx == -1 || ry == -1 {
		return false
	}
	if rx == ry {
		return true //本来就在一个集合里,也算合并成功
	}
	sx, sy := u.rank[rx], u.rank[ry]
	if sx > sy {
		u.parent[ry] = rx
	} else if sx < sy {
		u.parent[rx] = ry
	} else {
		u.parent[rx] = ry
		u.rank[ry] += 1
	}
	return true
}
