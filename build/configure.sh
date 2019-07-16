#!/bin/sh

set -e

if [ ! -f "build/configure.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

mkdir -p "$PWD/core/token"
sed 's/package core/package token/g' "$PWD/core/error.go" > "$PWD/core/token/error.go"
sed 's/package core/package token/g' "$PWD/core/events.go" > "$PWD/core/token/events.go"
sed 's/package core/package token/g' "$PWD/core/evm.go" > "$PWD/core/token/evm.go"
sed 's/package core/package token/g' "$PWD/core/gaspool.go" > "$PWD/core/token/gaspool.go"
sed 's/package core/package token/g' "$PWD/core/state_transition.go" > "$PWD/core/token/state_transition.go"
sed 's/package core/package token/g' "$PWD/core/tx_cacher.go" > "$PWD/core/token/tx_cacher.go"
sed 's/package core/package token/g' "$PWD/core/tx_journal.go" > "$PWD/core/token/tx_journal.go"
sed 's/package core/package token/g' "$PWD/core/tx_list.go" > "$PWD/core/token/tx_list.go"
sed 's/package core/package token/g' "$PWD/core/tx_pool.go" > "$PWD/core/token/tx_pool.go"
