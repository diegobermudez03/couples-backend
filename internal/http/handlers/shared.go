package handlers

import "net/http"

const limitFilter = "limit"
const pageFilter = "page"
const offsetFilter = "offset"

func getLimitPageOffset(r *http.Request)(string, string) {
	limit := r.URL.Query().Get(limitFilter)
	page := r.URL.Query().Get(pageFilter)
	return limit, page
}