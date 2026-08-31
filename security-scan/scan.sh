#!/usr/bin/env bash
# Runs Semgrep's public security rulesets (community-maintained, covers a broad set of
# CWEs -- log injection, CRLF injection, and more) against this repo's Lambda code.
#
# The inline Lambda handlers live as InlineCode: | blocks inside the CFN/SAM yaml
# templates, so they can't be scanned directly -- this pulls each block out into a
# real .py file (named <template>__<LogicalID>.py) first, then runs Semgrep on those.
#
# Requires network access to fetch rules from the Semgrep registry (works on
# GitHub-hosted runners; may be blocked by a local corporate proxy/SSL setup).
#
# Usage:
#   security-scan/scan.sh              # scan all templates in repo root
set -euo pipefail
cd "$(dirname "$0")/.."

command -v semgrep >/dev/null 2>&1 || { echo "semgrep not found. Install with: pipx install semgrep" >&2; exit 1; }

OUT_DIR="$(mktemp -d)"
python3 security-scan/extract_inline_lambdas.py . "$OUT_DIR"
echo
semgrep --config p/security-audit --metrics=off --error "$OUT_DIR"
