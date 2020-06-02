# GORM MySQL Dialector

## USAGE

```go
import (
  "gorm.io/dialector/mysql"
  "gorm.io/gorm"
)

// https://github.com/go-sql-driver/mysql
dsn := "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

Checkout [https://gorm.io](https://gorm.io) for details.
