package main

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
        "fmt"
)

type Config struct {
    BaseDir  string `toml:"base_dir"`
    ConfDir   string  `toml:"conf_dir"`
    RestServer RestServerConfig `toml:"rest_server"`
    DataBase  DataBaseConfig `toml:"database"`
    Packer    PackerConfig `toml:"packer"`
    Loging      LogingConfig    `toml:"loging"`
}

type RestServerConfig struct {
    Protocal  string `toml:"protocol"`
    Port      int `toml:"port"`
    FileApiPort int `toml:"file_api_port"`
    Streaming_TimeOut int `toml:"streaming_timeout"`
}

type DataBaseConfig struct {
    Host string "host"
    Port int     "port"
}

type PackerConfig struct {
    CacheDir  string "cache_dir"
    PackerDir   string "packer_dir"
    BuildDir    string  "build_dir"
}


type LogingConfig struct {
	File      string "file"
        Level     int    "level"
}

func ConfigFromFile(configPath string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

        config :=Config{}
	if _,err := toml.Decode(string(configBytes), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main(){
    config,err := ConfigFromFile("../conf/config.toml")
    if err != nil {
        panic (err.Error())
    }
    fmt.Println(config.Loging.File) 
}
