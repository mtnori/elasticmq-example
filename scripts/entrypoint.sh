#!/bin/bash

# バイナリのパス
binary_path="/app/cmd/lambda-invoker/main"

# バイナリが存在するかを確認し、存在しない場合はビルドを待機
if [ ! -f "$binary_path" ]; then
    echo "Binary does not exist, waiting for myapp binary..."
    while [ ! -f "$binary_path" ]; do
        sleep 1
    done
fi

echo "Starting myapp..."
"$binary_path" &


# バイナリのハッシュ値が変更されるまでループ
while true; do
    # 現在のバイナリのハッシュ値を取得
    current_hash=$(md5sum "$binary_path" | awk '{ print $1 }')

    # バイナリのハッシュ値が変更されたかどうかをチェック
    if [ "$initial_hash" != "$current_hash" ]; then
        echo "Binary updated, restarting myapp..."
        # myappを再起動
        pkill -f myapp
        "$binary_path" &
        # バイナリのハッシュ値を更新
        initial_hash="$current_hash"
    fi

    # バイナリのハッシュ値が変更されていない場合は待機
    sleep 5
done
