package configfile

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/url"
	"route_guide/db"
	"strings"
	"time"
)

// InitDB 初始化数据库连接
func InitDB(param string) (*db.MyDB, error) {
	if param == "" {
		return nil, errors.New("param is empty")
	}

	p, err := parseParam(param)
	if err != nil {
		return nil, errors.Wrapf(err, "bad param: %s", param)
	}

	c, err := mysql.ParseDSN(param)
	if err != nil {
		return nil, errors.Wrapf(err, "bad param: %s", param)
	}

	if p.Get("parseTime") == "" {
		c.ParseTime = true
	}

	if p.Get("loc") == "" {
		c.Loc = time.Local
	}

	if p.Get("timeout") == "" {
		c.Timeout = 3 * time.Second
	}

	param = c.FormatDSN()
	d, err := sqlx.Connect("mysql", param)
	if err != nil {
		return nil, errors.Wrap(err, "db Connect()")
	}

	return &db.MyDB{DB: d}, nil
}

func parseParam(param string) (url.Values, error) {
	index := strings.LastIndex(param, "?")
	if index == -1 {
		param = ""
	} else {
		param = param[index+1:]
	}

	m, err := url.ParseQuery(param)
	if err != nil {
		return nil, errors.Wrap(err, "parse param")
	}

	return m, nil
}
