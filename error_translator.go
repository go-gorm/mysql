package mysql

import (
	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

var errCodes = map[uint16]error{
	1062: gorm.ErrDuplicatedKey,
	1452: gorm.ErrForeignKeyViolated,
}

func (dialector Dialector) Translate(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if translatedErr, found := errCodes[mysqlErr.Number]; found {
			return translatedErr
		}
		return mysqlErr
	}

	return err
}
