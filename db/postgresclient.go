package db

import (
	blogProto "blog-service/rpc/blog"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/twitchtv/twirp"

	_ "github.com/lib/pq" // importing so drivers are registered with database/sql package, _ means we will not directly reference this package in code
)

/*
for reference:
https://pkg.go.dev/database/sql
https://www.calhoun.io/using-postgresql-with-go/
*/

type PostgresClient struct{}

var SqlDB *sql.DB

func NewPostgresClient() PostgresClient {
	return PostgresClient{}
}

func (p PostgresClient) Connect() error {
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
		return err
	}

	err = db.Ping() // verifies connection to db, establishes connection if necessary
	if err != nil {
		return err
	}

	SqlDB = db

	fmt.Println("Successfully connected to Postgres")
	return nil
}

func (p PostgresClient) CreateBlog(data *blogProto.CreateBlogRequest) (*blogProto.CreateBlogResponse, error) {

	sqlStatement := "INSERT INTO blogs (title, content) VALUES ($1, $2) RETURNING id"
	id := 0
	err := SqlDB.QueryRow(sqlStatement, data.Title, data.Content).Scan(&id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("There was an error creating a blog: %v", err))
	}

	return &blogProto.CreateBlogResponse{
		Id:      strconv.Itoa(id),
		Title:   data.Title,
		Content: data.Content,
	}, nil

}

func (p PostgresClient) GetBlog(data *blogProto.GetBlogRequest) (*blogProto.GetBlogResponse, error) {
	var title string
	var content string

	sqlStatement := "SELECT title, content FROM blogs WHERE id=$1"

	err := SqlDB.QueryRow(sqlStatement, data.Id).Scan(&title, &content)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("No documents were found for id: %v, err: %v", data.Id, err))
	}

	return &blogProto.GetBlogResponse{
		Id:      data.Id,
		Title:   title,
		Content: content,
	}, nil
}

func (p PostgresClient) UpdateBlog(data *blogProto.UpdateBlogRequest) (*blogProto.UpdateBlogResponse, error) {
	sqlStatement := "UPDATE blogs SET title=$2, content=$3 WHERE id=$1"
	result, err := SqlDB.Exec(sqlStatement, data.Id, data.Title, data.Content)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v, err: %v", data.Id, data, err))
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Blog id: %v could not be updated with %v, no matching rows", data.Id, data))
	}

	return &blogProto.UpdateBlogResponse{
		Id:      data.Id,
		Title:   data.Title,
		Content: data.Content,
	}, nil
}

func (p PostgresClient) DeleteBlog(data *blogProto.DeleteBlogRequest) (*blogProto.DeleteBlogResponse, error) {
	sqlStatement := "DELETE FROM blogs WHERE id=$1"
	result, err := SqlDB.Exec(sqlStatement, data.Id)
	if err != nil {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Unable to delete blog with ID: %v, err: %v", data.Id, err))
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, fmt.Sprintf("Unable to delete blog with ID: %v, no matching rows", data.Id))
	}

	return &blogProto.DeleteBlogResponse{
		Id: data.Id,
	}, nil
}

func (p PostgresClient) ListBlog(data *blogProto.ListBlogRequest) (*blogProto.ListBlogResponse, error) {
	sqlStatement := "SELECT * from blogs LIMIT $1"
	rows, err := SqlDB.Query(sqlStatement, data.Limit)
	if err != nil {
		return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error listing blogs: %v", err))
	}
	defer rows.Close()

	blogs := []*blogProto.CreateBlogResponse{}

	for rows.Next() {
		var id int
		var title string
		var content string
		err := rows.Scan(&id, &content, &title)
		if err != nil {
			return nil, twirp.NewError(twirp.NotFound, fmt.Sprintf("There was an error with blog id: %v, err: %v", id, err))
		}
		blog := blogProto.CreateBlogResponse{
			Id:      strconv.Itoa(id),
			Title:   title,
			Content: content,
		}
		blogs = append(blogs, &blog)
	}

	return &blogProto.ListBlogResponse{
		Blogs: blogs,
	}, nil
}
