package server

import (
	config "blog-service/config"
	blogProto "blog-service/rpc/blog"
	"context"
	"fmt"

	"github.com/twitchtv/twirp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct{}

// docs on bson struct tags and unmarshalling: https://docs.mongodb.com/drivers/go/current/fundamentals/bson/
type BlogItem struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

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
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	result := BlogItem{}

	// Decode() method unmarshals BSON into result
	unmarshal_err := config.Collection.FindOne(context.TODO(), filter).Decode(&result)
	if unmarshal_err != nil {
		if unmarshal_err == mongo.ErrNoDocuments {
			return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("No documents were found for id: %v", id))
		}
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error finding a blog with ID: %v \nError: %v", id, unmarshal_err))
	}

	return &blogProto.GetBlogResponse{
		Id:      id,
		Title:   result.Title,
		Content: result.Content,
	}, nil
}

func (*Server) UpdateBlog(ctx context.Context, req *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	newBlog := &blogProto.Blog{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	update := bson.D{{Key: "$set", Value: bson.M{"title": newBlog.Title, "content": newBlog.Content}}}

	result, update_err := config.Collection.UpdateOne(context.TODO(), filter, update)
	if update_err != nil || result.MatchedCount == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v", id, newBlog))
	}

	return &blogProto.UpdateBlogResponse{
		Id:      oid.Hex(),
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}, nil
}

func (*Server) DeleteBlog(ctx context.Context, req *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	fmt.Println("In DeleteBlog")
	return &blogProto.DeleteBlogResponse{}, nil
}

func (*Server) ListBlog(ctx context.Context, req *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	fmt.Println("In ListBlog")
	return &blogProto.ListBlogResponse{}, nil
}
