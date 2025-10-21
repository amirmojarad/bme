package repository

import (
	"fmt"
	"gorm.io/gorm"
)

func appendAppendStartsWith(query *gorm.DB, columnName, value string, toLowerCase bool) {
	if toLowerCase {
		columnName = fmt.Sprintf("lower(%s)", columnName)
	}

	query.Where(fmt.Sprintf(`%s like '%%%s%%'`, columnName, value))
}
