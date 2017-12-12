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
#### I. Data types
Like JSON, Yaml, **config** uses tags to define configurations. It supports the following golang data types:
  * bool
  * string
  * int8, int16, int, int32, int64
  * uint8, uint16, uint, uint32, uint64
  * float32, float64
  * slice type. e.g: []string, []int ...
  
#### II. Defines **default** values
Using **default** keyword in structure tags to define default value. Example codes:
```golang
  type Log struct {
    Path  string `default:"/var/logs"`
    Level string `default:"debug"`
  }
```

After that, calls ```config.ParseDefault(interface{})``` to set structure instance with default values. Example codes:
```golang
  logConfig := Log{}
  config.ParseDefault(&logConfig)
```

#### III. Reads configurations from JSON and Yaml files
##### 1. Using ```json``` to define configuration name
```golang
  type Database struct {
    Host string     `json:"host"`
    Port int        `json:"port"`
    Username string `json:"username" default:"admin"`
    Password string `json:"password" default:"admin"`
    Log Log         `json:"log"`
  }
```

##### 2. Using ```yaml``` to define configuration name
```golang
  type Database struct {
    Host string     `yaml:"host"`
    Port int        `yaml:"port"`
    Username string `yaml:"username" default:"admin"`
    Password string `yaml:"password" default:"admin"`
    Log Log         `yaml:"log"`
  }
```

Parsing configurations from JSON:
```golang
  dbConfig := Database{}
  config.ParseConfigFile(&dbConfig, "config.json")
```
