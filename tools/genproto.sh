#!/usr/bin/env bash
set -euo pipefail

# 用法：./tools/genproto.sh proto
# 默认扫描 proto 目录下所有 .proto
PROTO_DIR="${1:-proto}"

protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  $(find "${PROTO_DIR}" -name "*.proto")