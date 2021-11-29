package main

// import local packages with <module_name>/<package_name>
// in this case, module_name is blog-service (see go.mod)
import (
	config "blog-service/config"
	blogProto "blog-service/rpc/blog"
	"blog-service/server"
	"fmt"
	"net/http"
)

func startServer() {
	fmt.Println("Starting server")

	// assign server variable to the address of the Server struct in the server package
	server := &server.Server{}

	// assign handler variable to the TwirpServer generated by the NewBlogServiceServer function in service.twirp.go
	handler := blogProto.NewBlogServiceServer(server)

	fmt.Printf("Server listening on port: %v\n", config.Port)

	// format the port number to match expected argument format for http.ListenAndServe function
	listener := fmt.Sprintf(":%v", config.Port)

	// see docs for more information on http package: https://pkg.go.dev/net/http#example-ListenAndServe
	http.ListenAndServe(listener, handler)
}

func main() {
	config.ConnectToDB()
	startServer()
}
