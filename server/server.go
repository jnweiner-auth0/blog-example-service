package server

import (
	config "blog-service/config"
	blogProto "blog-service/rpc/blog"
	"context"
)

type Server struct{}

// use the BlogService interface generated in service.twirp.go as guideline to stub out the expected functions

func (*Server) CreateBlog(ctx context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	data := &blogProto.CreateBlogRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	res, err := config.DB.CreateBlog(data)
	return res, err
}

func (*Server) GetBlog(ctx context.Context, req *blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error) {
	data := &blogProto.GetBlogRequest{
		Id: req.GetId(),
	}

	res, err := config.DB.GetBlog(data)
	return res, err
}

func (*Server) UpdateBlog(ctx context.Context, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	data := &blogProto.UpdateBlogRequest{
		Id:      req.GetId(),
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	res, err := config.DB.UpdateBlog(data)
	return res, err
}

func (*Server) DeleteBlog(ctx context.Context, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	data := &blogProto.DeleteBlogRequest{
		Id: req.GetId(),
	}

	res, err := config.DB.DeleteBlog(data)
	return res, err
}

func (*Server) ListBlog(ctx context.Context, req *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	limit := int64(25)

	if req.GetLimit() > 0 {
		limit = req.GetLimit()
	}

	data := &blogProto.ListBlogRequest{
		Limit: limit,
	}

	res, err := config.DB.ListBlog(data)
	return res, err
}
