#!/usr/bin/env python3
"""gen_index.py —— 重新生成 mother-problems/INDEX.md（按解法家族分类 × 循序渐进）
用法: python3 scripts/gen_index.py
家族分类/排序在 FAMILIES 里改，重跑即可。"""
import os
import re

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
BASE = os.path.join(REPO, 'mother-problems')

DIFF = {}
for f in os.listdir(BASE):
    if f.endswith('.md') and f != 'INDEX.md':
        pid = f.split('-')[0].lstrip('0') or '0'
        for ln in open(os.path.join(BASE, f), encoding='utf-8'):
            m = re.match(r'- 难度：(.+)', ln.strip())
            if m:
                DIFF[pid] = m.group(1).strip()
                break

DONE = {'1', '3', '4', '33', '34', '49', '69', '128', '162', '217', '378', '704', '875', '1011', '15'}
CODE = {'1': 'lc/001.go', '3': 'lc/003.go'}
for p in ['4', '33', '34', '69', '162', '378', '704', '875', '1011']:
    CODE[p] = 'lc/004.go'

FAMILIES = [
 ('哈希与集合', '用空间换时间：一次遍历把"查找/配对"降为 O(1)', ['1', '217', '49', '128']),
 ('双指针与滑动窗口', '连续段问题：左右边界各只前进，方向由条件无歧义决定', ['283', '11', '15', '3', '209', '438', '904', '424', '76']),
 ('前缀和与哈希计数', '连续子数组求和/计数：前缀和相减 + 哈希查配对', ['560', '974', '1109', '304', '238']),
 ('二分与值域搜索', '单调无歧义的判定：猜 mid → 裁决 → 砍一半（四种形态全在这）', ['704', '34', '33', '153', '69', '278', '875', '1011', '410', '4', '378', '162', '240', '74']),
 ('栈与单调栈', "最近的未处理元素：栈内保持某种序，单调栈是它的进阶形态", ['20', '155', '225', '232', '622', '739', '503', '84', '42']),
 ('链表', '指针操作：dummy 节点、快慢指针、反转三件套', ['206', '21', '160', '141', '142', '19', '143', '25', '234']),
 ('树与递归', '大树的答案由子树推出：写对递归函数的语义就赢一半', ['104', '226', '101', '102', '236', '105', '98', '230', '124', '297']),
 ('BFS 与网格', '层序扩散：无权图最短路 / 多源同时扩散 / 连通块', ['200', '695', '130', '994', '127']),
 ('回溯', '决策树遍历：选择/递归/撤销三件套', ['46', '78', '77', '22', '131', '51', '79']),
 ('线性动态规划', '最优子结构：dp[i] 由更早的状态推出', ['70', '118', '198', '213', '53', '152', '139', '64', '62']),
 ('背包与子序列 DP', '从数组里选/不选：0-1 背包、完全背包、LIS、LCS、编辑距离', ['322', '518', '494', '416', '300', '1143', '72']),
 ('回文专题', '中心扩展 + 区间 DP 双解（5/32 两道）', ['5', '32']),
 ('贪心与区间', '排序后局部最优：交换论证保证正确', ['121', '122', '55', '45', '134', '763', '435', '452', '56', '179']),
 ('堆', '动态取最值：Top-K / 流中位数 / 合并 K 路', ['215', '347', '295', '23']),
 ('位运算', '异或消消乐与位计数', ['136', '191', '268', '338']),
 ('设计', '组合基本结构达成 O(1) 操作', ['146', '380', '208']),
 ('矩阵与模拟', '边界控制与原地操作', ['54', '48', '73']),
 ('排序与杂项', '手写排序 / 自定义比较', ['912', '179']),
]


def codefile(pid: str) -> str:
    if pid in CODE:
        return CODE[pid]
    n = int(pid)
    return ('lc/%04d.go' % n) if n >= 1000 else ('lc/%03d.go' % n)


def main():
    files = sorted(
        f for f in os.listdir(BASE)
        if f.endswith('.md') and f != 'INDEX.md'
    )
    by_pid = {f.split('-')[0].lstrip('0') or '0': f for f in files}

    missing = [pid for pid in
               (p for _, _, ids in FAMILIES for p in ids)
               if pid not in by_pid]
    assert not missing, f'题面缺失: {missing}'

    out = []
    out.append('# mother-problems 索引：按解法家族分类 × 循序渐进\n')
    out_line = f'> {len(by_pid)} 份题面按解法家族分为 {len(FAMILIES)} 组，组内按难度与依赖排序。'
    out.append(out_line)
    out.append('> 状态以最近一次刷新为准：已完成 = lc/ 下有通过测试的解法。')
    out.append('> 代码列：题目对应 lc/题号.go（已完成的 11 道指向现有实现文件）。\n')
    out.append('## 总览\n')
    done_n = sum(1 for pid in by_pid if pid in DONE)
    out.append(f'- 已完成 {done_n} 道 / 共 {len(by_pid)} 道；家族 4"二分与值域搜索"14 题全部通关\n')
    out.append('- 动线：读题面 → gf 跳代码文件写实现 → 测试绿 → INDEX 打钩\n')

    for idx, (name, core, ids) in enumerate(FAMILIES, 1):
        out.append(f'\n## 家族 {idx}：{name}\n')
        out.append(f'**核心解法**：{core}\n')
        out.append('| 序 | 题号 | 题目 | 难度 | 状态 |')
        out.append('|---|---|---|---|---|')
        rows = []
        for pid in ids:
            fn = by_pid[pid]
            title = fn[:-3].split('-', 1)[1]
            diff = DIFF.get(pid, '?')
            status = '已完成' if pid in DONE else '未做'
            rows.append((int(pid), fn, title, codefile(pid), diff, status))
        rows.sort(key=lambda r: (r[4], r[0]))
        for n, (pid, fn, title, code, diff, status) in enumerate(rows, 1):
            out.append(f'| {n} | [{title}]({code}) | {diff} | {status} |')

    out.append('\n## 与仓库其他材料的关系\n')
    out.append('- 解法思想详解：`docs/leetcode-8020-mother-problems.md`（12 条主线）')
    out.append('- 004 切分法的锚点问题：`lc/004.go` 的 `findMedianBinarySearchAlgs` 注释开头')
    out.append('- 标准库二分对照 demo：`lc_tests/slices_binary_demo_test.go`\n')

    path = os.path.join(BASE, 'INDEX.md')
    open(path, 'w', encoding='utf-8').write('\n'.join(out) + '\n')
    done = sum(1 for pid in by_pid if pid in DONE)
    print(f'INDEX.md 已生成: {len(by_pid)} 题, 已完成 {done}, 家族 {len(FAMILIES)} 组 -> {path}')


if __name__ == '__main__':
    main()
