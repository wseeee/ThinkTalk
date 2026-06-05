#!/usr/bin/env bash
set -euo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$PROJECT_ROOT"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC}  $*"; }
err()  { echo -e "${RED}[ERROR]${NC} $*"; }

BIN_DIR="$PROJECT_ROOT/bin"
LOG_DIR="$PROJECT_ROOT/logs"

# ---------- 服务定义 ----------
# 格式: "服务名|包路径|配置文件(相对于项目根目录)"
SERVICES=(
  # API
  "applet-api|./application/applet|application/applet/etc/applet-api.yaml"
  "article-api|./application/article/api|application/article/api/etc/article-api.yaml"
  "chat-api|./application/chat/api|application/chat/api/etc/chat-api.yaml"
  "qa-api|./application/qa/api|application/qa/api/etc/qa-api.yaml"
  # RPC
  "user-rpc|./application/user/rpc|application/user/rpc/etc/user.yaml"
  "article-rpc|./application/article/rpc|application/article/rpc/etc/article.yaml"
  "chat-rpc|./application/chat/rpc|application/chat/rpc/etc/chat.yaml"
  "concerned-rpc|./application/concerned/rpc|application/concerned/rpc/etc/concerned.yaml"
  "follow-rpc|./application/follow/rpc|application/follow/rpc/etc/follow.yaml"
  "like-rpc|./application/like/rpc|application/like/rpc/etc/like.yaml"
  "member-rpc|./application/member/rpc|application/member/rpc/etc/member.yaml"
  "message-rpc|./application/message/rpc|application/message/rpc/etc/message.yaml"
  "qa-rpc|./application/qa/rpc|application/qa/rpc/etc/qa.yaml"
  "reply-rpc|./application/reply/rpc|application/reply/rpc/etc/reply.yaml"
  "tag-rpc|./application/tag/rpc|application/tag/rpc/etc/tag.yaml"
  # MQ
  "article-mq|./application/article/mq|application/article/mq/etc/article.yaml"
  "chat-mq|./application/chat/mq|application/chat/mq/etc/chat.yaml"
  "concerned-mq|./application/concerned/mq|application/concerned/mq/etc/concerned.yaml"
  "like-mq|./application/like/mq|application/like/mq/etc/like.yaml"
  "member-mq|./application/member/mq|application/member/mq/etc/member.yaml"
  "message-mq|./application/message/mq|application/message/mq/etc/message.yaml"
  "qa-mq|./application/qa/mq|application/qa/mq/etc/qa.yaml"
  "reply-mq|./application/reply/mq|application/reply/mq/etc/reply.yaml"
)

total=${#SERVICES[@]}

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   ThinkTalk - 本地一键启动 ($total 个服务)${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# ---------- 前置检查 ----------
if ! command -v go &>/dev/null; then
    err "未找到 Go，请先安装 Go SDK"
    exit 1
fi

if [ ! -f .env ]; then
    warn ".env 文件不存在，将从 .env.example 复制"
    cp .env.example .env
    log "已创建 .env，请按需修改配置后重新运行"
    exit 0
fi

# Detect OS to determine executable extension
EXE_EXT=""
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    EXE_EXT=".exe"
fi

mkdir -p "$BIN_DIR" "$LOG_DIR"

# ---------- 编译 ----------
echo -e "${CYAN}[1/2] 编译服务...${NC}"
built=0
pids=()
for svc_def in "${SERVICES[@]}"; do
    IFS="|" read -r name pkg _ <<< "$svc_def"
    printf "  \033[0;36m编译 %-20s\033[0m %s\r" "$name" "..."
    go build -o "$BIN_DIR/${name}${EXE_EXT}" "$pkg" &
    pids+=($!)
    built=$((built + 1))
done

# 等待所有编译完成
failed=0
i=0
for svc_def in "${SERVICES[@]}"; do
    IFS="|" read -r name pkg _ <<< "$svc_def"
    wait "${pids[$i]}" || {
        err "编译失败: $name"
        failed=$((failed + 1))
    }
    i=$((i + 1))
done

if [ $failed -gt 0 ]; then
    err "$failed 个服务编译失败，请检查错误信息"
    exit 1
fi
log "全部 $built 个服务编译成功"
echo ""

# ---------- 启动 ----------
echo -e "${CYAN}[2/2] 启动服务...${NC}"
> "$LOG_DIR/pids.txt"

for svc_def in "${SERVICES[@]}"; do
    IFS="|" read -r name _ config <<< "$svc_def"
    exe="$BIN_DIR/${name}${EXE_EXT}"
    if [ ! -f "$exe" ]; then
        warn "$name: $exe 不存在，跳过"
        continue
    fi
    log_file="$LOG_DIR/${name}.log"
    "$exe" -f "$config" > "$log_file" 2>&1 &
    pid=$!
    echo "$pid $name" >> "$LOG_DIR/pids.txt"
    printf "  \033[0;32m启动 %-20s\033[0m PID=%s\n" "$name" "$pid"
done

echo ""
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   端口映射${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""
echo -e "  API 服务:"
echo -e "    applet-api   → http://localhost:8888"
echo -e "    article-api  → http://localhost:80"
echo -e "    chat-api     → http://localhost:8087"
echo -e "    qa-api       → http://localhost:8890"
echo ""
echo -e "  RPC 服务 (gRPC):"
echo -e "    user-rpc      :8080    article-rpc   :9090"
echo -e "    chat-rpc      :8087    concerned-rpc :8086"
echo -e "    follow-rpc    :8081    like-rpc      :8082"
echo -e "    member-rpc    :8091    message-rpc   :8085"
echo -e "    qa-rpc        :8090    reply-rpc     :8083"
echo -e "    tag-rpc       :8084"
echo ""
echo -e "${GREEN}全部启动完成!${NC}"
echo ""
echo -e "  查看单个日志:  tail -f logs/<服务名>.log"
echo -e "  查看全部日志:  tail -f logs/*.log"
echo -e "  停止所有服务:  ./stop-all.sh"
echo ""
