## Introduction
**config** is a simple golang library and designed to read configurations from JSON, Yaml files, environment variables and command line. **config** depends on [go-yaml](https://github.com/go-yaml/yaml) to anlayze Yaml file and uses built-in golang library to handle JSON file.

## Installation
1. Install [Yaml](https://github.com/go-yaml/yaml) library first:
```
go get gopkg.in/yaml.v2
```

2. Install **config** library:
```
go get github.com/eschao/config
```

## Usage
### I. Define configuration name in structure tags
Like JSON, Yaml, **config** uses tags to define configurations:

| Tag | Example | Function |
|-----|---------|------|
| json | Host string `json:"host"` | Maps `Host` to a JSON field: **host** |
| yaml | Host string `yaml:"host"` | Maps `Host` to a Yaml field: **host** |
| env | Host string `env:"HOST"` | Maps `Host` to a Environment variable: **HOST** |
| cli | Host string `cli:"host database host"` | Maps `Host` to a command line argument: **-host** or **--host** |
| default | Port int `default:"8080"` | Defines the port with default value: **8080** |
| separator | Path string `json:"path" separator:";"` | Separator is used to split string to a slice |


#### 1. Data types
 **config** supports the following golang data types:
  * bool
  * string
  * int8, int16, int, int32, int64
  * uint8, uint16, uint, uint32, uint64
  * float32, float64
  * slice type. e.g: []string, []int ...
  
#### 2. Defines **default** values
Using **default** keyword in structure tags to define default value:
```golang
  type Log struct {
    Path  string `default:"/var/logs"`
    Level string `default:"debug"`
  }
```

#### 3. Defines configruation name for JSON
Like parsing JSON object, using **json** keyword to define configuration name:
```golang
  type Database struct {
    Host string     `json:"host"`
    Port int        `json:"port"`
    Username string `json:"username" default:"admin"`
    Password string `json:"password" default:"admin"`
    Log Log         `json:"log"`
  }
```

Corresponding JSON file:
```json
 {
   "host": "test.db.hostname",
   "port": 8080,
   "username": "amdin",
   "password": "admin",
   "log": {
     "path": "/var/logs/db",
     "level": "debug"
   }
 }
 ```

#### 4. Defines configuration name for Yaml
Like parsing Yaml object, using **yaml** keyword to define configuration name
```golang
  type Database struct {
    Host string     `yaml:"host"`
    Port int        `yaml:"port"`
    Username string `yaml:"username" default:"admin"`
    Password string `yaml:"password" default:"admin"`
    Log Log         `yaml:"log"`
  }
```
Corresponding Yaml file:
```yaml
 host: test.db.hostname
 port: 8080
 username: amdin
 password: admin
 log:
   path: /var/logs/db
   level: debug
 ```
 
#### 5. Defines configuration name for Environment variable
Using **env** keyword to define configuration name
```golang
  type Database struct {
    Host string     `env:"DB_HOST"`
    Port int        `env:"DB_PORT"`
    Username string `env:"DB_USER" default:"admin"`
    Password string `env:"DB_PASSWORD" default:"admin"`
    Log Log         `env:"DB_LOG_"`
  }
```

Corresponding Environment variables:
```shell
 export DB_HOST=test.db.hostname
 export DB_PORT=8080
 export DB_USER=admin
 export DB_PASSWORD=admin
 export DB_LOG_PATH=/var/logs/db
 export DB_LOG_LEVEL=debug
```
Since the ```Log``` is a structure and nested in ```Database``` structure, the tag of ```Log``` and tags of its structure members will be combined to be an unique environment variable, for example: ```Path``` will be mapped to environment var: ```DB_LOG_PATH```. But if the ```Log``` has no tag definition, only tags of its structure members will be used, that means the ```Path``` will be mapped to ```PATH```.

#### 6. Defines configuration name for Command line
Using **cli** keyword to define configuration name
```golang
  type Database struct {
    Host string     `cli:"host database host name"`
    Port int        `cli:"port database port"`
    Username string `cli:"username database username" default:"admin"`
    Password string `cli:"password database password" default:"admin"`
    Log Log         `cli:"log database log configurations"`
  }
```
For **cli** definition, the string before the first space is command line argument, the rest string are the command line usage and will be oupputed when printing usage

Corresponding command line:
```shell
  ./main -host test.db.hostname -port 8080 -username admin -password admin log -path /var/logs/db -level debug
```
or
```shell
  ./main -host=test.db.hostname -port=8080 -username=admin -password=admin log -path=/var/logs/db -level=debug
```

#### 7. Defines configuration name as a slice type
Using **separator** to split string as a slice:
```golang
  type Log struct {
    Levels []string `env:"LEVELS" cli:"levels log levels" separator:";"`
  }
```

If the separator is not given, its default is **:**, The separator only works on **env** and **cli** tags
```golang
  logConfig := Log{}
  // export LEVELS=debug;error;info
  config.ParseEnv(&logConfig)
  // logConfig[0] == debug
  // logConfig[1] == error
  // logConfig[2] == info
```

### II. Parses configurations
#### 1. Parses default values
When default values are defined in tags, calls ```config.ParseDefault(interface{})``` to assign them to given structure instance **BEFORE** parsing any other configuration types:
```golang
  logConfig := Log{}
  config.ParseDefault(&logConfig)
```
>Note: Other parsing functions won't set structure instance with default values whatever if the configuration value is provided or not

#### 2. Parses from Environment variables
```golang
  dbConfig := Database{}
  config.ParseEnv(&dbConfig)
```

#### 3. Parses from Command line
```golang
  dbConfig := Database{}
  config.ParseCli(&dbConfig)
```

#### 4. Parses from default configuration files
Calls **ParseConfigFile(interface{}, string)** to parse given configuration file:
```golang
  dbConfig := Database{}
  config.ParseConfigFile(&dbConfig, "config.json")
```

If the configuration file is not given, the default configuration files: **config.json** and **config.yaml** will be located under the same folder with fixed searching order.

The **config.json** will be always first located, if it doesn't exist, then checks **config.yaml**. If all of them are not found, parsing will fail.
```golang
  dbConfig := Database{}
  config.ParseConfigFile(&dbConfig, "")
```

#### 4. Parses from configuration file specified by command line
Calls **ParseConfig(interface{}, string)** to parse the configuration file given by command line. The second parameter is a command line argument which is used to specifiy config file:
```golang
  dbConfig := Database{}
  config.ParseConfig(&dbConfig, "c")
```
Run application like:
```shell
  ./main -c config.json
```
**ParseConfig()** will analyze command line argument and extract **config.json** from argument **-c**

### III. Multi-configurations 
You can define all supported configuration tags in a structure and call corresponding functions in your desired order to parse.

Examples:
```golang
  type Log struct {
    Path   string `json:"path" yaml:"path" env:"PATH" cli:"path log path" default:"/var/logs"`
    Levels string `json:"levels" yaml:"levels" env:"LEVELS" cli:"levels log levels" default:"debug;error"`
  }
  
  type Database struct {
    Host     string `json:"host"   yaml:"host"   env:"DB_HOST"   cli:"host database host name"`
    Port     int    `json:"port"   yaml:"port"   env:"DB_PORT"   cli:"port database port"`
    Username string `json:"user"   yaml" user"   env:"DB_USER"   cli:"username database username" default:"admin"`
    Password string `json:"passwd" yaml:"passwd" env:"DB_PASSWD" cli:"password database password" default:"admin"`
    Log      Log    `json:"log"    yaml:"log"    env:"DB_LOG_"   cli:"log database log configurations"`
  }
```
Then, you can parse them like below:
```golang
 dbConfig := Database{}
 
 // parse default values
	if err := config.ParseDefault(&dbConfig); err != nil {
		// error handling
	}

	// parse configuration file from command line
	err := config.ParseConfig(&dbConfig, "c")
 
 // parse default configurations
 if err != nil {
	  err = config.ParseConfigFile(&dbConfig), "")
	}
 
 // parse environment variables
 if err != nil {
   err = config.ParseEnv(&dbConfig)
 }
 
 // parse command line
 if err != nil {
   err = config.ParseCli(&dbConfig)
 }
```

You don't need call all of them. Invokes parsing function that your need.

