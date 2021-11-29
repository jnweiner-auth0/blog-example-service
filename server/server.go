package server

import (
	config "blog-service/config"
	blogProto "blog-service/rpc/blog"
	"context"
	"fmt"

	"github.com/twitchtv/twirp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct{}

// use the BlogService interface generated in service.twirp.go as guideline to stub out the expected functions
// helpful Mongo driver docs: https://docs.mongodb.com/drivers/go/current/fundamentals/crud/

func (*Server) CreateBlog(ctx context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	data := &blogProto.CreateBlogRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	res, err := config.Collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("There was an error creating a blog: %v", err))
	}

	// type assertion that res.InsertedID is of type primitive.ObjectID
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Cannot convert to oid: %v", ok))
	}

	return &blogProto.CreateBlogResponse{
		Id:      oid.Hex(),
		Title:   data.Title,
		Content: data.Content,
	}, nil
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
