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
	"fmt"
	"os"
	"reflect"

	"github.com/eschao/config/utils"
)

// Parse parses given structure interface, extracts environment definitions
// from its tag and sets structure with defined environement variables
func Parse(i interface{}) error {
	return ParseWith(i, "")
}

// ParseWith parses with given structure interface and environment name prefix
// It is normally used in nested structure.
// For example: we have such structure
// type Database struct {
//   Host string `env:"HOST"`
// }

// type Server struct {
//   Server string `env:"SERVER"`
//   DB     Database `env:"DB_"`
// }
// The Server.DB.Host will be mapped to environment variable: DB_HOST which is
// concatenated from DB tag in Server struct and Host tag in Database struct
func ParseWith(i interface{}, prefix string) error {
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

	return parseValue(valueOfStruct, prefix)
}

// parseValue parses a reflect.Value object
func parseValue(v reflect.Value, prefix string) error {
	typeOfStruct := v.Type()
	var err error
	for i := 0; i < v.NumField() && err == nil; i++ {
		valueOfField := v.Field(i)
		kindOfField := valueOfField.Kind()
		structOfField := typeOfStruct.Field(i)

		// recursively unmarshal if value is ptr type
		if kindOfField == reflect.Ptr {
			if !valueOfField.IsNil() && valueOfField.CanSet() {
				err = ParseWith(valueOfField.Interface(),
					prefix+structOfField.Tag.Get("env"))
			} else {
				continue
			}
		} else if kindOfField == reflect.Struct {
			err = parseValue(valueOfField, prefix+structOfField.Tag.Get("env"))
		}

		err = setFieldValue(valueOfField, structOfField, prefix)
	}

	return err
}

// getEnvValue get environment value
func getEnvValue(envName string, f reflect.StructField) (string, bool) {
	//fmt.Printf("Lookup ENV: %s\n", envName)
	envValue, ok := os.LookupEnv(envName)
	if !ok {
		envValue, ok = f.Tag.Lookup("default")
	}

	return envValue, ok
}

// setFieldValue sets a reflect.Value with environment value
func setFieldValue(v reflect.Value, f reflect.StructField, prefix string) error {
	envName := f.Tag.Get("env")
	if envName == "" {
		return nil
	}

	envValue, ok := os.LookupEnv(prefix + envName)
	if !ok {
		return nil
	}

	if !v.CanSet() {
		return fmt.Errorf("%s: can't be set", f.Name)
	}

	var err error
	kind := v.Kind()
	switch kind {
	case reflect.Bool:
		err = utils.SetValueWithBool(v, envValue)
	case reflect.String:
		v.SetString(envValue)
	case reflect.Int8:
		err = utils.SetValueWithIntX(v, envValue, 8)
	case reflect.Int16:
		err = utils.SetValueWithIntX(v, envValue, 16)
	case reflect.Int, reflect.Int32:
		err = utils.SetValueWithIntX(v, envValue, 32)
	case reflect.Int64:
		err = utils.SetValueWithIntX(v, envValue, 64)
	case reflect.Uint8:
		err = utils.SetValueWithUintX(v, envValue, 8)
	case reflect.Uint16:
		err = utils.SetValueWithUintX(v, envValue, 16)
	case reflect.Uint, reflect.Uint32:
		err = utils.SetValueWithUintX(v, envValue, 32)
	case reflect.Uint64:
		err = utils.SetValueWithUintX(v, envValue, 64)
	case reflect.Float32:
		err = utils.SetValueWithFloatX(v, envValue, 32)
	case reflect.Float64:
		err = utils.SetValueWithFloatX(v, envValue, 64)

	case reflect.Slice:
		sp, ok := f.Tag.Lookup("separator")
		if !ok {
			sp = ":"
		}
		err = utils.SetValueWithSlice(v, envValue, sp)

	default:
		return fmt.Errorf("Can't support type: %s", kind.String())
	}

	if err != nil {
		return fmt.Errorf("%s: %s", f.Name, err.Error())
	}
	return nil
}
