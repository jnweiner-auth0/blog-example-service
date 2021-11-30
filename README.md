# blog-example-service
Simple example CRUD backend using Golang, gRPC, and protobufs.

## Initializing a local DB
You will need a MongoDB instance running on the default port with a db called `mydb` and a collection called `blog`.
<br>
Spin up with Docker:
```
$ docker run --name mongo-blog -d -p 27017:27017 mongo
$ docker exec -it mongo-blog mongo
$ use mydb
$ db.createCollection("blog")
```

## Testing
Start the server with `make serve` and then run the testing suite with `make test`.
