package paginator

import (
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

var PerPage int = 50

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		offset := (page - 1) * PerPage
		return db.Offset(offset).Limit(PerPage)
	}
}
