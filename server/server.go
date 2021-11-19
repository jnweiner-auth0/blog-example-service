package server

import (
	blogProto "blog-service/rpc/blog"
	"context"
	"fmt"
)

type Server struct{}

// use the BlogService interface generated in service.twirp.go as guideline to stub out the expected functions

func (*Server) CreateBlog(ctx context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	fmt.Printf("In CreateBlog \n req.Title: %v \n req.Content: %v \n", req.Title, req.Content)
	return &blogProto.CreateBlogResponse{}, nil
}

func (*Server) GetBlog(ctx context.Context, req *blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error) {
	fmt.Println("In GetBlog")
	return &blogProto.GetBlogResponse{}, nil
}

func (*Server) UpdateBlog(ctx context.Context, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	fmt.Println("In UpdateBlog")
	return &blogProto.UpdateBlogResponse{}, nil
}

func (*Server) DeleteBlog(ctx context.Context, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	fmt.Println("In DeleteBlog")
	return &blogProto.DeleteBlogResponse{}, nil
}

func (*Server) ListBlog(ctx context.Context, req *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	fmt.Println("In ListBlog")
	return &blogProto.ListBlogResponse{}, nil
}
