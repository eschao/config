package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/eschao/config/test"
	"github.com/stretchr/testify/assert"
)

const (
	LOGIN_USER        = "test-login-user"
	LOGIN_PASSWORD    = "test-login-passwd"
	SERVICE_HOST      = "test-service-host"
	SERVICE_PORT      = 8080
	SERVICE_LOG_PATH  = "/var/log/service"
	SERVICE_LOG_LEVEL = "debug"
	DB_HOST           = "test-db-host"
	DB_PORT           = 9090
	DB_USER           = "test-db-user"
	DB_PASSWORD       = "test-db-password"
	DB_LOG_PATH       = "/var/log/db"
	DB_LOG_LEVEL      = "error"
)

func TestDefaultValueConfig(t *testing.T) {
	conf := test.DefValueConfig{}
	assert := assert.New(t)
	assert.NoError(ParseDefault(&conf))

	assert.Equal(true, conf.BoolValue)
	assert.Equal(123, conf.IntValue)
	assert.Equal(float64(123.4567), conf.Float64Value)
	assert.Equal("default-string", conf.StrValue)
	assert.Equal(3, len(conf.SliceValue))
	assert.Equal("xx", conf.SliceValue[0])
	assert.Equal("yy", conf.SliceValue[1])
	assert.Equal("zz", conf.SliceValue[2])
	assert.Equal("", conf.NoDefValue)
}

func TestEnvConfig(t *testing.T) {
	dbLogPrefix := "LOG_"

	os.Setenv("HOST", DB_HOST)
	os.Setenv("PORT", strconv.Itoa(DB_PORT))
	os.Setenv("USER", DB_USER)
	os.Setenv("PASSWORD", DB_PASSWORD)
	os.Setenv(dbLogPrefix+"PATH", DB_LOG_PATH)
	os.Setenv(dbLogPrefix+"LEVEL", DB_LOG_LEVEL)

	defer os.Unsetenv("HOST")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("USER")
	defer os.Unsetenv("PASSWORD")
	defer os.Unsetenv(dbLogPrefix + "PATH")
	defer os.Unsetenv(dbLogPrefix + "LEVEL")

	conf := test.DBConfig{}
	assert := assert.New(t)
	assert.NoError(ParseEnv(&conf))
	assert.Equal(DB_HOST, conf.Host)
	assert.Equal(DB_PORT, conf.Port)
	assert.Equal(DB_USER, conf.User)
	assert.Equal(DB_PASSWORD, conf.Password)
	assert.Equal(DB_LOG_PATH, conf.Log.Path)
	assert.Equal(DB_LOG_LEVEL, conf.Log.Level)
}

func TestJSONConfigFile(t *testing.T) {
	_, curTestFile, _, _ := runtime.Caller(0)
	path := filepath.Dir(curTestFile)

	conf := test.DBConfig{}
	assert := assert.New(t)
	assert.NoError(ParseConfigFile(&conf, path+"/test/config.json"))
	assert.Equal(DB_HOST, conf.Host)
	assert.Equal(DB_PORT, conf.Port)
	assert.Equal(DB_USER, conf.User)
	assert.Equal(DB_PASSWORD, conf.Password)
	assert.Equal(DB_LOG_PATH, conf.Log.Path)
	assert.Equal(DB_LOG_LEVEL, conf.Log.Level)
}

func TestYamlConfigFile(t *testing.T) {
	_, curTestFile, _, _ := runtime.Caller(0)
	path := filepath.Dir(curTestFile)

	conf := test.DBConfig{}
	assert := assert.New(t)
	assert.NoError(ParseConfigFile(&conf, path+"/test/config.yaml"))
	assert.Equal(DB_HOST, conf.Host)
	assert.Equal(DB_PORT, conf.Port)
	assert.Equal(DB_USER, conf.User)
	assert.Equal(DB_PASSWORD, conf.Password)
	assert.Equal(DB_LOG_PATH, conf.Log.Path)
	assert.Equal(DB_LOG_LEVEL, conf.Log.Level)
}
