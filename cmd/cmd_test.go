package cmd

import (
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
	err := cmd.Init(&serviceConfig)
	if err != nil {
		t.Errorf("Can't init service command. %s", err.Error())
	}

	// assert login sub command
	loginCmd := cmd.SubCommands["login"]
	assert.NotNil(loginCmd, "service cmd should have {login} sub cmd")
	assert.NotNil(loginCmd.FlagSet.Lookup("user"),
		"login cmd should have {user} parameter")
	assert.NotNil(loginCmd.FlagSet.Lookup("password"),
		"login cmd should have {password} parameter")
}

func TestLoginCommand(t *testing.T) {
	loginConfig := loginConfig{}
	cmd := New("Login")
	if err := cmd.Init(&loginConfig); err != nil {
		t.Errorf("Can't init login command. %s", err.Error())
	}

	args := []string{"-user", "test", "-password", "pass", "log", "database"}
	if err := cmd.FlagSet.Parse(args); err != nil {
		t.Errorf("Can't parse login command. %s", err.Error())
	}

	uknArgs := cmd.FlagSet.Args()
	for i, arg := range uknArgs {
		t.Logf("arg[%d]=%s", i, arg)
	}
}
