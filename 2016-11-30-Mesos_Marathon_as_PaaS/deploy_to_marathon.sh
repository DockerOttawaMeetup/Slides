#!/usr/bin/env sh
MARATHON=http://$1:8080

echo "Installing $file..."
curl -X POST "$MARATHON/v2/apps" -d @"$2" -H "Content-type: application/json"
echo ""
