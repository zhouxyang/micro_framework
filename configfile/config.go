package configfile

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config 全局配置文件
type Config struct {
	Host               string
	Port               int
	ServerName         string
	ServerID           int
	LogPath            string
	EtcdHost           string
	EtcdPort           int
	FluentHost         string
	FluentPort         int
	ZipkinHTTPEndpoint string
	ConcurrencyLimit   int64
	MetricTime         int64

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
	host := os.Getenv("MY_POD_IP")
	if host != "" {
		conf.Host = host
	}
	return conf, nil

}
