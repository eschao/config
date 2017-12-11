package test

type DBConfig struct {
	Host     string    `json:"dbHost"     yaml:"dbHost"     env:"HOST"     cmd:"dbHost database server hostname"`
	Port     int       `json:"dbPort"     yaml:"dbPort"     env:"PORT"     cmd:"dbPort database server port"`
	User     string    `json:"dbUser"     yaml:"dbUser"     env:"USER"     cmd:"dbUser database username"`
	Password string    `json:"dbPassword" yaml:"dbPassword" env:"PASSWORD" cmd:"dbPassword database user password"`
	Log      LogConfig `json:"log"        yaml:"log"        env:"LOG_"     cmd:"log database log configuration"`
}

type LoginConfig struct {
	User     string `json:"user"     yaml:"user"     env:"USER"     cmd:"user login username"`
	Password string `json:"password" yaml:"password" env:"PASSWORD" cmd:"password login password"`
}

type LogConfig struct {
	Path  string `json:"path"  yaml:"path"  env:"PATH"  cmd:"path log path"`
	Level string `json:"level" yaml:"level" env:"LEVEL" cmd:"level log level {debug|warning|error}"`
}

type ServiceConfig struct {
	Host     string       `env:"CONFIG_TEST_SERVICE_HOST"   cmd:"hostname service hostname"`
	Port     int          `env:"CONFIG_TEST_SERVICE_PORT"   cmd:"port service port"`
	DBConfig DBConfig     `env:"CONFIG_TEST_SERVICE_DB_"    cmd:"database database configuration"`
	Login    *LoginConfig `env:"CONFIG_TEST_SERVICE_LOGIN_" cmd:"login login user and password"`
	Log      LogConfig    `env:"CONFIG_TEST_SERVICE_LOG_"   cmd:"log service log configuration"`
}

type TypesConfig struct {
	BoolValue    bool    `env:"CONFIG_TEST_BOOL"    cmd:"bool boolean value"`
	StrValue     string  `env:"CONFIG_TEST_STR"     cmd:"str string value"`
	Int8Value    int8    `env:"CONFIG_TEST_INT8"    cmd:"int8 int8 value"`
	Int16Value   int16   `env:"CONFIG_TEST_INT16"   cmd:"int16 int16 value"`
	IntValue     int     `env:"CONFIG_TEST_INT"     cmd:"int int value"`
	Int32Value   int32   `env:"CONFIG_TEST_INT32"   cmd:"int32 int32 value"`
	Int64Value   int64   `env:"CONFIG_TEST_INT64"   cmd:"int64 int64 value"`
	Uint8Value   uint8   `env:"CONFIG_TEST_UINT8"   cmd:"uint8 uint8 value"`
	Uint16Value  uint16  `env:"CONFIG_TEST_UINT16"  cmd:"uint16 uint16 value"`
	UintValue    uint    `env:"CONFIG_TEST_UINT"    cmd:"uint uint value"`
	Uint32Value  uint32  `env:"CONFIG_TEST_UINT32"  cmd:"uint32 uint32 value"`
	Uint64Value  uint64  `env:"CONFIG_TEST_UINT64"  cmd:"uint64 uint64 value"`
	Float32Value float32 `env:"CONFIG_TEST_FLOAT32" cmd:"float32 float32 value"`
	Float64Value float64 `env:"CONFIG_TEST_FLOAT64" cmd:"float64 float64 value"`
}

type DefValueConfig struct {
	BoolValue    bool     `env:"CONFIG_TEST_BOOL"        cmd:"bool boolean value" default:"true"`
	IntValue     int      `env:"CONFIG_TEST_INT"         cmd:"int int value" default:"123"`
	Float64Value float64  `env:"CONFIG_TEST_FLOAT64"     cmd:"float64 float64 value" default:"123.4567"`
	StrValue     string   `env:"CONFIG_TEST_STR"         cmd:"str string value" default:"default-string"`
	SliceValue   []string `env:"CONFIG_TEST_SLICE"       cmd:"slice slice values" default:"xx:yy:zz"`
	NoDefValue   string   `env:"CONFIG_TEST_NO_DEFVALUE" cmd:"nodefvalue no default value"`
}

type SlicesConfig struct {
	Paths  []string `env:"CONFIG_TEST_SLICES_PATHS"`
	Debugs []string `env:"CONFIG_TEST_SLICES_DEBUG" separator:";"`
	Values []int    `env:"CONFIG_TEST_SLICES_VALUES" separator:","`
}
