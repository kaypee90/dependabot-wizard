#!/bin/bash

set -e
npm run build && cp -fr build/* ../../cmd/depbot/static

echo "Done"
