package cli

import (
	"flag"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dbConfig struct {
	Host     string    `cmd:"dbHost database server hostname"`
	Port     int       `cmd:"dbPort database server port"`
	User     string    `cmd:"dbUser database username"`
	Password string    `cmd:"dbPassword database user password"`
	Log      logConfig `cmd:"log database log configuration"`
}

type loginConfig struct {
	User     string `cmd:"user login username"`
	Password string `cmd:"password login password"`
}

type logConfig struct {
	Path  string `cmd:"path log path"`
	Level string `cmd:"level log level {debug|warning|error}"`
}

type serviceConfig struct {
	Host     string       `cmd:"hostname service hostname"`
	Port     int          `cmd:"port service port"`
	DBConfig dbConfig     `cmd:"database database configuration"`
	Login    *loginConfig `cmd:"login login user and password"`
	Log      logConfig    `cmd:"log service log configuration"`
}

type typesConfig struct {
	BoolValue    bool    `cmd:"bool boolean value"`
	StrValue     string  `cmd:"str string value"`
	Int8Value    int8    `cmd:"int8 int8 value"`
	Int16Value   int16   `cmd:"int16 int16 value"`
	IntValue     int     `cmd:"int int value"`
	Int32Value   int32   `cmd:"int32 int32 value"`
	Int64Value   int64   `cmd:"int64 int64 value"`
	Uint8Value   uint8   `cmd:"uint8 uint8 value"`
	Uint16Value  uint16  `cmd:"uint16 uint16 value"`
	UintValue    uint    `cmd:"uint uint value"`
	Uint32Value  uint32  `cmd:"uint32 uint32 value"`
	Uint64Value  uint64  `cmd:"uint64 uint64 value"`
	Float32Value float32 `cmd:"float32 float32 value"`
	Float64Value float64 `cmd:"float64 float64 value"`
}

type defValueConfig struct {
	BoolValue bool `cmd:"bool boolean value" default:"true"`
}

func TestServiceCommand(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := serviceConfig{}
	cmd := New("Service")
	err := cmd.Init(&serviceConfig)
	if err != nil {
		t.Errorf("Can't init service command. %s", err.Error())
	}

	// assert service cmd
	assert.NotNil(cmd.FlagSet)
	assert.NotNil(cmd.FlagSet.Lookup("hostname"),
		"service cmd should have {hostname} parameter")
	assert.NotNil(cmd.FlagSet.Lookup("port"),
		"service cmd should have {port} parameter")
	assert.Equal(2, len(cmd.SubCommands),
		"service cmd should have 2 sub cmds")
	assert.Nil(cmd.SubCommands["login"],
		"service cmd shouldn't have {login} sub cmd")

	// assert database sub cmd
	dbCmd := cmd.SubCommands["database"]
	assert.NotNil(dbCmd, "service cmd should have {database} sub cmd")
	assert.NotNil(dbCmd.FlagSet.Lookup("dbHost"),
		"database cmd should have {dbHost} parameter")
	assert.NotNil(dbCmd.FlagSet.Lookup("dbPort"),
		"database cmd should have {dbPort} parameter")
	assert.NotNil(dbCmd.FlagSet.Lookup("dbUser"),
		"database cmd should have {dbUser} parameter")
	assert.NotNil(dbCmd.FlagSet.Lookup("dbPassword"),
		"database cmd should have {dbPassword} parameter")

	// assert database log sub cmd
	dbLogCmd := dbCmd.SubCommands["log"]
	assert.NotNil(dbCmd, "database cmd should have {log} sub cmd")
	assert.NotNil(dbLogCmd.FlagSet.Lookup("path"),
		"database log cmd should have {path} parameter")
	assert.NotNil(dbLogCmd.FlagSet.Lookup("level"),
		"database log cmd should have {level} parameter")
	assert.Equal(0, len(dbLogCmd.SubCommands),
		"database log cmd shouldn't have sub cmd")

	// assert log cmd
	logCmd := cmd.SubCommands["log"]
	assert.NotNil(logCmd, "service cmd should have {log} sub cmd")
	assert.NotNil(logCmd.FlagSet.Lookup("path"),
		"log cmd should have {path} parameter")
	assert.NotNil(logCmd.FlagSet.Lookup("level"),
		"log cmd should have {level} parameter")
}

func TestLoginSubCommand(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := serviceConfig{Login: &loginConfig{}}
	cmd := New("Service")
	assert.NoError(cmd.Init(&serviceConfig), "Can't init service command")

	// assert login sub command
	loginCmd := cmd.SubCommands["login"]
	assert.NotNil(loginCmd, "service cmd should have {login} sub cmd")
	assert.NotNil(loginCmd.FlagSet.Lookup("user"),
		"login cmd should have {user} parameter")
	assert.NotNil(loginCmd.FlagSet.Lookup("password"),
		"login cmd should have {password} parameter")
}

func TestLoginCommandWithValues(t *testing.T) {
	assert := assert.New(t)
	loginConfig := loginConfig{}
	cmd := New("Login")
	assert.NoError(cmd.Init(&loginConfig), "Can't init login command")

	username := "test-user"
	password := "test-passwd"
	args := []string{"-user", username, "--password", password}
	assert.NoError(cmd.Parse(args), "Can't parse login command")
	assert.Equal(username, loginConfig.User, "Failed to parse login command")
	assert.Equal(password, loginConfig.Password, "Failed to parse login command")
}

func TestServiceCommandWithValues(t *testing.T) {
	assert := assert.New(t)
	serviceConfig := serviceConfig{Login: &loginConfig{}}
	cmd := New("Service")
	assert.NoError(cmd.Init(&serviceConfig), "Can't init service command")

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
	assert.NoError(cmd.Parse(serviceArgs), "Can't parse service command")
	assert.Equal(serviceHost, serviceConfig.Host,
		"Service hostname is not equal")
	assert.Equal(servicePort, serviceConfig.Port,
		"Service port is not equal")
	assert.Equal(serviceLogPath, serviceConfig.Log.Path,
		"Service log path is not equal")
	assert.Equal(serviceLogLevel, serviceConfig.Log.Level,
		"Service log level is not equal")

	dbCmdArgs := []string{"database", "-dbHost", dbHost, "-dbPort",
		strconv.Itoa(dbPort), "-dbUser", dbUser, "-dbPassword", dbPassword}
	assert.NoError(cmd.Parse(dbCmdArgs), "Can't parse service database command")
	assert.Equal(dbHost, serviceConfig.DBConfig.Host,
		"Database hostname is not equal")
	assert.Equal(dbPort, serviceConfig.DBConfig.Port,
		"Database port is not equal")
	assert.Equal(dbUser, serviceConfig.DBConfig.User,
		"Database username is not equal")
	assert.Equal(dbPassword, serviceConfig.DBConfig.Password,
		"Database password is not equal")

	loginCmdArgs := []string{"login", "--user", loginUser, "-password",
		loginPassword}
	assert.NoError(cmd.Parse(loginCmdArgs), "Can't parse service login command")
	assert.Equal(loginUser, serviceConfig.Login.User,
		"Login username is not equal")
	assert.Equal(loginPassword, serviceConfig.Login.Password,
		"Login password is not equal")

	dbLogCmdArgs := []string{"database", "log", "-path", dbLogPath, "-level",
		dbLogLevel}
	assert.NoError(cmd.Parse(dbLogCmdArgs), "Can't parse database log command")
	assert.Equal(dbLogPath, serviceConfig.DBConfig.Log.Path,
		"Database log path is not equal")
	assert.Equal(dbLogLevel, serviceConfig.DBConfig.Log.Level,
		"Database log level is not equal")
}

func TestVariousTypeCommand(t *testing.T) {
	assert := assert.New(t)
	typesConfig := typesConfig{}
	cmd := NewWith("Types", flag.ContinueOnError, nil)
	assert.NoError(cmd.Init(&typesConfig))

	// bool value
	assert.NoError(cmd.Parse([]string{"-bool=true"}),
		"Can't parse bool value command")
	assert.Equal(true, typesConfig.BoolValue, "Bool value is not true")
	assert.NoError(cmd.Parse([]string{"-bool"}),
		"Can't parse bool value command")
	assert.Equal(true, typesConfig.BoolValue, "Bool value is not false")
	assert.Error(cmd.Parse([]string{"-bool=xxx"}),
		"Parsing string as bool should have an error")

	// string value
	assert.NoError(cmd.Parse([]string{"-str=xxx"}),
		"Can't parse string value command")
	assert.Equal("xxx", typesConfig.StrValue, "String value it not equal")
	assert.NoError(cmd.Parse([]string{"-str", "yyy"}),
		"Can't parse string value command")
	assert.Equal("yyy", typesConfig.StrValue, "String value is not equal")

	// int8 value
	assert.NoError(cmd.Parse([]string{"-int8=100"}),
		"Can't parse int8 value command")
	assert.Equal(int8(100), typesConfig.Int8Value, "Int8 value is not equal")
	assert.Error(cmd.Parse([]string{"-int8=xxx"}),
		"Parsing string as int8 should have an error")

	// int16 value
	assert.NoError(cmd.Parse([]string{"-int16=200"}),
		"Can't parse int16 value command")
	assert.Equal(int16(200), typesConfig.Int16Value, "Int16 value is not equal")
	assert.Error(cmd.Parse([]string{"-int16=xxx"}),
		"Parsing string as int16 should have an error")

	// int value
	assert.NoError(cmd.Parse([]string{"-int=300"}),
		"Can't parse int value command")
	assert.Equal(int(300), typesConfig.IntValue, "Int value is not equal")
	assert.Error(cmd.Parse([]string{"-int=xxx"}),
		"Parsing string as int should have an error")

	// int32 value
	assert.NoError(cmd.Parse([]string{"-int32=400"}),
		"Can't parse int32 value command")
	assert.Equal(int32(400), typesConfig.Int32Value, "Int32 value is not equal")
	assert.Error(cmd.Parse([]string{"-int32=xxx"}),
		"Parsing string as int32 should have an error")

	// int64 value
	assert.NoError(cmd.Parse([]string{"-int64=500"}),
		"Can't parse int64 value command")
	assert.Equal(int64(500), typesConfig.Int64Value, "Int64 value is not equal")
	assert.Error(cmd.Parse([]string{"-int64=xxx"}),
		"Parsing string as int64 should have an error")

	// uint8 value
	assert.NoError(cmd.Parse([]string{"-uint8=10"}),
		"Can't parse uint8 value command")
	assert.Equal(uint8(10), typesConfig.Uint8Value, "Uint8 value is not equal")
	assert.Error(cmd.Parse([]string{"-uint8=-10"}),
		"Parsing string as uint8 should have an error")

	// uint16 value
	assert.NoError(cmd.Parse([]string{"-uint16=1000"}),
		"Can't parse uint16 value command")
	assert.Equal(uint16(1000), typesConfig.Uint16Value,
		"Uint16 value is not equal")
	assert.Error(cmd.Parse([]string{"-uint16=xxx"}),
		"Parsing string as uint16 should have an error")

	// uint value
	assert.NoError(cmd.Parse([]string{"-uint=2000"}),
		"Can't parse uint value command")
	assert.Equal(uint(2000), typesConfig.UintValue, "Uint value is not equal")
	assert.Error(cmd.Parse([]string{"-uint=xxx"}),
		"Parsing string as uint should have an error")

	// uint32 value
	assert.NoError(cmd.Parse([]string{"-uint32=3000"}),
		"Can't parse uint32 value command")
	assert.Equal(uint32(3000), typesConfig.Uint32Value,
		"Uint32 value is not equal")
	assert.Error(cmd.Parse([]string{"-uint32=xxx"}),
		"Parsing string as uint32 should have an error")

	// uint64 value
	assert.NoError(cmd.Parse([]string{"-uint64=4000"}),
		"Can't parse uint64 value command")
	assert.Equal(uint64(4000), typesConfig.Uint64Value,
		"Uint64 value is not equal")
	assert.Error(cmd.Parse([]string{"-uint64=xxx"}),
		"Parsing string as uint64 should have an error")

	// float32 value
	assert.NoError(cmd.Parse([]string{"-float32=1.234"}),
		"Can't parse float32 value command")
	assert.Equal(float32(1.234), typesConfig.Float32Value,
		"Float32 value is not equal")
	assert.Error(cmd.Parse([]string{"-float32=xxx"}),
		"Parsing string as float32 should have an error")

	// float64 value
	assert.NoError(cmd.Parse([]string{"-float64=2.345"}),
		"Can't parse float64 value command")
	assert.Equal(float64(2.345), typesConfig.Float64Value,
		"Float64 value is not equal")
	assert.Error(cmd.Parse([]string{"-float64=xxx"}),
		"Parsing string as float64 should have an error")
}
