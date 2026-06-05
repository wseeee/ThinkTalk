#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_ROOT"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

PID_FILE="$PROJECT_ROOT/logs/pids.txt"

if [ ! -f "$PID_FILE" ] || [ ! -s "$PID_FILE" ]; then
    echo -e "${RED}[ERROR]${NC} 没有正在运行的服务记录 (logs/pids.txt 不存在或为空)"
    exit 1
fi

echo "[INFO] 正在停止所有 ThinkTalk 服务..."

killed=0
while IFS=" " read -r pid name; do
    [ -z "$pid" ] && continue
    if [[ "$pid" =~ ^[0-9]+$ ]]; then
        if kill "$pid" 2>/dev/null; then
            echo "  -> 已停止 $name (PID=$pid)"
            killed=$((killed + 1))
        else
            echo "  -> $name (PID=$pid) 已不存在"
        fi
    else
        # If it is a process name (e.g., applet-api, from .bat pids format)
        if taskkill //F //IM "${pid}.exe" &>/dev/null || pkill -f "${pid}" &>/dev/null || killall "${pid}" 2>/dev/null; then
            echo "  -> 已停止 $pid"
            killed=$((killed + 1))
        else
            echo "  -> $pid 已不存在"
        fi
    fi
done < "$PID_FILE"

rm -f "$PID_FILE"

echo ""
echo -e "${GREEN}[INFO]${NC} 共停止 $killed 个服务"
