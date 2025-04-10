package utils

import (
	"backend/internal/models"
	"net/http"
	"strconv"
)

const (
	defaultLimit  = 20
	maxLimit      = 100
	defaultPage   = 1
	defaultSortBy = "name"
	defaultOrder  = "asc"
)

var allowedSortBy = map[string]bool{
	"name":       true,
	"created_at": true,
}

func ParseExerciseFilter(r *http.Request) *models.ExerciseFilter {
	q := r.URL.Query()
	var f models.ExerciseFilter

	if v := q.Get("category_id"); v != "" {
		if id, err := strconv.Atoi(v); err == nil {
			f.CategoryID = &id
		}
	}

	if v := q.Get("search"); v != "" {
		f.Search = &v
	}

	if v := q.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l < maxLimit {
			f.Limit = l
		}
	}

	if f.Limit == 0 {
		f.Limit = defaultLimit
	}

	if v := q.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			f.Page = p
		}
	}

	if f.Page == 0 {
		f.Page = defaultPage
	}

	f.Offset = (f.Page - 1) * f.Limit

	if v := q.Get("sort_by"); allowedSortBy[v] {
		f.SortBy = v
	} else {
		f.SortBy = defaultSortBy
	}

	if v := q.Get("sort_order"); v == "desc" {
		f.SortOrder = v
	} else {
		f.SortOrder = defaultOrder
	}

	return &f

}
