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
    """Regenerate the module tables inside the `## Vul` section of README.md.

    The table CONTENT follows vul/README.md (the source of truth); the display
    ORDER of existing module subsections is preserved as curated in README.md,
    and modules that are new in vul/README.md are appended. Non-module content
    before the first subsection (intro text, legend table) is preserved.
    """
    with open(main_readme_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    # Find the "## Vul" section and its end (next level-2 header).
    section_start = None
    section_end = len(lines)

    for i, line in enumerate(lines):
        if section_start is None and line.startswith('## Vul'):
            section_start = i
            continue
        if section_start is not None and i > section_start and line.startswith('## '):
            section_end = i
            break

    if section_start is None:
        raise ValueError('Could not find a "## Vul" section in README.md')

    # Split the section body into the preamble (intro text, legend) and the
    # existing `### <module>` subsections, whose bodies are regenerated.
    preamble = []
    existing_order = []
    for line in lines[section_start + 1:section_end]:
        if line.startswith('### '):
            existing_order.append(line.strip()[4:].strip())
        elif not existing_order:
            preamble.append(line)

    module_map = {name: table_lines for name, table_lines in module_tables}
    dropped = [name for name in existing_order if name not in module_map]
    if dropped:
        print(
            f"Warning: dropping README modules absent from vul/README.md: "
            f"{', '.join(dropped)}",
            file=sys.stderr,
        )

    # Preserve the curated order for known modules; append new ones in source
    # order so a fresh module always lands in the generated section.
    order = [name for name in existing_order if name in module_map]
    order += [name for name in module_map if name not in order]

    new_module_blocks = []
    for module_name in order:
        new_module_blocks.append(f'### {module_name}\n\n')
        new_module_blocks.append('| vul | desc | check | exploit |\n')
        new_module_blocks.append('|-----|------|-------|---------|\n')

        for line in module_map[module_name]:
            if line.startswith('| vul') and 'desc' in line:
                continue  # 跳过表头
            if '---' in line:
                continue  # 跳过分隔线
            converted = convert_table_row(line)
            new_module_blocks.append(converted + '\n')

        new_module_blocks.append('\n')

    preamble_text = ''.join(preamble).rstrip('\n')
    new_content = ''.join(lines[:section_start + 1])
    new_content += preamble_text + '\n\n' if preamble_text else '\n'
    new_content += ''.join(new_module_blocks)
    new_content += ''.join(lines[section_end:])

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

