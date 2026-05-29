#!/bin/sh
set -e

CONFIG_DIR="/app/etc"
CONFIG_FILE="${CONFIG_DIR}/${CONFIG_FILE}"

if [ -d "$CONFIG_DIR" ] && [ -f "$CONFIG_FILE" ]; then
    echo "[entrypoint] Patching config for Docker environment..."

    for f in $(find "$CONFIG_DIR" -name "*.yaml"); do
        sed \
            -e "s|101\.42\.34\.232:2379|thinktalk-etcd:2379|g" \
            -e "s|101\.42\.34\.232:3306|thinktalk-mysql:3306|g" \
            -e "s|101\.42\.34\.232:6379|thinktalk-redis:6379|g" \
            -e "s|101\.42\.34\.232:9092|thinktalk-kafka:9092|g" \
            -e "s|101\.42\.34\.232:9000|thinktalk-minio:9000|g" \
            -e "s|localhost:9200|thinktalk-es:9200|g" \
            -e "s|127\.0\.0\.1:9200|thinktalk-es:9200|g" \
            -e "s|http://localhost:9200/|http://thinktalk-es:9200/|g" \
            -e "s|127\.0\.0\.1:14268|thinktalk-jaeger:14268|g" \
            "$f" > "$f.tmp" && mv "$f.tmp" "$f"
    done

    echo "[entrypoint] Config patched successfully."
fi

exec /app/service -f "$CONFIG_FILE"
