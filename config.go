package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
)

type Field struct {
	JsonName     string
	YamlName     string
	PropName     string
	EnvName      string
	CliName      string
	Value        reflect.Value
	DefaultValue string
	Separator    string
}

type Config struct {
	Fields []Field
}

func (this *Config) Init() *Config {
	if this.Fields == nil {
		this.Fields = []Field{}
	}

	return this
}

func (this *Config) ParseJSON(jsonFile string, data interface{}) error {
	raw, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return errors.New("Can't open json file. Err: " + err.Error())
	}

	err = json.Unmarshal(raw, data)
	if err != nil {
		return errors.New("Failed unmarshal json. Err: " + err.Error())
	}

	//fmt.Printf("Data: %v", *data.(*interface{}))

	return nil
}
