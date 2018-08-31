#!/usr/bin/env bash

echo "==> Checking that code complies with gofmt requirements..."

gofmt_files="$(\
    find -name "*.go" -not -path "./vendor/*" \
        -exec gofmt -l {} ';')"

if [[ -n ${gofmt_files} ]]; then 
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo
    exit 1
fi
