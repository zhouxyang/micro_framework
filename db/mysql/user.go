package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"route_guide/db"
)

// UserDao access to User
type UserDao struct {
	*db.MyDB
}

func (dao *UserDao) InsertUser(log *logrus.Entry, user *db.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `INSERT IGNORE INTO User(userid, name) VALUES (?, ?)`

	_, err := dao.Exec(log, query, user.UserID, user.Name)
	if err != nil {
		return errors.Wrap(err, "insert User")
	}
	return nil
}

func (dao *UserDao) GetUserByUserID(log *logrus.Entry, userid string) (*db.User, error) {
	query := `SELECT userid, name FROM User WHERE userid = ?`
	user := db.User{}
	err := dao.Get(log, &user, query, userid)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "db GetUserByUserID")
	}
	return &user, nil
}
