package svc

import (
	"ThinkTalk/application/applet/internal/config"
	"ThinkTalk/application/concerned/rpc/concerned"
	"ThinkTalk/application/follow/rpc/follow"
	like "ThinkTalk/application/like/rpc/like"
	"ThinkTalk/application/member/rpc/member"
	"ThinkTalk/application/message/rpc/message"
	"ThinkTalk/application/reply/rpc/reply"
	"ThinkTalk/application/tag/rpc/tag"
	"ThinkTalk/application/user/rpc/user"
	"ThinkTalk/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

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
