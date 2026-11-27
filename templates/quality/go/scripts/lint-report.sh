#!/bin/sh
# Print a per-severity summary of golangci-lint issues (project health).
# Severities come from the severity rules in .golangci.yml and are carried by
# the JSON output (.Issues[].Severity). Exits 1 when error-level issues exist.
set -eu

command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found (see .tool-versions)"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is required for the health report"; exit 1; }

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

golangci-lint run --output.text.path="$tmp/lint.txt" --output.json.path="$tmp/report.json" || true

jq -r '
  (.Issues // []) as $i
  | ($i | map(.Severity // "warning") | group_by(.)) as $g
  | "project health (golangci-lint, all linters):",
    ($g[] | "  \(.[0]): \(length) issue(s)"),
    "  total: \($i | length) issue(s)"
' "$tmp/report.json"

errors=$(jq '[.Issues[]? | select(.Severity == "error")] | length' "$tmp/report.json")
if [ "$errors" -gt 0 ]; then
	echo "failing: $errors error-level issue(s)"
	exit 1
fi
