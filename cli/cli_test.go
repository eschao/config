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
package cli

import (
	"flag"
	"strconv"
	"testing"

	"github.com/eschao/config/test"
	"github.com/stretchr/testify/assert"
)

func TestServiceCommand(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := test.ServiceConfig{}
	cmd := New("Service")
	err := cmd.Init(&serviceConfig)
	if err != nil {
		t.Errorf("Can't init service command. %s", err.Error())
	}

	// assert service cmd
	assert.NotNil(cmd.FlagSet)
	assert.NotNil(cmd.FlagSet.Lookup("hostname"))
	assert.NotNil(cmd.FlagSet.Lookup("port"))
	assert.Equal(2, len(cmd.SubCommands))
	assert.Nil(cmd.SubCommands["login"])

	// assert database sub cmd
	dbCmd := cmd.SubCommands["database"]
	assert.NotNil(dbCmd, "service cmd should have {database} sub cmd")
	assert.NotNil(dbCmd.FlagSet.Lookup("dbHost"))
	assert.NotNil(dbCmd.FlagSet.Lookup("dbPort"))
	assert.NotNil(dbCmd.FlagSet.Lookup("dbUser"))
	assert.NotNil(dbCmd.FlagSet.Lookup("dbPassword"))

	// assert database log sub cmd
	dbLogCmd := dbCmd.SubCommands["log"]
	assert.NotNil(dbCmd, "database cmd should have {log} sub cmd")
	assert.NotNil(dbLogCmd.FlagSet.Lookup("path"))
	assert.NotNil(dbLogCmd.FlagSet.Lookup("level"))
	assert.Equal(0, len(dbLogCmd.SubCommands))

	// assert log cmd
	logCmd := cmd.SubCommands["log"]
	assert.NotNil(logCmd, "service cmd should have {log} sub cmd")
	assert.NotNil(logCmd.FlagSet.Lookup("path"))
	assert.NotNil(logCmd.FlagSet.Lookup("level"))
}

func TestLoginSubCommand(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := test.ServiceConfig{Login: &test.LoginConfig{}}
	cmd := New("Service")
	assert.NoError(cmd.Init(&serviceConfig))

	// assert login sub command
	loginCmd := cmd.SubCommands["login"]
	assert.NotNil(loginCmd, "service cmd should have {login} sub cmd")
	assert.NotNil(loginCmd.FlagSet.Lookup("user"))
	assert.NotNil(loginCmd.FlagSet.Lookup("password"))
}

func TestLoginCommandWithValues(t *testing.T) {
	assert := assert.New(t)
	loginConfig := test.LoginConfig{}
	cmd := New("Login")
	assert.NoError(cmd.Init(&loginConfig), "Can't init login command")

	username := "test-user"
	password := "test-passwd"
	args := []string{"-user", username, "--password", password}
	assert.NoError(cmd.Parse(args))
	assert.Equal(username, loginConfig.User)
	assert.Equal(password, loginConfig.Password)
}

func TestServiceCommandWithValues(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := test.ServiceConfig{Login: &test.LoginConfig{}}
	cmd := New("Service")
	assert.NoError(cmd.Init(&serviceConfig))

	serviceHost := "service-hostname"
	servicePort := 8080
	serviceLogPath := "service-log-path"
	serviceLogLevel := "service-log-debug"

	dbHost := "database-hostname"
	dbPort := 9080
	dbUser := "database-user"
	dbPassword := "database-passwd"
	dbLogPath := "database-log-path"
	dbLogLevel := "database-log-error"

	loginUser := "login-user"
	loginPassword := "login-passwd"

	serviceArgs := []string{"--hostname", serviceHost, "--port",
		strconv.Itoa(servicePort), "log", "-path", serviceLogPath, "-level",
		serviceLogLevel}
	assert.NoError(cmd.Parse(serviceArgs))
	assert.Equal(serviceHost, serviceConfig.Host)
	assert.Equal(servicePort, serviceConfig.Port)
	assert.Equal(serviceLogPath, serviceConfig.Log.Path)
	assert.Equal(serviceLogLevel, serviceConfig.Log.Level)

	dbCmdArgs := []string{"database", "-dbHost", dbHost, "-dbPort",
		strconv.Itoa(dbPort), "-dbUser", dbUser, "-dbPassword", dbPassword}
	assert.NoError(cmd.Parse(dbCmdArgs))
	assert.Equal(dbHost, serviceConfig.DBConfig.Host)
	assert.Equal(dbPort, serviceConfig.DBConfig.Port)
	assert.Equal(dbUser, serviceConfig.DBConfig.User)
	assert.Equal(dbPassword, serviceConfig.DBConfig.Password)

	loginCmdArgs := []string{"login", "--user", loginUser, "-password",
		loginPassword}
	assert.NoError(cmd.Parse(loginCmdArgs))
	assert.Equal(loginUser, serviceConfig.Login.User)
	assert.Equal(loginPassword, serviceConfig.Login.Password)

	dbLogCmdArgs := []string{"database", "log", "-path", dbLogPath, "-level",
		dbLogLevel}
	assert.NoError(cmd.Parse(dbLogCmdArgs))
	assert.Equal(dbLogPath, serviceConfig.DBConfig.Log.Path)
	assert.Equal(dbLogLevel, serviceConfig.DBConfig.Log.Level)
}

func TestVariousTypeCommand(t *testing.T) {
	assert := assert.New(t)
	typesConfig := test.TypesConfig{}
	cmd := NewWith("Types", flag.ContinueOnError, func(cmd *Command) func() {
		return func() {
		}
	})
	assert.NoError(cmd.Init(&typesConfig))

	// bool value
	assert.NoError(cmd.Parse([]string{"-bool=true"}))
	assert.Equal(true, typesConfig.BoolValue)
	assert.NoError(cmd.Parse([]string{"-bool"}))
	assert.Equal(true, typesConfig.BoolValue)
	assert.Error(cmd.Parse([]string{"-bool=xxx"}))

	// string value
	assert.NoError(cmd.Parse([]string{"-str=xxx"}))
	assert.Equal("xxx", typesConfig.StrValue)
	assert.NoError(cmd.Parse([]string{"-str", "yyy"}))
	assert.Equal("yyy", typesConfig.StrValue)

	// int8 value
	assert.NoError(cmd.Parse([]string{"-int8=100"}))
	assert.Equal(int8(100), typesConfig.Int8Value)
	assert.Error(cmd.Parse([]string{"-int8=xxx"}))

	// int16 value
	assert.NoError(cmd.Parse([]string{"-int16=200"}))
	assert.Equal(int16(200), typesConfig.Int16Value)
	assert.Error(cmd.Parse([]string{"-int16=xxx"}))

	// int value
	assert.NoError(cmd.Parse([]string{"-int=300"}))
	assert.Equal(int(300), typesConfig.IntValue)
	assert.Error(cmd.Parse([]string{"-int=xxx"}))

	// int32 value
	assert.NoError(cmd.Parse([]string{"-int32=400"}))
	assert.Equal(int32(400), typesConfig.Int32Value)
	assert.Error(cmd.Parse([]string{"-int32=xxx"}))

	// int64 value
	assert.NoError(cmd.Parse([]string{"-int64=500"}))
	assert.Equal(int64(500), typesConfig.Int64Value)
	assert.Error(cmd.Parse([]string{"-int64=xxx"}))

	// uint8 value
	assert.NoError(cmd.Parse([]string{"-uint8=10"}))
	assert.Equal(uint8(10), typesConfig.Uint8Value)
	assert.Error(cmd.Parse([]string{"-uint8=-10"}))

	// uint16 value
	assert.NoError(cmd.Parse([]string{"-uint16=1000"}))
	assert.Equal(uint16(1000), typesConfig.Uint16Value)
	assert.Error(cmd.Parse([]string{"-uint16=xxx"}))

	// uint value
	assert.NoError(cmd.Parse([]string{"-uint=2000"}))
	assert.Equal(uint(2000), typesConfig.UintValue)
	assert.Error(cmd.Parse([]string{"-uint=xxx"}))

	// uint32 value
	assert.NoError(cmd.Parse([]string{"-uint32=3000"}))
	assert.Equal(uint32(3000), typesConfig.Uint32Value)
	assert.Error(cmd.Parse([]string{"-uint32=xxx"}))

	// uint64 value
	assert.NoError(cmd.Parse([]string{"-uint64=4000"}))
	assert.Equal(uint64(4000), typesConfig.Uint64Value)
	assert.Error(cmd.Parse([]string{"-uint64=xxx"}))

	// float32 value
	assert.NoError(cmd.Parse([]string{"-float32=1.234"}))
	assert.Equal(float32(1.234), typesConfig.Float32Value)
	assert.Error(cmd.Parse([]string{"-float32=xxx"}))

	// float64 value
	assert.NoError(cmd.Parse([]string{"-float64=2.345"}))
	assert.Equal(float64(2.345), typesConfig.Float64Value)
	assert.Error(cmd.Parse([]string{"-float64=xxx"}))
}

func TestCommandWithSlices(t *testing.T) {
	assert := assert.New(t)
	conf := test.SlicesConfig{}
	cmd := New("Slice")
	assert.NoError(cmd.Init(&conf), "Can't init slice command")

	paths := "/var:/home:/log"
	debugs := "error;info;debug"
	values := "100,200,300"
	args := []string{"-paths", paths, "-debugs", debugs, "-values", values}
	assert.NoError(cmd.Parse(args))

	assert.Equal(3, len(conf.Paths))
	assert.Equal("/var", conf.Paths[0])
	assert.Equal("/home", conf.Paths[1])
	assert.Equal("/log", conf.Paths[2])

	assert.Equal(3, len(conf.Debugs))
	assert.Equal("error", conf.Debugs[0])
	assert.Equal("info", conf.Debugs[1])
	assert.Equal("debug", conf.Debugs[2])

	assert.Equal(3, len(conf.Values))
	assert.Equal(100, conf.Values[0])
	assert.Equal(200, conf.Values[1])
	assert.Equal(300, conf.Values[2])
}
