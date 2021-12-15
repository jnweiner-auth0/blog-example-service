package db

import (
	blogProto "blog-service/rpc/blog"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/twitchtv/twirp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
for reference:
https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.7.2/mongo
https://docs.mongodb.com/drivers/go/current/quick-start/
https://docs.mongodb.com/drivers/go/current/fundamentals/bson/
https://docs.mongodb.com/drivers/go/current/fundamentals/crud/
*/

type MongoClient struct{}

var Collection *mongo.Collection

type BlogItem struct {
	Id      primitive.ObjectID `bson:"_id"`
	Title   string             `bson:"title"`
	Content string             `bson:"content"`
}

func NewMongoClient() MongoClient {
	return MongoClient{}
}

func (m MongoClient) Connect() error {
	fmt.Println("Connecting to MongoDB")

	// returned cancel function will cancel the created ctx and all associated resources, so ensures cleanup once db operations complete
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// mongo.Connect will create a new client and enable access to the MongoDB instance running on 27107
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	Collection = client.Database("mydb").Collection("blog")

	fmt.Println("Successfully connected to MongoDB")
	return nil
}

func (m MongoClient) CreateBlog(data *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {
	res, err := Collection.InsertOne(context.Background(), data)
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

func (m MongoClient) GetBlog(data *blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	result := BlogItem{}

	// Decode() method unmarshals BSON into result
	unmarshal_err := Collection.FindOne(context.TODO(), filter).Decode(&result)
	if unmarshal_err != nil {
		if unmarshal_err == mongo.ErrNoDocuments {
			return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("No documents were found for id: %v", data.Id))
		}
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error finding a blog with ID: %v \nError: %v", data.Id, unmarshal_err))
	}

	return &blogProto.GetBlogResponse{
		Id:      data.Id,
		Title:   result.Title,
		Content: result.Content,
	}, nil
}

func (m MongoClient) UpdateBlog(data *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update := bson.D{{Key: "$set", Value: bson.M{"title": data.Title, "content": data.Content}}}

	result, update_err := Collection.UpdateOne(context.TODO(), filter, update)
	if update_err != nil || result.MatchedCount == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v", data.Id, data))
	}

	return &blogProto.UpdateBlogResponse{
		Id:      oid.Hex(),
		Title:   data.Title,
		Content: data.Content,
	}, nil
}

func (m MongoClient) DeleteBlog(data *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	oid, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, "Invalid blog ID")
	}
	filter := bson.D{{Key: "_id", Value: oid}}

	result, delete_err := Collection.DeleteOne(context.TODO(), filter)
	if delete_err != nil || result.DeletedCount != 1 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Unable to delete blog with ID: %v", data.Id))
	}

	return &blogProto.DeleteBlogResponse{
		Id: data.Id,
	}, nil
}

func (m MongoClient) ListBlog(data *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	filter := bson.D{}

	options := &options.FindOptions{
		Limit: &data.Limit,
	}

	var results []BlogItem

	cursor, find_err := Collection.Find(context.TODO(), filter, options)
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
