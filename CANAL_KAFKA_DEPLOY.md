# Canal + Kafka 部署指南

## 架构概览

```
MySQL (binlog ROW) → Canal Server → Kafka (动态 topic 路由)
```

服务器：Ubuntu，IP: 101.42.34.232

### 组件版本

| 组件 | 镜像 | 端口 |
|------|------|------|
| MySQL 8.0 | mysql:8.0 | 3306 |
| Zookeeper | zookeeper:3.8 | 2181 |
| Kafka | wurstmeister/kafka:2.13-2.8.1 | 9092 |
| Canal Server | canal/canal-server:latest | 11111, 11112 |

---

## 1. MySQL binlog 配置

MySQL 容器内创建 `/etc/mysql/conf.d/binlog.cnf`：

```ini
[mysqld]
server-id=1
log-bin=mysql-bin
binlog-format=ROW
```

创建 Canal 用户（必须使用 `mysql_native_password`）：

```sql
CREATE USER 'canal'@'%' IDENTIFIED BY 'dsw123456';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
ALTER USER 'canal'@'%' IDENTIFIED WITH mysql_native_password BY 'dsw123456';
FLUSH PRIVILEGES;
```

验证 binlog：

```sql
SHOW VARIABLES LIKE 'log_bin';          -- ON
SHOW VARIABLES LIKE 'binlog_format';    -- ROW
SHOW MASTER STATUS;
```

---

## 2. Zookeeper 部署

`/opt/zookeeper/docker-compose.yml`：

```yaml
version: "3.8"
services:
  zookeeper:
    image: zookeeper:3.8
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      - ZOO_MY_ID=1
    restart: always
```

---

## 3. Kafka 部署

`/opt/kafka/docker-compose.yml`：

```yaml
version: "3.8"
services:
  kafka:
    image: wurstmeister/kafka:2.13-2.8.1
    container_name: kafka
    network_mode: host
    environment:
      - KAFKA_BROKER_ID=1
      - KAFKA_ZOOKEEPER_CONNECT=127.0.0.1:2181
      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://101.42.34.232:9092
      - KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092
      - KAFKA_LOG_RETENTION_HOURS=12
      - KAFKA_MESSAGE_MAX_BYTES=1048576
      - KAFKA_REPLICA_FETCH_MAX_BYTES=1048576
      - KAFKA_AUTO_CREATE_TOPICS_ENABLE=true
      - KAFKA_HEAP_OPTS=-Xmx256m -Xms128m
    restart: always
```

> **注意**：Kafka 使用 `network_mode: host`，Zookeeper 也必须在宿主机 2181 端口可达。

---

## 4. Canal 部署

### 4.1 docker-compose.yml

`/opt/canal/docker-compose.yml`：

```yaml
version: "3.8"
services:
  canal-server:
    image: canal/canal-server:latest
    container_name: canal-server
    ports:
      - "11111:11111"
      - "11112:11112"
    extra_hosts:
      - "host.docker.internal:host-gateway"
    restart: unless-stopped
```

### 4.2 修改 canal.properties（主配置）

```bash
# 设置 Kafka 模式
docker exec canal-server sed -i '0,/^canal.serverMode = tcp/s/canal.serverMode = tcp/canal.serverMode = kafka/' /home/admin/canal-server/conf/canal.properties

# 设置 Kafka 地址（容器内不能用 127.0.0.1，必须用 host.docker.internal）
docker exec canal-server sed -i 's/^kafka.bootstrap.servers.*/kafka.bootstrap.servers = host.docker.internal:9092/' /home/admin/canal-server/conf/canal.properties
```

验证：

```bash
docker exec canal-server grep -E "canal.serverMode|kafka.bootstrap.servers" /home/admin/canal-server/conf/canal.properties
```

应该只有：

```
canal.serverMode = kafka
kafka.bootstrap.servers = host.docker.internal:9092
```

> **踩坑**：`canal.serverMode` 只能有一行，多行时 Canal 读第一行。用 `sed -i '0,/.../s/.../.../'` 替换第一行。

### 4.3 编写 instance.properties（实例配置）

由于 `\\` 转义问题，**不要用 sed 修改**，用 heredoc + docker cp 写入：

```bash
cat > /tmp/instance.properties << 'ENDOFFILE'
canal.instance.gtidon=false
canal.instance.master.address=host.docker.internal:3306
canal.instance.master.journal.name=
canal.instance.master.position=
canal.instance.master.timestamp=
canal.instance.master.gtid=
canal.instance.tsdb.enable=true
canal.instance.dbUsername=canal
canal.instance.dbPassword=dsw123456
canal.instance.connectionCharset = UTF-8
canal.instance.enableDruid=false
canal.instance.filter.regex=thinktalk_like\\..*,thinktalk_article\\..*
canal.instance.filter.black.regex=mysql\\.slave_.*
canal.mq.topic=thinktalk_placeholder
canal.mq.dynamicTopic=thinktalk_like:thinktalk_like\\..*,thinktalk_article:thinktalk_article\\..*
canal.mq.partition=0
ENDOFFILE

docker cp /tmp/instance.properties canal-server:/home/admin/canal-server/conf/example/instance.properties
```

### 4.4 重启 Canal

```bash
docker restart canal-server
```

---

## 5. 关键配置说明

### 5.1 网络地址

| 场景 | 地址 |
|------|------|
| 容器内访问宿主机 | `host.docker.internal`（需 `extra_hosts`） |
| Kafka advertised listener | 宿主机公网 IP `101.42.34.232` |

### 5.2 filter.regex 语法

```
库名\\.表名
```

| 示例 | 含义 |
|------|------|
| `.*\\..*` | 所有库所有表 |
| `thinktalk_like\\..*` | thinktalk_like 库所有表 |
| `thinktalk_like\\.like_count` | 只监听一张表 |
| `thinktalk_like\\..*,thinktalk_article\\..*` | 多库用逗号分隔 |

### 5.3 动态 topic 路由

```properties
# 语法：topic名:库表正则
canal.mq.dynamicTopic=thinktalk_like:thinktalk_like\\..*,thinktalk_article:thinktalk_article\\..*
```

- `thinktalk_like` 库的变更 → Kafka topic `thinktalk_like`
- `thinktalk_article` 库的变更 → Kafka topic `thinktalk_article`

> **踩坑**：正则中必须用 `\\`（双反斜杠），sed 会吃掉一个，必须用 heredoc + docker cp 写入。

---

## 6. 常见问题排查

### 6.1 Canal 读取 binlog 但不发 Kafka

**现象**：meta.dat 中 position 在前进，但 Kafka topic 无消息。

**原因**：`canal.serverMode` 未设置为 `kafka`，或存在多行配置 Canal 读了第一行 `tcp`。

**解决**：

```bash
docker exec canal-server grep "canal.serverMode" /home/admin/canal-server/conf/canal.properties
# 确保只有一行：canal.serverMode = kafka
```

### 6.2 Access denied for user 'canal'

**现象**：Canal 日志报 `Access denied for user 'canal'@'172.17.0.1'`

**原因**：MySQL 8.0 默认用 `caching_sha2_password`，Canal 的 Java 客户端不支持。

**解决**：

```sql
ALTER USER 'canal'@'%' IDENTIFIED WITH mysql_native_password BY 'dsw123456';
```

### 6.3 Could not find first log file in binary log index

**现象**：`errno = 1236, Could not find first log file name in binary log index file`

**原因**：执行了 `RESET MASTER` 删除了旧 binlog，但 Canal 还保存着旧位置。

**解决**：

```bash
docker exec canal-server rm -f /home/admin/canal-server/conf/example/meta.dat
docker restart canal-server
```

### 6.4 Canal binlog 位置不前进

**现象**：Canal 连接 MySQL 正常（SHOW PROCESSLIST 有 Binlog Dump），但 position 不变。

**原因**：binlog dump 连接卡住，或 Kafka producer 连接失败但无日志。

**解决**：

```bash
# 删除 meta.dat 重置位置
docker exec canal-server rm -f /home/admin/canal-server/conf/example/meta.dat
docker restart canal-server

# 等待 15 秒后检查
sleep 15 && docker exec canal-server cat /home/admin/canal-server/logs/example/example.log | tail -10
```

### 6.5 Kafka 消费者 TimeoutException

**现象**：`kafka-console-consumer.sh` 报 TimeoutException，0 messages。

**原因**：
1. Canal 没有发送消息到该 topic（检查 topic offset）
2. 消费者没有加 `--from-beginning`，从最新 offset 开始读

**排查**：

```bash
# 检查 topic 是否有消息
docker exec kafka kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list 127.0.0.1:9092 --topic <topic名> --time -1

# 从头消费
docker exec kafka kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic <topic名> --from-beginning --timeout-ms 10000
```

### 6.6 动态 topic 不生效（消息去了默认 topic）

**现象**：消息出现在 `canal.mq.topic` 指定的默认 topic，而不是动态 topic。

**原因**：`canal.mq.dynamicTopic` 中正则的 `\\` 被 sed 吃掉，变成了单 `\`。

**解决**：用 heredoc + docker cp 重新写入配置文件，确保 `\\` 正确。

---

## 7. 验证命令速查

```bash
# 检查 Canal 是否连接 MySQL
docker exec -it mysql mysql -uroot -pdsw123456 -e "SHOW PROCESSLIST;" | grep -i canal

# 查看 Canal 实例日志
docker exec canal-server cat /home/admin/canal-server/logs/example/example.log | tail -20

# 查看 Canal 主日志
docker exec canal-server cat /home/admin/canal-server/logs/canal/canal.log | tail -20

# 查看 Canal binlog 位置
docker exec canal-server cat /home/admin/canal-server/conf/example/meta.dat

# 列出所有 Kafka topic
docker exec kafka kafka-topics.sh --bootstrap-server 127.0.0.1:9092 --list

# 查看 topic 消息数
docker exec kafka kafka-run-class.sh kafka.tools.GetOffsetShell --broker-list 127.0.0.1:9092 --topic <topic名> --time -1

# 消费 Kafka 消息
docker exec kafka kafka-console-consumer.sh --bootstrap-server 127.0.0.1:9092 --topic <topic名> --from-beginning --timeout-ms 10000

# 检查 MySQL binlog 状态
docker exec -it mysql mysql -uroot -pdsw123456 -e "SHOW MASTER STATUS;"
```

---

## 8. 修改监听范围

如需修改 Canal 监听的数据库：

1. 重新编写 `/tmp/instance.properties`（修改 `filter.regex` 和 `dynamicTopic`）
2. `docker cp` 到容器
3. 删除 `meta.dat` 并重启

```bash
# 示例：改为监听 thinktalk_user 和 thinktalk_follow
cat > /tmp/instance.properties << 'ENDOFFILE'
canal.instance.gtidon=false
canal.instance.master.address=host.docker.internal:3306
canal.instance.master.journal.name=
canal.instance.master.position=
canal.instance.master.timestamp=
canal.instance.master.gtid=
canal.instance.tsdb.enable=true
canal.instance.dbUsername=canal
canal.instance.dbPassword=dsw123456
canal.instance.connectionCharset = UTF-8
canal.instance.enableDruid=false
canal.instance.filter.regex=thinktalk_user\\..*,thinktalk_follow\\..*
canal.instance.filter.black.regex=mysql\\.slave_.*
canal.mq.topic=thinktalk_placeholder
canal.mq.dynamicTopic=thinktalk_user:thinktalk_user\\..*,thinktalk_follow:thinktalk_follow\\..*
canal.mq.partition=0
ENDOFFILE

docker cp /tmp/instance.properties canal-server:/home/admin/canal-server/conf/example/instance.properties
docker exec canal-server rm -f /home/admin/canal-server/conf/example/meta.dat
docker restart canal-server
```
