package config

import (
	"testing"
)

type TestConfig struct {
	Name string `json:"name" default:"test-name"`
	Path string `json:"path" default:"./"`
}

func TestJSONConfig(t *testing.T) {
	config := Config{}
	myConfig := TestConfig{}
	err := config.Init().ParseJSON("test.json", &myConfig)
	if err != nil {
		t.Errorf("JSON config test failed. ", err.Error())
	}

	if myConfig.Name != "jsonconfig" {
		t.Errorf("Name json value: %s != jsonconfig", myConfig.Name)
	}

	if myConfig.Path != "/var" {
		t.Errorf("Path json value: %s != /var", myConfig.Path)
	}
}
