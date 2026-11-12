package mysql

import (
	"strings"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

// The error codes to map mysql errors to gorm errors, here is the mysql error codes reference https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html.
var errCodes = map[uint16]error{
	1062: gorm.ErrDuplicatedKey,
	1451: gorm.ErrForeignKeyViolated,
	1452: gorm.ErrForeignKeyViolated,
	3819: gorm.ErrCheckConstraintViolated,
}

// The error codes that only apply to MariaDB. MariaDB reuses some numbers MySQL assigns to other errors
// (4025 is ER_INNODB_AUTOEXTEND_SIZE_OUT_OF_RANGE in MySQL), so they are checked only when the server is MariaDB,
// here is the mariadb error codes reference https://mariadb.com/kb/en/mariadb-error-code-reference/.
var mariadbErrCodes = map[uint16]error{
	4025: gorm.ErrCheckConstraintViolated,
}

func (dialector Dialector) Translate(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if translatedErr, found := errCodes[mysqlErr.Number]; found {
			return translatedErr
		}
		if dialector.Config != nil && strings.Contains(dialector.ServerVersion, "MariaDB") {
			if translatedErr, found := mariadbErrCodes[mysqlErr.Number]; found {
				return translatedErr
			}
		}
		return mysqlErr
	}

	return err
}
