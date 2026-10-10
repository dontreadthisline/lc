# scripts/ — 刷题工具集

## 用法总览

```bash
# 一键测试 + 进度统计
python3 scripts/check_progress.py

# 抓取一道新题的题面并入库 + 插入 Go 注释
python3 scripts/lcadd.py <slug>                    # 只抓题面
python3 scripts/lcadd.py <slug> <go文件> <函数名>   # 题面 + Go 注释

# 重新生成 INDEX.md（从 mother-problems/ 的题面元数据）
python3 scripts/gen_index.py
```

## 各脚本说明

### check_progress.py
做四件事：
1. 全量扫描 lc/ 目录（非 git diff），判断每道题的实现状态
2. 和 INDEX 状态交叉对比，报告不一致
3. 跑全套测试（排除已知挂起的 D 开头测试），输出通过/失败数
4. 按家族分组统计

**已知问题**：函数-题目映射有 bug。多个题共用一个 .go 文件时（如 003.go 含 10 道题），只找到第一个导出函数当作该文件所有题的主函数。导致"有实现未打钩"的统计不准确。**待修复**。

正确映射方式：需要维护一个 `题号 → 函数名` 的映射表，不能按文件自动检测。

### lcadd.py
从 leetcode.cn 抓取官方题面（通过 Playwright + GraphQL + cookies），转为 Markdown 存入 `mother-problems/`，可选插入 Go 文件的注释块。

依赖：
- `node` + `~/.config/playwright/lc-desc-one.js`
- cookies 有效（被拦时先跑 `python3 ~/.config/playwright/update-cookies.py`）

### gen_index.py
扫描 `mother-problems/*.md` 的元数据（难度、标签、所属家族），重新生成 `INDEX.md`。

**注意**：重新生成会覆盖手动修改的状态列（已完成/做错/非最优）。如果只更新个别题目状态，直接编辑 INDEX.md，不要重跑此脚本。
