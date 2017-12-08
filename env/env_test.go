package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type EnvConfig1 struct {
	Hostname string   `env:"CONFIG_TEST_HOSTNAME" default:"localhost"`
	Port     int      `env:"CONFIG_TEST_PORT"`
	User     string   `env:"CONFIG_TEST_USER"`
	Password string   `env:"CONFIG_TEST_PASSWORD"`
	Path1    []string `env:"CONFIG_TEST_PATH1"`
	Path2    []string `env:"CONFIG_TEST_PATH2" separator:";"`
	Home     string
}

const (
	TEST_HOSTNAME = "test-hostname"
	TEST_PORT     = 8080
	TEST_USER     = "test-user"
	TEST_PASSWORD = "test-password"
	TEST_PATH1    = "/usr:/var:/bin"
	TEST_PATH2    = "/root;/home;/tmp"
)

func setEnvConfig1() {
	os.Setenv("CONFIG_TEST_HOSTNAME", TEST_HOSTNAME)
	os.Setenv("CONFIG_TEST_PORT", strconv.Itoa(TEST_PORT))
	os.Setenv("CONFIG_TEST_USER", TEST_USER)
	os.Setenv("CONFIG_TEST_PASSWORD", TEST_PASSWORD)
	os.Setenv("CONFIG_TEST_PATH1", TEST_PATH1)
	os.Setenv("CONFIG_TEST_PATH2", TEST_PATH2)
}

func unsetEnvConfig1() {
	os.Unsetenv("CONFIG_TEST_HOSTNAME")
	os.Unsetenv("CONFIG_TEST_PORT")
	os.Unsetenv("CONFIG_TEST_USER")
	os.Unsetenv("CONFIG_TEST_PASSWORD")
}

func assertEqual(expected, actual []string) (bool, error) {
	if len(expected) != len(actual) {
		return false, fmt.Errorf("Expected length of array is %d, but actual is %d",
			len(expected), len(actual))
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != actual[i] {
			return false, fmt.Errorf("Expected array[%d]=%s, but acutal array[%d]=%s",
				i, expected[i], i, actual[i])
		}
	}

	return true, nil
}

func TestEnvConfig1(t *testing.T) {
	setEnvConfig1()
	defer unsetEnvConfig1()

	conf := EnvConfig1{}
	err := Unmarshal(&conf)
	if err != nil {
		t.Errorf("Can't unmarshal config1 from environemnt variables. %s",
			err.Error())
		return
	}

	if conf.Hostname != TEST_HOSTNAME {
		t.Errorf("Expect Hostname: %s, but got: %s", TEST_HOSTNAME, conf.Hostname)
	}

	if conf.Port != TEST_PORT {
		t.Errorf("Expect Port: %d, but got: %d", TEST_PORT, conf.Port)
	}

	if conf.User != TEST_USER {
		t.Errorf("Expect User: %s, but got: %s", TEST_USER, conf.User)
	}

	if conf.Password != TEST_PASSWORD {
		t.Errorf("Expect Password: %s, but got: %s", TEST_PASSWORD, conf.Password)
	}

	if conf.Home != "" {
		t.Errorf("Expect Home is empty, but got: %s", conf.Home)
	}

	expectedPath1 := strings.Split(TEST_PATH1, ":")
	if ok, err := assertEqual(expectedPath1, conf.Path1); !ok {
		t.Error(err.Error())
	}

	expectedPath2 := strings.Split(TEST_PATH2, ";")
	if ok, err := assertEqual(expectedPath2, conf.Path2); !ok {
		t.Error(err.Error())
	}
}

func TestEnvConfig1WithDefaultValue(t *testing.T) {
	os.Setenv("CONFIG_TEST_PORT", strconv.Itoa(TEST_PORT))
	os.Setenv("CONFIG_TEST_USER", TEST_USER)
	os.Setenv("CONFIG_TEST_PASSWORD", TEST_PASSWORD)
	defer unsetEnvConfig1()

	conf := EnvConfig1{}
	err := Unmarshal(&conf)
	if err != nil {
		t.Errorf("Can't unmarshal config1 from environemnt variables. %s",
			err.Error())
		return
	}

	if conf.Hostname != "localhost" {
		t.Errorf("Expect Hostname: localhost, bug got: %s", conf.Hostname)
	}

	if conf.Port != TEST_PORT {
		t.Errorf("Expect Port: %d, but got: %d", TEST_PORT, conf.Port)
	}

	if conf.User != TEST_USER {
		t.Errorf("Expect User: %s, but got: %s", TEST_USER, conf.User)
	}

	if conf.Password != TEST_PASSWORD {
		t.Errorf("Expect Password: %s, but got: %s", TEST_PASSWORD, conf.Password)
	}

	if conf.Home != "" {
		t.Errorf("Expect Home is empty, but got: %s", conf.Home)
	}
}

type EnvConfig2 struct {
	Config1 EnvConfig1
	Server  string `env:"CONFIG_ENV_TEST_SERVER"`
}

const (
	TEST_SERVER = "test-server"
)

func setEnvConfig2() {
	setEnvConfig1()
	os.Setenv("CONFIG_ENV_TEST_SERVER", TEST_SERVER)
}

func unsetEnvConfig2() {
	unsetEnvConfig1()
	os.Unsetenv("CONFIG_ENV_TEST_SERVER")
}

func TestEnvConfig2(t *testing.T) {
	setEnvConfig2()
	defer unsetEnvConfig2()

	conf := EnvConfig2{}
	err := Unmarshal(&conf)
	if err != nil {
		t.Errorf("Can't unmarshal config2 from environemnt variables. %s",
			err.Error())
		return
	}

	if conf.Config1.Hostname != TEST_HOSTNAME {
		t.Errorf("Expect Hostname: %s, but got: %s", TEST_HOSTNAME, conf.Config1.Hostname)
	}

	if conf.Config1.Port != TEST_PORT {
		t.Errorf("Expect Port: %d, but got: %d", TEST_PORT, conf.Config1.Port)
	}

	if conf.Config1.User != TEST_USER {
		t.Errorf("Expect User: %s, but got: %s", TEST_USER, conf.Config1.User)
	}

	if conf.Config1.Password != TEST_PASSWORD {
		t.Errorf("Expect Password: %s, but got: %s", TEST_PASSWORD, conf.Config1.Password)
	}

	if conf.Config1.Home != "" {
		t.Errorf("Expect Home is empty, but got: %s", conf.Config1.Home)
	}

	if conf.Server != TEST_SERVER {
		t.Errorf("Expect Server: %s, but got: %s", TEST_SERVER, conf.Server)
	}
}
