#!/usr/bin/env bash
set -e

rm -rf .gopath
mkdir -p .gopath/src/fourth.com/
ln -sf ../../.. .gopath/src/fourth.com/ratelimit
export GOPATH="$(pwd)/.gopath:$(pwd)/vendor"
go build -o ratelimit
