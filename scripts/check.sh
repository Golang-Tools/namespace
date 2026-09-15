#!/usr/bin/env bash
# namespace 本地校验脚本:gofmt / go vet / go test / race / Example / benchmark smoke。
# 每次改动后运行一次(本仓库不依赖 GitHub Actions)。
set -euo pipefail
cd "$(dirname "$0")/.."

echo "==> gofmt"
if [ -n "$(gofmt -l .)" ]; then
    gofmt -l .
    echo "FAIL: 存在未格式化的文件" >&2
    exit 1
fi

echo "==> go vet"
go vet ./...

echo "==> go test"
go test ./... -count=1

echo "==> go test -race"
go test ./... -race -count=1

echo "==> Example tests"
go test ./... -run '^Example' -count=1

echo "==> benchmark smoke"
go test ./... -run '^$' -bench . -benchtime=1x

echo "==> PASS"
