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
package test

type DBConfig struct {
	Host     string    `json:"dbHost"     yaml:"dbHost"     env:"HOST"     cli:"dbHost database server hostname"`
	Port     int       `json:"dbPort"     yaml:"dbPort"     env:"PORT"     cli:"dbPort database server port"`
	User     string    `json:"dbUser"     yaml:"dbUser"     env:"USER"     cli:"dbUser database username"`
	Password string    `json:"dbPassword" yaml:"dbPassword" env:"PASSWORD" cli:"dbPassword database user password"`
	Log      LogConfig `json:"log"        yaml:"log"        env:"LOG_"     cli:"log database log configuration"`
}

type LoginConfig struct {
	User     string `json:"user"     yaml:"user"     env:"USER"     prop:"user"     cli:"user login username"`
	Password string `json:"password" yaml:"password" env:"PASSWORD" prop:"password" cli:"password login password"`
}

type LogConfig struct {
	Path  string `json:"path"  yaml:"path"  env:"PATH"  prop:"path"  cli:"path log path"`
	Level string `json:"level" yaml:"level" env:"LEVEL" porp:"level" cli:"level log level {debug|warning|error}"`
}

type ServiceConfig struct {
	Host     string       `env:"CONFIG_TEST_SERVICE_HOST"   cli:"hostname service hostname"`
	Port     int          `env:"CONFIG_TEST_SERVICE_PORT"   cli:"port service port"`
	DBConfig DBConfig     `env:"CONFIG_TEST_SERVICE_DB_"    cli:"database database configuration"`
	Login    *LoginConfig `env:"CONFIG_TEST_SERVICE_LOGIN_" cli:"login login user and password"`
	Log      LogConfig    `env:"CONFIG_TEST_SERVICE_LOG_"   cli:"log service log configuration"`
}

type TypesConfig struct {
	BoolValue    bool    `env:"CONFIG_TEST_BOOL"    cli:"bool boolean value"`
	StrValue     string  `env:"CONFIG_TEST_STR"     cli:"str string value"`
	Int8Value    int8    `env:"CONFIG_TEST_INT8"    cli:"int8 int8 value"`
	Int16Value   int16   `env:"CONFIG_TEST_INT16"   cli:"int16 int16 value"`
	IntValue     int     `env:"CONFIG_TEST_INT"     cli:"int int value"`
	Int32Value   int32   `env:"CONFIG_TEST_INT32"   cli:"int32 int32 value"`
	Int64Value   int64   `env:"CONFIG_TEST_INT64"   cli:"int64 int64 value"`
	Uint8Value   uint8   `env:"CONFIG_TEST_UINT8"   cli:"uint8 uint8 value"`
	Uint16Value  uint16  `env:"CONFIG_TEST_UINT16"  cli:"uint16 uint16 value"`
	UintValue    uint    `env:"CONFIG_TEST_UINT"    cli:"uint uint value"`
	Uint32Value  uint32  `env:"CONFIG_TEST_UINT32"  cli:"uint32 uint32 value"`
	Uint64Value  uint64  `env:"CONFIG_TEST_UINT64"  cli:"uint64 uint64 value"`
	Float32Value float32 `env:"CONFIG_TEST_FLOAT32" cli:"float32 float32 value"`
	Float64Value float64 `env:"CONFIG_TEST_FLOAT64" cli:"float64 float64 value"`
}

type DefValueConfig struct {
	BoolValue    bool     `env:"CONFIG_TEST_BOOL"        cli:"bool boolean value" default:"true"`
	IntValue     int      `env:"CONFIG_TEST_INT"         cli:"int int value" default:"123"`
	Float64Value float64  `env:"CONFIG_TEST_FLOAT64"     cli:"float64 float64 value" default:"123.4567"`
	StrValue     string   `env:"CONFIG_TEST_STR"         cli:"str string value" default:"default-string"`
	SliceValue   []string `env:"CONFIG_TEST_SLICE"       cli:"slice slice values" default:"xx:yy:zz"`
	NoDefValue   string   `env:"CONFIG_TEST_NO_DEFVALUE" cli:"nodefvalue no default value"`
}

type SlicesConfig struct {
	Paths  []string `env:"CONFIG_TEST_SLICES_PATHS"  cli:"paths multiple path"`
	Debugs []string `env:"CONFIG_TEST_SLICES_DEBUG"  cli:"debugs multiple debug" separator:";"`
	Values []int    `env:"CONFIG_TEST_SLICES_VALUES" cli:"values multiple value" separator:","`
}
