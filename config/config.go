package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // importing so drivers are registered with database/sql package, _ means we will not directly reference this package in code

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Port = 5050
var Collection *mongo.Collection

// docs for reference:
// https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.7.2/mongo
// https://docs.mongodb.com/drivers/go/current/quick-start/

func ConnectToMongo() {
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

// for reference:
// https://pkg.go.dev/database/sql
// https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

func ConnectToPostgres() {
	fmt.Println("Connecting to Postgres")

	const (
		host     = "localhost"
		port     = 5432
		user     = "root"
		password = "password"
		dbname   = "root"
	)

	dbConnectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbConnectionString) // does not create connect to db, just validates arguments
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping() // verifies connection to db, establishes connection if necessary
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to Postgres")
}
