package connection

import (
	"errors"
	"testing"

	"github.com/AA55hex/golang_bootcamp/server/config"
	"github.com/stretchr/testify/require"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

func init() {
	config.LoadConfigs("../../configs.env")
	db_settings = &mysql.ConnectionURL{
		Database: config.MySQL.Database,
		Host:     config.MySQL.Host,
		User:     config.MySQL.User,
		Password: config.MySQL.Password,
	}
}

type mocked_adapter struct {
	isFailTest bool
}

var fail_error = errors.New("")

func (m *mocked_adapter) Open(connURL db.ConnectionURL) (db.Session, error) {
	if m.isFailTest {
		return nil, fail_error
	}
	sess, err := mysql.Open(connURL)
	return sess, err
}

var db_settings *mysql.ConnectionURL

func TestOpenSessionOnSuccess(t *testing.T) {
	adp := new(mocked_adapter)
	adapter = adp

	sess, err := OpenSession(db_settings, 1)
	defer func() {
		if session != nil {
			sess.Close()
			session = nil
		}
	}()
	require.NoError(t, err, "Unexpected error: %v", err)
	require.NotNil(t, sess, "Unexpected nil in session variable.")
}

func TestOpenSessionOnFail(t *testing.T) {
	adp := &mocked_adapter{isFailTest: true}
	adapter = adp

	sess, err := OpenSession(db_settings, 2)
	defer func() {
		if session != nil {
			sess.Close()
			session = nil
		}
	}()
	require.Error(t, err, "Expected error.")
	require.Nil(t, sess, "Unexpected data in session variable.")
}
