package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	BaseDir    string           `toml:"base_dir"`
	ConfDir    string           `toml:"conf_dir"`
	RestServer RestServerConfig `toml:"rest_server"`
	DataBase   DataBaseConfig   `toml:"database"`
	Packer     PackerConfig     `toml:"packer"`
	Loging     LogingConfig     `toml:"loging"`
	Fetcher    FetcherConfig    `toml:"fetcher"`
}

type RestServerConfig struct {
	Protocal          string `toml:"protocol"`
	Port              int    `toml:"port"`
	FileApiPort       int    `toml:"file_api_port"`
	Streaming_TimeOut int    `toml:"streaming_timeout"`
}

type DataBaseConfig struct {
	Host               string `toml:"host"`
	Port               int    `toml:"port"`
	User               string `toml:"user"`
	Password           string `toml:"password"`
	DataBase           string `toml:"database"`
	URL                string `toml:"url"`
	MaxOpenConnections int    `toml:"maxOpenConnections"`
	MaxIdleConnections int    `toml:"maxIdleConnections"`
}

type PackerConfig struct {
	CacheDir  string `toml:"cache_dir"`
	PackerDir string `toml:"packer_dir"`
	BuildDir  string `toml:"build_dir"`
}

type LogingConfig struct {
	File  string `toml:"file"`
	Level int    `toml:"level"`
}

type FetcherConfig struct {
	Timeout      string `toml:"timeout"`
	CheckRetry   int    `toml:"check_retry"`
	DataDir      string `toml:"data_dir"`
	CurlFormat   string `toml:"curl_format"`
	ScmTarFormat string `toml:"scm_tar_format"`
	ScmMD5Format string `toml:"scm_md5_format"`
	DbTable      string `toml:"db_table"`
}

func ConfigFromFile(configPath string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if _, err := toml.Decode(string(configBytes), &config); err != nil {
		return nil, err
	}

	return &config, nil
}
