package handler

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaginationParams struct {
	Limit  int    `form:"limit" validate:"min=1,max=100"`
	Sort   string `form:"sort" validate:"omitempty,oneof=name email created_at updated_at"`
	Order  string `form:"order" validate:"omitempty,oneof=asc desc"`
	Page   int    `form:"page" validate:"min=1"`
	Cursor string `form:"cursor"`
	Search string `form:"search"`
	Status string `form:"status" validate:"omitempty,oneof=active inactive deleted"`
	CreatedFrom string `form:"created_from" validate:"omitempty,datetime=2006-01-02"`
	CreatedTo   string `form:"created_to" validate:"omitempty,datetime=2006-01-02"`
}

type PaginationMeta struct {
	Page        int  `json:"page,omitempty"`
	Limit       int  `json:"limit"`
	Total       int  `json:"total,omitempty"`
	TotalPages  int  `json:"total_pages,omitempty"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
	NextCursor  string `json:"next_cursor,omitempty"`
	PrevCursor  string `json:"prev_cursor,omitempty"`
}

type Meta struct {
	Sort           string   `json:"sort,omitempty"`
	Order          string   `json:"order,omitempty"`
	FiltersApplied []string `json:"filters_applied,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{}     `json:"data"`
	Pagination PaginationMeta  `json:"pagination"`
	Meta       Meta            `json:"meta"`
}

type Cursor struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func encodeCursor(c Cursor) string {
	data, _ := json.Marshal(c)
	return base64.URLEncoding.EncodeToString(data)
}

func decodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return nil, nil
	}
	
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	
	var cursor Cursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, err
	}
	
	return &cursor, nil
}

func parsePaginationParams(c *gin.Context) PaginationParams {
	params := PaginationParams{
		Limit: 20,
		Sort:  "created_at",
		Order: "desc",
		Page:  1,
	}

	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 && limit <= 100 {
		params.Limit = limit
	}

	if sort := c.Query("sort"); sort != "" {
		params.Sort = sort
	}

	if order := c.Query("order"); order != "" {
		params.Order = order
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		params.Page = page
	}

	params.Cursor = c.Query("cursor")
	params.Search = c.Query("search")
	params.Status = c.Query("status")
	params.CreatedFrom = c.Query("created_from")
	params.CreatedTo = c.Query("created_to")

	return params
}

func buildPaginatedResponse(data interface{}, total int, params PaginationParams) *PaginatedResponse {
	totalPages := (total + params.Limit - 1) / params.Limit
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	var filtersApplied []string
	if params.Search != "" {
		filtersApplied = append(filtersApplied, "search")
	}
	if params.Status != "" {
		filtersApplied = append(filtersApplied, "status")
	}
	if params.CreatedFrom != "" || params.CreatedTo != "" {
		filtersApplied = append(filtersApplied, "created_date")
	}

	return &PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    hasNext,
			HasPrev:    hasPrev,
		},
		Meta: Meta{
			Sort:           params.Sort,
			Order:          params.Order,
			FiltersApplied: filtersApplied,
		},
	}
}

func buildCursorResponse(data interface{}, hasNext bool, nextCursor, prevCursor string, params PaginationParams) *PaginatedResponse {
	var filtersApplied []string
	if params.Search != "" {
		filtersApplied = append(filtersApplied, "search")
	}
	if params.Status != "" {
		filtersApplied = append(filtersApplied, "status")
	}
	if params.CreatedFrom != "" || params.CreatedTo != "" {
		filtersApplied = append(filtersApplied, "created_date")
	}

	return &PaginatedResponse{
		Data: data,
		Pagination: PaginationMeta{
			Limit:      params.Limit,
			HasNext:    hasNext,
			HasPrev:    prevCursor != "",
			NextCursor: nextCursor,
			PrevCursor: prevCursor,
		},
		Meta: Meta{
			Sort:           params.Sort,
			Order:          params.Order,
			FiltersApplied: filtersApplied,
		},
	}
}