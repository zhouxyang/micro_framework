package mysql

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"micro_framework/db"
)

// BalanceDao access to Balance
type BalanceDao struct {
	*db.MyDB
}

func (dao *BalanceDao) InsertBalance(log *logrus.Entry, balance *db.Balance) error {
	if balance == nil {
		return errors.New("balance is nil")
	}

	query := `INSERT IGNORE INTO Balance(userid, balance) VALUES (?, ?)`

	_, err := dao.Exec(log, query, balance.UserID, balance.Balance.String())
	if err != nil {
		return errors.Wrap(err, "insert Balance")
	}
	return nil
}

func (dao *BalanceDao) UpdateBalanceByUserID(log *logrus.Entry, userid string, balance decimal.Decimal) error {
	query := `UPDATE Balance SET balance = ? where userid = ?`
	_, err := dao.Exec(log, query, balance.String(), userid)

	switch {
	case err != nil:
		return errors.Wrap(err, "db UpdateBalanceByUserID")
	}
	return nil
}

func (dao *BalanceDao) GetBalanceByUserID(log *logrus.Entry, userid string) (*db.Balance, error) {
	query := `SELECT userid, balance FROM Balance WHERE userid = ?`
	balance := db.Balance{}
	err := dao.Get(log, &balance, query, userid)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, errors.Wrap(err, "db GetBalanceByUserID")
	}
	return &balance, nil

}
