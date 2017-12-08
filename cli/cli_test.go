package cli

import (
	"flag"
	"fmt"
	"testing"
)

type DBCli struct {
	Host     string `cli:"hostname database server hostname"`
	Port     string `cli:"port database server port"`
	User     string `cli:"user database username"`
	Password string `cli:"password database user password"`
}

type MyCli struct {
	Server string `cli:"-server server URL"`
	DBCli  DBCli  `cli:"database database information"`
}

func TestMain(m *testing.M) {
}

func TestCli(t *testing.T) {

	hostname := flag.String("hostanme", "127.0.0.1", "hostanme value")
	fmt.Printf("Hostname: %s", hostname)

	cli := DBCli{}
	root := CliFlag{}
	err := Parse(&cli, &root)

	if err != nil {
		t.Errorf("Can't parse cli. %s", err.Error())
	}
	flag.Parse()

}
