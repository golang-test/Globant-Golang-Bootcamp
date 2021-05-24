module github.com/SavinskiPavel/bookstore

go 1.16

replace github.com/SavinskiPavel/bookstore => ../bookstore

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/stretchr/testify v1.7.0
)
