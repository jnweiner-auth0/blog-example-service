package server

import (
	config "blog-service/config"
	blogProto "blog-service/rpc/blog"
	"context"
	"fmt"
	"strconv"

	"github.com/twitchtv/twirp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
// helpful Postgres walkthrough: https://www.calhoun.io/using-postgresql-with-go/ and docs https://pkg.go.dev/database/sql

func (*Server) CreateBlog(ctx context.Context, req *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	data := &blogProto.CreateBlogRequest{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	// postgres variant

	if config.Database == "postgres" {
		sqlStatement := "INSERT INTO blogs (title, content) VALUES ($1, $2) RETURNING id"
		id := 0
		err := config.SqlDB.QueryRow(sqlStatement, data.Title, data.Content).Scan(&id)
		if err != nil {
			return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("There was an error creating a blog: %v", err))
		}

		return &blogProto.CreateBlogResponse{
			Id:      strconv.Itoa(id),
			Title:   data.Title,
			Content: data.Content,
		}, nil
	}

	// mongo variant

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

	// postgres variant

	if config.Database == "postgres" {
		var title string
		var content string

		sqlStatement := "SELECT title, content FROM blogs WHERE id=$1"

		err := config.SqlDB.QueryRow(sqlStatement, id).Scan(&title, &content)
		if err != nil {
			return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("No documents were found for id: %v, err: %v", id, err))
		}

		return &blogProto.GetBlogResponse{
			Id:      id,
			Title:   title,
			Content: content,
		}, nil
	}

	// mongo variant

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

	newBlog := &blogProto.Blog{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	}

	// postgres variant

	if config.Database == "postgres" {
		sqlStatement := "UPDATE blogs SET title=$2, content=$3 WHERE id=$1"
		result, err := config.SqlDB.Exec(sqlStatement, id, newBlog.Title, newBlog.Content)
		if err != nil {
			return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v, err: %v", id, newBlog, err))
		}
		rows, err := result.RowsAffected()
		if err != nil || rows == 0 {
			return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v, no matching rows", id, newBlog))
		}

		return &blogProto.UpdateBlogResponse{
			Id:      id,
			Title:   newBlog.Title,
			Content: newBlog.Content,
		}, nil
	}

	// mongo variant

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

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
	id := req.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}
	filter := bson.D{{Key: "_id", Value: oid}}

	result, delete_err := config.Collection.DeleteOne(context.TODO(), filter)
	if delete_err != nil || result.DeletedCount != 1 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Unable to delete blog with ID: %v", id))
	}

	return &blogProto.DeleteBlogResponse{
		Id: id,
	}, nil
}

func (*Server) ListBlog(ctx context.Context, req *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	filter := bson.D{}
	limit := int64(25)

	if req.GetLimit() > 0 {
		limit = req.GetLimit()
	}

	options := &options.FindOptions{
		Limit: &limit,
	}

	var results []BlogItem

	cursor, find_err := config.Collection.Find(context.TODO(), filter, options)
	if find_err != nil {
		if find_err == mongo.ErrNoDocuments {
			return nil, twirp.NewError(twirp.NotFound, "No documents were found")
		}
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error listing blogs: %v", find_err))
	}

	if find_err := cursor.All(context.TODO(), &results); find_err != nil {
		fmt.Printf("error: %v", find_err)
	}

	blogs := []*blogProto.CreateBlogResponse{}

	for _, result := range results {
		blog := blogProto.CreateBlogResponse{
			Id:      result.Id.Hex(),
			Title:   result.Title,
			Content: result.Content,
		}
		blogs = append(blogs, &blog)
	}

	return &blogProto.ListBlogResponse{
		Blogs: blogs,
	}, nil
}
