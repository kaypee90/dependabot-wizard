#!/bin/bash
#
# Builds the React app and replaces the copy embedded in the Go binary.
# The destination is cleared first: copying over it left every previous
# bundle behind, and //go:embed shipped all of them.

set -euo pipefail

DESTINATION="../../cmd/depbot/static"

npm run build

rm -rf "$DESTINATION"
mkdir -p "$DESTINATION"
cp -R build/. "$DESTINATION"

echo "Embedded $(find "$DESTINATION" -type f | wc -l | tr -d ' ') files into $DESTINATION"
