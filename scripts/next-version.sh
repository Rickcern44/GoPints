#!/usr/bin/env bash
# Prints the next stable semver for a package (gopints|web) based on the
# Conventional Commits that touched its directory since its last
# <package>-v* tag. Prints nothing when there is nothing to release.
#
#   breaking change (type! or BREAKING CHANGE:) -> major
#   feat                                        -> minor
#   fix                                         -> patch
#
# Run from anywhere inside a clone with full history and tags.
set -euo pipefail

pkg="${1:?usage: next-version.sh <gopints|web>}"
prefix="${pkg}-v"

cd "$(git rev-parse --show-toplevel)"
last_tag="$(git tag --list "${prefix}*" --sort=-v:refname | head -1)"
range="${last_tag:+${last_tag}..}HEAD"
# The retired K8s manifests lived under <pkg>/infrastructure/; commits that
# only touched them must not trigger a release.
commits="$(git log "$range" --format='%B' -- "${pkg}/" ":(exclude)${pkg}/infrastructure")"

if grep -qE '^[a-z]+(\([^)]*\))?!:|BREAKING CHANGE:' <<<"$commits"; then
  bump=major
elif grep -qE '^feat(\([^)]*\))?:' <<<"$commits"; then
  bump=minor
elif grep -qE '^fix(\([^)]*\))?:' <<<"$commits"; then
  bump=patch
else
  exit 0
fi

current="${last_tag#"$prefix"}"
current="${current%%-*}"
IFS=. read -r major minor patch <<<"${current:-0.0.0}"

case "$bump" in
  major) echo "$((major + 1)).0.0" ;;
  minor) echo "${major}.$((minor + 1)).0" ;;
  patch) echo "${major}.${minor}.$((patch + 1))" ;;
esac
