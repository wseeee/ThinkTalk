# 隐私配置下沉到 .env 实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 将所有微服务 YAML 配置文件中的敏感隐私信息配置到项目根目录的 `.env` 中，以提高安全性与可移植性。

**架构：**
1. 引入 `joho/godotenv` 并实现一个 `pkg/env` 模块，它会在启动时自适应向上追溯查找 `.env` 文件进行加载。
2. 每一个微服务的主入口函数（共 23 个）都会在加载 YAML 前执行 `env.LoadEnv()`。
3. YAML 配置中的敏感数据全部替换为 go-zero 原生支持的 `${VAR}` 占位符。
4. 在 `docker-compose.yml` 中配置 `env_file: - .env`，并在项目根目录提供 `.env.example`。

**技术栈：** Go 1.26, go-zero 1.10.1, github.com/joho/godotenv

---

### 任务 1：安装依赖包

**文件：**
- 修改：`go.mod`
- 修改：`go.sum`

- [ ] **步骤 1：添加 go 依赖**
  运行：
  ```bash
  go get github.com/joho/godotenv
  ```
  预期：依赖成功下载，并且 `go.mod` 增加了 `github.com/joho/godotenv` 条目。

- [ ] **步骤 2：执行 go mod tidy**
  运行：
  ```bash
  go mod tidy
  ```
  预期：退出码为 0，依赖整洁无误。

- [ ] **步骤 3：验证已有测试依旧通过**
  运行：
  ```bash
  go test ./pkg/encrypt/...
  ```
  预期：PASS

- [ ] **步骤 4：Commit**
  ```bash
  git add go.mod go.sum
  git commit -m "chore: add joho/godotenv dependency"
  ```

---

### 任务 2：实现自适应环境载入包 (`pkg/env`)

**文件：**
- 创建：`pkg/env/env.go`
- 创建：`pkg/env/env_test.go`

- [ ] **步骤 1：编写失败的单元测试**
  在 `pkg/env/env_test.go` 中编写以下内容，测试向上追溯加载 `.env` 的逻辑：
  ```go
  package env

  import (
  	"os"
  	"path/filepath"
  	"testing"
  )

  func TestLoadEnv(t *testing.T) {
  	cwd, err := os.Getwd()
  	if err != nil {
  		t.Fatalf("failed to get wd: %v", err)
  	}

  	tempDir, err := os.MkdirTemp("", "test_env")
  	if err != nil {
  		t.Fatalf("failed to create temp dir: %v", err)
  	}
  	defer os.RemoveAll(tempDir)

  	envFile := filepath.Join(tempDir, ".env")
  	err = os.WriteFile(envFile, []byte("TEST_DOTENV_VAR=dotenv_value\n"), 0644)
  	if err != nil {
  		t.Fatalf("failed to write temp .env: %v", err)
  	}

  	subDir := filepath.Join(tempDir, "sub_service")
  	err = os.Mkdir(subDir, 0755)
  	if err != nil {
  		t.Fatalf("failed to create sub dir: %v", err)
  	}

  	err = os.Chdir(subDir)
  	if err != nil {
  		t.Fatalf("failed to change dir: %v", err)
  	}
  	defer os.Chdir(cwd)

  	os.Unsetenv("TEST_DOTENV_VAR")

  	LoadEnv()

  	val := os.Getenv("TEST_DOTENV_VAR")
  	if val != "dotenv_value" {
  		t.Errorf("expected TEST_DOTENV_VAR to be 'dotenv_value', got '%s'", val)
  	}
  }
  ```

- [ ] **步骤 2：运行测试验证失败**
  运行：
  ```bash
  go test ./pkg/env/...
  ```
  预期：编译失败，因为 `pkg/env` 目录或其中的 `LoadEnv` 函数不存在。

- [ ] **步骤 3：编写最小实现代码**
  在 `pkg/env/env.go` 中实现追溯查找逻辑：
  ```go
  package env

  import (
  	"os"
  	"path/filepath"
  	"github.com/joho/godotenv"
  )

  func LoadEnv() {
  	dir, err := os.Getwd()
  	if err != nil {
  		_ = godotenv.Load()
  		return
  	}

  	for {
  		envPath := filepath.Join(dir, ".env")
  		if _, err := os.Stat(envPath); err == nil {
  			_ = godotenv.Load(envPath)
  			return
  		}
  		parent := filepath.Dir(dir)
  		if parent == dir {
  			break
  		}
  		dir = parent
  	}

  	_ = godotenv.Load()
  }
  ```

- [ ] **步骤 4：运行测试验证通过**
  运行：
  ```bash
  go test ./pkg/env/...
  ```
  预期：PASS

- [ ] **步骤 5：Commit**
  ```bash
  git add pkg/env/env.go pkg/env/env_test.go
  git commit -m "feat: implement pkg/env with upward recursive dotenv lookup"
  ```

---

### 任务 3：创建 `.env.example` 和本地 `.env`

**文件：**
- 创建：`.env.example`
- 创建：`.env`

- [ ] **步骤 1：创建 `.env.example`**
  在项目根目录创建 `.env.example`，内容为公开的默认结构：
  ```env
  # 数据库配置
  DB_USER=root
  DB_PASS=dsw123456
  DB_HOST=101.42.34.232
  DB_PORT=3306

  # Redis 配置
  REDIS_HOST=101.42.34.232
  REDIS_PORT=6379
  REDIS_PASS=dsw123456

  # Etcd 配置
  ETCD_HOST=101.42.34.232
  ETCD_PORT=2379

  # Kafka 配置
  KAFKA_HOST=101.42.34.232
  KAFKA_PORT=9092

  # JWT 密钥配置
  JWT_ACCESS_SECRET=ThinkTalk
  JWT_REFRESH_SECRET=ThinkTalk
  ```

- [ ] **步骤 2：创建本地使用的 `.env`**
  拷贝 `.env.example` 到 `.env`：
  ```bash
  cp .env.example .env
  ```
  预期：存在 `.env` 且内容与模板一致。

- [ ] **步骤 3：在 `.gitignore` 中追加 `.env` 规则（如果还没有）**
  打开 `.gitignore`，在文件尾部追加 `.env`。
  预期：`.gitignore` 包含 `.env` 屏蔽。

- [ ] **步骤 4：Commit**
  ```bash
  git add .env.example .gitignore
  git commit -m "chore: add .env.example template and update .gitignore"
  ```

---

### 任务 4：微服务入口集成

**文件：**
- 修改：`application/applet/applet.go`
- 修改：`application/user/rpc/user.go`
- 修改其它 21 个主入口文件。

- [ ] **步骤 1：率先集成并验证核心服务 (applet, user/rpc)**
  在 `application/applet/applet.go` 的 `main()` 函数中导入并加入 `env.LoadEnv()`。
  在 `application/user/rpc/user.go` 的 `main()` 函数中导入并加入 `env.LoadEnv()`。

  修改示例如下：
  ```go
  import (
  	"ThinkTalk/pkg/env"
      // ...其他
  )

  func main() {
  	flag.Parse()

  	env.LoadEnv()

  	var c config.Config
  	conf.MustLoad(*configFile, &c)
      // ...
  }
  ```

- [ ] **步骤 2：集成其余 21 个入口**
  修改列表：
  * `application/article/api/article.go`
  * `application/article/mq/main.go`
  * `application/article/rpc/article.go`
  * `application/chat/api/chat.go`
  * `application/chat/mq/main.go`
  * `application/chat/rpc/chat.go`
  * `application/concerned/mq/main.go`
  * `application/concerned/rpc/concerned.go`
  * `application/follow/rpc/follow.go`
  * `application/like/mq/main.go`
  * `application/like/rpc/like.go`
  * `application/member/mq/main.go`
  * `application/member/rpc/member.go`
  * `application/message/mq/main.go`
  * `application/message/rpc/message.go`
  * `application/qa/api/qa.go`
  * `application/qa/mq/main.go`
  * `application/qa/rpc/qa.go`
  * `application/reply/mq/main.go`
  * `application/reply/rpc/reply.go`
  * `application/tag/rpc/tag.go`

- [ ] **步骤 3：验证编译是否正常**
  运行：
  ```bash
  go build ./application/applet
  go build ./application/user/rpc
  ```
  预期：编译成功，没有编译错误。

- [ ] **步骤 4：Commit**
  ```bash
  git add application/
  git commit -m "feat: integrate env.LoadEnv in all microservice main entrypoints"
  ```

---

### 任务 5：重构所有服务 YAML 配置以使用环境变量

**文件：**
- 修改：所有 23 个服务对应的 YAML 配置文件。

- [ ] **步骤 1：重构 `applet` 与 `user` 的 YAML 并完成验证**
  - 修改 `application/applet/etc/applet-api.yaml`，将：
    * `AccessSecret: ThinkTalk` 替换为 `${JWT_ACCESS_SECRET}`
    * `RefreshSecret: ThinkTalk` 替换为 `${JWT_REFRESH_SECRET}`
    * `Hosts: - 101.42.34.232:2379` 等替换为 `- ${ETCD_HOST}:${ETCD_PORT}`
    * `Host: 101.42.34.232:6379` 替换为 `${REDIS_HOST}:${REDIS_PORT}`
    * `Pass: dsw123456` 替换为 `${REDIS_PASS}`
  - 修改 `application/user/rpc/etc/user.yaml`，将：
    * `DataSource` 替换为 `${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/thinktalk_user?parseTime=true`
    * `Host: 101.42.34.232:6379` 替换为 `${REDIS_HOST}:${REDIS_PORT}`
    * `Pass: dsw123456` 替换为 `${REDIS_PASS}`
    * `Hosts: - 101.42.34.232:2379` 替换为 `- ${ETCD_HOST}:${ETCD_PORT}`

  - **启动服务验证**：
    1. 运行 `go run application/user/rpc/user.go -f application/user/rpc/etc/user.yaml`
    2. 运行 `go run application/applet/applet.go -f application/applet/etc/applet-api.yaml`
    预期：两服务均能成功读取环境变量并正常启动。

- [ ] **步骤 2：重构其余微服务的 YAML 配置文件**
  逐个将所有敏感的 MySQL 连接串、Redis 主机与密码、Kafka 主机与端口、Etcd 主机与端口替换为对应的 `${VAR}`。

  - 微服务 YAML 列表：
    * `application/article/api/etc/article-api.yaml`
    * `application/article/mq/etc/article.yaml`
    * `application/article/rpc/etc/article.yaml`
    * `application/chat/api/etc/chat-api.yaml`
    * `application/chat/mq/etc/chat.yaml`
    * `application/chat/rpc/etc/chat.yaml`
    * `application/concerned/mq/etc/concerned.yaml`
    * `application/concerned/rpc/etc/concerned.yaml`
    * `application/follow/rpc/etc/follow.yaml`
    * `application/like/mq/etc/like.yaml`
    * `application/like/rpc/etc/like.yaml`
    * `application/member/mq/etc/member.yaml`
    * `application/member/rpc/etc/member.yaml`
    * `application/message/mq/etc/message.yaml`
    * `application/message/rpc/etc/message.yaml`
    * `application/qa/api/etc/qa-api.yaml`
    * `application/qa/mq/etc/qa.yaml`
    * `application/qa/rpc/etc/qa.yaml`
    * `application/reply/mq/etc/reply.yaml`
    * `application/reply/rpc/etc/reply.yaml`
    * `application/tag/rpc/etc/tag.yaml`

- [ ] **步骤 3：验证微服务是否可编译**
  运行：
  ```bash
  go vet ./...
  ```
  预期：没有静态检查错误。

- [ ] **步骤 4：Commit**
  ```bash
  git add application/
  git commit -m "refactor: configure sensitive credentials as env variables in all yaml files"
  ```

---

### 任务 6：在 `docker-compose.yml` 中注入 `.env`

**文件：**
- 修改：`docker-compose.yml`

- [ ] **步骤 1：为每一个服务容器加上 `env_file`**
  修改 `docker-compose.yml`，给里面每一个服务组件（例如 `applet-api`、`user-rpc`、`chat-rpc` 等共 23 个服务）添加配置：
  ```yaml
      env_file:
        - .env
  ```
  确保其在容器化启动时能自动注入根目录下的 `.env` 变量。

- [ ] **步骤 2：Commit**
  ```bash
  git add docker-compose.yml
  git commit -m "deploy: integrate env_file in docker-compose.yml"
  ```

---

### 任务 7：全局最终验证

- [ ] **步骤 1：运行所有模块的 Go 单元测试**
  运行：
  ```bash
  go test ./...
  ```
  预期：所有测试全部通过。

- [ ] **步骤 2：全量服务编译验证**
  运行：
  ```bash
  go build ./application/applet
  go build ./application/user/rpc
  go build ./application/article/api
  ```
  预期：顺利通过，无编译错误。
