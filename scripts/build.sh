#!/bin/bash

# ビルド
cd /app/cmd/lambda-invoker || exit
go build -o main .
