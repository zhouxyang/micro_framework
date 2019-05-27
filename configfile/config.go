package configfile

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config 全局配置文件
type Config struct {
	Host       string
	Port       int
	ServerName string
	ServerID   int
	LogPath    string
	EtcdHost   string
	EtcdPort   int

	JSONDBFile string
	MyService  struct {
		ExternalServiceName string
	}
}

// GetConfig 解析并获取配置文件
func GetConfig(filename string) (*Config, error) {
	conf := &Config{}
	if _, err := toml.DecodeFile(filename, conf); err != nil {
		werr := errors.Wrap(err, "ParseConfig error")
		return nil, werr
	}
	return conf, nil

}
