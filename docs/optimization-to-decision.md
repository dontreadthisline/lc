# 优化 → 判定: 形式化理论

> 全符号化版本。结构: 定义 → 引理 → 定理 → 证明 → 算法不变量 → 松弛理论 → 实例。
> 实例题: 410 及其 max-min 变式, 2226, 875/1011, 合集 13 题。

## 1. 基本定义

**定义 1.1 (离散优化问题).** 三元组 $\mathcal{P} = \langle X, f, \mathbb{D} \rangle$, 其中

- $X$: 有限非空集合, 称**决策空间** (410: 所有恰好 $k$ 段切法);
- $f: X \to \mathbb{D}$, $\mathbb{D} \subseteq \mathbb{Z}$, 称**目标映射** (410: $f(P) = \max_{S \in P} \text{sum}(S)$);
- 目标: 求 $\mathrm{OPT} := \min_{x \in X} f(x)$。

**定义 1.2 (下水平集).** $\forall t \in \mathbb{Z}: \quad S_t := \{\, x \in X \mid f(x) \le t \,\}$。

**定义 1.3 (判定族).** $\forall t \in \mathbb{Z}: \quad Q(t) := [\, S_t \ne \varnothing \,] \in \{0,1\}$,
即 $Q(t) = 1 \iff \exists x \in X:\ f(x) \le t$。

**定义 1.4 (交界点).** $T^\* := \min \{\, t \mid Q(t) = 1 \,\}$。
存在性: 取任意 $x_0 \in X$, $Q(f(x_0)) = 1$, 集合非空; $\mathbb{Z}$ 良序, min 存在。

**定义 1.5 (值域与包围盒).** $f(X) := \{f(x) \mid x \in X\}$ (值空间中的**稀疏**子集);
$[\ell, u] \subseteq \mathbb{Z}$ 为包围盒, 若 $\ell \le \mathrm{OPT} \le u$。

## 2. 基本定理

**引理 2.1 (单调性).** $\forall t_1 \le t_2:\ Q(t_1) \le Q(t_2)$ (布尔序)。

证. $f(x) \le t_1 \implies f(x) \le t_2$, 故 $S_{t_1} \subseteq S_{t_2}$;
$S_{t_1} \ne \varnothing \implies S_{t_2} \ne \varnothing$。∎

由引理 2.1, $Q$ 在 $\mathbb{Z}$ 上呈 $0\cdots0\,1\cdots1$ 形, 交界点唯一且即 $T^\*$。

**定理 2.2 (两夹定理: 判定-优化等价).** $T^\* = \mathrm{OPT}$。

证.
(i) $X$ 有限非空, 取 $x^\* \in \arg\min_X f$, 则 $f(x^\*) = \mathrm{OPT} \le \mathrm{OPT}$,
故 $x^\* \in S_{\mathrm{OPT}}$, $Q(\mathrm{OPT}) = 1$。
由定义 1.4 (最小性), $T^\* \le \mathrm{OPT}$。
(ii) $Q(T^\* - 1) = 0$ (定义 1.4 最小性, 整数域保证前元存在), 即
$\neg\exists x: f(x) \le T^\*-1$, 故 $\forall x \in X: f(x) \ge T^\*$, 即 $\mathrm{OPT} \ge T^\*$。
(i)(ii) 合并: $T^\* = \mathrm{OPT}$。∎

**定理 2.3 (见证紧性).** $\forall x \in S_{T^\*}:\ f(x) = T^\*$。

证. $x \in S_{T^\*} \implies f(x) \le T^\*$。若 $f(x) \le T^\* - 1$,
则 $x \in S_{T^\*-1} \implies Q(T^\*-1) = 1$, 与 $T^\*$ 最小性矛盾。故 $f(x) = T^\*$。∎

**推论 2.4.** $S_{T^\*} \ne \varnothing$ 且 $S_{T^\*} \subseteq \arg\min_X f$;
特别地 $T^\* \in f(X)$ —— **交界点落在值域中是定理的结论, 不是算法的前提**。

**定理 2.5 (对偶形式, max-min).** 设 $h: X \to \mathbb{Z}$, 求 $\max_{x} h(x)$。
定义上水平集 $S'^t := \{x \mid h(x) \ge t\}$, 判定 $Q'(t) := [S'^t \ne \varnothing]$,
$T'^\* := \max\{t \mid Q'(t) = 1\}$。则引理 2.1 (反向单调)、定理 2.2、2.3 逐条镜像成立,
$T'^\* = \max_X h$。

## 3. 算法与不变量 (二分正确性)

**算法 A (在包围盒上二分).** 输入 $[\ell, u]$ 与 oracle $Q$ (约定 $Q(u) = 1$):

```python
while ℓ < u:
    mid := ⌊(ℓ+u)/2⌋
    if Q(mid) = 1: u := mid
    else:          ℓ := mid + 1
return ℓ
```

**不变量 INV.** $Q(\ell - 1) = 0 \ \wedge\ Q(u) = 1 \ \wedge\ \ell \le T^\* \le u$。

- 初始: 以 $\ell_0 := \ell - 1$ 为虚哨兵 ($Q(\ell_0) = 0$ 由 $\ell \le \mathrm{OPT}$ 保证);
  $Q(u) = 1$ 由 $u \ge \mathrm{OPT}$ 保证。
- 保持: $Q(\text{mid}) = 1 \Rightarrow u' = \text{mid}$, $Q(u') = 1$ 保持;
  $Q(\text{mid}) = 0 \Rightarrow \ell' = \text{mid}+1$, $Q(\ell' - 1) = Q(\text{mid}) = 0$ 保持。
- 终止: $\ell = u$, 由 INV 夹死 $\ell = u = T^\*$。

**复杂度.** oracle 调用 $\lceil \log_2(u - \ell) \rceil$ 次:
$$\text{优化} \ \le\ \text{判定} \times \log(\text{包围盒}) \tag{3.1}$$

**定理 2.6 (自归约: 由判定构造解).** 设 $X \subseteq \Sigma^n$ (决策为长度 $n$ 的串),
且 $Q$ 支持带前缀约束的变体 $Q_p(t) := [\exists x \in X: x_{1..|p|} = p \wedge f(x) \le t]$。
则最优解可由 $O(n \cdot |\Sigma|)$ 次约束判定构造: 从 $p = \varepsilon$ 起,
逐位贪心延伸使 $Q_p(T^\*) = 1$ 的前缀。∎ (构造性证明)

## 4. 松弛理论 (第③步的一般化)

实践中 oracle 常用**松弛版** $\tilde{Q}$ (更易判定)。何时无害:

**定义 4.1 (松弛).** $\tilde{Q}$ 是 $Q$ 的松弛, 若 $\forall t:\ Q(t) = 1 \implies \tilde{Q}(t) = 1$。
记对应交界点 $\tilde{T}^\*$。一般地 $Q \implies \tilde{Q}$ 给出 $\tilde{T}^\* \le T^\*$ (松弛可能**过头**)。

**定理 4.2 (松弛无害条件).** 若 $\forall t:\ \tilde{Q}(t) = 1 \implies Q(t) = 1$ (反向亦真, 即 $\tilde{Q} \equiv Q$),
则 $\tilde{T}^\* = T^\*$, 用 $\tilde{Q}$ 二分正确。∎ (平凡)

**定义 4.3 (划分格与细化算子).** 记 $\Pi_k$ = 将 $[n]$ 分成恰好 $k$ 个连续区间的所有划分。
定义细化偏序: $P \sqsubseteq P'$ ($P'$ 由 $P$ 继续切刀得到)。$(\Pi, \sqsubseteq)$ 构成格;
$P' \sqsupseteq P$ 时段数增加。

**命题 4.4 (切割单调性, min-max 侧).** 元素非负时:
$$P' \sqsupseteq P \implies f(P') = \max_{S \in P'} \text{sum}(S) \le \max_{S \in P} \text{sum}(S) = f(P) \tag{4.1}$$
即 $f$ 沿细化**单调不增** (切刀不亏)。
证. $P'$ 的每段含于 $P$ 的某段, 非负性 $\implies$ 子段和 $\le$ 母段和, 逐段取 max。∎

**命题 4.5 (合并单调性, max-min 侧).** 元素非负时, $g(P) := \min_{S \in P} \text{sum}(S)$ 沿细化**单调不增**
(即沿合并单调不减, 合并不亏)。证. 细化后某段被拆小, min 只可能降。∎

**定理 4.6 (410 的松弛等价, 形式化).** 记
$Q(t) = [\exists P \in \Pi_k: f(P) \le t]$ (恰好 $k$ 段),
$\tilde{Q}(t) = [\exists P \in \bigcup_{j \le k} \Pi_j: f(P) \le t]$ ($\le k$ 段)。
元素非负且 $k \le n$ 时, $Q \equiv \tilde{Q}$。

证. "$\implies$": $\Pi_k \subseteq \bigcup_{j\le k}\Pi_j$, 平凡。
"$\impliedby$": 设 $P \in \Pi_j$, $j < k$, $f(P) \le t$。因 $j < k \le n$, $P$ 中存在长度 $\ge 2$ 的段,
对它切一刀得 $P' \in \Pi_{j+1}$; 由命题 4.4, $f(P') \le f(P) \le t$。
归纳至 $j = k$。∎

**定理 4.7 (变式的松弛等价).** $g$-问题 (max-min) 中,
$Q(t) = [\exists P \in \Pi_k: g(P) \ge t]$ 与 $\tilde{Q}(t) = [\exists P \in \bigcup_{j \ge k} \Pi_j: g(P) \ge t]$ ($\ge k$ 段) 等价。
"$\impliedby$" 方向: $j > k$ 时对 $P$ 合并相邻段 (命题 4.5, $g$ 不减) 至 $k$ 段。∎

**注 4.8 (松弛方向的生成规则).**
$$\text{min-max} \xrightarrow{\text{切刀不亏 (4.1)}} \text{松弛为 } \le k \quad\quad \text{max-min} \xrightarrow{\text{合并不亏 (4.5)}} \text{松弛为 } \ge k$$
方向由 $f$ 在划分格上的单调方向**推导**, 非人为选择。

## 5. 复杂性注记

- (3.1) 是**归约**不是免费午餐: $Q$ 本身可以 NP-hard (TSP 判定)。
  二分保证 "优化不比判定难 $\log$ 因子", 不降低内在难度。
- 判定 $\le_p$ 优化: 已知 $\mathrm{OPT}$, $Q(t) = [t \ge \mathrm{OPT}]$;
  优化 $\le_p$ 判定: (3.1) + 定理 2.6。故在 NP 类上二者多项式等价。
- 判定便宜 (贪心/DP 多项式) 时二分大赚; 判定也难时转 branch & bound
  (LP 松弛给下界 + 剪枝, 即引理 2.1 的工程精装)、状压 DP ($n \le 20$)、近似算法。

## 6. 实例: 410 全套形式化

$$X = \Pi_k = \{P = (S_1, \dots, S_k)\ :\ S_i \text{ 连续非空},\ \bigsqcup S_i = [n]\}$$
$$f(P) = \max_i \sum_{j \in S_i} \text{nums}_j, \qquad \mathbb{D} = [\max_i \text{nums}_i,\ \sum_j \text{nums}_j]$$
$$\tilde{Q}(t): \text{贪心 (最晚塞)} \mapsto \text{最少段数 } c(t) \le k\ ? \qquad \text{(贪心最优性: 交换论证)}$$
$$\text{算法 A 于 } \tilde{Q} \implies \text{返回 } T^\* = \mathrm{OPT} \quad \text{(定理 4.6 + 2.2)}$$

贪心 $c(t)$ 的最优性 (它数出的确是最少段数), 交换论证:
设任一可行划分 $P$ 各段右端点 $r_1 < r_2 < \cdots$; 归纳可证贪心的第 $i$ 段右端点
$r_i' \le r_i$ (贪心每段塞到装不下为止, 不晚于任何可行划分), 故贪心段数 $\ge$ 任何可行划分的段数取到最少。∎

## 7. 一页总结

| 概念 | 符号 | 410 实例 |
|---|---|---|
| 决策空间 | $X$ | 恰好 $k$ 段切法 $\Pi_k$ |
| 目标映射 | $f: X \to \mathbb{Z}$ | $\max$ 段和 |
| 下水平集 | $S_t = \{f \le t\}$ | 各段都 $\le t$ 的切法 |
| 判定 | $Q(t) = [S_t \ne \varnothing]$ | 贪心数段 $\le k$ (松弛 $\tilde{Q}$) |
| 单调性 | $S_{t_1} \subseteq S_{t_2}$ | 同一见证 |
| 交界点 | $T^\* = \min\{Q = 1\}$ | 二分返回值 |
| 两夹定理 | $T^\* = \mathrm{OPT}$ | 判定值 = 最优值 |
| 紧性 | $x \in S_{T^\*} \Rightarrow f(x) = T^\*$ | 见证恰好压线 |
| 松弛无害 | $\tilde{Q} \equiv Q$ | 切刀不亏 $\implies$ $\le k \equiv$ 恰好 $k$ |
| 对偶 | 上水平集, $\max$ 交界 | 变式: $\ge k \equiv$ 恰好 $k$ |

## 8. 判定的两型: 参数扫描 vs 存在判定(875 与 410 的难度墙)

**定义 8.1 (参数扫描型).** 若 $X \subseteq \mathbb{Z}$ (决策本身是数值参数:
速度/天数/容量/半径), 则二分直接在 $X$ 上进行, $Q(x) = [f(x) \le h]$ 为
**一阶评估**(喂进模拟器直接算), 无构造。

**定义 8.2 (存在判定型).** 若 $X$ 为组合结构(切法/分配/划分), 二分在值空间
包围盒上进行, $Q(t) = [\exists P \in X: f(P) \le t]$ 含**存在量词**。
check 的实现 = 构造性存在证明(贪心/匹配), 且常需松弛等价(定理 4.6/4.7)。

| | 参数扫描 (875) | 存在判定 (410) |
|---|---|---|
| 二分空间 | 决策空间 $X \subseteq \mathbb{Z}$ 本身 | 值空间包围盒 $[\ell, u]$ |
| check 逻辑形式 | $f(x) \le h$ (一阶) | $\exists P: f(P) \le t$ (二阶) |
| 额外负担 | 无 | 构造贪心 + 松弛等价证明 |

**中间形态**: $B6/B7$ 的 check 为 $\sum_i \lceil b_i/T \rceil \le n$ ——
形式是存在型("能否分配"), 实质退化为算术(构造有闭式)。$B8$ 双 check 皆可
(贪心摊平 vs 前缀均值公式), 是横跨两型的标本。

**识别流程**: 先问 mid 是操作参数还是质量界限 → 参数型直接写模拟;
存在型再问构造贪心的方向与松弛的代价单调性(切/合不亏)。

## 9. 附注: 浮点边界的"恰好相等"陷阱 (1870 实证)

当最优解落在 total == hour 的精确边界时, 正向模拟与反向推导是两条不同的舍入路径:

```
正向: total = 2.0 + double(0.01)  == double(2.01)   ✓ (恰好相等, Go 实测)
反向: slack = double(2.01) - 2.0 < 0.01 (小 1 ulp)
      需求速度 = 100000/slack ≈ 10^7 + ε > 上界       ✗
```

同一题、两种都自洽的判定, 给出相反答案(10^7 vs -1)。逻辑排他: 官方
expected=-1 ⟹ 官方判题器必非纯正向 double 语义; 具体机制(反向推导 /
整数化截断 / 恰等判负)三种皆可自洽, 未实证确认, 待 LeetCode 提交验证。
防御: 判定全程只用一条路径(正向累加), 不做减法反推; 边界语义写进测试注释。
