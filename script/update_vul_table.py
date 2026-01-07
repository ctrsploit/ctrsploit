#!/usr/bin/env python3
"""
Script to update the vulnerability table in README.md from vul/README.md

This script:
1. Extracts the table from vul/README.md (8 columns)
2. Converts it to 4 columns (vul, desc, check, exploit)
3. Fixes link paths to use ./vul/ prefix
4. Updates the table in README.md
"""

import re
import sys
from pathlib import Path


def extract_module_tables_from_vul_readme(vul_readme_path):
    """
    从 vul/README.md 的 `## module` 章节下按模块拆分出各自的表格。

    返回值示例:
    [
        ("config", [table_line1, table_line2, ...]),
        ("runc", [table_line1, ...]),
        ...
    ]
    """
    with open(vul_readme_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    module_section_start = None
    module_section_end = len(lines)

    # 找到 "## module" 章节范围
    for i, line in enumerate(lines):
        if module_section_start is None and line.startswith('## module'):
            module_section_start = i
            continue
        if module_section_start is not None and line.startswith('## ') and not line.startswith('## module'):
            module_section_end = i
            break

    if module_section_start is None:
        raise ValueError('Could not find "## module" section in vul/README.md')

    modules = []
    current_module_name = None
    current_table_lines = []
    collecting_table = False

    for line in lines[module_section_start + 1:module_section_end]:
        # 模块小节，如 "### config"
        if line.startswith('### '):
            # 收尾前一个模块
            if current_module_name and current_table_lines:
                modules.append((current_module_name, current_table_lines))
            current_module_name = line.strip()[4:]  # 去掉 "### "
            current_table_lines = []
            collecting_table = False
            continue

        if current_module_name is None:
            # 还没进入第一个模块小节
            continue

        stripped = line.strip()

        # 表格行
        if stripped.startswith('|'):
            collecting_table = True
            current_table_lines.append(stripped)
            continue

        # 已经在收集表格，遇到非表格行则结束当前表格
        if collecting_table:
            collecting_table = False
            continue

    # 收尾最后一个模块
    if current_module_name and current_table_lines:
        modules.append((current_module_name, current_table_lines))

    if not modules:
        raise ValueError('Could not find any module tables in vul/README.md')

    return modules


def convert_table_row(row):
    """Convert a table row from 8 columns to 4 columns and fix links"""
    # Split by | and clean up
    parts = [p.strip() for p in row.split('|')]
    # Remove empty first and last parts (markdown table format)
    # But keep empty middle parts (empty columns)
    if parts and not parts[0]:
        parts = parts[1:]
    if parts and not parts[-1]:
        parts = parts[:-1]
    
    if len(parts) < 4:
        return row
    
    # Get first 4 columns
    vul_col = parts[0]
    desc_col = parts[1] if len(parts) > 1 else ''
    check_col = parts[2] if len(parts) > 2 else ''
    exploit_col = parts[3] if len(parts) > 3 else ''
    
    # Fix links in vul column: [text](path) -> [text](./vul/path)
    def fix_link(match):
        text = match.group(1)
        path = match.group(2)
        # If path doesn't start with ./ or http, add ./vul/
        if not path.startswith('./') and not path.startswith('http'):
            path = f"./vul/{path}"
        return f"[{text}]({path})"
    
    # Fix markdown links
    vul_col = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', fix_link, vul_col)
    
    # Build the new row
    new_row = f"| {vul_col} | {desc_col} | {check_col} | {exploit_col} |"
    
    return new_row


def update_readme_table(main_readme_path, module_tables):
    """根据 vul/README.md 的模块划分，更新 README.md 中的 `### module` 章节"""
    with open(main_readme_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    # Find the "### module" section and its end (next "###" header)
    module_header_idx = None
    next_header_idx = None

    for i, line in enumerate(lines):
        if module_header_idx is None and line.startswith('### module'):
            module_header_idx = i
            continue
        if module_header_idx is not None and line.startswith('### ') and not line.startswith('### module'):
            next_header_idx = i
            break

    if module_header_idx is None:
        raise ValueError('Could not find "### module" section in README.md')

    if next_header_idx is None:
        # If there's no following section header, treat end of file as section end
        next_header_idx = len(lines)

    # Preserve existing content inside module section *before* any auto-generated
    # content (旧表格 / 旧模块小节)，方便多次运行脚本。
    module_block = lines[module_header_idx + 1:next_header_idx]
    preserved_block = []
    auto_section_started = False

    for line in module_block:
        stripped = line.strip()
        # 旧的自动生成内容可能以表格头或 "#### 模块名" 开始
        if not auto_section_started and (
            stripped.startswith('| vul') and 'desc' in stripped
            or line.startswith('#### ')
        ):
            auto_section_started = True
            continue
        if auto_section_started:
            # 丢弃旧的自动生成部分
            continue
        preserved_block.append(line)

    # 构建新的模块章节内容：按模块拆分，每个模块一个子标题 + 4 列表
    new_module_blocks = []

    for module_name, table_lines in module_tables:
        # 使用 "####" 作为 README 中 module 章节下的子标题
        new_module_blocks.append(f'#### {module_name}\n\n')
        new_module_blocks.append('| vul | desc | check | exploit |\n')
        new_module_blocks.append('|-----|------|-------|---------|\n')

        for line in table_lines:
            if line.strip().startswith('| vul') and 'desc' in line:
                continue  # 跳过表头
            if '---' in line:
                continue  # 跳过分隔线
            converted = convert_table_row(line)
            new_module_blocks.append(converted + '\n')

        new_module_blocks.append('\n')

    # Replace the content inside the "module" section between the preserved
    # description (if any) and the next "###" header with the new module blocks.
    new_content = (
        ''.join(lines[:module_header_idx + 1]) +
        ''.join(preserved_block) +
        ''.join(new_module_blocks) +
        ''.join(lines[next_header_idx:])
    )
    
    # Write back
    with open(main_readme_path, 'w', encoding='utf-8') as f:
        f.write(new_content)
    
    print(f"Successfully updated table in {main_readme_path}")


def main():
    """Main function"""
    # Get script directory and repo root
    script_dir = Path(__file__).parent
    repo_root = script_dir.parent
    
    vul_readme = repo_root / 'vul' / 'README.md'
    main_readme = repo_root / 'README.md'
    
    # Check if files exist
    if not vul_readme.exists():
        print(f"Error: {vul_readme} not found", file=sys.stderr)
        sys.exit(1)
    
    if not main_readme.exists():
        print(f"Error: {main_readme} not found", file=sys.stderr)
        sys.exit(1)
    
    try:
        # 从 vul/README.md 中按模块拆分表格
        module_tables = extract_module_tables_from_vul_readme(vul_readme)

        # 更新 README.md 中的 module 章节
        update_readme_table(main_readme, module_tables)
        
        print("Table update completed successfully!")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()

