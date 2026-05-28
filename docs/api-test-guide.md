# ThinkTalk 接口测试文档

> 生成日期: 2026-05-27
> 基础地址: `http://localhost:8888` (applet), `http://localhost:<port>` (article), `http://localhost:<port>` (chat), `http://localhost:8890` (qa)

---

## 目录

1. [认证说明](#1-认证说明)
2. [用户服务 (applet-api)](#2-用户服务-applet-api)
3. [文章服务 (article-api)](#3-文章服务-article-api)
4. [聊天服务 (chat-api)](#4-聊天服务-chat-api)
5. [gRPC 内部服务](#5-grpc-内部服务)
6. [问答服务 (qa-api)](#6-问答服务-qa-api)

---

## 1. 认证说明

| 认证方式            | 说明                                    |
| --------------- | ------------------------------------- |
| 无认证             | 公开接口，无需 Token                         |
| JWT             | 请求头携带 `Authorization: Bearer <token>` |
| JWT + Signature | JWT 基础上增加签名校验（由中间件处理）                 |

**获取 Token 流程**: 先调用 `/v1/verification` 发送验证码，再调用 `/v1/login` 获取 `access_token`。

---

## 2. 用户服务 (applet-api)

### 2.1 用户注册

- **接口**: `POST /v1/register`
- **认证**: 无

**请求参数**:

| 参数                | 类型     | 必填  | 说明  |
| ----------------- | ------ | --- | --- |
| name              | string | 是   | 用户名 |
| mobile            | string | 是   | 手机号 |
| password          | string | 是   | 密码  |
| verification_code | string | 是   | 验证码 |

**请求示例**:

```json
{
  "name": "testuser",
  "mobile": "13800138000",
  "password": "123456",
  "verification_code": "123456"
}
```

**响应示例**:

```json
{
  "user_id": 10001,
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "access_expire": 1717000000
  }
}
```

**测试用例**:

| 用例     | 输入     | 期望结果                    |
| ------ | ------ | ----------------------- |
| 正常注册   | 合法参数   | 200, 返回 user_id 和 token |
| 手机号已存在 | 重复手机号  | 错误码, 提示已注册              |
| 验证码错误  | 错误验证码  | 错误码, 提示验证码无效            |
| 参数为空   | 缺少必填字段 | 400, 参数校验失败             |

---

### 2.2 发送验证码

- **接口**: `POST /v1/verification`
- **认证**: 无

**请求参数**:

| 参数     | 类型     | 必填  | 说明  |
| ------ | ------ | --- | --- |
| mobile | string | 是   | 手机号 |

**请求示例**:

```json
{
  "mobile": "13800138000"
}
```

**响应**: 成功返回空 body (200)

**测试用例**:

| 用例      | 输入       | 期望结果        |
| ------- | -------- | ----------- |
| 正常发送    | 合法手机号    | 200         |
| 频率限制    | 短时间内重复发送 | 错误码, 提示稍后重试 |
| 手机号格式错误 | 非法手机号    | 400, 参数校验失败 |

---

### 2.3 用户登录

- **接口**: `POST /v1/login`
- **认证**: 无

**请求参数**:

| 参数                | 类型     | 必填  | 说明  |
| ----------------- | ------ | --- | --- |
| mobile            | string | 是   | 手机号 |
| verification_code | string | 是   | 验证码 |

**请求示例**:

```json
{
  "mobile": "13800138000",
  "verification_code": "123456"
}
```

**响应示例**:

```json
{
  "userId": 10001,
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "access_expire": 1717000000
  }
}
```

**测试用例**:

| 用例    | 输入            | 期望结果          |
| ----- | ------------- | ------------- |
| 正常登录  | 合法手机号 + 正确验证码 | 200, 返回 token |
| 验证码错误 | 错误验证码         | 错误码           |
| 用户不存在 | 未注册手机号        | 错误码           |

---

### 2.4 获取用户信息

- **接口**: `GET /v1/user/info`
- **认证**: JWT + Signature

**请求参数**: 无 (从 Token 中获取 userId)

**响应示例**:

```json
{
  "user_id": 10001,
  "username": "testuser",
  "avatar": "https://example.com/avatar.jpg"
}
```

**测试用例**:

| 用例       | 输入       | 期望结果        |
| -------- | -------- | ----------- |
| 正常获取     | 有效 Token | 200, 返回用户信息 |
| 未登录      | 无 Token  | 401         |
| Token 过期 | 过期 Token | 401         |

---

### 2.5 点赞/踩

- **接口**: `POST /v1/like/thumbup`
- **认证**: JWT + Signature

**请求参数**:

| 参数        | 类型     | 必填  | 说明                  |
| --------- | ------ | --- | ------------------- |
| biz_id    | string | 是   | 业务 ID (如 "article") |
| obj_id    | int64  | 是   | 目标对象 ID             |
| like_type | int32  | 是   | 1=点赞, 2=踩           |

**请求示例**:

```json
{
  "biz_id": "article",
  "obj_id": 1001,
  "like_type": 1
}
```

**响应示例**:

```json
{
  "biz_id": "article",
  "obj_id": 1001,
  "like_num": 42,
  "dislike_num": 3
}
```

**测试用例**:

| 用例   | 输入          | 期望结果               |
| ---- | ----------- | ------------------ |
| 点赞   | like_type=1 | 200, like_num+1    |
| 踩    | like_type=2 | 200, dislike_num+1 |
| 重复点赞 | 已点赞再次点赞     | 取消点赞, like_num-1   |
| 未登录  | 无 Token     | 401                |

---

### 2.6 查询点赞状态

- **接口**: `GET /v1/like/thumbup`
- **认证**: JWT + Signature

**请求参数** (Query):

| 参数        | 类型     | 必填  | 说明      |
| --------- | ------ | --- | ------- |
| biz_id    | string | 是   | 业务 ID   |
| target_id | int64  | 是   | 目标对象 ID |

**响应示例**:

```json
{
  "like_type": 1,
  "thumbup_time": 1716900000
}
```

**测试用例**:

| 用例   | 输入        | 期望结果             |
| ---- | --------- | ---------------- |
| 已点赞  | 有效参数      | 200, like_type=1 |
| 未点赞  | 有效参数      | 200, like_type=0 |
| 参数缺失 | 缺少 biz_id | 400              |

---

### 2.7 关注用户

- **接口**: `POST /v1/follow/follow`
- **认证**: JWT + Signature

**请求参数**:

| 参数               | 类型    | 必填  | 说明       |
| ---------------- | ----- | --- | -------- |
| followed_user_id | int64 | 是   | 被关注用户 ID |

**请求示例**:

```json
{
  "followed_user_id": 10002
}
```

**测试用例**:

| 用例    | 输入      | 期望结果    |
| ----- | ------- | ------- |
| 正常关注  | 有效用户 ID | 200     |
| 关注自己  | 自己的 ID  | 错误码     |
| 重复关注  | 已关注用户   | 幂等, 200 |
| 用户不存在 | 无效 ID   | 错误码     |

---

### 2.8 取消关注

- **接口**: `POST /v1/follow/unfollow`
- **认证**: JWT + Signature

**请求参数**:

| 参数               | 类型    | 必填  | 说明         |
| ---------------- | ----- | --- | ---------- |
| followed_user_id | int64 | 是   | 被取消关注用户 ID |

**测试用例**:

| 用例     | 输入    | 期望结果    |
| ------ | ----- | ------- |
| 正常取消   | 已关注用户 | 200     |
| 未关注时取消 | 未关注用户 | 幂等, 200 |

---

### 2.9 关注列表

- **接口**: `GET /v1/follow/list`
- **认证**: JWT + Signature

**请求参数** (Query):

| 参数        | 类型    | 必填  | 说明          |
| --------- | ----- | --- | ----------- |
| cursor    | int64 | 否   | 分页游标, 默认 0  |
| page_size | int64 | 否   | 每页数量, 默认 10 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 1,
      "followed_user_id": 10002,
      "create_time": 1716900000,
      "fans_count": 128
    }
  ],
  "cursor": 10002,
  "is_end": false
}
```

**测试用例**:

| 用例  | 输入           | 期望结果                       |
| --- | ------------ | -------------------------- |
| 首页  | cursor=0     | 200, 返回列表                  |
| 翻页  | cursor=上次返回值 | 200, 返回下一页                 |
| 空列表 | 新用户          | 200, items=[], is_end=true |

---

### 2.10 粉丝列表

- **接口**: `GET /v1/follow/fans`
- **认证**: JWT + Signature

**请求参数**: 同关注列表

**测试用例**: 同关注列表

---

### 2.11 通知列表

- **接口**: `GET /v1/message/list`
- **认证**: JWT

**请求参数** (Query):

| 参数         | 类型    | 必填  | 说明     |
| ---------- | ----- | --- | ------ |
| notif_type | int32 | 否   | 通知类型过滤 |
| cursor     | int64 | 否   | 分页游标   |
| page_size  | int64 | 否   | 每页数量   |

**响应示例**:

```json
{
  "items": [
    {
      "id": 1,
      "type": 1,
      "title": "有人点赞了你的文章",
      "content": "...",
      "is_read": false,
      "trigger_user_id": 10002,
      "ref_id": 1001,
      "create_time": 1716900000
    }
  ],
  "cursor": 1,
  "is_end": true
}
```

**测试用例**:

| 用例    | 输入                 | 期望结果        |
| ----- | ------------------ | ----------- |
| 全部通知  | 不传 notif_type      | 200, 返回所有类型 |
| 按类型过滤 | notif_type=1       | 200, 只返回该类型 |
| 分页    | cursor + page_size | 200, 正确分页   |

---

### 2.12 未读通知数

- **接口**: `GET /v1/message/unread`
- **认证**: JWT

**响应示例**:

```json
{
  "total": 5,
  "type_counts": {
    "1": 3,
    "2": 2
  }
}
```

---

### 2.13 标记单条通知已读

- **接口**: `POST /v1/message/read`
- **认证**: JWT

**请求参数**:

| 参数              | 类型    | 必填  | 说明    |
| --------------- | ----- | --- | ----- |
| notification_id | int64 | 是   | 通知 ID |

---

### 2.14 标记全部通知已读

- **接口**: `POST /v1/message/readall`
- **认证**: JWT

**请求参数**:

| 参数         | 类型    | 必填  | 说明            |
| ---------- | ----- | --- | ------------- |
| notif_type | int32 | 否   | 通知类型, 不传则全部标记 |

---

### 2.15 收藏

- **接口**: `POST /v1/concerned/add`
- **认证**: JWT + Signature

**请求参数**:

| 参数     | 类型     | 必填  | 说明      |
| ------ | ------ | --- | ------- |
| biz_id | string | 是   | 业务 ID   |
| obj_id | int64  | 是   | 目标对象 ID |

---

### 2.16 取消收藏

- **接口**: `POST /v1/concerned/cancel`
- **认证**: JWT + Signature

**请求参数**: 同收藏

---

### 2.17 检查收藏状态

- **接口**: `GET /v1/concerned/check`
- **认证**: JWT + Signature

**请求参数** (Query):

| 参数     | 类型     | 必填  | 说明      |
| ------ | ------ | --- | ------- |
| biz_id | string | 是   | 业务 ID   |
| obj_id | int64  | 是   | 目标对象 ID |

**响应示例**:

```json
{
  "is_concerned": true
}
```

---

### 2.18 收藏列表

- **接口**: `GET /v1/concerned/list`
- **认证**: JWT + Signature

**请求参数** (Query):

| 参数        | 类型     | 必填  | 说明    |
| --------- | ------ | --- | ----- |
| biz_id    | string | 是   | 业务 ID |
| cursor    | int64  | 否   | 分页游标  |
| page_size | int64  | 否   | 每页数量  |

---

### 2.19 会员信息

- **接口**: `GET /v1/member/info`
- **认证**: JWT

**响应示例**:

```json
{
  "user_id": 10001,
  "level": 2,
  "level_name": "黄金会员",
  "expire_time": 1720000000,
  "status": 1
}
```

---

### 2.20 会员权益检查

- **接口**: `GET /v1/member/right`
- **认证**: JWT

**请求参数** (Query):

| 参数        | 类型     | 必填  | 说明   |
| --------- | ------ | --- | ---- |
| right_key | string | 是   | 权益标识 |

**响应示例**:

```json
{
  "has_right": true,
  "level": 2
}
```

### 2.21 创建标签

- **接口**: `POST /v1/tag/create`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_name | string | 是 | 标签名 |
| tag_desc | string | 否 | 标签描述 |

**请求示例**:

```json
{
  "tag_name": "Go语言",
  "tag_desc": "Go 相关技术讨论"
}
```

**响应示例**:

```json
{
  "tag_id": 1001
}
```

---

### 2.22 更新标签

- **接口**: `POST /v1/tag/update`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_id | int64 | 是 | 标签 ID |
| tag_name | string | 是 | 标签名 |
| tag_desc | string | 否 | 标签描述 |

---

### 2.23 删除标签

- **接口**: `POST /v1/tag/delete`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_id | int64 | 是 | 标签 ID |

---

### 2.24 标签详情

- **接口**: `GET /v1/tag/detail`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_id | int64 | 是 | 标签 ID |

**响应示例**:

```json
{
  "tag_id": 1001,
  "tag_name": "Go语言",
  "tag_desc": "Go 相关技术讨论",
  "resource_count": 42,
  "create_time": 1716900000
}
```

---

### 2.25 标签列表

- **接口**: `GET /v1/tag/list`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "tag_id": 1001,
      "tag_name": "Go语言",
      "tag_desc": "Go 相关技术讨论",
      "resource_count": 42,
      "create_time": 1716900000
    }
  ],
  "cursor": 1001,
  "is_end": false
}
```

---

### 2.26 热门标签

- **接口**: `GET /v1/tag/hot`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int32 | 否 | 数量上限 |

---

### 2.27 打标签

- **接口**: `POST /v1/tag/resource/add`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID (如 "article") |
| target_id | int64 | 是 | 资源 ID |
| tag_id | int64 | 是 | 标签 ID |

---

### 2.28 去标签

- **接口**: `POST /v1/tag/resource/remove`
- **认证**: JWT + Signature

**请求参数**: 同打标签

---

### 2.29 标签下的资源列表

- **接口**: `GET /v1/tag/resource/list`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_id | int64 | 是 | 标签 ID |
| biz_id | string | 否 | 业务 ID 过滤 |
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "target_id": 1001,
      "biz_id": "article",
      "create_time": 1716900000
    }
  ],
  "cursor": 1001,
  "is_end": false
}
```

---

### 2.30 资源的标签列表

- **接口**: `GET /v1/tag/resource/tags`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID |
| target_id | int64 | 是 | 资源 ID |

---

### 2.31 发表评论

- **接口**: `POST /v1/reply/create`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID |
| target_id | int64 | 是 | 评论目标 ID |
| be_reply_user_id | int64 | 否 | 被回复用户 ID |
| parent_id | int64 | 否 | 父评论 ID (0=根评论) |
| content | string | 是 | 评论内容 |

**请求示例**:

```json
{
  "biz_id": "article",
  "target_id": 1001,
  "content": "写得很不错!"
}
```

**响应示例**:

```json
{
  "reply_id": 2001
}
```

---

### 2.32 删除评论

- **接口**: `POST /v1/reply/delete`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reply_id | int64 | 是 | 评论 ID |

---

### 2.33 评论详情

- **接口**: `GET /v1/reply/detail`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reply_id | int64 | 是 | 评论 ID |

**响应示例**:

```json
{
  "reply": {
    "reply_id": 2001,
    "biz_id": "article",
    "target_id": 1001,
    "reply_user_id": 10001,
    "be_reply_user_id": 0,
    "parent_id": 0,
    "content": "写得很不错!",
    "like_num": 5,
    "create_time": 1716900000,
    "sub_replies": []
  }
}
```

---

### 2.34 评论列表

- **接口**: `GET /v1/reply/list`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID |
| target_id | int64 | 是 | 目标 ID |
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |
| sort_type | int32 | 否 | 排序: 0=时间倒序, 1=点赞倒序 |

---

### 2.35 评论计数

- **接口**: `GET /v1/reply/count`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID |
| target_id | int64 | 是 | 目标 ID |

**响应示例**:

```json
{
  "reply_num": 28,
  "reply_root_num": 10
}
```

---

### 2.36 收藏计数

- **接口**: `GET /v1/concerned/count`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| biz_id | string | 是 | 业务 ID |
| obj_id | int64 | 是 | 目标对象 ID |

**响应示例**:

```json
{
  "concerned_num": 128
}
```

---

### 2.37 升级会员

- **接口**: `POST /v1/member/upgrade`
- **认证**: JWT + Signature

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| level | int32 | 是 | 会员等级 |
| duration_days | int32 | 是 | 购买天数 |
| transaction_id | string | 是 | 交易流水号 |
| amount | int64 | 是 | 支付金额(分) |
| pay_channel | string | 是 | 支付渠道 |

**请求示例**:

```json
{
  "level": 2,
  "duration_days": 365,
  "transaction_id": "TXN20260527001",
  "amount": 29900,
  "pay_channel": "wechat"
}
```

---

### 2.38 会员订单列表

- **接口**: `GET /v1/member/orders`
- **认证**: JWT

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 1,
      "user_id": 10001,
      "level": 2,
      "duration_days": 365,
      "amount": 29900,
      "pay_channel": "wechat",
      "status": 1,
      "create_time": 1716900000
    }
  ],
  "cursor": 1,
  "is_end": true
}
```

---

## 3. 文章服务 (article-api)

### 3.1 上传文章封面

- **接口**: `POST /v1/article/upload/cover`
- **认证**: JWT
- **Content-Type**: `multipart/form-data`

**请求参数**:

| 参数   | 类型   | 必填  | 说明   |
| ---- | ---- | --- | ---- |
| file | file | 是   | 图片文件 |

**响应示例**:

```json
{
  "cover_url": "https://cdn.example.com/covers/xxx.jpg"
}
```

**测试用例**:

| 用例    | 输入         | 期望结果        |
| ----- | ---------- | ----------- |
| 正常上传  | JPG/PNG 图片 | 200, 返回 URL |
| 文件过大  | 超过限制大小     | 错误码         |
| 非图片文件 | txt/pdf 文件 | 错误码         |

---

### 3.2 发布文章

- **接口**: `POST /v1/article/publish`
- **认证**: JWT

**请求参数**:

| 参数          | 类型     | 必填  | 说明       |
| ----------- | ------ | --- | -------- |
| title       | string | 是   | 文章标题     |
| content     | string | 是   | 文章内容     |
| description | string | 否   | 文章摘要     |
| cover       | string | 否   | 封面图片 URL |

**请求示例**:

```json
{
  "title": "Go-Zero 微服务实践",
  "content": "# Go-Zero\n\n这是一篇关于 Go-Zero 的文章...",
  "description": "Go-Zero 微服务框架入门指南",
  "cover": "https://cdn.example.com/covers/xxx.jpg"
}
```

**响应示例**:

```json
{
  "article_id": 1001
}
```

**测试用例**:

| 用例   | 输入         | 期望结果               |
| ---- | ---------- | ------------------ |
| 正常发布 | 完整参数       | 200, 返回 article_id |
| 标题为空 | title=""   | 400, 参数校验失败        |
| 内容为空 | content="" | 400, 参数校验失败        |

---

### 3.3 文章详情

- **接口**: `GET /v1/article/detail`
- **认证**: JWT

**请求参数** (Query):

| 参数         | 类型    | 必填  | 说明    |
| ---------- | ----- | --- | ----- |
| article_id | int64 | 是   | 文章 ID |

**响应示例**:

```json
{
  "title": "Go-Zero 微服务实践",
  "content": "...",
  "description": "...",
  "cover": "...",
  "author_id": "10001",
  "author_name": "testuser"
}
```

---

### 3.4 文章列表

- **接口**: `GET /v1/article/list`
- **认证**: JWT

**请求参数** (Query):

| 参数         | 类型    | 必填  | 说明    |
| ---------- | ----- | --- | ----- |
| author_id  | int64 | 否   | 按作者过滤 |
| cursor     | int64 | 否   | 分页游标  |
| page_size  | int64 | 否   | 每页数量  |
| sort_type  | int32 | 否   | 排序方式  |
| article_id | int64 | 否   | 文章 ID |

---

### 3.5 搜索文章

- **接口**: `GET /v1/article/search`
- **认证**: JWT

**请求参数** (Query):

| 参数        | 类型     | 必填  | 说明    |
| --------- | ------ | --- | ----- |
| keyword   | string | 是   | 搜索关键词 |
| cursor    | int64  | 否   | 分页游标  |
| page_size | int64  | 否   | 每页数量  |

**响应示例**:

```json
{
  "articles": [
    {
      "article_id": 1001,
      "title": "Go-Zero 实践",
      "description": "...",
      "cover": "...",
      "author_id": 10001,
      "author_name": "testuser",
      "like_num": 42,
      "comment_num": 8,
      "publish_time": 1716900000
    }
  ],
  "cursor": 1001,
  "is_end": false
}
```

### 3.6 删除文章

- **接口**: `POST /v1/article/delete`
- **认证**: JWT

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int64 | 是 | 文章 ID |

**测试用例**:

| 用例 | 输入 | 期望结果 |
|------|------|----------|
| 正常删除 | 自己发布的文章 | 200 |
| 删除他人文章 | 非自己的文章 | 错误码 |
| 文章不存在 | 无效 ID | 错误码 |

---

## 4. 聊天服务 (chat-api)

### 4.1 会话列表

- **接口**: `GET /v1/chat/conversations`
- **认证**: JWT

**请求参数** (Query):

| 参数        | 类型    | 必填  | 说明   |
| --------- | ----- | --- | ---- |
| cursor    | int64 | 否   | 分页游标 |
| page_size | int64 | 否   | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 1,
      "target_user_id": 10002,
      "last_message": "你好",
      "last_message_time": 1716900000,
      "unread_count": 3
    }
  ],
  "cursor": 1,
  "is_end": false
}
```

---

### 4.2 消息列表

- **接口**: `GET /v1/chat/messages`
- **认证**: JWT

**请求参数** (Query):

| 参数              | 类型    | 必填  | 说明    |
| --------------- | ----- | --- | ----- |
| conversation_id | int64 | 是   | 会话 ID |
| cursor          | int64 | 否   | 分页游标  |
| page_size       | int64 | 否   | 每页数量  |

---

### 4.3 发送消息 (REST)

- **接口**: `POST /v1/chat/send`
- **认证**: JWT

**请求参数**:

| 参数          | 类型     | 必填  | 说明          |
| ----------- | ------ | --- | ----------- |
| receiver_id | int64  | 是   | 接收者 ID      |
| content     | string | 是   | 消息内容        |
| msg_type    | int32  | 是   | 消息类型 (1=文本) |

---

### 4.4 标记会话已读

- **接口**: `POST /v1/chat/markread`
- **认证**: JWT

**请求参数**:

| 参数              | 类型    | 必填  | 说明    |
| --------------- | ----- | --- | ----- |
| conversation_id | int64 | 是   | 会话 ID |

---

### 4.5 未读消息数

- **接口**: `GET /v1/chat/unread`
- **认证**: JWT

**响应示例**:

```json
{
  "total": 5
}
```

---

### 4.6 WebSocket 连接

- **接口**: `GET /v1/chat/ws`
- **认证**: JWT
- **协议**: WebSocket

**连接方式**:

```
ws://localhost:<port>/v1/chat/ws?token=<jwt_token>
```

**客户端发送消息格式**:

```json
{
  "type": "message",
  "receiver_id": 10002,
  "content": "你好",
  "msg_type": 1
}
```

**服务端推送消息格式**:

```json
{
  "type": "message",
  "sender_id": 10001,
  "content": "你好",
  "msg_type": 1,
  "create_time": 1716900000
}
```

**服务端确认格式 (ACK)**:

```json
{
  "type": "ack",
  "message": "sent"
}
```

**测试用例**:

| 用例       | 输入       | 期望结果            |
| -------- | -------- | --------------- |
| 建立连接     | 有效 Token | WebSocket 握手成功  |
| 发送消息     | 正确格式     | 收到 ack          |
| 接收消息     | 对方发送消息   | 收到 message 类型推送 |
| Token 无效 | 过期 Token | 连接被拒绝           |
| 格式错误     | 错误 JSON  | 收到 error 类型     |

---

## 5. gRPC 内部服务

> 以下为内部微服务 RPC 接口，通过 etcd 服务发现调用，不直接对外暴露。

### 5.1 User RPC (端口 8080)

| 方法           | 请求                                 | 响应                               | 说明       |
| ------------ | ---------------------------------- | -------------------------------- | -------- |
| Register     | username, mobile, avatar, password | userId                           | 用户注册     |
| FindById     | userId                             | userId, username, mobile, avatar | 按 ID 查用户 |
| FindByMobile | mobile                             | userId, username, mobile, avatar | 按手机号查用户  |
| SendSms      | userId, mobile                     | -                                | 发送短信     |

### 5.2 Article RPC

| 方法             | 请求                                            | 响应                        | 说明   |
| -------------- | --------------------------------------------- | ------------------------- | ---- |
| Publish        | userId, title, content, description, cover    | articleId                 | 发布文章 |
| Articles       | userId, cursor, pageSize, sortType, articleId | articles[], isEnd, cursor | 文章列表 |
| ArticleDelete  | userId, articleId                             | -                         | 删除文章 |
| ArticleDetail  | articleId                                     | article                   | 文章详情 |
| SearchArticles | keyword, cursor, pageSize                     | items[], cursor, isEnd    | 搜索文章 |

### 5.3 Follow RPC

| 方法         | 请求                       | 响应                     | 说明   |
| ---------- | ------------------------ | ---------------------- | ---- |
| Follow     | userId, followedUserId   | -                      | 关注   |
| UnFollow   | userId, followedUserId   | -                      | 取消关注 |
| FollowList | userId, cursor, pageSize | items[], cursor, isEnd | 关注列表 |
| FansList   | userId, cursor, pageSize | items[], cursor, isEnd | 粉丝列表 |

### 5.4 Like RPC

| 方法        | 请求                             | 响应                                | 说明     |
| --------- | ------------------------------ | --------------------------------- | ------ |
| Thumbup   | bizId, objId, userId, likeType | bizId, objId, likeNum, dislikeNum | 点赞/踩   |
| IsThumbup | bizId, targetId, userId        | userThumbups map                  | 查询点赞状态 |

### 5.5 Tag RPC

| 方法             | 请求                             | 响应                                     | 说明     |
| -------------- | ------------------------------ | -------------------------------------- | ------ |
| CreateTag      | tagName, tagDesc               | tagId                                  | 创建标签   |
| UpdateTag      | tagId, tagName, tagDesc        | -                                      | 更新标签   |
| DeleteTag      | tagId                          | -                                      | 删除标签   |
| TagDetail      | tagId                          | tagId, tagName, tagDesc, resourceCount | 标签详情   |
| TagList        | cursor, pageSize               | items[], cursor, isEnd                 | 标签列表   |
| HotTags        | limit                          | items[]                                | 热门标签   |
| TagResource    | bizId, targetId, tagId, userId | -                                      | 打标签    |
| UntagResource  | bizId, targetId, tagId, userId | -                                      | 移除标签   |
| ResourcesByTag | tagId, bizId, cursor, pageSize | items[], cursor, isEnd                 | 按标签查资源 |
| TagsByResource | bizId, targetId                | items[]                                | 按资源查标签 |

### 5.6 Reply RPC

| 方法          | 请求                                                             | 响应                     | 说明   |
| ----------- | -------------------------------------------------------------- | ---------------------- | ---- |
| CreateReply | bizId, targetId, replyUserId, beReplyUserId, parentId, content | replyId                | 发表回复 |
| DeleteReply | replyId, userId                                                | -                      | 删除回复 |
| ReplyDetail | replyId                                                        | reply                  | 回复详情 |
| ReplyList   | bizId, targetId, cursor, pageSize, sortType                    | items[], cursor, isEnd | 回复列表 |
| ReplyCount  | bizId, targetId                                                | replyNum, replyRootNum | 回复数量 |

### 5.7 Message RPC

| 方法               | 请求                             | 响应                     | 说明   |
| ---------------- | ------------------------------ | ---------------------- | ---- |
| NotificationList | userId, cursor, pageSize, type | items[], cursor, isEnd | 通知列表 |
| UnreadCount      | userId                         | total, typeCounts      | 未读数  |
| MarkRead         | userId, notificationId         | -                      | 标记已读 |
| MarkAllRead      | userId, type                   | -                      | 全部已读 |

### 5.8 Concerned RPC

| 方法              | 请求                              | 响应                     | 说明   |
| --------------- | ------------------------------- | ---------------------- | ---- |
| AddConcerned    | bizId, objId, userId            | -                      | 添加收藏 |
| CancelConcerned | bizId, objId, userId            | -                      | 取消收藏 |
| IsConcerned     | bizId, objId, userId            | isConcerned            | 检查收藏 |
| ConcernedList   | userId, bizId, cursor, pageSize | items[], cursor, isEnd | 收藏列表 |
| ConcernedCount  | bizId, objId                    | concernedNum           | 收藏数量 |

### 5.9 QA RPC

| 方法              | 请求                                 | 响应                     | 说明   |
| --------------- | ---------------------------------- | ---------------------- | ---- |
| PublishQuestion | userId, title, content, tagIds     | questionId             | 发布问题 |
| AnswerQuestion  | questionId, userId, content        | answerId               | 回答问题 |
| AcceptAnswer    | questionId, answerId, userId       | -                      | 采纳回答 |
| Questions       | userId, cursor, pageSize, sortType | items[], cursor, isEnd | 问题列表 |
| QuestionDetail  | questionId                         | question               | 问题详情 |
| QuestionDelete  | userId, questionId                 | -                      | 删除问题 |
| AnswerList      | questionId, cursor, pageSize       | items[], cursor, isEnd | 回答列表 |
| AnswerDelete    | userId, answerId                   | -                      | 删除回答 |
| SearchQuestions | keyword, cursor, pageSize          | items[], cursor, isEnd | 搜索问题 |

### 5.10 Chat RPC

| 方法            | 请求                                     | 响应                     | 说明   |
| ------------- | -------------------------------------- | ---------------------- | ---- |
| SendMessage   | senderId, receiverId, content, msgType | messageId              | 发送消息 |
| Conversations | userId, cursor, pageSize               | items[], cursor, isEnd | 会话列表 |
| Messages      | conversationId, cursor, pageSize       | items[], cursor, isEnd | 消息列表 |
| MarkRead      | userId, conversationId                 | -                      | 标记已读 |
| UnreadCount   | userId                                 | total                  | 未读数  |

### 5.11 Member RPC

| 方法               | 请求                                                             | 响应                                           | 说明   |
| ---------------- | -------------------------------------------------------------- | -------------------------------------------- | ---- |
| MemberInfo       | userId                                                         | userId, level, levelName, expireTime, status | 会员信息 |
| UpgradeMember    | userId, level, durationDays, transactionId, amount, payChannel | -                                            | 升级会员 |
| CheckMemberRight | userId, rightKey                                               | hasRight, level                              | 检查权益 |
| MemberOrderList  | userId, cursor, pageSize                                       | items[], cursor, isEnd                       | 订单列表 |

---

---

## 6. 问答服务 (qa-api)

- **基础地址**: `http://localhost:8890`
- **认证**: 全部接口使用 JWT

### 6.1 发布问题

- **接口**: `POST /v1/qa/publish`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 问题标题 |
| content | string | 是 | 问题内容 |
| tag_ids | string | 否 | 标签 ID 列表 |

**请求示例**:

```json
{
  "title": "Go-Zero 如何集成 Redis?",
  "content": "刚开始学 go-zero，想知道怎么配置 Redis...",
  "tag_ids": "1,2,3"
}
```

**响应示例**:

```json
{
  "question_id": 3001
}
```

**测试用例**:

| 用例 | 输入 | 期望结果 |
|------|------|----------|
| 正常发布 | 完整参数 | 200, 返回 question_id |
| 标题为空 | title="" | 400 |

---

### 6.2 回答问题

- **接口**: `POST /v1/qa/answer`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int64 | 是 | 问题 ID |
| content | string | 是 | 回答内容 |

**响应示例**:

```json
{
  "answer_id": 4001
}
```

---

### 6.3 采纳回答

- **接口**: `POST /v1/qa/accept`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int64 | 是 | 问题 ID |
| answer_id | int64 | 是 | 回答 ID |

**测试用例**:

| 用例 | 输入 | 期望结果 |
|------|------|----------|
| 正常采纳 | 自己发布的问题 | 200 |
| 采纳他人问题的回答 | 非自己发布 | 错误码 |

---

### 6.4 问题列表

- **接口**: `GET /v1/qa/list`

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |
| sort_type | int32 | 否 | 排序: 0=时间, 1=热度 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 3001,
      "title": "Go-Zero 如何集成 Redis?",
      "content": "刚开始学 go-zero...",
      "author_id": 10001,
      "answer_num": 5,
      "view_num": 128,
      "tag_ids": "1,2,3",
      "create_time": 1716900000
    }
  ],
  "cursor": 3001,
  "is_end": false
}
```

---

### 6.5 问题详情

- **接口**: `GET /v1/qa/detail`

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int64 | 是 | 问题 ID |

**响应示例**:

```json
{
  "question": {
    "id": 3001,
    "title": "Go-Zero 如何集成 Redis?",
    "content": "刚开始学 go-zero...",
    "author_id": 10001,
    "answer_num": 5,
    "view_num": 128,
    "tag_ids": "1,2,3",
    "create_time": 1716900000
  }
}
```

---

### 6.6 删除问题

- **接口**: `POST /v1/qa/delete`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int64 | 是 | 问题 ID |

---

### 6.7 回答列表

- **接口**: `GET /v1/qa/answers`

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question_id | int64 | 是 | 问题 ID |
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 4001,
      "question_id": 3001,
      "author_id": 10002,
      "content": "redigo 是最简单的选择...",
      "is_accepted": true,
      "like_num": 15,
      "reply_num": 3,
      "create_time": 1716900000
    }
  ],
  "cursor": 4001,
  "is_end": false
}
```

---

### 6.8 删除回答

- **接口**: `POST /v1/qa/answer/delete`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| answer_id | int64 | 是 | 回答 ID |

---

### 6.9 搜索问题

- **接口**: `GET /v1/qa/search`

**请求参数** (Query):

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| cursor | int64 | 否 | 分页游标 |
| page_size | int64 | 否 | 每页数量 |

**响应示例**:

```json
{
  "items": [
    {
      "id": 3001,
      "title": "Go-Zero 如何集成 Redis?",
      "content": "刚开始学 go-zero...",
      "author_id": 10001,
      "answer_num": 5,
      "tag_ids": "1,2,3",
      "create_time": 1716900000
    }
  ],
  "cursor": 3001,
  "is_end": false
}
```

---

## 附录: 接口统计

| 类别                              | 数量     |
| ------------------------------- | ------ |
| REST 接口 (applet-api)            | 38     |
| REST 接口 (article-api)           | 6      |
| REST 接口 (chat-api, 含 WebSocket) | 6      |
| REST 接口 (qa-api)                | 9      |
| **REST 接口合计**                   | **59** |
| gRPC 接口 (11 个服务)                | 57     |
| **总接口数**                        | **116** |
