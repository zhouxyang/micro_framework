package configfile

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config 全局配置文件
type Config struct {
	Host                   string
	Port                   int
	ServerName             string
	ServerID               int
	LogPath                string
	EtcdHost               string
	EtcdPort               int
	ZipkinHTTPEndpoint     string
	ZipkinRecorderHostPort string

	OrderService struct {
		OrderDB        string
		UserService    string
		ProductService string
		BalanceService string
	}
	UserService struct {
		UserDB string
	}
	ProductService struct {
		ProductDB string
	}
	BalanceService struct {
		BalanceDB string
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
