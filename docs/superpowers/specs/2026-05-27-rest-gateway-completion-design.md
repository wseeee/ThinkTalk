# ThinkTalk REST 网关补齐设计文档

> 日期: 2026-05-27
> 状态: 待审查

---

## 1. 背景与目标

ThinkTalk 项目现有 11 个 gRPC 微服务，但只有 3 个 REST API 网关 (applet-api、article-api、chat-api)，共 31 个对外接口。3 个 gRPC 服务 (tag/rpc、reply/rpc、qa/rpc) 完全没有 REST 网关，另有 4 个 gRPC 方法虽有对应网关但遗漏了 REST 端点。

**目标**: 补齐缺失的 REST 网关层，使所有 gRPC 服务均有对应的 REST 对外接口。

**补齐后**: 从 31 个 REST 接口增至 59 个。

---

## 2. 架构决策

### 2.1 部署模式: 混合模式

| 模块 | 部署方式 | 理由 |
|------|----------|------|
| Tag (10 接口) | 集成到 applet-api | 标签是通用基础设施，与 like/follow 同级 |
| Reply (5 接口) | 集成到 applet-api | 评论是通用基础设施 |
| QA (9 接口) | 独立网关 qa-api | 问答是独立业务模块，参照 article-api 模式 |
| ArticleDelete (1) | 集成到 article-api | 与现有 article 路由同属一个服务 |
| ConcernedCount (1) | 集成到 applet-api | 与现有 concerned 路由同属一个服务 |
| UpgradeMember (1) | 集成到 applet-api | 与现有 member 路由同属一个服务 |
| MemberOrderList (1) | 集成到 applet-api | 与现有 member 路由同属一个服务 |

### 2.2 认证策略: 按操作类型区分

- 写操作 (POST 创建/修改/删除): JWT + Signature
- 读操作 (GET 查询): JWT
- userId 始终从 JWT 中提取，不暴露在请求体中

### 2.3 URL 前缀

| 模块 | 前缀 |
|------|------|
| Tag | `/v1/tag` |
| Reply | `/v1/reply` |
| QA | `/v1/qa` |
| ArticleDelete | `/v1/article` (已有) |
| ConcernedCount | `/v1/concerned` (已有) |
| UpgradeMember/MemberOrderList | `/v1/member` (已有) |

### 2.4 实现方式

全部手写 Go 代码，不使用 goctl 生成，与 applet-api 中 like/follow/message 等模块保持一致的编码模式。

---

## 3. 接口定义

### 3.1 Tag 模块 — 集成到 applet-api

前缀: `/v1/tag`

**写操作 (JWT + Signature)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| POST | `/v1/tag/create` | 创建标签 | CreateTag |
| POST | `/v1/tag/update` | 更新标签 | UpdateTag |
| POST | `/v1/tag/delete` | 删除标签 | DeleteTag |
| POST | `/v1/tag/resource/add` | 打标签 | TagResource |
| POST | `/v1/tag/resource/remove` | 去标签 | UntagResource |

**读操作 (JWT)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| GET | `/v1/tag/detail` | 标签详情 | TagDetail |
| GET | `/v1/tag/list` | 标签列表 | TagList |
| GET | `/v1/tag/hot` | 热门标签 | HotTags |
| GET | `/v1/tag/resource/list` | 标签下的资源 | ResourcesByTag |
| GET | `/v1/tag/resource/tags` | 资源的标签 | TagsByResource |

**请求/响应类型**:

```go
// TagCreateRequest — 创建标签
type TagCreateRequest struct {
    TagName string `json:"tag_name"`
    TagDesc string `json:"tag_desc"`
}
type TagCreateResponse struct {
    TagId int64 `json:"tag_id"`
}

// TagUpdateRequest — 更新标签
type TagUpdateRequest struct {
    TagId   int64  `json:"tag_id"`
    TagName string `json:"tag_name"`
    TagDesc string `json:"tag_desc"`
}

// TagDeleteRequest — 删除标签
type TagDeleteRequest struct {
    TagId int64 `json:"tag_id"`
}

// TagDetailRequest — 标签详情
type TagDetailRequest struct {
    TagId int64 `form:"tag_id"`
}

// TagItem — 标签信息(公共)
type TagItem struct {
    TagId         int64  `json:"tag_id"`
    TagName       string `json:"tag_name"`
    TagDesc       string `json:"tag_desc"`
    ResourceCount int64  `json:"resource_count"`
    CreateTime    int64  `json:"create_time"`
}
type TagDetailResponse = TagItem

// TagListRequest/Response — 标签列表
type TagListRequest struct {
    Cursor   int64 `form:"cursor"`
    PageSize int64 `form:"page_size"`
}
type TagListResponse struct {
    Items  []*TagItem `json:"items"`
    Cursor int64      `json:"cursor"`
    IsEnd  bool       `json:"is_end"`
}

// HotTagsRequest/Response — 热门标签
type HotTagsRequest struct {
    Limit int32 `form:"limit"`
}
type HotTagsResponse struct {
    Items []*TagItem `json:"items"`
}

// TagResourceRequest — 打标签
type TagResourceRequest struct {
    BizId    string `json:"biz_id"`
    TargetId int64  `json:"target_id"`
    TagId    int64  `json:"tag_id"`
}

// UntagResourceRequest — 去标签
type UntagResourceRequest struct {
    BizId    string `json:"biz_id"`
    TargetId int64  `json:"target_id"`
    TagId    int64  `json:"tag_id"`
}

// ResourcesByTagRequest/Response — 标签下的资源
type ResourcesByTagRequest struct {
    TagId    int64  `form:"tag_id"`
    BizId    string `form:"biz_id,optional"`
    Cursor   int64  `form:"cursor"`
    PageSize int64  `form:"page_size"`
}
type ResourceItem struct {
    TargetId   int64  `json:"target_id"`
    BizId      string `json:"biz_id"`
    CreateTime int64  `json:"create_time"`
}
type ResourcesByTagResponse struct {
    Items  []*ResourceItem `json:"items"`
    Cursor int64           `json:"cursor"`
    IsEnd  bool            `json:"is_end"`
}

// TagsByResourceRequest/Response — 资源的标签
type TagsByResourceRequest struct {
    BizId    string `form:"biz_id"`
    TargetId int64  `form:"target_id"`
}
type TagsByResourceResponse struct {
    Items []*TagItem `json:"items"`
}
```

### 3.2 Reply 模块 — 集成到 applet-api

前缀: `/v1/reply`

**写操作 (JWT + Signature)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| POST | `/v1/reply/create` | 发表评论 | CreateReply |
| POST | `/v1/reply/delete` | 删除评论 | DeleteReply |

**读操作 (JWT)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| GET | `/v1/reply/detail` | 评论详情 | ReplyDetail |
| GET | `/v1/reply/list` | 评论列表 | ReplyList |
| GET | `/v1/reply/count` | 评论计数 | ReplyCount |

**请求/响应类型**:

```go
// ReplyCreateRequest — 发表评论
type ReplyCreateRequest struct {
    BizId         string `json:"biz_id"`
    TargetId      int64  `json:"target_id"`
    BeReplyUserId int64  `json:"be_reply_user_id,optional"`
    ParentId      int64  `json:"parent_id,optional"`
    Content       string `json:"content"`
}
type ReplyCreateResponse struct {
    ReplyId int64 `json:"reply_id"`
}

// ReplyDeleteRequest — 删除评论
type ReplyDeleteRequest struct {
    ReplyId int64 `json:"reply_id"`
}

// ReplyDetailRequest — 评论详情
type ReplyDetailRequest struct {
    ReplyId int64 `form:"reply_id"`
}

// ReplyItem — 评论信息(公共)
type ReplyItem struct {
    ReplyId       int64        `json:"reply_id"`
    BizId         string       `json:"biz_id"`
    TargetId      int64        `json:"target_id"`
    ReplyUserId   int64        `json:"reply_user_id"`
    BeReplyUserId int64        `json:"be_reply_user_id"`
    ParentId      int64        `json:"parent_id"`
    Content       string       `json:"content"`
    LikeNum       int64        `json:"like_num"`
    CreateTime    int64        `json:"create_time"`
    SubReplies    []*ReplyItem `json:"sub_replies,omitempty"`
}
type ReplyDetailResponse struct {
    Reply *ReplyItem `json:"reply"`
}

// ReplyListRequest/Response — 评论列表
type ReplyListRequest struct {
    BizId    string `form:"biz_id"`
    TargetId int64  `form:"target_id"`
    Cursor   int64  `form:"cursor"`
    PageSize int64  `form:"page_size"`
    SortType int32  `form:"sort_type,optional"`
}
type ReplyListResponse struct {
    Items  []*ReplyItem `json:"items"`
    Cursor int64        `json:"cursor"`
    IsEnd  bool         `json:"is_end"`
}

// ReplyCountRequest/Response — 评论计数
type ReplyCountRequest struct {
    BizId    string `form:"biz_id"`
    TargetId int64  `form:"target_id"`
}
type ReplyCountResponse struct {
    ReplyNum     int64 `json:"reply_num"`
    ReplyRootNum int64 `json:"reply_root_num"`
}
```

### 3.3 ArticleDelete — 集成到 article-api

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| POST | `/v1/article/delete` | 删除文章 | ArticleDelete |

**请求类型**:

```go
type ArticleDeleteRequest struct {
    ArticleId int64 `json:"article_id"`
}
```

路由注册到 article-api 的 JWT 路由组中 (`/v1/article` 前缀)。

### 3.4 ConcernedCount — 集成到 applet-api

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| GET | `/v1/concerned/count` | 收藏计数 | ConcernedCount |

**请求/响应类型**:

```go
type ConcernedCountRequest struct {
    BizId string `form:"biz_id"`
    ObjId int64  `form:"obj_id"`
}
type ConcernedCountResponse struct {
    ConcernedNum int64 `json:"concerned_num"`
}
```

路由注册到 applet-api 的 JWT 路由组中 (`/v1/concerned` 前缀)。

### 3.5 UpgradeMember + MemberOrderList — 集成到 applet-api

**写操作 (JWT + Signature)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| POST | `/v1/member/upgrade` | 升级会员 | UpgradeMember |

**读操作 (JWT)**

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| GET | `/v1/member/orders` | 会员订单列表 | MemberOrderList |

**请求/响应类型**:

```go
// UpgradeMemberRequest — 升级会员
type UpgradeMemberRequest struct {
    Level         int32  `json:"level"`
    DurationDays  int64  `json:"duration_days"`
    TransactionId string `json:"transaction_id"`
    Amount        int64  `json:"amount"`
    PayChannel    string `json:"pay_channel"`
}

// MemberOrderListRequest/Response — 订单列表
type MemberOrderListRequest struct {
    Cursor   int64 `form:"cursor"`
    PageSize int64 `form:"page_size"`
}
type MemberOrderItem struct {
    Id           int64  `json:"id"`
    UserId       int64  `json:"user_id"`
    Level        int32  `json:"level"`
    DurationDays int64  `json:"duration_days"`
    Amount       int64  `json:"amount"`
    PayChannel   string `json:"pay_channel"`
    Status       int32  `json:"status"`
    CreateTime   int64  `json:"create_time"`
}
type MemberOrderListResponse struct {
    Items  []*MemberOrderItem `json:"items"`
    Cursor int64              `json:"cursor"`
    IsEnd  bool               `json:"is_end"`
}
```

### 3.6 QA 模块 — 独立网关 qa-api

前缀: `/v1/qa`, 认证: JWT, 端口: 8890

| 方法 | 路径 | 说明 | gRPC |
|------|------|------|------|
| POST | `/v1/qa/publish` | 发布问题 | PublishQuestion |
| POST | `/v1/qa/answer` | 回答问题 | AnswerQuestion |
| POST | `/v1/qa/accept` | 采纳回答 | AcceptAnswer |
| GET | `/v1/qa/list` | 问题列表 | Questions |
| GET | `/v1/qa/detail` | 问题详情 | QuestionDetail |
| POST | `/v1/qa/delete` | 删除问题 | QuestionDelete |
| GET | `/v1/qa/answers` | 回答列表 | AnswerList |
| POST | `/v1/qa/answer/delete` | 删除回答 | AnswerDelete |
| GET | `/v1/qa/search` | 搜索问题 | SearchQuestions |

**请求/响应类型**:

```go
// PublishQuestion — 发布问题
type PublishQuestionRequest struct {
    Title   string `json:"title"`
    Content string `json:"content"`
    TagIds  string `json:"tag_ids"`
}
type PublishQuestionResponse struct {
    QuestionId int64 `json:"question_id"`
}

// AnswerQuestion — 回答问题
type AnswerQuestionRequest struct {
    QuestionId int64  `json:"question_id"`
    Content    string `json:"content"`
}
type AnswerQuestionResponse struct {
    AnswerId int64 `json:"answer_id"`
}

// AcceptAnswer — 采纳回答
type AcceptAnswerRequest struct {
    QuestionId int64 `json:"question_id"`
    AnswerId   int64 `json:"answer_id"`
}

// QuestionList — 问题列表
type QuestionListRequest struct {
    Cursor   int64 `form:"cursor"`
    PageSize int64 `form:"page_size"`
    SortType int32 `form:"sort_type,optional"`
}

type QuestionItem struct {
    Id         int64  `json:"id"`
    Title      string `json:"title"`
    Content    string `json:"content"`
    AuthorId   int64  `json:"author_id"`
    AnswerNum  int64  `json:"answer_num"`
    ViewNum    int64  `json:"view_num"`
    TagIds     string `json:"tag_ids"`
    CreateTime int64  `json:"create_time"`
}
type QuestionListResponse struct {
    Items  []*QuestionItem `json:"items"`
    Cursor int64           `json:"cursor"`
    IsEnd  bool            `json:"is_end"`
}

// QuestionDetail — 问题详情
type QuestionDetailRequest struct {
    QuestionId int64 `form:"question_id"`
}
type QuestionDetailResponse struct {
    Question *QuestionItem `json:"question"`
}

// QuestionDelete — 删除问题
type QuestionDeleteRequest struct {
    QuestionId int64 `json:"question_id"`
}

// AnswerList — 回答列表
type AnswerListRequest struct {
    QuestionId int64 `form:"question_id"`
    Cursor     int64 `form:"cursor"`
    PageSize   int64 `form:"page_size"`
}

type AnswerItem struct {
    Id         int64  `json:"id"`
    QuestionId int64  `json:"question_id"`
    AuthorId   int64  `json:"author_id"`
    Content    string `json:"content"`
    IsAccepted bool   `json:"is_accepted"`
    LikeNum    int64  `json:"like_num"`
    ReplyNum   int64  `json:"reply_num"`
    CreateTime int64  `json:"create_time"`
}
type AnswerListResponse struct {
    Items  []*AnswerItem `json:"items"`
    Cursor int64         `json:"cursor"`
    IsEnd  bool          `json:"is_end"`
}

// AnswerDelete — 删除回答
type AnswerDeleteRequest struct {
    AnswerId int64 `json:"answer_id"`
}

// SearchQuestions — 搜索问题
type SearchQuestionsRequest struct {
    Keyword  string `form:"keyword"`
    Cursor   int64  `form:"cursor"`
    PageSize int64  `form:"page_size"`
}

type SearchQuestionItem struct {
    Id         int64  `json:"id"`
    Title      string `json:"title"`
    Content    string `json:"content"`
    AuthorId   int64  `json:"author_id"`
    AnswerNum  int64  `json:"answer_num"`
    TagIds     string `json:"tag_ids"`
    CreateTime int64  `json:"create_time"`
}
type SearchQuestionsResponse struct {
    Items  []*SearchQuestionItem `json:"items"`
    Cursor int64                 `json:"cursor"`
    IsEnd  bool                  `json:"is_end"`
}
```

---

## 4. 路由注册

### 4.1 applet-api routes.go 追加

```go
// Tag 读操作 — JWT
server.AddRoutes([]rest.Route{
    {Method: http.MethodGet, Path: "/detail",        Handler: TagDetailHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/list",          Handler: TagListHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/hot",           Handler: HotTagsHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/resource/list", Handler: ResourcesByTagHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/resource/tags", Handler: TagsByResourceHandler(serverCtx)},
}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
   rest.WithPrefix("/v1/tag"))

// Tag 写操作 — JWT + Signature
server.AddRoutes([]rest.Route{
    {Method: http.MethodPost, Path: "/create",        Handler: TagCreateHandler(serverCtx)},
    {Method: http.MethodPost, Path: "/update",        Handler: TagUpdateHandler(serverCtx)},
    {Method: http.MethodPost, Path: "/delete",        Handler: TagDeleteHandler(serverCtx)},
    {Method: http.MethodPost, Path: "/resource/add",  Handler: TagResourceHandler(serverCtx)},
    {Method: http.MethodPost, Path: "/resource/remove", Handler: UntagResourceHandler(serverCtx)},
}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
   rest.WithSignature(serverCtx.Config.Signature),
   rest.WithPrefix("/v1/tag"))

// Reply 读操作 — JWT
server.AddRoutes([]rest.Route{
    {Method: http.MethodGet, Path: "/detail", Handler: ReplyDetailHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/list",   Handler: ReplyListHandler(serverCtx)},
    {Method: http.MethodGet, Path: "/count",  Handler: ReplyCountHandler(serverCtx)},
}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
   rest.WithPrefix("/v1/reply"))

// Reply 写操作 — JWT + Signature
server.AddRoutes([]rest.Route{
    {Method: http.MethodPost, Path: "/create", Handler: ReplyCreateHandler(serverCtx)},
    {Method: http.MethodPost, Path: "/delete", Handler: ReplyDeleteHandler(serverCtx)},
}, rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
   rest.WithSignature(serverCtx.Config.Signature),
   rest.WithPrefix("/v1/reply"))

// ConcernedCount — 追加到已有 concerned JWT 路由组
// {Method: http.MethodGet, Path: "/count", Handler: ConcernedCountHandler(serverCtx)},

// UpgradeMember — 追加到已有 member JWT+Sig 路由组(需改为混合组或新增)
// MemberOrderList — 追加到已有 member JWT 路由组
```

### 4.2 article-api routes.go 追加

```go
// 追加到已有 JWT 路由组
// {Method: http.MethodPost, Path: "/delete", Handler: ArticleDeleteHandler(serverCtx)},
```

### 4.3 qa-api routes.go

```go
func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
    server.AddRoutes([]rest.Route{
        {Method: http.MethodPost, Path: "/publish",       Handler: PublishQuestionHandler(serverCtx)},
        {Method: http.MethodPost, Path: "/answer",        Handler: AnswerQuestionHandler(serverCtx)},
        {Method: http.MethodPost, Path: "/accept",        Handler: AcceptAnswerHandler(serverCtx)},
        {Method: http.MethodGet,  Path: "/list",          Handler: QuestionsHandler(serverCtx)},
        {Method: http.MethodGet,  Path: "/detail",        Handler: QuestionDetailHandler(serverCtx)},
        {Method: http.MethodPost, Path: "/delete",        Handler: QuestionDeleteHandler(serverCtx)},
        {Method: http.MethodGet,  Path: "/answers",       Handler: AnswerListHandler(serverCtx)},
        {Method: http.MethodPost, Path: "/answer/delete", Handler: AnswerDeleteHandler(serverCtx)},
        {Method: http.MethodGet,  Path: "/search",        Handler: SearchQuestionsHandler(serverCtx)},
    }, rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
       rest.WithPrefix("/v1/qa"))
}
```

---

## 5. ServiceContext 与配置变更

### 5.1 applet-api

**config.go 追加**:

```go
type Config struct {
    // ... 已有字段 ...
    TagRpc   zrpc.RpcClientConf
    ReplyRpc zrpc.RpcClientConf
}
```

**servicecontext.go 追加**:

```go
type ServiceContext struct {
    // ... 已有字段 ...
    TagRPC   tag.Tag
    ReplyRPC reply.Reply
}

func NewServiceContext(c config.Config) *ServiceContext {
    // ... 已有初始化 ...
    tagRPC := zrpc.MustNewClient(c.TagRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
    replyRPC := zrpc.MustNewClient(c.ReplyRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
    return &ServiceContext{
        // ...
        TagRPC:   tag.NewTag(tagRPC),
        ReplyRPC: reply.NewReply(replyRPC),
    }
}
```

**applet-api.yaml 追加**:

```yaml
TagRpc:
    Etcd:
        Hosts:
            - 101.42.34.232:2379
        Key: tag.rpc
    NonBlock: true
ReplyRpc:
    Etcd:
        Hosts:
            - 101.42.34.232:2379
        Key: reply.rpc
    NonBlock: true
```

### 5.2 qa-api 新建

**config.go**:

```go
package config

import (
    "github.com/zeromicro/go-zero/rest"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    QaRpc zrpc.RpcClientConf
}
```

**servicecontext.go**:

```go
package svc

import (
    "ThinkTalk/application/qa/api/internal/config"
    "ThinkTalk/application/qa/rpc/qa"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
    Config config.Config
    QaRPC  qa.QA
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        QaRPC:  qa.NewQA(zrpc.MustNewClient(c.QaRpc)),
    }
}
```

**qa-api.yaml**:

```yaml
Name: qa-api
Host: 0.0.0.0
Port: 8890
Auth:
    AccessSecret: ThinkTalk
    AccessExpire: 604800
QaRpc:
    Etcd:
        Hosts:
            - 101.42.34.232:2379
        Key: qa.rpc
    NonBlock: true
```

---

## 6. 文件改动清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `application/applet/internal/types/types.go` | 编辑 | 追加 Tag/Reply/ConcernedCount/UpgradeMember/MemberOrderList 类型 |
| `application/applet/internal/handler/routes.go` | 编辑 | 追加 Tag/Reply/ConcernedCount/UpgradeMember/MemberOrderList 路由 |
| `application/applet/internal/handler/taghandler.go` | 新建 | Tag 10 个 handler |
| `application/applet/internal/handler/replyhandler.go` | 新建 | Reply 5 个 handler |
| `application/applet/internal/handler/concernedhandler.go` | 编辑 | 追加 ConcernedCountHandler |
| `application/applet/internal/handler/memberhandler.go` | 编辑 | 追加 UpgradeMemberHandler, MemberOrderListHandler |
| `application/applet/internal/logic/taglogic.go` | 新建 | Tag 10 个 logic |
| `application/applet/internal/logic/replylogic.go` | 新建 | Reply 5 个 logic |
| `application/applet/internal/logic/concernedlogic.go` | 编辑 | 追加 ConcernedCount |
| `application/applet/internal/logic/memberlogic.go` | 编辑 | 追加 UpgradeMember, MemberOrderList |
| `application/applet/internal/svc/servicecontext.go` | 编辑 | 注入 TagRPC, ReplyRPC |
| `application/applet/internal/config/config.go` | 编辑 | 追加 TagRpc, ReplyRpc |
| `application/applet/etc/applet-api.yaml` | 编辑 | 追加 TagRpc, ReplyRpc 配置 |
| `application/article/api/internal/types/types.go` | 编辑 | 追加 ArticleDeleteRequest, 修复 AuthorId 类型 |
| `application/article/api/internal/handler/routes.go` | 编辑 | 追加 POST /delete 路由 |
| `application/article/api/internal/handler/detailhandler.go` 等 | 编辑 | 追加 ArticleDeleteHandler |
| `application/article/api/internal/logic/articledetaillogic.go` 等 | 编辑 | 追加 ArticleDelete logic, 修复 AuthorId string→int64 |
| `application/qa/api/qa.go` | 新建 | QA 网关入口 |
| `application/qa/api/etc/qa-api.yaml` | 新建 | QA 配置 |
| `application/qa/api/internal/config/config.go` | 新建 | QA Config |
| `application/qa/api/internal/svc/servicecontext.go` | 新建 | QA ServiceContext |
| `application/qa/api/internal/types/types.go` | 新建 | QA 类型 |
| `application/qa/api/internal/handler/qahandler.go` | 新建 | QA handler |
| `application/qa/api/internal/logic/qalogic.go` | 新建 | QA logic |

---

## 7. 附带修复

| 文件 | 问题 | 修复 |
|------|------|------|
| `article/api/internal/types/types.go:15` | `ArticleDetailResponse.AuthorId` 类型为 `string`，proto 为 `int64` | 改为 `int64` |
| `article/api/internal/logic/articledetaillogic.go:55` | 使用 `strconv.FormatInt()` 转换 | 去掉转换，直接赋值 int64 |

---

## 8. 接口统计

| 类别 | 已有 | 新增 | 补齐后 |
|------|------|------|--------|
| applet-api | 20 | 18 | 38 |
| article-api | 5 | 1 | 6 |
| chat-api | 6 | 0 | 6 |
| qa-api (新建) | 0 | 9 | 9 |
| **合计** | **31** | **28** | **59** |

---

## 9. 跨字段校验记录

已逐字段校验以下映射关系:

- Tag: 10 个 proto 方法的 request/response 字段与 REST 类型一一对应，userId 从 JWT 获取
- Reply: 5 个 proto 方法字段映射正确，replyUserId/beReplyUserId 仅 CreateReply 中保留于请求体（业务语义需要），DeleteReply 的 userId 从 JWT 获取
- QA: 9 个 proto 方法字段映射正确，userId 从 JWT 获取，tagIds 保持为 proto 的 string 格式
- ArticleDelete: 确认归属 article-api，articleId 从 JSON body 获取，userId 从 JWT 获取
- ConcernedCount: bizId/objId 从 query 获取，proto 不包含 userId
- UpgradeMember: userId 从 JWT 获取，level/durationDays/transactionId/amount/payChannel 从 JSON body 获取
- MemberOrderList: userId 从 JWT 获取，cursor/pageSize 从 query 获取

无多写、漏写字段。

---

## 10. 实施偏离记录

执行过程中发现设计阶段未覆盖的问题，已在实施时修正：

### 10.1 zrpc 客户端包装器缺失

Tag、Reply、QA 三个 gRPC 服务的 proto 只生成了 `pb` 包（protobuf 定义），缺少 go-zero 框架所需的 zrpc 客户端包装器。参照 `concerned/rpc/concerned/concerned.go` 模式，新建了三个包装器：

| 文件 | 说明 |
|------|------|
| `application/tag/rpc/tag/tag.go` | Tag zrpc 客户端 (10 个方法) |
| `application/reply/rpc/reply/reply.go` | Reply zrpc 客户端 (5 个方法) |
| `application/qa/rpc/qa/qa.go` | QA zrpc 客户端 (9 个方法) |

每个包装器包含: 类型别名(proto→client)、接口定义、构造函数、方法实现(通过 pb 生成的 gRPC client)。

### 10.2 Member DurationDays 类型修正

`UpgradeMemberRequest.DurationDays` 和 `MemberOrderItem.DurationDays` 在 proto 中定义为 `int32`，设计文档中误写为 `int64`，实施时修正为 `int32`。

### 10.3 ReplyCreateRequest json tag 修正

`ReplyCreateRequest` 中 `BeReplyUserId` 和 `ParentId` 的 json tag 从 `,optional` 改为 `,omitempty`，因为 Go 标准库 `encoding/json` 不支持 `optional` 选项。

### 10.4 SignatureConf 不含 RewriteThreshold

go-zero v1.10.1 的 `SignatureConf` 字段为 `Strict`、`Expiry`、`PrivateKeys`，无 `RewriteThreshold`。新路由通过 `rest.WithSignature(serverCtx.Config.Signature)` 使用与已有 Like/Follow 路由相同的签名配置，无需额外设置。
