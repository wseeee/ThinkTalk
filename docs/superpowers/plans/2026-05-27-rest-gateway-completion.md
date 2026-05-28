# REST 网关补齐实现计划

> **面向 AI 代理的工作者：** 必需子技能：使用 superpowers:subagent-driven-development（推荐）或 superpowers:executing-plans 逐任务实现此计划。步骤使用复选框（`- [ ]`）语法来跟踪进度。

**目标：** 补齐 28 个 REST 接口，涉及 Tag(10)、Reply(5)、QA(9) + ArticleDelete(1) + ConcernedCount/UpgradeMember/MemberOrderList(3)，全部手写 Go。

**架构：** Tag 和 Reply 集成到现有 applet-api，QA 新建独立网关 qa-api（端口 8890），ArticleDelete 集成到 article-api，3 个遗漏接口集成到 applet-api。写操作 JWT+Sig，读操作 JWT。

**技术栈：** Go 1.26, go-zero v1.10.1, gRPC, etcd

**参考文档：** `docs/superpowers/specs/2026-05-27-rest-gateway-completion-design.md`

---

### 任务 1：applet-api config + ServiceContext + yaml 注入 TagRPC/ReplyRPC

**文件：**
- 修改：`application/applet/internal/config/config.go`
- 修改：`application/applet/internal/svc/servicecontext.go`
- 修改：`application/applet/etc/applet-api.yaml`

- [ ] **步骤 1：追加 TagRpc/ReplyRpc 到 Config**

在 `config.go` 的 `Config` 结构体中追加两行：

```go
type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
	}
	Redis       redis.RedisConf
	UserRpc     zrpc.RpcClientConf
	LikeRpc     zrpc.RpcClientConf
	FollowRpc   zrpc.RpcClientConf
	MessageRpc  zrpc.RpcClientConf
	ConcernedRpc zrpc.RpcClientConf
	MemberRpc   zrpc.RpcClientConf
	TagRpc      zrpc.RpcClientConf
	ReplyRpc    zrpc.RpcClientConf
}
```

- [ ] **步骤 2：追加 imports 和 ServiceContext 字段 + 初始化**

`servicecontext.go` — 在 import 块追加：

```go
import (
	"ThinkTalk/application/applet/internal/config"
	"ThinkTalk/application/concerned/rpc/concerned"
	"ThinkTalk/application/follow/rpc/follow"
	like "ThinkTalk/application/like/rpc/like"
	"ThinkTalk/application/member/rpc/member"
	"ThinkTalk/application/message/rpc/message"
	"ThinkTalk/application/tag/rpc/tag"
	"ThinkTalk/application/reply/rpc/reply"
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)
```

`ServiceContext` 结构体追加 `TagRPC` 和 `ReplyRPC`：

```go
type ServiceContext struct {
	Config       config.Config
	RDB          *redis.Redis
	UserRPC      user.User
	LikeRPC      like.Like
	FollowRPC    follow.Follow
	MessageRPC   message.Message
	ConcernedRPC concerned.Concerned
	MemberRPC    member.Member
	TagRPC       tag.Tag
	ReplyRPC     reply.Reply
}
```

`NewServiceContext` 追加初始化：

```go
func NewServiceContext(c config.Config) *ServiceContext {
	rdb, _ := redis.NewRedis(c.Redis)
	userRPC := zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	likeRPC := zrpc.MustNewClient(c.LikeRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	followRPC := zrpc.MustNewClient(c.FollowRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	messageRPC := zrpc.MustNewClient(c.MessageRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	concernedRPC := zrpc.MustNewClient(c.ConcernedRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	memberRPC := zrpc.MustNewClient(c.MemberRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	tagRPC := zrpc.MustNewClient(c.TagRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	replyRPC := zrpc.MustNewClient(c.ReplyRpc, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:       c,
		RDB:          rdb,
		UserRPC:      user.NewUser(userRPC),
		LikeRPC:      like.NewLike(likeRPC),
		FollowRPC:    follow.NewFollow(followRPC),
		MessageRPC:   message.NewMessage(messageRPC),
		ConcernedRPC: concerned.NewConcerned(concernedRPC),
		MemberRPC:    member.NewMember(memberRPC),
		TagRPC:       tag.NewTag(tagRPC),
		ReplyRPC:     reply.NewReply(replyRPC),
	}
}
```

- [ ] **步骤 3：追加 applet-api.yaml 配置**

在 `applet-api.yaml` 末尾追加：

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

- [ ] **步骤 4：编译验证**

运行：`go build ./application/applet/...`
预期：编译成功（types 未定义会报错，但 config+svc 层本身正确）

---

### 任务 2：applet-api types.go 追加 Tag/Reply/遗漏接口类型

**文件：**
- 修改：`application/applet/internal/types/types.go`

- [ ] **步骤 1：在 types.go 末尾追加所有新类型**

```go
// ========== Tag 模块 ==========

type TagCreateRequest struct {
	TagName string `json:"tag_name"`
	TagDesc string `json:"tag_desc"`
}

type TagCreateResponse struct {
	TagId int64 `json:"tag_id"`
}

type TagUpdateRequest struct {
	TagId   int64  `json:"tag_id"`
	TagName string `json:"tag_name"`
	TagDesc string `json:"tag_desc"`
}

type TagDeleteRequest struct {
	TagId int64 `json:"tag_id"`
}

type TagDetailRequest struct {
	TagId int64 `form:"tag_id"`
}

type TagItem struct {
	TagId         int64  `json:"tag_id"`
	TagName       string `json:"tag_name"`
	TagDesc       string `json:"tag_desc"`
	ResourceCount int64  `json:"resource_count"`
	CreateTime    int64  `json:"create_time"`
}

type TagDetailResponse = TagItem

type TagListRequest struct {
	Cursor   int64 `form:"cursor"`
	PageSize int64 `form:"page_size"`
}

type TagListResponse struct {
	Items  []*TagItem `json:"items"`
	Cursor int64      `json:"cursor"`
	IsEnd  bool       `json:"is_end"`
}

type HotTagsRequest struct {
	Limit int32 `form:"limit"`
}

type HotTagsResponse struct {
	Items []*TagItem `json:"items"`
}

type TagResourceRequest struct {
	BizId    string `json:"biz_id"`
	TargetId int64  `json:"target_id"`
	TagId    int64  `json:"tag_id"`
}

type UntagResourceRequest struct {
	BizId    string `json:"biz_id"`
	TargetId int64  `json:"target_id"`
	TagId    int64  `json:"tag_id"`
}

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

type TagsByResourceRequest struct {
	BizId    string `form:"biz_id"`
	TargetId int64  `form:"target_id"`
}

type TagsByResourceResponse struct {
	Items []*TagItem `json:"items"`
}

// ========== Reply 模块 ==========

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

type ReplyDeleteRequest struct {
	ReplyId int64 `json:"reply_id"`
}

type ReplyDetailRequest struct {
	ReplyId int64 `form:"reply_id"`
}

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

type ReplyCountRequest struct {
	BizId    string `form:"biz_id"`
	TargetId int64  `form:"target_id"`
}

type ReplyCountResponse struct {
	ReplyNum     int64 `json:"reply_num"`
	ReplyRootNum int64 `json:"reply_root_num"`
}

// ========== ConcernedCount ==========

type ConcernedCountRequest struct {
	BizId string `form:"biz_id"`
	ObjId int64  `form:"obj_id"`
}

type ConcernedCountResponse struct {
	ConcernedNum int64 `json:"concerned_num"`
}

// ========== Member 遗漏 ==========

type UpgradeMemberRequest struct {
	Level         int32  `json:"level"`
	DurationDays  int64  `json:"duration_days"`
	TransactionId string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	PayChannel    string `json:"pay_channel"`
}

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

- [ ] **步骤 2：编译验证**

运行：`go build ./application/applet/...`
预期：报错 — handler 和 logic 文件未创建（正常，下一步创建）

---

### 任务 3：applet-api Tag handler + logic

**文件：**
- 创建：`application/applet/internal/handler/taghandler.go`
- 创建：`application/applet/internal/logic/taglogic.go`

- [ ] **步骤 1：创建 taghandler.go**

```go
package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

// ========== 写操作 ==========

func TagCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.CreateTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagUpdateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.UpdateTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTag(uid, &req)
		writeJSON(w, resp, err)
	}
}

func TagResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.TagResourceRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagResource(uid, &req)
		writeJSON(w, resp, err)
	}
}

func UntagResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.UntagResourceRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.UntagResource(uid, &req)
		writeJSON(w, resp, err)
	}
}

// ========== 读操作 ==========

func TagDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagDetailRequest
		if v := r.URL.Query().Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TagId)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagDetail(&req)
		writeJSON(w, resp, err)
	}
}

func TagListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagListRequest
		parseQuery(r, &req)
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagList(&req)
		writeJSON(w, resp, err)
	}
}

func HotTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HotTagsRequest
		if v := r.URL.Query().Get("limit"); v != "" {
			json.Unmarshal([]byte(v), &req.Limit)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.HotTags(&req)
		writeJSON(w, resp, err)
	}
}

func ResourcesByTagHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResourcesByTagRequest
		q := r.URL.Query()
		if v := q.Get("tag_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TagId)
		}
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.ResourcesByTag(&req)
		writeJSON(w, resp, err)
	}
}

func TagsByResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagsByResourceRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		l := logic.NewTagLogic(r.Context(), svcCtx)
		resp, err := l.TagsByResource(&req)
		writeJSON(w, resp, err)
	}
}
```

- [ ] **步骤 2：创建 taglogic.go**

```go
package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/tag/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagLogic {
	return &TagLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *TagLogic) CreateTag(userId int64, req *types.TagCreateRequest) (*types.TagCreateResponse, error) {
	resp, err := l.svcCtx.TagRPC.CreateTag(l.ctx, &pb.CreateTagRequest{
		TagName: req.TagName,
		TagDesc: req.TagDesc,
	})
	if err != nil {
		l.Errorf("[CreateTag] rpc err: %v", err)
		return nil, err
	}
	return &types.TagCreateResponse{TagId: resp.TagId}, nil
}

func (l *TagLogic) UpdateTag(userId int64, req *types.TagUpdateRequest) (*types.TagUpdateResponse, error) {
	_, err := l.svcCtx.TagRPC.UpdateTag(l.ctx, &pb.UpdateTagRequest{
		TagId:   req.TagId,
		TagName: req.TagName,
		TagDesc: req.TagDesc,
	})
	if err != nil {
		l.Errorf("[UpdateTag] rpc err: %v", err)
		return nil, err
	}
	return &types.TagUpdateResponse{}, nil
}

func (l *TagLogic) DeleteTag(userId int64, req *types.TagDeleteRequest) (*types.TagDeleteResponse, error) {
	_, err := l.svcCtx.TagRPC.DeleteTag(l.ctx, &pb.DeleteTagRequest{
		TagId: req.TagId,
	})
	if err != nil {
		l.Errorf("[DeleteTag] rpc err: %v", err)
		return nil, err
	}
	return &types.TagDeleteResponse{}, nil
}

func (l *TagLogic) TagDetail(req *types.TagDetailRequest) (*types.TagDetailResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagDetail(l.ctx, &pb.TagDetailRequest{
		TagId: req.TagId,
	})
	if err != nil {
		l.Errorf("[TagDetail] rpc err: %v", err)
		return nil, err
	}
	return &types.TagDetailResponse{
		TagId:         resp.TagId,
		TagName:       resp.TagName,
		TagDesc:       resp.TagDesc,
		ResourceCount: resp.ResourceCount,
		CreateTime:    resp.CreateTime,
	}, nil
}

func (l *TagLogic) TagList(req *types.TagListRequest) (*types.TagListResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagList(l.ctx, &pb.TagListRequest{
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[TagList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.TagItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.TagItem{
			TagId:         item.TagId,
			TagName:       item.TagName,
			TagDesc:       item.TagDesc,
			ResourceCount: item.ResourceCount,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.TagListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *TagLogic) HotTags(req *types.HotTagsRequest) (*types.HotTagsResponse, error) {
	resp, err := l.svcCtx.TagRPC.HotTags(l.ctx, &pb.HotTagsRequest{
		Limit: req.Limit,
	})
	if err != nil {
		l.Errorf("[HotTags] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.TagItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.TagItem{
			TagId:         item.TagId,
			TagName:       item.TagName,
			TagDesc:       item.TagDesc,
			ResourceCount: item.ResourceCount,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.HotTagsResponse{Items: items}, nil
}

func (l *TagLogic) TagResource(userId int64, req *types.TagResourceRequest) (*types.TagResourceResponse, error) {
	_, err := l.svcCtx.TagRPC.TagResource(l.ctx, &pb.TagResourceRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		TagId:    req.TagId,
		UserId:   userId,
	})
	if err != nil {
		l.Errorf("[TagResource] rpc err: %v", err)
		return nil, err
	}
	return &types.TagResourceResponse{}, nil
}

func (l *TagLogic) UntagResource(userId int64, req *types.UntagResourceRequest) (*types.UntagResourceResponse, error) {
	_, err := l.svcCtx.TagRPC.UntagResource(l.ctx, &pb.UntagResourceRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		TagId:    req.TagId,
		UserId:   userId,
	})
	if err != nil {
		l.Errorf("[UntagResource] rpc err: %v", err)
		return nil, err
	}
	return &types.UntagResourceResponse{}, nil
}

func (l *TagLogic) ResourcesByTag(req *types.ResourcesByTagRequest) (*types.ResourcesByTagResponse, error) {
	resp, err := l.svcCtx.TagRPC.ResourcesByTag(l.ctx, &pb.ResourcesByTagRequest{
		TagId:    req.TagId,
		BizId:    req.BizId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[ResourcesByTag] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.ResourceItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.ResourceItem{
			TargetId:   item.TargetId,
			BizId:      item.BizId,
			CreateTime: item.CreateTime,
		})
	}
	return &types.ResourcesByTagResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *TagLogic) TagsByResource(req *types.TagsByResourceRequest) (*types.TagsByResourceResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagsByResource(l.ctx, &pb.TagsByResourceRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[TagsByResource] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.TagItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.TagItem{
			TagId:         item.TagId,
			TagName:       item.TagName,
			TagDesc:       item.TagDesc,
			ResourceCount: item.ResourceCount,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.TagsByResourceResponse{Items: items}, nil
}
```

注意: 新增空的 response 类型 `TagUpdateResponse`、`TagDeleteResponse`、`TagResourceResponse`、`UntagResourceResponse` — 需在 types.go 补充。

- [ ] **步骤 3：补充 types.go 中遗漏的空 response 类型**

在 types.go 中补充：

```go
type TagUpdateResponse struct{}
type TagDeleteResponse struct{}
type TagResourceResponse struct{}
type UntagResourceResponse struct{}
```

- [ ] **步骤 4：编译验证**

运行：`go build ./application/applet/...`
预期：routes.go 报错 — Tag handler 函数引用未声明的路由（正常，下一步补充路由）

---

### 任务 4：applet-api Reply handler + logic

**文件：**
- 创建：`application/applet/internal/handler/replyhandler.go`
- 创建：`application/applet/internal/logic/replylogic.go`

- [ ] **步骤 1：创建 replyhandler.go**

```go
package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/applet/internal/logic"
	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
)

func ReplyCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ReplyCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.CreateReply(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ReplyDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.ReplyDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.DeleteReply(uid, &req)
		writeJSON(w, resp, err)
	}
}

func ReplyDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyDetailRequest
		if v := r.URL.Query().Get("reply_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ReplyId)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyDetail(&req)
		writeJSON(w, resp, err)
	}
}

func ReplyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyListRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		if v := q.Get("sort_type"); v != "" {
			json.Unmarshal([]byte(v), &req.SortType)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyList(&req)
		writeJSON(w, resp, err)
	}
}

func ReplyCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReplyCountRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("target_id"); v != "" {
			json.Unmarshal([]byte(v), &req.TargetId)
		}
		l := logic.NewReplyLogic(r.Context(), svcCtx)
		resp, err := l.ReplyCount(&req)
		writeJSON(w, resp, err)
	}
}
```

- [ ] **步骤 2：创建 replylogic.go**

```go
package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/reply/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReplyLogic {
	return &ReplyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ReplyLogic) CreateReply(userId int64, req *types.ReplyCreateRequest) (*types.ReplyCreateResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.CreateReply(l.ctx, &pb.CreateReplyRequest{
		BizId:         req.BizId,
		TargetId:      req.TargetId,
		ReplyUserId:   userId,
		BeReplyUserId: req.BeReplyUserId,
		ParentId:      req.ParentId,
		Content:       req.Content,
	})
	if err != nil {
		l.Errorf("[CreateReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCreateResponse{ReplyId: resp.ReplyId}, nil
}

func (l *ReplyLogic) DeleteReply(userId int64, req *types.ReplyDeleteRequest) (*types.ReplyDeleteResponse, error) {
	_, err := l.svcCtx.ReplyRPC.DeleteReply(l.ctx, &pb.DeleteReplyRequest{
		ReplyId: req.ReplyId,
		UserId:  userId,
	})
	if err != nil {
		l.Errorf("[DeleteReply] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyDeleteResponse{}, nil
}

func (l *ReplyLogic) ReplyDetail(req *types.ReplyDetailRequest) (*types.ReplyDetailResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyDetail(l.ctx, &pb.ReplyDetailRequest{
		ReplyId: req.ReplyId,
	})
	if err != nil {
		l.Errorf("[ReplyDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Reply == nil {
		return &types.ReplyDetailResponse{}, nil
	}
	return &types.ReplyDetailResponse{Reply: convertReplyItem(resp.Reply)}, nil
}

func convertReplyItem(pb *pb.ReplyItem) *types.ReplyItem {
	if pb == nil {
		return nil
	}
	item := &types.ReplyItem{
		ReplyId:       pb.ReplyId,
		BizId:         pb.BizId,
		TargetId:      pb.TargetId,
		ReplyUserId:   pb.ReplyUserId,
		BeReplyUserId: pb.BeReplyUserId,
		ParentId:      pb.ParentId,
		Content:       pb.Content,
		LikeNum:       pb.LikeNum,
		CreateTime:    pb.CreateTime,
	}
	for _, sub := range pb.SubReplies {
		item.SubReplies = append(item.SubReplies, convertReplyItem(sub))
	}
	return item
}

func (l *ReplyLogic) ReplyList(req *types.ReplyListRequest) (*types.ReplyListResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyList(l.ctx, &pb.ReplyListRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[ReplyList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.ReplyItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, convertReplyItem(item))
	}
	return &types.ReplyListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *ReplyLogic) ReplyCount(req *types.ReplyCountRequest) (*types.ReplyCountResponse, error) {
	resp, err := l.svcCtx.ReplyRPC.ReplyCount(l.ctx, &pb.ReplyCountRequest{
		BizId:    req.BizId,
		TargetId: req.TargetId,
	})
	if err != nil {
		l.Errorf("[ReplyCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ReplyCountResponse{
		ReplyNum:     resp.ReplyNum,
		ReplyRootNum: resp.ReplyRootNum,
	}, nil
}
```

注意: 需要在 types.go 补充 `ReplyDeleteResponse{}`。

- [ ] **步骤 3：补充 types.go 遗漏类型**

```go
type ReplyDeleteResponse struct{}
```

---

### 任务 5：applet-api 遗漏接口 — ConcernedCount/UpgradeMember/MemberOrderList

**文件：**
- 修改：`application/applet/internal/handler/concernedhandler.go`
- 修改：`application/applet/internal/logic/concernedlogic.go`
- 修改：`application/applet/internal/handler/memberhandler.go`
- 修改：`application/applet/internal/logic/memberlogic.go`

- [ ] **步骤 1：在 concernedhandler.go 末尾追加 ConcernedCountHandler**

```go
func ConcernedCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConcernedCountRequest
		q := r.URL.Query()
		if v := q.Get("biz_id"); v != "" {
			req.BizId = v
		}
		if v := q.Get("obj_id"); v != "" {
			json.Unmarshal([]byte(v), &req.ObjId)
		}
		l := logic.NewConcernedLogic(r.Context(), svcCtx)
		resp, err := l.ConcernedCount(&req)
		writeJSON(w, resp, err)
	}
}
```

- [ ] **步骤 2：在 concernedlogic.go 末尾追加 ConcernedCount 方法**

```go
func (l *ConcernedLogic) ConcernedCount(req *types.ConcernedCountRequest) (*types.ConcernedCountResponse, error) {
	resp, err := l.svcCtx.ConcernedRPC.ConcernedCount(l.ctx, &pb.ConcernedCountRequest{
		BizId: req.BizId,
		ObjId: req.ObjId,
	})
	if err != nil {
		l.Errorf("[ConcernedCount] rpc err: %v", err)
		return nil, err
	}
	return &types.ConcernedCountResponse{ConcernedNum: resp.ConcernedNum}, nil
}
```

- [ ] **步骤 3：在 memberhandler.go 末尾追加 UpgradeMemberHandler + MemberOrderListHandler**

```go
func UpgradeMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.UpgradeMemberRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.UpgradeMember(uid, &req)
		writeJSON(w, resp, err)
	}
}

func MemberOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.MemberOrderListRequest
		parseQuery(r, &req)
		l := logic.NewMemberLogic(r.Context(), svcCtx)
		resp, err := l.MemberOrderList(uid, &req)
		writeJSON(w, resp, err)
	}
}
```

- [ ] **步骤 4：在 memberlogic.go 末尾追加 UpgradeMember + MemberOrderList 方法**

```go
func (l *MemberLogic) UpgradeMember(userId int64, req *types.UpgradeMemberRequest) (*types.UpgradeMemberResponse, error) {
	_, err := l.svcCtx.MemberRPC.UpgradeMember(l.ctx, &pb.UpgradeMemberRequest{
		UserId:        userId,
		Level:         req.Level,
		DurationDays:  req.DurationDays,
		TransactionId: req.TransactionId,
		Amount:        req.Amount,
		PayChannel:    req.PayChannel,
	})
	if err != nil {
		l.Errorf("[UpgradeMember] rpc err: %v", err)
		return nil, err
	}
	return &types.UpgradeMemberResponse{}, nil
}

func (l *MemberLogic) MemberOrderList(userId int64, req *types.MemberOrderListRequest) (*types.MemberOrderListResponse, error) {
	resp, err := l.svcCtx.MemberRPC.MemberOrderList(l.ctx, &pb.MemberOrderListRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[MemberOrderList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.MemberOrderItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.MemberOrderItem{
			Id:           item.Id,
			UserId:       item.UserId,
			Level:        item.Level,
			DurationDays: item.DurationDays,
			Amount:       item.Amount,
			PayChannel:   item.PayChannel,
			Status:       item.Status,
			CreateTime:   item.CreateTime,
		})
	}
	return &types.MemberOrderListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
```

注意: 需要在 types.go 补充 `UpgradeMemberResponse{}`。

- [ ] **步骤 5：补充 types.go**

```go
type UpgradeMemberResponse struct{}
```

---

### 任务 6：applet-api routes.go 追加所有新路由

**文件：**
- 修改：`application/applet/internal/handler/routes.go`

- [ ] **步骤 1：在 RegisterHandlers 末尾追加 Tag/Reply/ConcernedCount/Member 扩展路由**

在 `RegisterHandlers` 函数末尾（`member` 路由组之后，闭合大括号之前）追加：

```go
// Tag 读操作 — JWT
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodGet, Path: "/detail", Handler: TagDetailHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/list", Handler: TagListHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/hot", Handler: HotTagsHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/resource/list", Handler: ResourcesByTagHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/resource/tags", Handler: TagsByResourceHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithPrefix("/v1/tag"),
)

// Tag 写操作 — JWT + Signature
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodPost, Path: "/create", Handler: TagCreateHandler(serverCtx)},
		{Method: http.MethodPost, Path: "/update", Handler: TagUpdateHandler(serverCtx)},
		{Method: http.MethodPost, Path: "/delete", Handler: TagDeleteHandler(serverCtx)},
		{Method: http.MethodPost, Path: "/resource/add", Handler: TagResourceHandler(serverCtx)},
		{Method: http.MethodPost, Path: "/resource/remove", Handler: UntagResourceHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithSignature(serverCtx.Config.Signature),
	rest.WithPrefix("/v1/tag"),
)

// Reply 读操作 — JWT
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodGet, Path: "/detail", Handler: ReplyDetailHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/list", Handler: ReplyListHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/count", Handler: ReplyCountHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithPrefix("/v1/reply"),
)

// Reply 写操作 — JWT + Signature
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodPost, Path: "/create", Handler: ReplyCreateHandler(serverCtx)},
		{Method: http.MethodPost, Path: "/delete", Handler: ReplyDeleteHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithSignature(serverCtx.Config.Signature),
	rest.WithPrefix("/v1/reply"),
)
```

- [ ] **步骤 2：修改 existing concerned 路由组，追加 ConcernedCount**

在 concerned 路由组中追加一行：

```go
{
	Method:  http.MethodGet,
	Path:    "/count",
	Handler: ConcernedCountHandler(serverCtx),
},
```

- [ ] **步骤 3：修改 existing member 路由组并新增 UpgradeMember 路由组**

将 member 的 JWT 路由组从：

```go
server.AddRoutes(
	[]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/info",
			Handler: MemberInfoHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/right",
			Handler: MemberRightHandler(serverCtx),
		},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithPrefix("/v1/member"),
)
```

改为追加 MemberOrderList：

```go
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodGet, Path: "/info", Handler: MemberInfoHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/right", Handler: MemberRightHandler(serverCtx)},
		{Method: http.MethodGet, Path: "/orders", Handler: MemberOrderListHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithPrefix("/v1/member"),
)
```

然后在该 JWT 组之后追加新的 UpgradeMember 路由组：

```go
// Member 写操作 — JWT + Signature
server.AddRoutes(
	[]rest.Route{
		{Method: http.MethodPost, Path: "/upgrade", Handler: UpgradeMemberHandler(serverCtx)},
	},
	rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	rest.WithSignature(serverCtx.Config.Signature),
	rest.WithPrefix("/v1/member"),
)
```

- [ ] **步骤 4：编译验证**

运行：`go build ./application/applet/...`
预期：编译成功，无报错

---

### 任务 7：article-api 修复 AuthorId + 新增 ArticleDelete

**文件：**
- 修改：`application/article/api/internal/types/types.go`
- 修改：`application/article/api/internal/logic/articledetaillogic.go`
- 创建：`application/article/api/internal/handler/articledeletehandler.go`
- 创建：`application/article/api/internal/logic/articledeletelogic.go`
- 修改：`application/article/api/internal/handler/routes.go`

- [ ] **步骤 1：修复 types.go 中 AuthorId 类型**

将 `ArticleDetailResponse.AuthorId` 类型从 `string` 改为 `int64`：

```go
type ArticleDetailResponse struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	AuthorId    int64  `json:"author_id"`
	AuthorName  string `json:"author_name"`
}
```

- [ ] **步骤 2：追加 ArticleDeleteRequest 到 types.go 末尾**

```go
type ArticleDeleteRequest struct {
	ArticleId int64 `json:"article_id"`
}
```

- [ ] **步骤 3：修复 articledetaillogic.go**

将第 55 行：
```go
AuthorId: strconv.FormatInt(articleInfo.Article.AuthorId, 10),
```
改为：
```go
AuthorId: articleInfo.Article.AuthorId,
```

注意：确保 `"strconv"` 不再被使用时可从 import 中移除。

- [ ] **步骤 4：创建 articledeletehandler.go**

```go
package handler

import (
	"net/http"

	"ThinkTalk/application/article/api/internal/logic"
	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ArticleDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleDeleteRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewArticleDeleteLogic(r.Context(), svcCtx)
		err := l.ArticleDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, nil)
		}
	}
}
```

- [ ] **步骤 5：创建 articledeletelogic.go**

```go
package logic

import (
	"context"
	"encoding/json"

	"ThinkTalk/application/article/api/internal/svc"
	"ThinkTalk/application/article/api/internal/types"
	"ThinkTalk/application/article/rpc/article"
	"ThinkTalk/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleDeleteLogic {
	return &ArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleDeleteLogic) ArticleDelete(req *types.ArticleDeleteRequest) error {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return xcode.NoLogin
	}
	_, err = l.svcCtx.ArticleRPC.ArticleDelete(l.ctx, &article.ArticleDeleteRequest{
		UserId:    userId,
		ArticleId: req.ArticleId,
	})
	if err != nil {
		logx.Errorf("ArticleDelete req: %v userId: %d error: %v", req, userId, err)
		return err
	}
	return nil
}
```

- [ ] **步骤 6：在 routes.go 中追加 ArticleDelete 路由**

在 article-api 的 JWT 路由组中追加：

```go
{
	Method:  http.MethodPost,
	Path:    "/delete",
	Handler: ArticleDeleteHandler(serverCtx),
},
```

- [ ] **步骤 7：编译验证**

运行：`go build ./application/article/api/...`
预期：编译成功

---

### 任务 8：qa-api 新建独立网关

**文件：**
- 创建：`application/qa/api/qa.go`
- 创建：`application/qa/api/etc/qa-api.yaml`
- 创建：`application/qa/api/internal/config/config.go`
- 创建：`application/qa/api/internal/svc/servicecontext.go`
- 创建：`application/qa/api/internal/types/types.go`

- [ ] **步骤 1：创建目录结构**

```bash
mkdir -p application/qa/api/etc
mkdir -p application/qa/api/internal/config
mkdir -p application/qa/api/internal/handler
mkdir -p application/qa/api/internal/logic
mkdir -p application/qa/api/internal/svc
mkdir -p application/qa/api/internal/types
```

- [ ] **步骤 2：创建 qa.go 入口文件**

```go
package main

import (
	"flag"
	"fmt"

	"ThinkTalk/application/qa/api/internal/config"
	"ThinkTalk/application/qa/api/internal/handler"
	"ThinkTalk/application/qa/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/qa-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
```

- [ ] **步骤 3：创建 qa-api.yaml**

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

- [ ] **步骤 4：创建 config.go**

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

- [ ] **步骤 5：创建 servicecontext.go**

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

- [ ] **步骤 6：创建 types.go**

```go
package types

// ========== PublishQuestion ==========

type PublishQuestionRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	TagIds  string `json:"tag_ids"`
}

type PublishQuestionResponse struct {
	QuestionId int64 `json:"question_id"`
}

// ========== AnswerQuestion ==========

type AnswerQuestionRequest struct {
	QuestionId int64  `json:"question_id"`
	Content    string `json:"content"`
}

type AnswerQuestionResponse struct {
	AnswerId int64 `json:"answer_id"`
}

// ========== AcceptAnswer ==========

type AcceptAnswerRequest struct {
	QuestionId int64 `json:"question_id"`
	AnswerId   int64 `json:"answer_id"`
}

// ========== QuestionDetail ==========

type QuestionDetailRequest struct {
	QuestionId int64 `form:"question_id"`
}

// ========== QuestionList ==========

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

type QuestionDetailResponse struct {
	Question *QuestionItem `json:"question"`
}

// ========== QuestionDelete ==========

type QuestionDeleteRequest struct {
	QuestionId int64 `json:"question_id"`
}

// ========== AnswerList ==========

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

// ========== AnswerDelete ==========

type AnswerDeleteRequest struct {
	AnswerId int64 `json:"answer_id"`
}

// ========== SearchQuestions ==========

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

### 任务 9：qa-api handler + logic + routes

**文件：**
- 创建：`application/qa/api/internal/handler/qahandler.go`
- 创建：`application/qa/api/internal/logic/qalogic.go`
- 创建：`application/qa/api/internal/handler/routes.go`

- [ ] **步骤 1：创建 qahandler.go**

```go
package handler

import (
	"encoding/json"
	"net/http"

	"ThinkTalk/application/qa/api/internal/logic"
	"ThinkTalk/application/qa/api/internal/svc"
	"ThinkTalk/application/qa/api/internal/types"
)

func getUserID(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(json.Number)
	uid, _ := userId.Int64()
	return uid
}

func writeJSON(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(data)
}

func PublishQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.PublishQuestionRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.PublishQuestion(uid, &req)
		writeJSON(w, resp, err)
	}
}

func AnswerQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AnswerQuestionRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AnswerQuestion(uid, &req)
		writeJSON(w, resp, err)
	}
}

func AcceptAnswerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AcceptAnswerRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AcceptAnswer(uid, &req)
		writeJSON(w, resp, err)
	}
}

func QuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.QuestionListRequest
		q := r.URL.Query()
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		if v := q.Get("sort_type"); v != "" {
			json.Unmarshal([]byte(v), &req.SortType)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.Questions(uid, &req)
		writeJSON(w, resp, err)
	}
}

func QuestionDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionDetailRequest
		if v := r.URL.Query().Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.QuestionDetail(&req)
		writeJSON(w, resp, err)
	}
}

func QuestionDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.QuestionDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.QuestionDelete(uid, &req)
		writeJSON(w, resp, err)
	}
}

func AnswerListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AnswerListRequest
		q := r.URL.Query()
		if v := q.Get("question_id"); v != "" {
			json.Unmarshal([]byte(v), &req.QuestionId)
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AnswerList(&req)
		writeJSON(w, resp, err)
	}
}

func AnswerDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := getUserID(r)
		var req types.AnswerDeleteRequest
		json.NewDecoder(r.Body).Decode(&req)
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.AnswerDelete(uid, &req)
		writeJSON(w, resp, err)
	}
}

func SearchQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchQuestionsRequest
		q := r.URL.Query()
		if v := q.Get("keyword"); v != "" {
			req.Keyword = v
		}
		if v := q.Get("cursor"); v != "" {
			json.Unmarshal([]byte(v), &req.Cursor)
		}
		if v := q.Get("page_size"); v != "" {
			json.Unmarshal([]byte(v), &req.PageSize)
		}
		l := logic.NewQALogic(r.Context(), svcCtx)
		resp, err := l.SearchQuestions(&req)
		writeJSON(w, resp, err)
	}
}
```

- [ ] **步骤 2：创建 qalogic.go**

```go
package logic

import (
	"context"

	"ThinkTalk/application/qa/api/internal/svc"
	"ThinkTalk/application/qa/api/internal/types"
	"ThinkTalk/application/qa/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QALogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQALogic(ctx context.Context, svcCtx *svc.ServiceContext) *QALogic {
	return &QALogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QALogic) PublishQuestion(userId int64, req *types.PublishQuestionRequest) (*types.PublishQuestionResponse, error) {
	resp, err := l.svcCtx.QaRPC.PublishQuestion(l.ctx, &pb.PublishQuestionRequest{
		UserId:  userId,
		Title:   req.Title,
		Content: req.Content,
		TagIds:  req.TagIds,
	})
	if err != nil {
		l.Errorf("[PublishQuestion] rpc err: %v", err)
		return nil, err
	}
	return &types.PublishQuestionResponse{QuestionId: resp.QuestionId}, nil
}

func (l *QALogic) AnswerQuestion(userId int64, req *types.AnswerQuestionRequest) (*types.AnswerQuestionResponse, error) {
	resp, err := l.svcCtx.QaRPC.AnswerQuestion(l.ctx, &pb.AnswerQuestionRequest{
		QuestionId: req.QuestionId,
		UserId:     userId,
		Content:    req.Content,
	})
	if err != nil {
		l.Errorf("[AnswerQuestion] rpc err: %v", err)
		return nil, err
	}
	return &types.AnswerQuestionResponse{AnswerId: resp.AnswerId}, nil
}

func (l *QALogic) AcceptAnswer(userId int64, req *types.AcceptAnswerRequest) error {
	_, err := l.svcCtx.QaRPC.AcceptAnswer(l.ctx, &pb.AcceptAnswerRequest{
		QuestionId: req.QuestionId,
		AnswerId:   req.AnswerId,
		UserId:     userId,
	})
	if err != nil {
		l.Errorf("[AcceptAnswer] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) Questions(userId int64, req *types.QuestionListRequest) (*types.QuestionListResponse, error) {
	resp, err := l.svcCtx.QaRPC.Questions(l.ctx, &pb.QuestionsRequest{
		UserId:   userId,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
		SortType: req.SortType,
	})
	if err != nil {
		l.Errorf("[Questions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.QuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.QuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			ViewNum:    item.ViewNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.QuestionListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *QALogic) QuestionDetail(req *types.QuestionDetailRequest) (*types.QuestionDetailResponse, error) {
	resp, err := l.svcCtx.QaRPC.QuestionDetail(l.ctx, &pb.QuestionDetailRequest{
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDetail] rpc err: %v", err)
		return nil, err
	}
	if resp == nil || resp.Question == nil {
		return &types.QuestionDetailResponse{}, nil
	}
	q := resp.Question
	return &types.QuestionDetailResponse{
		Question: &types.QuestionItem{
			Id:         q.Id,
			Title:      q.Title,
			Content:    q.Content,
			AuthorId:   q.AuthorId,
			AnswerNum:  q.AnswerNum,
			ViewNum:    q.ViewNum,
			TagIds:     q.TagIds,
			CreateTime: q.CreateTime,
		},
	}, nil
}

func (l *QALogic) QuestionDelete(userId int64, req *types.QuestionDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.QuestionDelete(l.ctx, &pb.QuestionDeleteRequest{
		UserId:     userId,
		QuestionId: req.QuestionId,
	})
	if err != nil {
		l.Errorf("[QuestionDelete] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) AnswerList(req *types.AnswerListRequest) (*types.AnswerListResponse, error) {
	resp, err := l.svcCtx.QaRPC.AnswerList(l.ctx, &pb.AnswerListRequest{
		QuestionId: req.QuestionId,
		Cursor:     req.Cursor,
		PageSize:   req.PageSize,
	})
	if err != nil {
		l.Errorf("[AnswerList] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.AnswerItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.AnswerItem{
			Id:         item.Id,
			QuestionId: item.QuestionId,
			AuthorId:   item.AuthorId,
			Content:    item.Content,
			IsAccepted: item.IsAccepted,
			LikeNum:    item.LikeNum,
			ReplyNum:   item.ReplyNum,
			CreateTime: item.CreateTime,
		})
	}
	return &types.AnswerListResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}

func (l *QALogic) AnswerDelete(userId int64, req *types.AnswerDeleteRequest) error {
	_, err := l.svcCtx.QaRPC.AnswerDelete(l.ctx, &pb.AnswerDeleteRequest{
		UserId:   userId,
		AnswerId: req.AnswerId,
	})
	if err != nil {
		l.Errorf("[AnswerDelete] rpc err: %v", err)
		return err
	}
	return nil
}

func (l *QALogic) SearchQuestions(req *types.SearchQuestionsRequest) (*types.SearchQuestionsResponse, error) {
	resp, err := l.svcCtx.QaRPC.SearchQuestions(l.ctx, &pb.SearchQuestionsRequest{
		Keyword:  req.Keyword,
		Cursor:   req.Cursor,
		PageSize: req.PageSize,
	})
	if err != nil {
		l.Errorf("[SearchQuestions] rpc err: %v", err)
		return nil, err
	}
	items := make([]*types.SearchQuestionItem, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, &types.SearchQuestionItem{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			AuthorId:   item.AuthorId,
			AnswerNum:  item.AnswerNum,
			TagIds:     item.TagIds,
			CreateTime: item.CreateTime,
		})
	}
	return &types.SearchQuestionsResponse{Items: items, Cursor: resp.Cursor, IsEnd: resp.IsEnd}, nil
}
```

- [ ] **步骤 3：创建 routes.go**

```go
package handler

import (
	"net/http"

	"ThinkTalk/application/qa/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{Method: http.MethodPost, Path: "/publish", Handler: PublishQuestionHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/answer", Handler: AnswerQuestionHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/accept", Handler: AcceptAnswerHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/list", Handler: QuestionsHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/detail", Handler: QuestionDetailHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/delete", Handler: QuestionDeleteHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/answers", Handler: AnswerListHandler(serverCtx)},
			{Method: http.MethodPost, Path: "/answer/delete", Handler: AnswerDeleteHandler(serverCtx)},
			{Method: http.MethodGet, Path: "/search", Handler: SearchQuestionsHandler(serverCtx)},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/v1/qa"),
	)
}
```

---

### 任务 10：全部编译验证

- [ ] **步骤 1：构建所有相关服务**

```bash
go build ./application/applet/...
go build ./application/article/api/...
go build ./application/qa/api/...
```

预期：全部编译成功，无报错。

- [ ] **步骤 2：运行 go vet 检查**

```bash
go vet ./application/applet/...
go vet ./application/article/api/...
go vet ./application/qa/api/...
```

预期：无警告或错误。

- [ ] **步骤 3：检查未使用的 import**

手动确认每个新文件没有未使用的 import。如果 articledetaillogic.go 中移除 `strconv.FormatInt` 后 `"strconv"` 不再使用，需从 import 中移除。

---

### 任务 11：Commit

- [ ] **步骤 1：验证改动清单**

```bash
git diff --stat
```

- [ ] **步骤 2：提交**

```bash
git add application/applet/internal/config/config.go \
        application/applet/internal/svc/servicecontext.go \
        application/applet/etc/applet-api.yaml \
        application/applet/internal/types/types.go \
        application/applet/internal/handler/taghandler.go \
        application/applet/internal/logic/taglogic.go \
        application/applet/internal/handler/replyhandler.go \
        application/applet/internal/logic/replylogic.go \
        application/applet/internal/handler/concernedhandler.go \
        application/applet/internal/logic/concernedlogic.go \
        application/applet/internal/handler/memberhandler.go \
        application/applet/internal/logic/memberlogic.go \
        application/applet/internal/handler/routes.go \
        application/article/api/internal/types/types.go \
        application/article/api/internal/logic/articledetaillogic.go \
        application/article/api/internal/handler/articledeletehandler.go \
        application/article/api/internal/logic/articledeletelogic.go \
        application/article/api/internal/handler/routes.go \
        application/qa/api/

git commit -m "feat: add REST gateways for Tag, Reply, QA services and fill missing endpoints

- Add Tag (10 endpoints) and Reply (5 endpoints) to applet-api
- Create qa-api (9 endpoints) as independent gateway on port 8890
- Add ArticleDelete to article-api
- Add ConcernedCount, UpgradeMember, MemberOrderList to applet-api
- Fix ArticleDetailResponse.AuthorId type string -> int64
- Total: 28 new REST endpoints, 31 -> 59"
```

---

## 实施偏离记录

执行过程中发现以下问题，已记录在设计文档第 10 节：

1. **zrpc 客户端包装器缺失** — Tag/Reply/QA 只有 pb 包，缺少 go-zero zrpc 客户端，新建 3 个包装器文件
2. **Member DurationDays 类型** — proto 为 `int32`，设计文档误写为 `int64`，实施时修正
3. **ReplyCreateRequest json tag** — `,optional` 改为 `,omitempty`（Go 标准库不支持 optional）
4. **SignatureConf 无 RewriteThreshold** — go-zero v1.10.1 的签名配置无此字段，现有 Like/Follow 已正常工作
```
