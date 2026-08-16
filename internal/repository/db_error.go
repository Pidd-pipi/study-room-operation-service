package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/go-sql-driver/mysql"
)

// dbOrTx 返回事务连接（若 tx 非空）或默认连接。
func dbOrTx(db, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return db
}

// isDuplicateMySQL 识别 MySQL 唯一键冲突。
func isDuplicateMySQL(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

// isDuplicate 兼容 MySQL 唯一键错误文本。
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	if isDuplicateMySQL(err) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
