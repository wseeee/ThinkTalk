# 隐私配置下沉到 .env 设计方案

本项目是一个基于 go-zero 框架的多微服务社交平台。为了提高项目的安全性和灵活性，需要将所有微服务 YAML 配置文件中的敏感信息（如数据库 DSN 密码、Redis 密码、JWT 密钥、Etcd 地址等）移出代码库，统一配置在根目录的 `.env` 文件中。

## 需求与目的

1. **安全性**：防止敏感的数据库密码、外部 IP 和密钥被硬编码并提交至 Git 仓库。
2. **灵活性**：同一套微服务代码在不同的环境（本地开发、测试、生产）下，只需修改外部环境配置即可，无需修改代码。
3. **开发体验**：Go 进程支持本地自动检测并载入项目根目录下的 `.env`，而在 Docker 环境中能够完美兼容宿主机注入的系统环境变量。

## 方案设计

### 1. 配置项清单

创建两个配置文件：
* `d:\Go_Zero\ThinkTalk\.env`：包含真实的敏感参数（需加入 `.gitignore`）。
* `d:\Go_Zero\ThinkTalk\.env.example`：示例配置文件模板。

提取配置项如下：

| 配置键名 | 说明 | 示例值 |
|---|---|---|
| `DB_USER` | MySQL 用户名 | `root` |
| `DB_PASS` | MySQL 密码 | `dsw123456` |
| `DB_HOST` | MySQL 主机 | `101.42.34.232` |
| `DB_PORT` | MySQL 端口 | `3306` |
| `REDIS_HOST` | Redis 主机 | `101.42.34.232` |
| `REDIS_PORT` | Redis 端口 | `6379` |
| `REDIS_PASS` | Redis 密码 | `dsw123456` |
| `ETCD_HOST` | Etcd 主机 | `101.42.34.232` |
| `ETCD_PORT` | Etcd 端口 | `2379` |
| `KAFKA_HOST` | Kafka 主机 | `101.42.34.232` |
| `KAFKA_PORT` | Kafka 端口 | `9092` |
| `JWT_ACCESS_SECRET` | JWT Access Token 签名密钥 | `ThinkTalk` |
| `JWT_REFRESH_SECRET` | JWT Refresh Token 签名密钥 | `ThinkTalk` |

### 2. 追溯式环境载入包 (`pkg/env`)

新增包 `ThinkTalk/pkg/env`，内部实现从当前目录向上追溯查找 `.env` 文件并使用 `joho/godotenv` 加载。

* **接口定义**：
  ```go
  package env

  // LoadEnv 向上寻找 .env 并注入环境变量，找不到则忽略
  func LoadEnv()
  ```

### 3. 微服务入口集成

项目中的微服务入口分布在 `application/` 目录下的各个 Go 程序中（包括 `api`、`rpc`、`mq` 等，共有 23 个）。
在每个服务的 `main()` 函数启动时，在 `conf.MustLoad` 执行前调用 `env.LoadEnv()`：

```diff
 func main() {
 	flag.Parse()
 
+	env.LoadEnv()
+
 	var c config.Config
 	conf.MustLoad(*configFile, &c)
    ...
```

### 4. 配置文件 YAML 改写

对以下服务的配置文件进行改写，将其中的敏感硬编码字符串替换为 `${VAR}` 占位符：

1. **applet**: `application/applet/etc/applet-api.yaml`
2. **article**:
   * `application/article/api/etc/article-api.yaml`
   * `application/article/mq/etc/article.yaml`
   * `application/article/rpc/etc/article.yaml`
3. **chat**:
   * `application/chat/api/etc/chat-api.yaml`
   * `application/chat/mq/etc/chat.yaml`
   * `application/chat/rpc/etc/chat.yaml`
4. **concerned**:
   * `application/concerned/mq/etc/concerned.yaml`
   * `application/concerned/rpc/etc/concerned.yaml`
5. **follow**: `application/follow/rpc/etc/follow.yaml`
6. **like**:
   * `application/like/mq/etc/like.yaml`
   * `application/like/rpc/etc/like.yaml`
7. **member**:
   * `application/member/mq/etc/member.yaml`
   * `application/member/rpc/etc/member.yaml`
8. **message**:
   * `application/message/mq/etc/message.yaml`
   * `application/message/rpc/etc/message.yaml`
9. **qa**:
   * `application/qa/api/etc/qa-api.yaml`
   * `application/qa/mq/etc/qa.yaml`
   * `application/qa/rpc/etc/qa.yaml`
10. **reply**:
    * `application/reply/mq/etc/reply.yaml`
    * `application/reply/rpc/etc/reply.yaml`
11. **tag**: `application/tag/rpc/etc/tag.yaml`
12. **user**: `application/user/rpc/etc/user.yaml`

#### 替换规则示例：
* 数据库: `DataSource: ${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/thinktalk_xxx?parseTime=true`
* Redis: `Host: ${REDIS_HOST}:${REDIS_PORT}`，`Pass: ${REDIS_PASS}`
* Etcd: `- ${ETCD_HOST}:${ETCD_PORT}`
* Kafka: `- ${KAFKA_HOST}:${KAFKA_PORT}`
* JWT: `AccessSecret: ${JWT_ACCESS_SECRET}`，`RefreshSecret: ${JWT_REFRESH_SECRET}`

### 5. Docker Compose 集成

修改根目录的 `docker-compose.yml`，给每一个服务容器添加 `env_file` 配置：
```yaml
    env_file:
      - .env
```
以此确保容器能够直接拿到 `.env` 中定义的变量。

---

## 验证计划

1. **依赖添加验证**：运行 `go get github.com/joho/godotenv` 并执行 `go mod tidy`。
2. **本地启动验证**：
   * 创建 `.env` 并填入真实配置。
   * 启动 `user` rpc 服务，观察其是否成功解析环境变量并正常运行（监听 `8080` 端口）。
   * 启动 `applet` api 服务，调用 API 接口或验证 RPC 联通性。
3. **编译构建验证**：运行 `go build ./application/applet` 等命令检查无编译错误。
4. **代码规范与验证**：使用 `go vet ./...` 检查代码规范。
