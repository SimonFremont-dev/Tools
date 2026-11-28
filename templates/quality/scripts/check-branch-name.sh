#!/bin/sh
# Fail when a branch does not follow conventional branch names
# (https://conventionalbranch.org/): <type>/<description>, types from
# conventional commits. main/master are exempt.
# Usage: check-branch-name.sh [branch]   (default: current branch)
set -eu

branch="${1:-$(git branch --show-current)}"

case "$branch" in
  main|master) exit 0 ;;
esac

if echo "$branch" | grep -Eq '^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert|hotfix)(/[a-z0-9][a-z0-9-]*)?$'; then
	exit 0
fi

echo "branch '$branch' does not follow conventional branch names (https://conventionalbranch.org/)"
echo "expected <type>/<description>, e.g. feat/user-login or fix/crash-on-startup"
exit 1
