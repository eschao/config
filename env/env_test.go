/*
 * Copyright (C) 2017 eschao <esc.chao@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package env

import (
	"os"
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

func TestLoginConfigEnv(t *testing.T) {
	os.Setenv("USER", LOGIN_USER)
	os.Setenv("PASSWORD", LOGIN_PASSWORD)
	defer os.Unsetenv("USER")
	defer os.Unsetenv("PASSWORD")

	assert := assert.New(t)
	loginConfig := test.LoginConfig{}
	assert.NoError(Parse(&loginConfig))

	assert.Equal(LOGIN_USER, loginConfig.User)
	assert.Equal(LOGIN_PASSWORD, loginConfig.Password)
}

func TestLoginConfigEnvWithPrefix(t *testing.T) {
	os.Setenv("DB_USER", LOGIN_USER)
	os.Setenv("DB_PASSWORD", LOGIN_PASSWORD)
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")

	assert := assert.New(t)
	loginConfig := test.LoginConfig{}
	assert.NoError(ParseWith(&loginConfig, "DB_"))
	assert.Equal(LOGIN_USER, loginConfig.User)
	assert.Equal(LOGIN_PASSWORD, loginConfig.Password)
}

func TestServiceConfigEnv(t *testing.T) {
	servicePrefix := "CONFIG_TEST_SERVICE_"
	serviceLogPrefix := servicePrefix + "LOG_"
	dbPrefix := servicePrefix + "DB_"
	dbLogPrefix := dbPrefix + "LOG_"

	os.Setenv(servicePrefix+"HOST", SERVICE_HOST)
	os.Setenv(servicePrefix+"PORT", strconv.Itoa(SERVICE_PORT))
	os.Setenv(serviceLogPrefix+"PATH", SERVICE_LOG_PATH)
	os.Setenv(serviceLogPrefix+"LEVEL", SERVICE_LOG_LEVEL)
	os.Setenv(dbPrefix+"HOST", DB_HOST)
	os.Setenv(dbPrefix+"PORT", strconv.Itoa(DB_PORT))
	os.Setenv(dbPrefix+"USER", DB_USER)
	os.Setenv(dbPrefix+"PASSWORD", DB_PASSWORD)
	os.Setenv(dbLogPrefix+"PATH", DB_LOG_PATH)
	os.Setenv(dbLogPrefix+"LEVEL", DB_LOG_LEVEL)

	defer os.Unsetenv(servicePrefix + "HOST")
	defer os.Unsetenv(servicePrefix + "PORT")
	defer os.Unsetenv(serviceLogPrefix + "PATH")
	defer os.Unsetenv(serviceLogPrefix + "LEVEL")
	defer os.Unsetenv(dbPrefix + "HOST")
	defer os.Unsetenv(dbPrefix + "PORT")
	defer os.Unsetenv(dbPrefix + "USER")
	defer os.Unsetenv(dbPrefix + "PASSWORD")
	defer os.Unsetenv(dbLogPrefix + "PATH")
	defer os.Unsetenv(dbLogPrefix + "LEVEL")

	assert := assert.New(t)
	serviceConfig := test.ServiceConfig{}
	assert.NoError(Parse(&serviceConfig))
	assert.Equal(SERVICE_HOST, serviceConfig.Host)
	assert.Equal(SERVICE_PORT, serviceConfig.Port)
	assert.Equal(SERVICE_LOG_PATH, serviceConfig.Log.Path)
	assert.Equal(SERVICE_LOG_LEVEL, serviceConfig.Log.Level)
	assert.Equal(DB_HOST, serviceConfig.DBConfig.Host)
	assert.Equal(DB_PORT, serviceConfig.DBConfig.Port)
	assert.Equal(DB_USER, serviceConfig.DBConfig.User)
	assert.Equal(DB_PASSWORD, serviceConfig.DBConfig.Password)
	assert.Equal(DB_LOG_PATH, serviceConfig.DBConfig.Log.Path)
	assert.Equal(DB_LOG_LEVEL, serviceConfig.DBConfig.Log.Level)
}

func TestServiceLoginConfigEnv(t *testing.T) {
	serviceLoginPrefix := "CONFIG_TEST_SERVICE_LOGIN_"
	os.Setenv(serviceLoginPrefix+"USER", LOGIN_USER)
	os.Setenv(serviceLoginPrefix+"PASSWORD", LOGIN_PASSWORD)
	defer os.Unsetenv(serviceLoginPrefix + "USER")
	defer os.Unsetenv(serviceLoginPrefix + "PASSWORD")

	assert := assert.New(t)
	serviceConfig := test.ServiceConfig{Login: &test.LoginConfig{}}
	assert.NoError(Parse(&serviceConfig))
	assert.Equal(LOGIN_USER, serviceConfig.Login.User)
	assert.Equal(LOGIN_PASSWORD, serviceConfig.Login.Password)
}

func TestTypesConfigEnv(t *testing.T) {
	typesPrefix := "CONFIG_TEST_"
	os.Setenv(typesPrefix+"BOOL", "true")
	os.Setenv(typesPrefix+"STR", "test-string")
	os.Setenv(typesPrefix+"INT8", "100")
	os.Setenv(typesPrefix+"INT16", "1000")
	os.Setenv(typesPrefix+"INT", "10000")
	os.Setenv(typesPrefix+"INT32", "100000")
	os.Setenv(typesPrefix+"INT64", "1000000")
	os.Setenv(typesPrefix+"UINT8", "200")
	os.Setenv(typesPrefix+"UINT16", "2000")
	os.Setenv(typesPrefix+"UINT", "20000")
	os.Setenv(typesPrefix+"UINT32", "200000")
	os.Setenv(typesPrefix+"UINT64", "2000000")
	os.Setenv(typesPrefix+"FLOAT32", "1.234")
	os.Setenv(typesPrefix+"FLOAT64", "2222.33333")

	defer os.Unsetenv(typesPrefix + "BOOL")
	defer os.Unsetenv(typesPrefix + "STR")
	defer os.Unsetenv(typesPrefix + "INT8")
	defer os.Unsetenv(typesPrefix + "INT16")
	defer os.Unsetenv(typesPrefix + "INT")
	defer os.Unsetenv(typesPrefix + "INT32")
	defer os.Unsetenv(typesPrefix + "INT64")
	defer os.Unsetenv(typesPrefix + "UINT8")
	defer os.Unsetenv(typesPrefix + "UINT16")
	defer os.Unsetenv(typesPrefix + "UINT")
	defer os.Unsetenv(typesPrefix + "UINT32")
	defer os.Unsetenv(typesPrefix + "UINT64")
	defer os.Unsetenv(typesPrefix + "FLOAT32")
	defer os.Unsetenv(typesPrefix + "FLOAT64")

	assert := assert.New(t)
	typesConfig := test.TypesConfig{}
	assert.NoError(Parse(&typesConfig))
	assert.Equal(true, typesConfig.BoolValue)
	assert.Equal("test-string", typesConfig.StrValue)
	assert.Equal(int8(100), typesConfig.Int8Value)
	assert.Equal(int16(1000), typesConfig.Int16Value)
	assert.Equal(10000, typesConfig.IntValue)
	assert.Equal(int32(100000), typesConfig.Int32Value)
	assert.Equal(int64(1000000), typesConfig.Int64Value)
	assert.Equal(uint8(200), typesConfig.Uint8Value)
	assert.Equal(uint16(2000), typesConfig.Uint16Value)
	assert.Equal(uint(20000), typesConfig.UintValue)
	assert.Equal(uint32(200000), typesConfig.Uint32Value)
	assert.Equal(uint64(2000000), typesConfig.Uint64Value)
	assert.Equal(float32(1.234), typesConfig.Float32Value)
	assert.Equal(float64(2222.33333), typesConfig.Float64Value)
}

func TestTypesConfigWithErrorEnv(t *testing.T) {
	assert := assert.New(t)
	typesConfig := test.TypesConfig{}
	typesPrefix := "CONFIG_TEST_"
	os.Setenv(typesPrefix+"BOOL", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "BOOL")

	os.Setenv(typesPrefix+"INT8", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "INT8")

	os.Setenv(typesPrefix+"INT16", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "INT16")

	os.Setenv(typesPrefix+"INT", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "INT")

	os.Setenv(typesPrefix+"INT32", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "INT32")

	os.Setenv(typesPrefix+"INT64", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "INT64")

	os.Setenv(typesPrefix+"UINT8", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "UINT8")

	os.Setenv(typesPrefix+"UINT16", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "UINT16")

	os.Setenv(typesPrefix+"UINT", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "UINT")

	os.Setenv(typesPrefix+"UINT32", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "UINT32")

	os.Setenv(typesPrefix+"UINT64", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "UINT64")

	os.Setenv(typesPrefix+"FLOAT32", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "FLOAT32")

	os.Setenv(typesPrefix+"FLOAT64", "xxx")
	assert.Error(Parse(&typesConfig))
	os.Unsetenv(typesPrefix + "FLOAT64")

	defer os.Unsetenv(typesPrefix + "BOOL")
	defer os.Unsetenv(typesPrefix + "INT8")
	defer os.Unsetenv(typesPrefix + "INT16")
	defer os.Unsetenv(typesPrefix + "INT")
	defer os.Unsetenv(typesPrefix + "INT32")
	defer os.Unsetenv(typesPrefix + "INT64")
	defer os.Unsetenv(typesPrefix + "UINT8")
	defer os.Unsetenv(typesPrefix + "UINT16")
	defer os.Unsetenv(typesPrefix + "UINT")
	defer os.Unsetenv(typesPrefix + "UINT32")
	defer os.Unsetenv(typesPrefix + "UINT64")
	defer os.Unsetenv(typesPrefix + "FLOAT32")
	defer os.Unsetenv(typesPrefix + "FLOAT64")
}

func TestSlicesConfigEnv(t *testing.T) {
	prefix := "CONFIG_TEST_SLICES_"
	os.Setenv(prefix+"PATHS", "/var:/usr:/home")
	os.Setenv(prefix+"DEBUG", "/root;/log;/opt")
	os.Setenv(prefix+"VALUES", "1,2,4,5")

	defer os.Unsetenv(prefix + "PATHS")
	defer os.Unsetenv(prefix + "DEBUG")
	defer os.Unsetenv(prefix + "VALUES")

	assert := assert.New(t)
	conf := test.SlicesConfig{}
	assert.NoError(Parse(&conf))
	assert.Equal(3, len(conf.Paths))
	assert.Equal("/var", conf.Paths[0])
	assert.Equal("/usr", conf.Paths[1])
	assert.Equal("/home", conf.Paths[2])
	assert.Equal(3, len(conf.Debugs))
	assert.Equal("/root", conf.Debugs[0])
	assert.Equal("/log", conf.Debugs[1])
	assert.Equal("/opt", conf.Debugs[2])
	assert.Equal(4, len(conf.Values))
	assert.Equal(1, conf.Values[0])
	assert.Equal(2, conf.Values[1])
	assert.Equal(4, conf.Values[2])
	assert.Equal(5, conf.Values[3])
}
