#!/usr/bin/env python3
"""lcadd.py —— 力扣题面一条龙（抓取 → 入库 → 插入 go 注释），本会话沉淀的固定工具。

用法:
  python3 scripts/lcadd.py <slug>                          # 只抓题面入库 mother-problems/
  python3 scripts/lcadd.py <slug> <go文件> <函数名>         # 再把题面注释插到该函数上方

依赖:
  - node + ~/.config/playwright/lc-desc-one.js（真实浏览器抓官方 GraphQL，带 cookie）
  - 首次/被 WAF 拦时：Chrome 打开一次 leetcode.cn，再跑
    python3 ~/.config/playwright/update-cookies.py

效果:
  1) 题面转 Markdown 存 mother-problems/<题号>-<标题>.md（自动去 HTML 标签）
  2) go 文件里生成/替换注释块（带 [题面-begin/end] 标记，可重复执行不重复插入）
"""
import json
import os
import re
import subprocess
import sys
from html.parser import HTMLParser

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DESC_DIR = os.path.join(REPO, 'mother-problems')
FETCHER = os.path.expanduser('~/.config/playwright/lc-desc-one.js')
DIFF = {'EASY': '简单', 'MEDIUM': '中等', 'HARD': '困难',
        'Easy': '简单', 'Medium': '中等', 'Hard': '困难'}


class LC2MD(HTMLParser):
    """力扣题面 HTML → Markdown（与历史 106 题转换规则一致，勿改语义）"""

    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.out = []
        self.pre = 0

    def handle_starttag(self, tag, attrs):
        if tag == 'pre':
            self.pre += 1
            self.out.append('\n```\n')
        elif tag in ('p', 'ul', 'ol'):
            self.out.append('\n')
        elif tag == 'li':
            self.out.append('\n- ')
        elif tag == 'br':
            self.out.append('\n')
        elif tag == 'hr':
            self.out.append('\n---\n')
        elif tag == 'strong' and not self.pre:
            self.out.append('**')
        elif tag == 'em' and not self.pre:
            self.out.append('*')
        elif tag == 'code' and not self.pre:
            self.out.append('`')
        elif tag == 'sup':
            self.out.append('^')
        elif tag == 'sub':
            self.out.append('_')

    def handle_endtag(self, tag):
        if tag == 'pre':
            self.pre = max(0, self.pre - 1)
            self.out.append('\n```\n')
        elif tag in ('p', 'ul', 'ol', 'li'):
            self.out.append('\n')
        elif tag in ('strong', 'em') and not self.pre:
            self.out.append('**' if tag == 'strong' else '*')
        elif tag == 'code' and not self.pre:
            self.out.append('`')

    def handle_data(self, data):
        self.out.append(data)


def html2md(html: str) -> str:
    p = LC2MD()
    p.feed(html)
    md = ''.join(p.out).replace('\xa0', ' ')
    md = re.sub(r'[ \t]+\n', '\n', md)
    md = re.sub(r'```\n+', '```\n', md)
    md = re.sub(r'\n+```', '\n```', md)
    md = re.sub(r'\n{3,}', '\n\n', md)
    return md.strip()


def fetch(slug: str) -> dict:
    subprocess.run(
        ['node', FETCHER, slug],
        check=True, timeout=180,
        stdout=subprocess.DEVNULL,
    )
    return json.load(open(f'/tmp/lc-desc-{slug}.json', encoding='utf-8'))


def to_md(qd: dict, track: str = '未分组') -> str:
    title = qd.get('translatedTitle') or qd['title']
    diff = DIFF.get(qd['difficulty'], qd['difficulty'])
    tags = '、'.join((t.get('translatedName') or t['name']) for t in qd.get('topicTags') or [])
    body = html2md(qd['translatedContent'])
    return (f"# {qd['questionFrontendId']}. {title}\n\n"
            f"- 难度：{diff}\n- 标签：{tags}\n"
            f"- 链接：https://leetcode.cn/problems/{qd['titleSlug']}/\n"
            f"- 所属：{track}\n\n---\n\n{body}\n")


def md2comments(md_text: str) -> str:
    """markdown → 纯文本注释块（去 md 装饰：围栏/加粗/反引号，neovim 直读友好）"""
    lines = md_text.strip().splitlines()
    out, keep = [], False
    for ln in lines:
        if ln.strip() == '---':
            keep = True
            continue
        if not keep:
            continue
        if ln.strip() == '```':
            continue  # 去围栏行
        ln = ln.rstrip().replace('**', '').replace('`', '').replace('* *', ' ')
        out.append('// ' + ln if ln.strip() else '//')
    txt = '\n'.join(out)
    txt = re.sub(r'[ \t]+\n', '\n', txt)
    return re.sub(r'(//\n){2,}', '//\n', txt)


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        sys.exit(1)
    slug = sys.argv[1]
    go_file = sys.argv[2] if len(sys.argv) > 2 else None
    func = sys.argv[3] if len(sys.argv) > 3 else None

    os.makedirs(DESC_DIR, exist_ok=True)
    qd = fetch(slug)
    title = qd.get('translatedTitle') or qd['title']
    md = to_md(qd)
    clean_title = title.strip().replace(' ', '')
    md_path = os.path.join(DESC_DIR, f"{int(qd['questionFrontendId']):04d}-{clean_title}.md")
    with open(md_path, 'w', encoding='utf-8') as f:
        f.write(md)
    print(f'[1/2] 题面入库: {md_path}')

    if not (go_file and func):
        print('[2/2] 未指定 go 文件，跳过插入')
        return

    block = (f'// [题面-begin: {slug}]\n'
             + md2comments(md)
             + f'\n// [题面-end: {slug}]')
    src = open(go_file, encoding='utf-8').read()
    pat = re.compile(rf'// \[题面-begin: {re.escape(slug)}\].*?// \[题面-end: {re.escape(slug)}\]\n?', re.S)
    if pat.search(src):
        src = pat.sub(block + '\n', src)
        action = '已替换既有题面块'
    else:
        src = src.replace(f'func {func}(', block + '\nfunc ' + func + '(', 1)
        action = '已插入题面块'
    open(go_file, 'w', encoding='utf-8').write(src)
    print(f'[2/2] {action}: {go_file} -> func {func}')


if __name__ == '__main__':
    main()
