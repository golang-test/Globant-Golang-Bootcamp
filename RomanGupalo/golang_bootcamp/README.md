# Go server
---

* [About project](#about-project)
* [Install](#install)
* [Using](#using)
* [Requests](#requests)

---

## About project

Project consists of the following packages:
* handlers
* config
* connection
* entity
* main

The handler package consists of router handlers using gorilla-mux.
This packet is responsible for all the routing logic.

The config package is used to load and use .env files for server configuration.

The connection package is used for creation session between server and database and running database migrations.

The entity package is used for database interaction. All CRUD requests go through this package.

The main package is used for configure and starting server.

All requirements are met.

The migration scripts are located in the `/server/db/migrations`. They are checked automatically every time the server starts.

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/AA55hex/golang_bootcamp
```
## Using

There is Makefile for using project. Use `make` for help.

### Launch

Use `make go-run` to run the docker-compose in attachment mode.
There you can see server-related log information (after conteiners initialization)

### Testing

Use `make go-test` to test project packages
Tests are available for router handlers, connections, configuration packages.

### Clearing

Use `make clear` to down docker-compose conteiners

## Requests
### Get book by id

Use GET method for `localhost:3000/books/<book_id>` to get json-information about book.
If there is no book with such id It's response `404 Not found`.

CURL request:
```bash
curl -XGET 'localhost:3000/books/<book_id>'
```

### Create book

Use POST method for `localhost:3000/books/new` with json body to create book.
Request can pass only if the `Content-Type` header is `application/json`.

CURL request example:
```bash
curl -XPOST -H "Content-type: application/json" -d '{ "name": "book1", "price": 9999, "genre": 1, "amount": 9999 }' 'localhost:3000/books/new'
```

### Update book

Use PUT method for `localhost:3000/books/<book_id>` with json body to update book.
Request can pass only if the `Content-Type` header is `application/json`.

CURL request example:
```bash
curl -XPUT -H "Content-type: application/json" -d '{ "name": "book1", "price": 15, "genre": 2, "amount": 249 }' 'localhost:3000/books/<book_id>'
```

### Delete book by id

Use DELETE method for `localhost:3000/books/<book_id>` to delete book by id.
If there is no book with such id It's response StatusNotFound with error discription in body.

CURL request:
```bash
curl -XDELETE 'localhost:3000/books/<book_id>'
```

### Get books by filter

Use GET method for: 
```http
localhost:3000/books?name=<name>&minPrice=<price>&maxPrice=<price>&genre=<genre_id>
```
You can skip parameters or take it with empty value if you want. 
All requirements are met.

CURL request:
```bash
curl -XGET 'localhost:3000/books?name=<name>&minPrice=<price>&maxPrice=<price>&genre=<genre_id>'
```