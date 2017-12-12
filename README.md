## Introduction
**config** is a simple golang library and designed to read configurations from JSON, Yaml files, environment variables and command line. **config** depends on [go-yaml](https://github.com/go-yaml/yaml) to anlayze Yaml file and use built-in golang library to handle JSON file.

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
#### Defines **default** values
**config** library supports defining a default value for structure members by using **default** keyword in structure tags
```golang
  type Database struct {
    Username string `default:"admin"`
    Password string `default:"admin"`
  }
```

After specified default value in tags, calls ```config.ParseDefault(interface{})``` to set it on structure instance, example codes as the below:
```golang
  dbConfig := Database{}
  config.ParseDefault(&dbConfig)
```
