#!/usr/bin/env python3
"""一键测试 + 进度统计"""
import re, os, glob, subprocess, sys
from collections import defaultdict

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
LC_DIR = os.path.join(ROOT, 'lc')
TEST_DIR = os.path.join(ROOT, 'lc_tests')
INDEX = os.path.join(ROOT, 'mother-problems', 'INDEX.md')

# ---- 显式映射: 题号 -> 主函数名(该题的完成依据) ----
# 必须映射的场景:
#   1. 多题共文件: lc/001.go lc/003.go(家族2十题) lc/004.go(家族4九题)
#   2. 文件内含非本题的辅助函数: 046 的 Factorial / 450 的 SearchNode
#   3. 多解法文件中需指定以哪个为准
# 未列出的单题文件走自动推断(任一非 New/Constructor 函数非桩即视为有实现)
MAIN_FUNC = {
    # lc/001.go
    '0001': 'TwoSum',
    # lc/003.go 家族2 十题共文件
    '0003': 'NoRepeateLongestSubStr',
    '0209': 'MinSubArrayLen',
    '0424': 'CharacterReplacement',
    '0438': 'FindAnagrams',
    '0904': 'TotalFruit',
    '1004': 'LongestOnes',
    '1493': 'LongestSubarray',
    '0713': 'NumSubarrayProductLessThanK',
    '0930': 'NumSubarraysWithSum',
    '0395': 'LongestSubstring',
    # lc/004.go 家族4 九题共文件
    '0004': 'FindMedian',
    '0033': 'SearchRotated',
    '0034': 'SearchRange',
    '0069': 'SqrtSearchAlgs',
    '0704': 'BinarySearch',
    '0162': 'SearchPeakBinarySearchAlgs',
    '0875': 'EtaFuckBananas',
    '1011': 'Transmite',
    '0378': 'SearchMatrixSearchWay',
    # 多解法/命名不规则的单题文件
    '0049': 'GroupAnagrams',
    '0015': 'ThreeSum',
    '0019': 'RemoveNthFromEnd',
    '0128': 'LongestConsecutive',
    '0283': 'MoveZeroes',
    '0992': 'GoodArrays',
    # 含辅助函数易误判
    '0046': 'Permute',            # Factorial 是辅助, 不能当完成依据
    '0450': 'DeleteNode',         # SearchNode 有实现但题目本体是桩
    # 显式声明主解法
    '0101': 'IsSymmetric',
    '0121': 'MaxProfit',
    '0217': 'ContainsDuplicate',
    '0230': 'KthSmallest',
    '0435': 'EraseOverlapIntervals',
    '0056': 'MergeRange',
    '0503': 'NextGreaterElements',
    '0142': 'DetectCycle',
    '0124': 'MaxPathSum',
    '0994': 'OrangesRotting',
    '0127': 'LadderLength',
    # 测试先行题(未做, 桩 + 红测试等着实现; 反向检测依赖这些映射提示打钩)
    '0042': 'Trap',
    '0084': 'LargestRectangleArea',
    '0143': 'ReorderList',
    '0025': 'ReverseKGroup',
    '0105': 'BuildTree',
    '0337': 'RobIII',
    '0022': 'GenerateParenthesis',
    '0077': 'Combine',
    '0078': 'Subsets',
    '0079': 'Exist',
    '0131': 'Partition',
    '0051': 'SolveNQueens',
    '0053': 'MaxSubArray',
    '0062': 'UniquePaths',
    '0064': 'MinPathSum',
    '0139': 'WordBreak',
    '0152': 'MaxProduct',
    '0198': 'Rob',
    '0213': 'RobII',
    '0070': 'ClimbStairs',
    '0118': 'Generate',
    '0072': 'MinDistance',
    '0300': 'LengthOfLIS',
    '0322': 'CoinChange',
    '0416': 'CanPartition',
    '0494': 'FindTargetSumWays',
    '0518': 'Change',
    '1143': 'LongestCommonSubsequence',
    '0032': 'LongestValidParentheses',
    '0045': 'Jump',
    '0055': 'CanJump',
    '0122': 'MaxProfitII',
    '0134': 'CanCompleteCircuit',
    '0452': 'FindMinArrowShots',
    '0763': 'PartitionLabels',
    '0054': 'SpiralOrder',
    '0073': 'SetZeroes',
    '0048': 'Rotate',
    '1526': 'MinNumberOperations',
    '0410': 'SplitArray',
}

def read_index():
    src = open(INDEX).read()
    rows = []
    for ln in src.split('\n'):
        if not ln.startswith('|') or 'lc/' not in ln or '[' not in ln:
            continue
        num_m = re.search(r'\| \d+ \| (\d+) \|', ln)
        link_m = re.search(r'\(lc/(\d+)\.go\)', ln)
        title_m = re.search(r'\[([^\]]+)\]', ln)
        if not (num_m and link_m and title_m):
            continue
        num, link, title = num_m.group(1), link_m.group(1), title_m.group(1)
        if '已完成' in ln and '非最优' in ln: status = 'suboptimal'
        elif '已完成' in ln: status = 'done'
        elif '做错' in ln: status = 'wrong'
        else: status = 'todo'
        rows.append({'num': num, 'link': link, 'title': title, 'index_status': status})
    return rows

def find_go_file(link):
    for c in [f'{link}.go', f'{int(link):03d}.go', f'{int(link):04d}.go']:
        p = os.path.join(LC_DIR, c)
        if os.path.exists(p): return p
    return None

def is_stub_body(body):
    """判断函数体是否是桩"""
    lines = [l.strip() for l in body.split('\n') if l.strip() and not l.strip().startswith('//')]
    if not lines: return True
    if any('panic' in l or 'TODO' in l for l in lines): return True
    if len(lines) == 1 and re.match(r'^return\s+-?\w+$', lines[0]): return True
    return False

def analyze_impl(path, num):
    if path is None: return 'missing'
    src = open(path).read()
    main = MAIN_FUNC.get(num)
    if main:
        m = re.search(r'^func (?:\(\w+ \*?\w+\) )?' + re.escape(main) + r'\([^)]*\)[^{]*\{([^}]*)\}', src, re.M)
        if not m: return 'no_func'
        return 'stub' if is_stub_body(m.group(1)) else 'has_impl'
    funcs = list(re.finditer(r'^func (?:\(\w+ \*\w+\) )?([A-Z]\w*)\([^)]*\)[^{]*\{([^}]*)\}', src, re.M))
    if not funcs: return 'no_func'
    has_real = False
    for m in funcs:
        fname = m.group(1)
        if fname.startswith('New') or fname == 'Constructor': continue
        if not is_stub_body(m.group(2)):
            has_real = True
            break
    return 'has_impl' if has_real else 'stub'

def collect_test_calls():
    """测试函数名 -> 其函数体内调用的 lc.XXX 集合(用于红色测试反查题号)"""
    calls = {}
    for f in glob.glob(os.path.join(TEST_DIR, '*_test.go')):
        src = open(f).read()
        for m in re.finditer(r'^func (Test\w+)\(t \*testing\.T\) \{(.*?)^\}', src, re.M | re.S):
            calls[m.group(1)] = set(re.findall(r'\blc\.([A-Z]\w*)\b', m.group(2)))
    return calls

def main():
    rows = read_index()
    stats = []
    discrepancies = []
    for row in rows:
        path = find_go_file(row['link'])
        row['impl_status'] = analyze_impl(path, row['num'])
        if row['impl_status'] == 'has_impl' and row['index_status'] == 'todo':
            discrepancies.append(f"  {row['num']} {row['title']}: 有实现但标未做")
        if row['impl_status'] != 'has_impl' and row['index_status'] in ('done','suboptimal'):
            discrepancies.append(f"  {row['num']} {row['title']}: 标已完成但主函数是桩/缺失 (依据 MAIN_FUNC['{row['num']}'])")
        if path is None:
            discrepancies.append(f"  {row['num']} {row['title']}: 文件不存在")
        stats.append(row)

    done = sum(1 for r in stats if r['index_status'] in ('done','suboptimal'))
    wrong = sum(1 for r in stats if r['index_status'] == 'wrong')
    todo = sum(1 for r in stats if r['index_status'] == 'todo')
    subopt = sum(1 for r in stats if r['index_status'] == 'suboptimal')

    # 按家族
    family_stats = defaultdict(lambda: {'done':0,'wrong':0,'todo':0,'subopt':0,'total':0})
    src = open(INDEX).read()
    current_family = ''
    fam_names = {}
    for ln in src.split('\n'):
        m = re.match(r'## 家族 (\d+)：(.+)', ln)
        if m:
            current_family = m.group(1)
            fam_names[current_family] = m.group(2).split('——')[0].strip()
            continue
        if ln.startswith('|') and 'lc/' in ln and '[' in ln:
            family_stats[current_family]['total'] += 1
            if '非最优' in ln: family_stats[current_family]['subopt'] += 1
            elif '已完成' in ln: family_stats[current_family]['done'] += 1
            elif '做错' in ln: family_stats[current_family]['wrong'] += 1
            else: family_stats[current_family]['todo'] += 1

    print(f"{'家族':<4} {'名称':<16} {'最优':>4} {'非最优':>5} {'做错':>4} {'未做':>4} {'合计':>4}")
    print("-" * 58)
    for fam in sorted(family_stats.keys(), key=int):
        s = family_stats[fam]
        name = fam_names.get(fam,'')[:14]
        print(f"{fam:<4} {name:<16} {s['done']:>4} {s['subopt']:>5} {s['wrong']:>4} {s['todo']:>4} {s['total']:>4}")

    print("-" * 58)
    print(f"{'合计':<21} {done-subopt:>4} {subopt:>5} {wrong:>4} {todo:>4} {len(stats):>4}")

    wrong_list = [f"{r['num']}" for r in stats if r['index_status']=='wrong']
    subopt_list = [f"{r['num']}" for r in stats if r['index_status']=='suboptimal']

    if wrong_list: print(f"\n做错({len(wrong_list)}): {', '.join(wrong_list)}")
    if subopt_list: print(f"非最优({len(subopt_list)}): {', '.join(subopt_list)}")

    impl_gap = sum(1 for r in stats if r['impl_status']=='has_impl' and r['index_status']=='todo')
    if impl_gap: print(f"⚠ 有实现未打钩: {impl_gap} 题")

    if discrepancies:
        print(f"\n{'!'*50}")
        print(f"INDEX 不一致 ({len(discrepancies)} 处):")
        for d in discrepancies: print(d)
    else:
        print("\n✓ INDEX 与代码完全一致")

    # 测试
    print(f"\n{'='*70}")
    print("运行测试...")
    print("="*70)
    failed_tests = []
    passed_tests = set()
    try:
        # TestDetectCycle 挂起已修复(lc/types.go 的 String() 加防环截断), 不再需要 -skip
        # -skip: 栈溢出(fatal 级, recover 拦不住)的测试会杀掉整个进程
        # TestPartition(131 半成品, 无限递归) —— 待 131 写完后移除此跳过
        r = subprocess.run(['go','test','./lc_tests/','-skip','^TestPartition$','-timeout','120s','-v'],
                          capture_output=True, text=True, timeout=150, cwd=ROOT)
        output = r.stdout + r.stderr
        passed = len(re.findall(r'^\s*--- PASS', output, re.M))
        failed = len(re.findall(r'^\s*--- FAIL', output, re.M))
        failed_tests = list(dict.fromkeys(re.findall(r'^--- FAIL: (\S+)', output, re.M)))
        passed_tests = set(re.findall(r'^--- PASS: (\S+)', output, re.M))
        print(f"通过: {passed}  失败: {failed}")
        for ft in failed_tests: print(f"  ✗ {ft}")
    except subprocess.TimeoutExpired:
        print("⚠ 测试超时")

    # 红色测试 <-> INDEX 状态交叉验证
    if failed_tests:
        test_calls = collect_test_calls()
        func_to_num = {}
        for num, fn in MAIN_FUNC.items():
            func_to_num.setdefault(fn, []).append(num)
        num_status = {r['num']: r['index_status'] for r in stats}
        num_title = {r['num']: r['title'] for r in stats}
        conflicts, expected_red = [], []
        for ft in failed_tests:
            for fn in test_calls.get(ft, set()):
                for num in func_to_num.get(fn, []):
                    st = num_status.get(num)
                    if st in ('done', 'suboptimal'):
                        conflicts.append((num, num_title.get(num,''), ft))
                    elif st == 'wrong':
                        expected_red.append((num, ft))
        if conflicts:
            print(f"\n⚠ 假阳性警报: 标已完成但测试红 ({len(conflicts)} 题):")
            seen = set()
            for num, title, ft in conflicts:
                if num in seen: continue
                seen.add(num)
                print(f"  {num} {title} ← {ft}")
        if expected_red:
            seen = set()
            line = []
            for num, ft in expected_red:
                if num not in seen:
                    seen.add(num); line.append(num)
            print(f"预期红色(做错题): {', '.join(sorted(seen))}")

        # 反向验证: 调用主函数的测试全部 PASS, 但 INDEX 仍标做错/未做 -> 该打钩了
        test_calls = collect_test_calls()  # 复用已收集的调用表
        stale = []
        for num, fn in MAIN_FUNC.items():
            related = [t for t, calls in test_calls.items() if fn in calls]
            if not related: continue  # 没有任何测试锁定该题, 无法判断
            if all(t in passed_tests for t in related):
                st = num_status.get(num)
                if st in ('wrong', 'todo'):
                    stale.append((num, num_title.get(num, ''), st))
        if stale:
            print(f"\n⚠ 可打钩未打: 测试全绿但 INDEX 状态过时 ({len(stale)} 题):")
            for num, title, st in sorted(stale):
                print(f"  {num} {title} 当前={st} -> 应为已完成")

    print(f"\n{'='*70}")
    print(f"总结: {done-subopt} 最优 / {subopt} 非最优 / {wrong} 做错 / {todo} 未做 / 共 {len(stats)}")

if __name__ == '__main__':
    main()
