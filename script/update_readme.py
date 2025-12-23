#!/usr/bin/env python3
"""
Script to automatically update README.md

This script updates:
1. Command help sections (env, vul, exploit, checksec) by running ctrsploit commands
2. Vulnerability table (by calling update_vul_table.py or its logic)
"""

import re
import subprocess
import sys
from pathlib import Path


# Path to ctrsploit binary
CTRSPLOIT_BIN = Path(__file__).parent.parent / "bin" / "latest" / "ctrsploit_linux_amd64"

# Commands to update help for
COMMANDS = {
    "env": {
        "section": "### env",
        "command": "env",
        "prompt": "$ ctrsploit env"
    },
    "vul": {
        "section": "### vul",
        "command": "vul",
        "prompt": "$ ctrsploit vul"
    },
    "exploit": {
        "section": "### exploit",
        "command": "exploit",
        "prompt": "$ ctrsploit exploit"
    },
    "checksec": {
        "section": "### checksec",
        "command": "checksec",
        "prompt": "$ ctrsploit checksec",
        # Note: checksec auto output is environment-dependent and contains warnings/errors
        # It's better to keep the example output manually curated rather than auto-updated
        # "examples": [
        #     {
        #         "command": ["--colorful", "checksec", "auto"],
        #         "prompt": "$ ctrsploit --colorful checksec auto"
        #     }
        # ]
    },
}


def get_command_help(command):
    """Run ctrsploit command and get help output"""
    if not CTRSPLOIT_BIN.exists():
        raise FileNotFoundError(f"ctrsploit binary not found at {CTRSPLOIT_BIN}")
    
    try:
        result = subprocess.run(
            [str(CTRSPLOIT_BIN), command, "--help"],
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        raise RuntimeError(f"Failed to run ctrsploit {command} --help: {e.stderr}")


def get_command_output(command_args):
    """Run ctrsploit command with given arguments and get output"""
    if not CTRSPLOIT_BIN.exists():
        raise FileNotFoundError(f"ctrsploit binary not found at {CTRSPLOIT_BIN}")
    
    # Don't use check=True, so we can capture output even if command fails
    result = subprocess.run(
        [str(CTRSPLOIT_BIN)] + command_args,
        capture_output=True,
        text=True,
        check=False
    )
    # Return stdout even if command failed (checksec auto might fail but still produce output)
    return result.stdout if result.stdout else ""


def update_command_section_in_content(content, command_name, command_info):
    """Update a command help section in README.md content"""
    # Get the help output
    help_output = get_command_help(command_info["command"])
    
    # Build the new code block content (everything inside ```shell ... ```)
    prompt = command_info["prompt"]
    help_output_clean = help_output.rstrip()
    new_code_block_content = f"{prompt}\n{help_output_clean}"
    
    # Find the section header first to locate the right code block
    section_header = command_info['section']
    section_pattern = rf"({re.escape(section_header)})"
    section_match = re.search(section_pattern, content)
    
    if not section_match:
        print(f"Warning: Could not find {command_name} section header")
        return None
    
    # Find the code block after the section (first ```shell after section header)
    section_start = section_match.end()
    remaining_content = content[section_start:]
    
    # Match the first code block: ```shell + prompt (with optional trailing spaces) + content + ```
    # Allow trailing spaces after prompt line
    prompt_escaped = re.escape(prompt.rstrip())
    prompt_pattern = prompt_escaped + r'\s*'
    code_block_pattern = rf"(```shell\n{prompt_pattern}.*?\n```)"
    
    code_block_match = re.search(code_block_pattern, remaining_content, re.DOTALL)
    if code_block_match:
        # Replace the code block content
        new_code_block = f"```shell\n{new_code_block_content}\n```"
        new_content = (
            content[:section_start + code_block_match.start(1)] +
            new_code_block +
            content[section_start + code_block_match.end(1):]
        )
        print(f"Updated {command_name} section")
        return new_content
    else:
        print(f"Warning: Could not find code block for {command_name} section")
        return None


def update_example_output_in_content(content, command_name, command_info):
    """Update example output code blocks for a command"""
    if "examples" not in command_info:
        return content
    
    section_header = command_info['section']
    
    # Find the section header
    section_pattern = rf"({re.escape(section_header)})"
    section_match = re.search(section_pattern, content)
    if not section_match:
        return content
    
    # Update each example output
    for example in command_info["examples"]:
        # Skip if explicitly disabled
        if example.get("skip", False):
            print(f"Skipping example (disabled): {example.get('prompt', 'unknown')}")
            continue
        
        example_prompt = example["prompt"]
        example_command = example["command"]
        
        # Get the example output
        example_output = get_command_output(example_command)
        if not example_output:
            print(f"Warning: No output for example: {example_prompt}")
            continue
        
        # Filter output if filter function is provided
        if "filter" in example and callable(example["filter"]):
            example_output = example["filter"](example_output)
        
        example_output_clean = example_output.rstrip()
        new_example_code_block = f"```shell\n{example_prompt}\n{example_output_clean}\n```"
        
        # Find the code block with this prompt in the entire content
        # Escape prompt but allow optional trailing spaces
        example_prompt_escaped = re.escape(example_prompt.rstrip()) + r'\s*'
        example_pattern = rf"(```shell\n{example_prompt_escaped}.*?\n```)"
        example_match = re.search(example_pattern, content, re.DOTALL)
        
        if example_match:
            # Replace the example code block
            content = (
                content[:example_match.start(1)] +
                new_example_code_block +
                content[example_match.end(1):]
            )
            print(f"Updated example output: {example_prompt}")
        else:
            print(f"Warning: Could not find example code block: {example_prompt}")
    
    return content


def update_all_command_sections(readme_path):
    """Update all command help sections"""
    with open(readme_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    for command_name, command_info in COMMANDS.items():
        print(f"Updating {command_name} help section...")
        updated_content = update_command_section_in_content(content, command_name, command_info)
        if updated_content:
            content = updated_content
        else:
            print(f"Warning: Failed to update {command_name} section, skipping...")
        
        # Update example outputs if any
        if "examples" in command_info:
            print(f"Updating {command_name} example outputs...")
            content = update_example_output_in_content(content, command_name, command_info)
    
    # Write back all updates at once
    with open(readme_path, 'w', encoding='utf-8') as f:
        f.write(content)


def update_vul_table(readme_path):
    """Update vulnerability table by calling update_vul_table.py"""
    script_path = Path(__file__).parent / "update_vul_table.py"
    if not script_path.exists():
        print("Warning: update_vul_table.py not found, skipping vulnerability table update")
        return
    
    try:
        result = subprocess.run(
            [sys.executable, str(script_path)],
            check=True,
            capture_output=True,
            text=True
        )
        print(result.stdout)
    except subprocess.CalledProcessError as e:
        print(f"Error updating vulnerability table: {e.stderr}", file=sys.stderr)
        raise


def main():
    """Main function"""
    repo_root = Path(__file__).parent.parent
    readme_path = repo_root / "README.md"
    
    if not readme_path.exists():
        print(f"Error: {readme_path} not found", file=sys.stderr)
        sys.exit(1)
    
    if not CTRSPLOIT_BIN.exists():
        print(f"Error: ctrsploit binary not found at {CTRSPLOIT_BIN}", file=sys.stderr)
        print("Please build the binary first: make binary", file=sys.stderr)
        sys.exit(1)
    
    try:
        # Update command help sections
        print("Updating command help sections...")
        update_all_command_sections(readme_path)
        
        # Update vulnerability table
        print("\nUpdating vulnerability table...")
        update_vul_table(readme_path)
        
        print("\nREADME.md update completed successfully!")
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == '__main__':
    main()

