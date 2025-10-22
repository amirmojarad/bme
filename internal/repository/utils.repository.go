package repository

import (
	"bme/internal/service"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"math"
)

const (
	MaxPaginationSize     = 50
	DefaultPaginationSize = 20
)

func appendAppendStartsWith(query *gorm.DB, columnName, value string, toLowerCase bool) {
	if toLowerCase {
		columnName = fmt.Sprintf("lower(%s)", columnName)
	}

	query.Where(fmt.Sprintf(`%s like '%s%%'`, columnName, value))
}

func Paginate(paginator *service.PaginationMeta, cl clause.Interface) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := paginator.CurrentPage
		pageSize := paginator.PerPage

		switch {
		case pageSize > MaxPaginationSize:
			pageSize = MaxPaginationSize
		case pageSize <= 0:
			pageSize = DefaultPaginationSize
		}

		if page <= 0 {
			page = 1
		}

		offset := (page - 1) * pageSize
		if offset < 0 {
			offset = 0
		}

		countSession := db.Session(&gorm.Session{Initialized: true}).Model(db.Statement.Model)

		if cl != nil {
			countSession.Statement.AddClause(cl)
		}

		var count int64
		_ = countSession.Count(&count)

		paginator.Total = int(count)
		paginator.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))

		paginator.CurrentPage = page
		paginator.PerPage = pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

func WithFullQ(q *gorm.DB) clause.Interface {
	if ex, ok := q.Statement.Clauses["WHERE"]; ok {
		expression, ok := ex.Expression.(clause.Interface)
		if ok {
			return expression
		}
	}

	return nil
}
