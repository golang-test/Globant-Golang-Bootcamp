package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvStringOnSuccess(t *testing.T) {
	var test_buffer string
	os.Setenv("TESTING_ENV_VARIABLE", "SOME TEXT")
	req := os.Getenv("TESTING_ENV_VARIABLE")

	// Check on success return
	setEnvString(&test_buffer, "TESTING_ENV_VARIABLE", "DEFAULT")
	require.True(t, test_buffer == req,
		"Want `%v` value.\nReturns `%v` value.", req, test_buffer)
}

func TestEnvStringOnDefault(t *testing.T) {
	var test_buffer string
	os.Setenv("TESTING_ENV_VARIABLE", "SOME TEXT")
	req := "DEFAULT"

	// Check on default return
	setEnvString(&test_buffer, "//NOT EXISTS\\", req)
	require.True(t, test_buffer == "DEFAULT",
		"Want `%v` value.\nReturns `%v` value.", req, test_buffer)
}

func TestEnvIntOnSuccess(t *testing.T) {
	var test_buffer int
	os.Setenv("TESTING_ENV_VARIABLE", "404")
	req, _ := strconv.Atoi(os.Getenv("TESTING_ENV_VARIABLE"))

	setEnvInt(&test_buffer, "TESTING_ENV_VARIABLE", 0)
	require.True(t, test_buffer == req,
		"Want `%v` value.\nReturns `%v` value.", req, test_buffer)

}

func TestEnvIntOnDefault(t *testing.T) {
	var test_buffer int
	os.Setenv("TESTING_ENV_VARIABLE", "404")
	req := 0

	setEnvInt(&test_buffer, "//NOT EXISTS\\", req)
	require.True(t, test_buffer == req,
		"Want `%v` value.\nReturns `%v` value.", req, test_buffer)

}

func TestEnvIntOnBadEnv(t *testing.T) {
	var test_buffer int
	os.Setenv("TESTING_ENV_VARIABLE", "NOT INT")
	req := 0

	setEnvInt(&test_buffer, "TESTING_ENV_VARIABLE", req)
	require.True(t, test_buffer == req,
		"Want `%v` value.\nReturns `%v` value.", req, test_buffer)

}

func TestLoadConfigsOnSuccess(t *testing.T) {
	err := LoadConfigs("../configs.env")
	require.NoError(t, err, err)
}

func TestLoadConfigOnBadPath(t *testing.T) {
	err := LoadConfigs("//not_exist\\.env")
	require.Error(t, err, "Expected error")
}
