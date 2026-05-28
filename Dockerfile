FROM golang:1.26-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 一次性顺序编译全部 23 个服务
RUN set -e; \
    for pkg in \
        application/applet \
        application/article/api \
        application/article/rpc \
        application/article/mq \
        application/chat/api \
        application/chat/rpc \
        application/chat/mq \
        application/concerned/rpc \
        application/concerned/mq \
        application/follow/rpc \
        application/like/rpc \
        application/like/mq \
        application/member/rpc \
        application/member/mq \
        application/message/rpc \
        application/message/mq \
        application/qa/api \
        application/qa/rpc \
        application/qa/mq \
        application/reply/rpc \
        application/reply/mq \
        application/tag/rpc \
        application/user/rpc \
    ; do \
        name=$(echo "$pkg" | sed 's|application/||; s|/|-|g'); \
        [ "$name" = "applet" ] && name="applet-api"; \
        echo ">>> building $pkg -> /app/$name"; \
        CGO_ENABLED=0 GOOS=linux GOMAXPROCS=4 go build -p 4 -o /app/$name ./$pkg; \
    done

FROM alpine:3.20

ENV TZ=Asia/Shanghai

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app /app
