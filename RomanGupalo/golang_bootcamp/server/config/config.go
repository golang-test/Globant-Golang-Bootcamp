package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var MySQL struct {
	User               string
	Password           string
	Database           string
	Host               string
	ConnectionTryCount int
	MigrationSource    string
}

var Server struct {
	Address string
}

// setting env variable into string variable
func setEnvString(env_var *string, env_string string, default_string string) {
	var ok bool
	*env_var, ok = os.LookupEnv(env_string)
	if !ok {
		*env_var = default_string
		log.Println(env_string, ":\t\tdefault.")
		return
	}
	log.Println(env_string, ":\t\t", *env_var)
}

// setting env variable into int variable
func setEnvInt(env_var *int, env_string string, default_int int) {
	str, ok := os.LookupEnv(env_string)
	buff, err := strconv.ParseInt(str, 10, 32)
	switch {
	case !ok:
		*env_var = default_int
		log.Println(env_string, ": default.")
		return
	case err != nil:
		*env_var = default_int
		log.Println("Bad ", env_string, ": ", buff, ". Set to default.")
	default:
		*env_var = int(buff)
		log.Println(env_string, ":\t\t", buff)
	}
}

// LoadConfigs load environment variables from file
// to MySQL and Server sttructures
func LoadConfigs(path string) error {
	log.Println("Loading ", path)
	if err := godotenv.Load(path); err != nil {
		return err
	}

	setEnvString(&MySQL.Host, "MYSQL_HOST", "localhost:3306")
	setEnvString(&MySQL.Database, "MYSQL_DATABASE", "")
	setEnvString(&MySQL.User, "MYSQL_USER", "root")
	setEnvString(&MySQL.Password, "MYSQL_PASSWORD", "")
	setEnvString(&MySQL.MigrationSource, "MIGRATION_SOURCE", "file:///migrations/")
	setEnvString(&Server.Address, "SERVER_ADDRESS", "localhost:8000")

	setEnvInt(&MySQL.ConnectionTryCount, "CONNECTION_TRY_COUNT", 50)
	return nil
}
