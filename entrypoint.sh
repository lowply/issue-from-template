#!/bin/sh

if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN is missing"
    exit 1
fi

/go/bin/main
