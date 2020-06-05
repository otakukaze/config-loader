# Config Loader

## Install
```
go get -u git.trj.tw/golang/config-loader
```

## Usage

### struct tag
> default : default value define
> length : use for slice type min default length

### support type
- string
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- bool
- struct
- slice
  - []string
  - []int, []int8, []int16, []int32, []int64
  - []uint, []uint8, []uint16, []uint32, []uint64
  - []float32, []float64
  - []bool
  - []struct

### example
```go
package main

import (
  "fmt"
  "log"

  confloader "git.trj.tw/golang/config-loader"
)

type Server struct {
  Port int `json:"port" default:"10230"`
}
type Config struct {
  Server     Server
  Switchs    []int  `json:"switchs" default:"1" length:"2"`
  URL        string `json:"url" default:"https://google.com"`
  OutputFile bool   `json:"output_file" default:"true"`
}

func main() {
  conf := &Config{}
  if err := confloader.Load(conf, nil); err != nil {
    log.Fatal(err)
  }

  fmt.Println(conf)
}
```
