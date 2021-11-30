package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type BlogResponse struct {
	Id      string `json:"_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var SubjectId string
var BaseURL string = "http://localhost:5050/twirp/service.BlogService"

// note: NOT unit tests, more to demonstrate how to write tests in Golang
func TestServer_Create_And_Delete(t *testing.T) {
	t.Parallel()
	values := map[string]string{"title": "Test title", "content": "Test content"}

	json_data, err := json.Marshal(values)
	if err != nil {
		log.Fatalf("There was an error marshalling JSON: %v", err)
	}

	createRes, err := http.Post(fmt.Sprintf("%s/CreateBlog", BaseURL), "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatalf("There was an error POSTing JSON to localhost:5050: %v", err)
	}

	var createResDecoded BlogResponse

	json.NewDecoder(createRes.Body).Decode(&createResDecoded)
	SubjectId = createResDecoded.Id

	require.NoError(t, err)

	blog_delete := map[string]string{"id": SubjectId}
	json_delete, err := json.Marshal(blog_delete)
	if err != nil {
		log.Fatalf("There was an error marshalling JSON: %v", err)
	}

	deleteRes, err := http.Post(fmt.Sprintf("%s/DeleteBlog", BaseURL), "application/json", bytes.NewBuffer(json_delete))
	if err != nil {
		log.Fatalf("There was an error deleting the blog post: %v", err)
	}

	var deleteResDecoded BlogResponse
	json.NewDecoder(deleteRes.Body).Decode(&deleteResDecoded)

	require.True(t, deleteResDecoded.Id == SubjectId)
	require.NoError(t, err)
}
