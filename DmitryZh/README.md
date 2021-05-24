At first: make some folder in Terminal and make this command:
    git clone https://github.com/dmitriyzhevnov/library.git

============================================================

For Run tests make comands:
    cd library/src/service
    go test

============================================================

For Run Application you must: open the root folder "library" in the terminal and then run the command:
    docker-compose up --build

============================================================

Then go to Postman and make Get request for get all books:
    http://localhost:8080/api/books

============================================================

For get the book by Id enter the Get request:
    http://localhost:8080/api/books/2

============================================================
If you enter in URL invalid Id, for example:
	http://localhost:8080/api/books/2gd
...you will see Error message :
sdfsdfsf
============================================================

For filter by price:
    http://localhost:8080/api/books/price/10/20

============================================================
If you enter in URL only one price parameter, for example:
	http://localhost:8080/api/books/price/10
...you will see Error message :
sdfsdfsf

============================================================

For filter by genre:
    http://localhost:8080/api/books/genre/1

============================================================

For find book by name:
    http://localhost:8080/api/books/name/book3

============================================================

For add new book make Post request at URL:
	http://localhost:8080/api/books
...with raw JSON Body:
{
    "name":"New Book",
    "price": 10,
    "genre": 2,
    "amount": 100
}

============================================================

For edit some book make Put request at URL:
	http://localhost:8080/api/books/2
...with raw JSON Body:
 {
        "name": "EditedBook",
        "price": 50,
        "genre": 1,
        "amount": 50
 }

============================================================

For delete book make Delete request at URL:
	http://localhost:8080/api/books/2	

============================================================
