#!/bin/bash
set -euo pipefail

# This is used by the linter action.
# Recursively finds all directories with a go.mod file and creates
# a GitHub Actions JSON output option.

echo "Resolving modules in $(pwd)"

PROJECT_HOME=$(
	cd "$(dirname "${BASH_SOURCE[0]}")" &&
		cd .. &&
		pwd
)

source "${PROJECT_HOME}/.hack/util.sh"

FAILURE_FILE=${PROJECT_HOME}/.hack/.lintcheck_failures

all_modules=$(util::find_modules)
failing_modules=()
while IFS='' read -r line; do failing_modules+=("$line"); done < <(cat "$FAILURE_FILE")

echo "Ignored failing modules:"
echo "${failing_modules[*]}"
echo

PATHS=""

for mod in $all_modules; do
	echo "Checking module: $mod"
	util::array_contains "$mod" "${failing_modules[*]}" && in_failing=$? || in_failing=$?
	if [[ "$in_failing" -ne "0" ]]; then
		PATHS+=$(printf '{"workdir":"%s"},' ${mod})
	fi
done

echo "matrix={\"include\":[${PATHS%?}]}" >>$GITHUB_OUTPUT

MODULES=$(find . -name "go.mod" -type f -not -path "*/vendor/*" -not -path "*/\.*/*" | sort | xargs -I{} dirname {})
JSON='{"include":['

JSON="${JSON}{\"workdir\":\".\", \"module\":\"main\"}"

for MODULE in $MODULES; do
	if [ "$MODULE" != "." ]; then
		MODULE_NAME=$(basename "$MODULE")
		JSON="${JSON},{\"workdir\":\"$MODULE\", \"module\":\"$MODULE_NAME\"}"
	fi
done

JSON="${JSON}]}"

echo "matrix=$JSON" >>$GITHUB_OUTPUT

echo "$JSON"
