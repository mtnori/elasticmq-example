#!/bin/bash

set -Ceu

# ビルド
(
  cd /app/cmd/lambda-invoker
  go build -o main .
)

sleep 20
