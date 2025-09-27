#!/bin/bash

PLUGIN_BIN="./protoc-gen-dict"

# 使用 dlv 来启动插件，并连接到 GoLand 的远程调试服务器
dlv exec --headless --listen=:2345 --api-version=2 --accept-multiclient \
  --log-dest=2 \
  "$PLUGIN_BIN" -- "$@"