package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Port = 5050
var Collection *mongo.Collection

// docs for reference:
// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.7.2/mongo
// https://docs.mongodb.com/drivers/go/current/quick-start/

func ConnectToDB() {
	fmt.Println("Connecting to MongoDB")

	// returned cancel function will cancel the created ctx and all associated resources, so ensures cleanup once db operations complete
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// mongo.Connect will create a new client and enable access to the MongoDB instance running on 27107
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database("mydb").Collection("blog")
}
