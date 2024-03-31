#!/bin/bash

set -Ceu

# ビルド
(
  cd /app/cmd/lambda-invoker
  rm main
  go build -o main .
)

sleep 20
