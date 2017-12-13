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
   "password": "admin"
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
Since the ```Log``` is a structure and nested in ```Database``` structure, the tag of ```Log``` and tags of its structure members will be combined to be an unique environment variable, for example: ```Path``` will be mapped to environment var: ```DB_LOG_PATH```. But if the ```Log``` has no tag defination, only tags of its structure members will be used, that means the ```Path``` will be mapped to ```PATH```.

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
>Note: Other parsing functions won't set structure instance with default values whatever if the configuration is provided or not

#### 2. Parses configuration files:

Parsing configurations from JSON:
```golang
  dbConfig := Database{}
  config.ParseConfigFile(&dbConfig, "config.json")
```
