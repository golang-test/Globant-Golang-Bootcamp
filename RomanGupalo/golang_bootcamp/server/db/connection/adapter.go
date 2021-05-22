package connection

import (
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

// created for testing
type i_adapter interface {
	Open(connURL db.ConnectionURL) (db.Session, error)
}

type mysql_adapter struct{}

func (mysql_adapter) Open(connURL db.ConnectionURL) (db.Session, error) {
	return mysql.Open(connURL)
}

var adapter i_adapter

func init() {
	adapter = mysql_adapter{}
}
