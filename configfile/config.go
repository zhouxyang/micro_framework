package configfile

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	Host     string
	Port     int
	ServerID int
	LogPath  string
	EtcdHost string
	EtcdPort int
}

func GetConfig(filename string) (*Config, error) {
	conf := &Config{}
	if _, err := toml.DecodeFile(filename, conf); err != nil {
		werr := errors.Wrap(err, "ParseConfig error")
		return nil, werr
	}
	return conf, nil

}
