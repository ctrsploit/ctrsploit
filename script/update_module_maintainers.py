#!/usr/bin/env python3
"""
Update module README maintainer front matter from git history.

The module scope is the README parent directory. Maintainers are ranked by
commit count within that directory, then by most recent commit.
"""

import argparse
import collections
import json
import subprocess
import sys
from pathlib import Path


DEFAULT_MAINTAINER = "ssst0n3"

EMAIL_TO_MAINTAINER = {
    "ssst0n3@gmail.com": "ssst0n3",
    "st0n3@debian.home": "ssst0n3",
    "r081n4k@gmail.com": "r0binak",
    "hpdoger@d0g3.cn": "Hpd0ger",
    "axsl666@foxmail.com": "Axsl666",
    "normalbe@163.com": "normalbe",
    "test@example.com": "noirfate",
}

NAME_TO_MAINTAINER = {
    "LEI WANG": "ssst0n3",
    "Lei Wang": "ssst0n3",
    "st0n3": "ssst0n3",
    "Sergey K.": "r0binak",
}


def repo_root():
    return Path(__file__).resolve().parent.parent


def load_module_docs(modules_json):
    with modules_json.open("r", encoding="utf-8") as f:
        catalog = json.load(f)

    docs = set()

    def collect(module):
        doc = module.get("doc")
        if doc:
            doc = doc.split("#", 1)[0]
            docs.add(doc.removeprefix("./"))
        for child in module.get("children", []):
            collect(child)

    for module_list in catalog["modules"].values():
        for module in module_list:
            collect(module)

    return sorted(docs)


def git_contributors(root, scope):
    result = subprocess.run(
        [
            "git",
            "log",
            "--format=%aI%x00%aN%x00%aE",
            "--",
            str(scope),
        ],
        cwd=root,
        text=True,
        capture_output=True,
        check=True,
    )

    counts = collections.Counter()
    latest_index = {}
    for index, line in enumerate(result.stdout.splitlines()):
        if not line:
            continue
        date, name, email = line.split("\0")
        maintainer = identity_to_maintainer(name, email)
        counts[maintainer] += 1
        latest_index.setdefault(maintainer, index)

    return [
        maintainer
        for maintainer, _ in sorted(
            counts.items(),
            key=lambda item: (-item[1], latest_index[item[0]], item[0].lower()),
        )
    ]


def identity_to_maintainer(name, email):
    maintainer = EMAIL_TO_MAINTAINER.get(email.lower())
    if maintainer:
        return maintainer
    maintainer = NAME_TO_MAINTAINER.get(name)
    if maintainer:
        return maintainer
    if name:
        return name
    if email:
        return email.split("@", 1)[0]
    return DEFAULT_MAINTAINER


def render_maintainer_block(maintainers):
    lines = ["maintainer:\n"]
    for maintainer in maintainers:
        lines.append(f"    - {maintainer}\n")
    return lines


def is_front_matter_delimiter(line):
    return line.strip() == "---"


def split_front_matter(lines):
    if not lines or not is_front_matter_delimiter(lines[0]):
        return [], lines

    for index, line in enumerate(lines[1:], start=1):
        if is_front_matter_delimiter(line):
            return lines[1:index], lines[index + 1 :]

    return [], lines


def top_level_key(line):
    if line.startswith((" ", "\t", "-")):
        return ""
    if ":" not in line:
        return ""
    return line.split(":", 1)[0].strip()


def remove_key_block(front_matter, key):
    result = []
    index = 0
    while index < len(front_matter):
        if top_level_key(front_matter[index]) != key:
            result.append(front_matter[index])
            index += 1
            continue

        index += 1
        while index < len(front_matter):
            line = front_matter[index]
            if line.strip() and top_level_key(line):
                break
            index += 1

    return result


def block_end(front_matter, start):
    index = start + 1
    while index < len(front_matter):
        line = front_matter[index]
        if line.strip() and top_level_key(line):
            break
        index += 1
    return index


def maintainer_insert_index(front_matter):
    for preferred_key in ("author", "tags"):
        for index, line in enumerate(front_matter):
            if top_level_key(line) == preferred_key:
                return block_end(front_matter, index)
    return 0


def update_readme_content(content, maintainers):
    lines = content.splitlines(keepends=True)
    front_matter, body = split_front_matter(lines)
    front_matter = remove_key_block(front_matter, "maintainer")

    insert_at = maintainer_insert_index(front_matter)
    new_front_matter = (
        front_matter[:insert_at]
        + render_maintainer_block(maintainers)
        + front_matter[insert_at:]
    )

    return "".join(["---\n"] + new_front_matter + ["---\n"] + body)


def update_readme(path, maintainers, check):
    content = path.read_text(encoding="utf-8")
    updated = update_readme_content(content, maintainers)
    if updated == content:
        return False
    if check:
        return True
    path.write_text(updated, encoding="utf-8")
    return True


def main():
    parser = argparse.ArgumentParser(
        description="Update module README maintainer front matter from git history.",
    )
    parser.add_argument(
        "--modules-json",
        default="modules.json",
        help="module catalog path relative to the repository root",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="fail if any README maintainer metadata is not up to date",
    )
    args = parser.parse_args()

    root = repo_root()
    modules_json = root / args.modules_json
    changed = []

    for doc in load_module_docs(modules_json):
        readme = root / doc
        if not readme.exists():
            print(f"warning: skip missing README: {doc}", file=sys.stderr)
            continue

        maintainers = git_contributors(root, Path(doc).parent)
        if not maintainers:
            maintainers = [DEFAULT_MAINTAINER]

        if update_readme(readme, maintainers, args.check):
            changed.append(doc)
            action = "would update" if args.check else "updated"
            print(f"{action} {doc}: {', '.join(maintainers)}")

    if args.check and changed:
        print(
            f"error: {len(changed)} README maintainer entries are out of date",
            file=sys.stderr,
        )
        sys.exit(1)

    suffix = "out of date" if args.check else "changed"
    print(f"module maintainer update complete: {len(changed)} {suffix}")


if __name__ == "__main__":
    main()
