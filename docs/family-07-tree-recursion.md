# 家族 7：树与递归 —— 攻坚手册

> 沉淀到 docs 下。按依赖顺序编排，每道题先给"从暴力怎么想到"，再给代码，再给"你的错题在哪"。

## 0. 树题的全部思维模式（先看这张图）

```
树上的信息只有两个流动方向：

    自顶向下（传参）              自底向上（返回值）
    祖先 → 后代                  后代 → 祖先
    传"约束"或"路径状态"           返回"子问题的答案"

    ┌─────────────┐             ┌─────────────┐
    │ 098 BST 验证  │             │ 236 LCA      │
    │ 传 (lo, hi)   │             │ 左右回报见没见到│
    ├─────────────┤             ├─────────────┤
    │ 437 路径总和III│             │ 124 最大路径和 │
    │ 传 path 前缀和 │             │ 返回单链，全局记拐弯│
    ├─────────────┤             ├─────────────┤
    │ 129 根到叶数字和│             │ 337 打家劫舍III│
    │ 传累积数字     │             │ 返回(选,不选)两种状态│
    └─────────────┘             └─────────────┘

    另有一个横切面（不是流向，是遍历方式）：
    层序（BFS）：102/103/107/199/515 —— 按层访问，队列驱动
    中序（BST 专用）：230/98/173/538 —— BST 的中序 = 有序序列
    构造：105/106/108/1008 —— 从序列重建树
    重构：114/226/116 —— 边遍历边改结构
```

**拿到任何树题，先问一句：信息往哪流？** 答案决定了骨架。

---

## 1. 基础层：单节点递归（已做，跳过）

### 104 最大深度 / 226 翻转 / 101 对称

这三道的共同点：递归函数的返回值 = 答案，没有任何分离。

```go
// 104: 返回值 = 以我为根的最大深度
func MaxDepth(root *TreeNode) int {
    if root == nil { return 0 }
    return 1 + max(MaxDepth(root.Left), MaxDepth(root.Right))
}
```

---

## 2. 层序遍历（BFS on Tree）

### 102 层序遍历（已做 ✅）

队列 + 层快照。已通过测试。

### 199 右视图（变体，未做）

每层只取最后一个节点。在 102 的骨架上改一行：`res = append(res, level[len(level)-1])`。

**生成规则**：103 锯齿（奇偶层翻转）、107 倒序（最后 reverse）、515 每层最大值（level 里取 max）——全是 102 的换装。

---

## 3. 自顶向下：传约束

### 098 验证 BST（做错 → 修复中）

**从暴力推**：

暴力是每个节点检查全部后代。O(n²)。反转方向：后代自己检查自己，祖先的要求压缩成 (lo, hi) 传下来。

**建议先写中序版**（更直观）：

```go
func IsValidBST(root *TreeNode) bool {
    prev := math.MinInt
    var inorder func(*TreeNode) bool
    inorder = func(node *TreeNode) bool {
        if node == nil { return true }
        if !inorder(node.Left) { return false }
        if node.Val <= prev { return false }
        prev = node.Val
        return inorder(node.Right)
    }
    return inorder(root)
}
```

**中序版的直觉**：BST 的中序遍历严格递增。遍历一遍，每步检查"我比前一个大"。

**上下界版为什么等价**：中序遍历的"prev"（我的中序前驱）就是沿路径往左走的最近祖先——上下界版的 hi 就是它的显式追踪。两种写法是同一个信息（"我的所有左拐祖先"）的不同编码。

### 235 BST 的最近公共祖先（未做，235 → 236 的前置）

BST 的有序性把 LCA 简化成一次走树：

```go
func LowestCommonAncestorBST(root, p, q *TreeNode) *TreeNode {
    for root != nil {
        if p.Val < root.Val && q.Val < root.Val {
            root = root.Left       // p、q 都在左子树
        } else if p.Val > root.Val && q.Val > root.Val {
            root = root.Right      // p、q 都在右子树
        } else {
            return root            // 分岔了，我就是 LCA
        }
    }
    return nil
}
```

**直觉**：如果 p 和 q 在 root 的两侧（或其中一个就是 root），root 就是分岔点。BST 的有序性让你一次比较就知道往哪边走。

---

## 4. 自顶向下：传路径状态

### 437 路径总和 III（未做，你刚踩过的坑）

**题面**：数"从某个祖先往下走到某个后代"的和为 target 的路径条数。

**这就是你给 124 写的那套骨架的正确归宿**——维护根到当前节点的路径，算后缀和。但因为树上有负数，暴力后缀和是 O(n²)。

**优化**：把家族 3 的前缀和 + 哈希原样搬上树：

```go
func PathSum(root *TreeNode, targetSum int) int {
    count := make(map[int64]int)
    count[0] = 1  // 种子：空前缀
    var dfs func(node *TreeNode, prefix int64) int
    dfs = func(node *TreeNode, prefix int64) int {
        if node == nil { return 0 }
        prefix += int64(node.Val)
        res := count[prefix - int64(targetSum)]  // 查"之前有多少前缀 = 我 - target"
        count[prefix]++                          // 入账
        res += dfs(node.Left, prefix)
        res += dfs(node.Right, prefix)
        count[prefix]--                          // 回溯：离开这条路，撤销记账
        return res
    }
    return dfs(root, 0)
}
```

**和 560 的对照**：完全同一套配对键恒等的查账逻辑，只是"从左往右扫数组"变成了"从根往叶走树"——多了一步**回溯**（离开一条路径后要撤销它的前缀，否则左子树的前缀会污染右子树的查询）。

**为什么你给 124 写的骨架用不到 437 上也不行**：437 只需竖直路径（祖先→后代），你的 path + 后缀和思路方向对了，但实现应该是"前缀哈希"而不是"逐个算后缀"（O(n) vs O(n²)）。

---

## 5. 自底向上：聚合（返回值 = 答案，一步到位）

### 236 最近公共祖先（未做）

后序聚合：左右子树各回报"见没见到 p/q"，两边都有 → 我是分岔点。

```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    if root == nil || root == p || root == q { return root }
    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)
    if left != nil && right != nil { return root }  // 分岔
    if left != nil { return left }
    return right
}
```

**递归语义**："在以我为根的子树里找 p 和 q 的 LCA；如果只找到一个，返回那个；一个没找到返回 nil。"

---

## 6. 自底向上：单链 / 全局分离（124 的核心难点）

### 124 最大路径和（做错，算法是直链 → 需要改成拐弯）

**你的错误代码在做什么**：

```go
// 你的版本：维护 path（根到当前节点），算后缀和
path = append(path, root.Val)
for i := len(path)-1; i >= 0; i-- { sum += path[i]; res = max(res, sum) }
```

这搜到的是**竖直路径**的最大和——"从某个祖先往下走到我"。代码本身没有 bug（Left/Right 修好后是正确的竖直路径搜索），但 124 的答案可能是**拐弯路径**：

```
    -10
   /   \
  9     20      ← 答案 = 15 + 20 + 7 = 42，路径在 20 处拐弯
       /  \
      15   7
```

15 和 7 分别在 20 的两侧。**竖直路径最多包含其中一边**，你的搜索空间里根本没有 42。

**为什么必须自底向上**：竖直路径可以自顶向下传（你的 path 参数就是）；拐弯路径的顶点在中间，顶点需要**同时知道左链最优和右链最优**——这两个信息只能从下往上聚合，不能从上往下传（上面不知道下面的内部结构）。

**两本账的正确代码**：

```go
func MaxPathSum(root *TreeNode) int {
    res := math.MinInt
    var dfs func(*TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil { return 0 }
        left  := max(0, dfs(node.Left))   // 左子树给我的最好单链（负不如不要）
        right := max(0, dfs(node.Right))  // 右子树给我的最好单链
        res = max(res, node.Val + left + right)  // 拐弯账：左链+我+右链
        return node.Val + max(left, right)        // 单链账：只能选一边，报给父
    }
    dfs(root)
    return res
}
```

**为什么返回值只能选一边**：父节点把你接在它的某条链上。你的返回值代表"以我为顶点、往下延伸的一条线"。如果你把左右都带上去，父节点那条线就拐了两次——不再是线。

**为什么拐弯要在自己这层记**：拐弯的顶点就是你自己。父节点不需要也不应该知道你内部拐了弯——它只关心你作为一段链能贡献多少。两本账各管各的。

**你的暴力版和正确版的关系**：你的 path + 后缀和算的是"单链账"的一部分（只是从祖先视角看，而不是从当前节点视角看）。改造的关键不是修 bug，是**加一个自底向上的通道**（返回值），让拐弯路径也可达。

---

## 7. 自底向上：树形 DP（选/不选两状态）

### 337 打家劫舍 III（未做，124 的进阶）

**题面**：二叉树上打家劫舍，不能同时选父节点和子节点，求最大收益。

**和 124 的关系**：124 的两本账是"单链 vs 拐弯"两种**形状**；337 的两本账是"选我 vs 不选我"两种**决策**。骨架一模一样，只是状态含义换了。

```go
func RobIII(root *TreeNode) int {
    var rob func(*TreeNode) (int, int)  // 返回 (选我的最大, 不选我的最大)
    rob = func(node *TreeNode) (int, int) {
        if node == nil { return 0, 0 }
        lRob, lSkip := rob(node.Left)
        rRob, rSkip := rob(node.Right)
        // 选我：左右子节点都不能选
        take := node.Val + lSkip + rSkip
        // 不选我：左右子节点各选最优
        skip := max(lRob, lSkip) + max(rRob, rSkip)
        return take, skip
    }
    take, skip := rob(root)
    return max(take, skip)
}
```

**递归语义**：返回两个数——"偷我的最优收益"和"不偷我的最优收益"。父节点根据这两个数做自己的决策。**状态从 124 的一个返回值（单链和）扩展到两个返回值（两种决策）**——这就是树形 DP 的全部。

---

## 8. BST 结构操作

### 450 删除 BST 中的节点（未做）

三种情况（按目标节点子节点数分）：

1. **叶节点**：直接删
2. **一个子节点**：用子节点替代
3. **两个子节点**：找右子树的最小节点（右子树一路向左），用它的值替代目标，然后在右子树中删除那个最小节点

情况 3 是难点：**为什么用右子树最小替代？** 因为它 > 左子树所有人（保持 BST），且 < 右子树其他人（删除后仍有序）。它是"能坐在我的位置上、不打乱两边的唯一候选人"。

---

## 9. 树的重构

### 114 二叉树展开为链表（未做）

**要求**：前序遍历的顺序，全部变成右链。

```
    1               1
   / \               \
  2   5        →      2
 / \   \               \
3   4   6               3
                         \
                          4
                           \
                            5
                             \
                              6
```

**思路（后序重构）**：先递归展开左右子树，然后把左子树（已展开为链）接到右边，原右子树接到新链的末尾。

```go
func Flatten(root *TreeNode) {
    if root == nil { return }
    Flatten(root.Left)
    Flatten(root.Right)
    right := root.Right
    root.Right = root.Left
    root.Left = nil
    // 走到新右链的末端，接上原右子树
    for root.Right != nil { root = root.Right }
    root.Right = right
}
```

**关键**：先递归（子问题已解决），再拼接。和 226 翻转的骨架一样——后序遍历中做结构调整。

---

## 10. 构造

### 105 前序 + 中序构造二叉树（未做）

前序第一个是根；在中序里定位根，左段是左子树，右段是右子树；递归。

```go
func buildTree(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 { return nil }
    rootVal := preorder[0]
    idx := 0
    for i, v := range inorder {
        if v == rootVal { idx = i; break }
    }
    return &TreeNode{
        Val:   rootVal,
        Left:  buildTree(preorder[1:idx+1], inorder[:idx]),
        Right: buildTree(preorder[idx+1:], inorder[idx+1:]),
    }
}
```

**优化**：哈希预存中序值→下标，省去每次线性扫。

---

## 11. 序列化

### 297 序列化与反序列化（未做）

前序遍历 + 空节点标记（`#,`），序列化和反序列化互为镜像。

```go
// 序列化：前序遍历，空节点输出 "#,"
// 反序列化：逐 token 消费，遇到 "#" 返回 nil，否则建节点递归左右
```

---

## 练习清单（从浅入深）

按此顺序做，每道题给出一句话自检（做完口述）：

| 序 | 题号 | 题目 | 难度 | 前置 | 自检口述 |
|---|---|---|---|---|---|
| 1 | 104 | 最大深度 | 简单 | — | "递归函数返回什么的语义" |
| 2 | 226 | 翻转二叉树 | 简单 | — | 同上 |
| 3 | 101 | 对称二叉树 | 简单 | — | "两次递归怎么配对" |
| 4 | 102 | 层序遍历 | 中等 | — | "层边界怎么记" |
| 5 | 098 | 验证 BST | 中等 | — | "中序版：为什么中序递增 ⟺ BST" |
| 6 | 235 | BST 的 LCA | 中等 | 098 | "BST 的有序性怎么把 LCA 变成一次走树" |
| 7 | 236 | 二叉树的 LCA | 中等 | 235 | "后序聚合：左右各回报什么，什么时候判定分岔" |
| 8 | **124** | **最大路径和** | **困难** | 236 | **"返回值和答案为什么不是一回事，单链账为什么只能选一边"** |
| 9 | 437 | 路径总和 III | 中等 | 124 | "前缀和+哈希怎么从数组搬到树上，为什么要回溯" |
| 10 | 337 | 打家劫舍 III | 中等 | 124 | "两本账从'形状'变成'决策'，返回值从 1 个变 2 个" |
| 11 | 450 | 删除 BST 节点 | 中等 | 098/235 | "两个子节点时为什么用右子树最小替代" |
| 12 | 114 | 展开为链表 | 中等 | 226 | "后序重构：先递归再拼接的顺序为什么不能反" |
| 13 | 105 | 前中序构造 | 中等 | — | "前序定根，中序切分的对应关系" |
| 14 | 297 | 序列化 | 困难 | 105 | "空标记为什么必须显式" |
| 15 | 230 | BST 第 K 小 | 中等 | 098 | "中序遍历的第 k 个就是答案" |

**做 124 之前确保 236 过了**（LCA 是"后序聚合"的最简形态，先在简单题上体会"左右回报→父决策"的节奏，再去打 124 的"两本账"）。

**做 437 之前确保 124 过了**（437 的前缀哈希 + 回溯需要你先理解"树上路径"和"数组子数组"的对应，而 124 会让你真正看清"树上路径的形状"）。

---

## 本仓库已有错题

| 题号 | 错误 | 根因 | 修复方向 |
|---|---|---|---|
| 098 | 只查直接子节点 | 未意识到孙辈约束 | 先写中序版 |
| 124 | 算法是直链 | 混淆了 124 和 437 的路径形态 | 重写为两本账 |
