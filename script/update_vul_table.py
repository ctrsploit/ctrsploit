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


def extract_table_from_vul_readme(vul_readme_path):
    """Extract the table from vul/README.md"""
    with open(vul_readme_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    # Find the table start (line with | vul | desc | ...)
    table_start = None
    for i, line in enumerate(lines):
        if line.strip().startswith('| vul') and 'desc' in line:
            table_start = i
            break
    
    if table_start is None:
        raise ValueError("Could not find table start in vul/README.md")
    
    # Extract table lines (skip the header separator line)
    table_lines = []
    for i in range(table_start, len(lines)):
        line = lines[i].strip()
        if not line or not line.startswith('|'):
            break
        if '---' in line:  # Skip separator line
            continue
        table_lines.append(line)
    
    return table_lines


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


def update_readme_table(main_readme_path, new_table_lines):
    """Update the table in README.md"""
    with open(main_readme_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    # Find the table start and end
    table_start = None
    table_end = None
    
    for i, line in enumerate(lines):
        if table_start is None and line.strip().startswith('| vul') and 'desc' in line:
            table_start = i
        elif table_start is not None:
            # Table ends when we hit a line that doesn't start with | or is empty followed by ###
            if not line.strip().startswith('|'):
                # Check if next non-empty line is a section header
                j = i
                while j < len(lines) and not lines[j].strip():
                    j += 1
                if j < len(lines) and lines[j].startswith('###'):
                    table_end = i
                    break
    
    if table_start is None:
        raise ValueError("Could not find table start in README.md")
    
    if table_end is None:
        raise ValueError("Could not find table end in README.md")
    
    # Build the new table with optimized header
    new_table_lines_list = []
    
    # Add optimized header (compact format)
    new_table_lines_list.append('| vul | desc | check | exploit |\n')
    new_table_lines_list.append('|-----|------|-------|---------|\n')
    
    # Add all data rows (skip header from new_table_lines)
    for line in new_table_lines:
        if line.strip().startswith('| vul') and 'desc' in line:
            continue  # Skip header
        if '---' in line:
            continue  # Skip separator
        converted = convert_table_row(line)
        new_table_lines_list.append(converted + '\n')
    
    # Replace the table section
    new_content = (
        ''.join(lines[:table_start]) +
        ''.join(new_table_lines_list) +
        ''.join(lines[table_end:])
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
        # Extract table from vul/README.md
        table_lines = extract_table_from_vul_readme(vul_readme)
        
        # Update README.md
        update_readme_table(main_readme, table_lines)
        
        print("Table update completed successfully!")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()

