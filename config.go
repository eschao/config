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
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/eschao/config/cli"
	"github.com/eschao/config/env"
	"github.com/eschao/config/utils"
	"gopkg.in/yaml.v2"
)

// Default configuration file
const (
	DefaultJSONConfig = "config.json"
	DefaultYamlConfig = "config.yaml"
	DefaultPropConfig = "config.properties"
)

const (
	JSONConfigType = "json"
	YamlConfigType = "yaml"
	PropConfigType = "properties"
)

// ParseDefault parses the given structure, extract default value from its tag
// and set structure with these values.
// Normally, ParseDefault should be called before any other parsing functions
// to set default values for structure.
func ParseDefault(i interface{}) error {
	ptrRef := reflect.ValueOf(i)

	if ptrRef.IsNil() || ptrRef.Kind() != reflect.Ptr {
		return fmt.Errorf("Expect a structure pointer type instead of %s",
			ptrRef.Kind().String())
	}

	valueOfStruct := ptrRef.Elem()
	if valueOfStruct.Kind() != reflect.Struct {
		return fmt.Errorf("Expect a structure pointer type instead of %s",
			valueOfStruct.Kind().String())
	}

	return parseValue(valueOfStruct)
}

func parseValue(v reflect.Value) error {
	typeOfStruct := v.Type()
	var err error
	for i := 0; i < v.NumField() && err == nil; i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)

		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				err = ParseDefault(valueOfField.Interface())
			} else {
				continue
			}
		} else if kindOfField == reflect.Struct {
			err = parseValue(valueOfField)
		}

		defValue, ok := structOfField.Tag.Lookup("default")
		if !ok {
			continue
		}

		kind := valueOfField.Kind()
		switch kind {
		case reflect.Bool:
			err = utils.SetValueWithBool(valueOfField, defValue)
		case reflect.String:
			valueOfField.SetString(defValue)
		case reflect.Int8:
			err = utils.SetValueWithIntX(valueOfField, defValue, 8)
		case reflect.Int16:
			err = utils.SetValueWithIntX(valueOfField, defValue, 16)
		case reflect.Int, reflect.Int32:
			err = utils.SetValueWithIntX(valueOfField, defValue, 32)
		case reflect.Int64:
			err = utils.SetValueWithIntX(valueOfField, defValue, 64)
		case reflect.Uint8:
			err = utils.SetValueWithUintX(valueOfField, defValue, 8)
		case reflect.Uint16:
			err = utils.SetValueWithUintX(valueOfField, defValue, 16)
		case reflect.Uint, reflect.Uint32:
			err = utils.SetValueWithUintX(valueOfField, defValue, 32)
		case reflect.Uint64:
			err = utils.SetValueWithUintX(valueOfField, defValue, 64)
		case reflect.Float32:
			err = utils.SetValueWithFloatX(valueOfField, defValue, 32)
		case reflect.Float64:
			err = utils.SetValueWithFloatX(valueOfField, defValue, 64)
		case reflect.Slice:
			sp, ok := structOfField.Tag.Lookup("separator")
			if !ok {
				sp = ":"
			}
			err = utils.SetValueWithSlice(valueOfField, defValue, sp)

		default:
			return fmt.Errorf("Can't support type: %s", kind.String())
		}
	}

	return err
}

// ParseEnv parses given structure interface and set it with corresponding
// environment values
func ParseEnv(i interface{}) error {
	return env.ParseWith(i, "")
}

// ParseCli parses given structure interface and set it with command line input
func ParseCli(i interface{}) error {
	cli := cli.New(os.Args[0])
	if err := cli.Init(i); err != nil {
		return err
	}
	if err := cli.Parse(os.Args[1:]); err != nil {
		return err
	}
	return nil
}

// ParseConfig parses given structure interface and set it with default
// configuration file.
// configFlag is a command line flag to tell where to locate configure file.
// If the config file doesn't exist, the default config fill will be searched
// under the same folder with the fixed order: config.json, config.yaml and
// config.properties
func ParseConfig(i interface{}, configFlag string) error {
	configFile := flag.String(configFlag, "", "Specifiy configuration file")
	flag.Parse()
	return ParseConfigFile(i, *configFile)
}

// ParseConfigFile parses given structure interface and set its value with
// the specified configuration file
func ParseConfigFile(i interface{}, configFile string) error {
	var err error
	if configFile == "" {
		configFile, err = getDefaultConfigFile()
		if err != nil {
			return err
		}
	}

	configType, err := getConfigFileType(configFile)
	if err != nil {
		return err
	}

	switch configType {
	case JSONConfigType:
		return parseJSON(i, configFile)
	case YamlConfigType:
		return parseYaml(i, configFile)
	case PropConfigType:
		return parseProp(i, configFile)
	default:
		return fmt.Errorf("Can't support config file: %s", configFile)
	}

	return nil
}

// parseJSON parses JSON file and set structure with its value
func parseJSON(i interface{}, jsonFile string) error {
	raw, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("Can't open json config file. %s", err.Error())
	}

	return json.Unmarshal(raw, i)
}

// parseYaml parses Yaml file and set structure with its value
func parseYaml(i interface{}, yamlFile string) error {
	raw, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("Can't open yaml config file. %s", err.Error())
	}

	return yaml.Unmarshal(raw, i)
}

// parseProp parses Properties file and set structure with its value
func parseProp(i interface{}, propFile string) error {
	return fmt.Errorf("Properties config has not implemented!")
}

// getDefaultConfigFile returns a existing default config file. The checking
// order is fixed with beginning from: config.json to config.yaml and
// config.properties
func getDefaultConfigFile() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Can't find default config file. %s", err.Error())
	}

	path := filepath.Dir(exe) + string(filepath.Separator)

	// check json config
	jsonConfig := path + DefaultJSONConfig
	if _, err := os.Stat(jsonConfig); err == nil {
		return jsonConfig, nil
	}

	// check yaml config
	yamlConfig := path + DefaultYamlConfig
	if _, err := os.Stat(yamlConfig); err == nil {
		return yamlConfig, nil
	}

	// check prop config
	propConfig := path + DefaultPropConfig
	if _, err := os.Stat(propConfig); err == nil {
		return propConfig, nil
	}

	return "", fmt.Errorf("No default config file found in path: %s", path)
}

// getConfigFileType analyzes config file extension name and return
// corresponding type: json, yaml or properties
func getConfigFileType(configFile string) (string, error) {
	ext := filepath.Ext(configFile)
	if ext == ".json" {
		return JSONConfigType, nil
	} else if ext == ".yaml" || ext == ".yml" {
		return YamlConfigType, nil
	} else if ext == ".properties" || ext == ".prop" {
		return PropConfigType, nil
	}

	return "", fmt.Errorf("Can't support file type: %s", configFile)
}
