#!/bin/bash
#
# Guards the embedded copy of the React build against the two ways it has
# broken before:
#
#   1. index.html referencing asset paths the Go file server cannot resolve
#      (the committed file had been hand-edited and no rebuild reproduced it).
#   2. Stale bundles accumulating, because the copy step never cleared the
#      destination and //go:embed shipped every past build.

set -euo pipefail

STATIC_DIR="${1:-cmd/depbot/static}"
INDEX="$STATIC_DIR/index.html"

if [ ! -f "$INDEX" ]; then
  echo "FAIL: $INDEX does not exist. Run 'make ui'."
  exit 1
fi

status=0

# 1. Every asset index.html references must exist at that path in the tree
#    the binary embeds, since the app is served from the root of it.
while read -r asset; do
  [ -z "$asset" ] && continue
  case "$asset" in
    http*|data:*|"#"*) continue ;;
  esac
  if [ ! -f "$STATIC_DIR/${asset#/}" ]; then
    echo "FAIL: index.html references $asset but $STATIC_DIR/${asset#/} is missing."
    status=1
  fi
done < <(grep -oE '(src|href)="[^"]*"' "$INDEX" | sed -E 's/.*="([^"]*)".*/\1/')

# 2. Exactly one application bundle should be present.
bundle_count=$(find "$STATIC_DIR" -name 'main.*.js' -type f | wc -l | tr -d ' ')
if [ "$bundle_count" -ne 1 ]; then
  echo "FAIL: expected 1 main bundle in $STATIC_DIR, found $bundle_count. Run 'make ui'."
  find "$STATIC_DIR" -name 'main.*.js' -type f | sed 's/^/  /'
  status=1
fi

if [ "$status" -eq 0 ]; then
  echo "OK: embedded assets are self-consistent ($bundle_count bundle)."
fi

exit $status
