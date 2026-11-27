#!/bin/sh
# Print a per-severity summary of golangci-lint issues (project health).
# Severities come from the severity rules in .golangci.yml and are carried
# by the SARIF output. Exits 1 when at least one error-level issue is found.
set -eu

command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found (see .tool-versions)"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "jq is required for the health report"; exit 1; }

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

golangci-lint run --output.text.path="$tmp/lint.txt" --output.sarif.path="$tmp/report.sarif" || true

jq -r '
  (.runs[0].results // []) as $r
  | "project health (golangci-lint, all linters):",
    (["error", "warning", "note"] | map(. as $l | "\($l): \([$r[] | select((.level // "note") == $l)] | length)") | .[]),
    "total: \($r | length)"
' "$tmp/report.sarif"

errors=$(jq '[.runs[0].results[]? | select((.level // "note") == "error")] | length' "$tmp/report.sarif")
if [ "$errors" -gt 0 ]; then
	echo "failing: $errors error-level issue(s)"
	exit 1
fi
