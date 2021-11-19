gen:
		rm -rf rpc
		mkdir rpc
		protoc ./proto/* --go_out=. --go-grpc_out=. --twirp_out=.

# for more info about the protoc CLI, see docs: https://grpc.io/docs/languages/go/quickstart/

serve:
		go run main.go