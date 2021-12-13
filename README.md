# blog-example-service
Simple example CRUD backend using Golang, gRPC, and protobufs.

## Using with MongoDB as the database (default)
You will need a MongoDB instance running on the default port with a db called `mydb` and a collection called `blog`.
<br>
Spin up with Docker:
```
$ docker run --name mongo-blog -d -p 27017:27017 mongo
$ docker exec -it mongo-blog mongo
$ use mydb
$ db.createCollection("blog")
```
Start the server with `make serve` (will default to Mongo) or `make mongo` (to be explicit).

## Using with Postgres as the database
You will need a Postgres instance running on the default port with a db called `root` and a table called `blogs`.
<br>
Spin up with Docker:
```
$ docker run --name postgres-blog -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres
$ docker exec -it postgres-blog psql
$ create table blogs ( id SERIAL PRIMARY KEY, title TEXT, content TEXT );
```
Start the server with `make postgres`.

## Testing
Start the server (`make mongo` or `make postgres`) and then run the testing suite with `make test`.
