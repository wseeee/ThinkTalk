package logic

import (
	"context"

	"ThinkTalk/application/applet/internal/svc"
	"ThinkTalk/application/applet/internal/types"
	"ThinkTalk/application/tag/rpc/tag"

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
	resp, err := l.svcCtx.TagRPC.CreateTag(l.ctx, &tag.CreateTagRequest{
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
	_, err := l.svcCtx.TagRPC.UpdateTag(l.ctx, &tag.UpdateTagRequest{
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
	_, err := l.svcCtx.TagRPC.DeleteTag(l.ctx, &tag.DeleteTagRequest{
		TagId: req.TagId,
	})
	if err != nil {
		l.Errorf("[DeleteTag] rpc err: %v", err)
		return nil, err
	}
	return &types.TagDeleteResponse{}, nil
}

func (l *TagLogic) TagDetail(req *types.TagDetailRequest) (*types.TagDetailResponse, error) {
	resp, err := l.svcCtx.TagRPC.TagDetail(l.ctx, &tag.TagDetailRequest{
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
	resp, err := l.svcCtx.TagRPC.TagList(l.ctx, &tag.TagListRequest{
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
	resp, err := l.svcCtx.TagRPC.HotTags(l.ctx, &tag.HotTagsRequest{
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
	_, err := l.svcCtx.TagRPC.TagResource(l.ctx, &tag.TagResourceRequest{
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
	_, err := l.svcCtx.TagRPC.UntagResource(l.ctx, &tag.UntagResourceRequest{
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
	resp, err := l.svcCtx.TagRPC.ResourcesByTag(l.ctx, &tag.ResourcesByTagRequest{
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
	resp, err := l.svcCtx.TagRPC.TagsByResource(l.ctx, &tag.TagsByResourceRequest{
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
