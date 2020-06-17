# GORM MySQL Driver

## Quick Start

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

// https://github.com/go-sql-driver/mysql
dsn := "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

## Configuration

```go
import (
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
  DisableDatetimePrecision: true, // disable datetime precision support (added since mysql 5.6)
  DefaultStringSize: 256, // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
}), &gorm.Config{})
```

Checkout [https://gorm.io](https://gorm.io) for details.
