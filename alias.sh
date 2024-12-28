#!/bin/bash

PROJECT_DIR=$(dirname "$(realpath "$0")")
ORIGINAL_DIR="$(pwd)"
cd $PROJECT_DIR || exit
zig build run -- "$ORIGINAL_DIR" "$@"
