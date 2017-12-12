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
#### 1. Defines **default** values
**config** library supports defining a default value for structure members by using **default** keyword in structure tags
```golang
  type Log struct {
    Path  string `default:"/var/logs"`
    Level string `default:"debug"`
  }
```

After that, calls ```config.ParseDefault(interface{})``` to set it on structure instance, example codes as the below:
```golang
  logConfig := Log{}
  config.ParseDefault(&logConfig)
```

#### 2. Reads configurations from JSON file
Like analyzing JSON object from file, reading configurations from JSON file is also simple:
```golang
  type Database struct {
    Host string     `json:"host"`
    Port int        `json:"port"`
    Username string `json:"username" default:"admin"`
    Password string `json:"password" default:"admin"`
    Log Log         `json:"log"`
  }
```

Parsing configurations from JSON:
```golang
  dbConfig := Database{}
  config.ParseConfigFile(&dbConfig, "config.json")
```
