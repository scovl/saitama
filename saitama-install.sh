#!/bin/bash
set -e

SAITAMA_URL="https://github.com/lobocode/saitama/raw/master/bin/saitama"
SAITAMA_PATH="/usr/local/bin/saitama"

wget -O "$SAITAMA_PATH" "$SAITAMA_URL"
chmod 755 "$SAITAMA_PATH"