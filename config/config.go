package config

import (
	db "blog-service/db"
	blogProto "blog-service/rpc/blog"
	"log"
)

var DB DBClient
var Port = 5050

type DBClient interface {
	Connect() error
	CreateBlog(*blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error)
	GetBlog(*blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error)
	UpdateBlog(*blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error)
	DeleteBlog(*blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error)
	ListBlog(*blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error)
}

func SetDB(dbToUse string) {
	if dbToUse == "postgres" {
		DB = db.NewPostgresClient()
	} else {
		DB = db.NewMongoClient()
	}
	err := DB.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
