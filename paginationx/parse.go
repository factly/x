package paginationx

import (
	"net/url"
	"strconv"
)

// Parse - pagination
func Parse(query url.Values) (int, int) {

	page, _ := strconv.Atoi(query.Get("page"))     // page number
	perPage, _ := strconv.Atoi(query.Get("limit")) // perPage number

	offset := 0 // no. of records to skip
	limit := 10 // limit

	if perPage > 0 && perPage <= 40 {
		limit = perPage
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return offset, limit

}
