package mysql

import (
	"encoding/json"
	"gorm.io/gorm"
)

var errCodes = map[string]int{
	"uniqueConstraint": 1062,
}

type ErrMessage struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

func (dialector Dialector) Translate(err error) error {
	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if errMsg.Number == errCodes["uniqueConstraint"] {
		return gorm.ErrDuplicatedKey
	}

	return err
}
