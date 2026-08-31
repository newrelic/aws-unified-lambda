#!/usr/bin/env python3
"""Extract InlineCode: | python blocks from CFN/SAM yaml templates into real .py files,
named after the CloudFormation Logical ID, so SAST tools (Semgrep) can parse them."""
import re
import sys
from pathlib import Path

SRC_DIR = Path(sys.argv[1] if len(sys.argv) > 1 else ".")
OUT_DIR = Path(sys.argv[2] if len(sys.argv) > 2 else "/tmp/cwe-scan/extracted")
OUT_DIR.mkdir(parents=True, exist_ok=True)

resource_re = re.compile(r"^  (\w+):\s*$")
inline_re = re.compile(r"^(\s+)InlineCode:\s*\|\s*$")

count = 0
for yaml_file in sorted(SRC_DIR.glob("*.yaml")) + sorted(SRC_DIR.glob("*.yml")):
    lines = yaml_file.read_text().splitlines()
    current_logical_id = None
    i = 0
    while i < len(lines):
        line = lines[i]
        m = resource_re.match(line)
        if m:
            current_logical_id = m.group(1)
        m2 = inline_re.match(line)
        if m2:
            indent = len(m2.group(1)) + 2  # InlineCode's own indent + block scalar indent guess
            body = []
            j = i + 1
            base_indent = None
            while j < len(lines):
                l = lines[j]
                if l.strip() == "":
                    body.append("")
                    j += 1
                    continue
                cur_indent = len(l) - len(l.lstrip(" "))
                if base_indent is None:
                    base_indent = cur_indent
                if cur_indent < base_indent:
                    break
                body.append(l[base_indent:])
                j += 1
            name = current_logical_id or f"unknown_{count}"
            out_path = OUT_DIR / f"{yaml_file.stem}__{name}.py"
            out_path.write_text("\n".join(body) + "\n")
            print(f"extracted {out_path} (from {yaml_file.name}, Logical ID: {name})")
            count += 1
            i = j
            continue
        i += 1

print(f"\n{count} inline Lambda code block(s) extracted to {OUT_DIR}")
