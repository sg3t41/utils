package handler

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testTime = time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)

func TestParsePaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		query    string
		expected PaginationParams
	}{
		{
			name:  "default values",
			query: "",
			expected: PaginationParams{
				Limit: 20,
				Sort:  "created_at",
				Order: "desc",
				Page:  1,
			},
		},
		{
			name:  "custom pagination",
			query: "limit=10&page=2&sort=name&order=asc",
			expected: PaginationParams{
				Limit: 10,
				Sort:  "name",
				Order: "asc",
				Page:  2,
			},
		},
		{
			name:  "with filters",
			query: "search=john&status=active&created_from=2023-01-01&created_to=2023-12-31",
			expected: PaginationParams{
				Limit:       20,
				Sort:        "created_at",
				Order:       "desc",
				Page:        1,
				Search:      "john",
				Status:      "active",
				CreatedFrom: "2023-01-01",
				CreatedTo:   "2023-12-31",
			},
		},
		{
			name:  "with cursor",
			query: "cursor=eyJpZCI6IjEyMyJ9&limit=5",
			expected: PaginationParams{
				Limit:  5,
				Sort:   "created_at",
				Order:  "desc",
				Page:   1,
				Cursor: "eyJpZCI6IjEyMyJ9",
			},
		},
		{
			name:  "invalid limit uses default",
			query: "limit=200&page=0",
			expected: PaginationParams{
				Limit: 20,
				Sort:  "created_at",
				Order: "desc",
				Page:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tt.query != "" {
				values, _ := url.ParseQuery(tt.query)
				c.Request = httptest.NewRequest("GET", "/?"+tt.query, nil)
				c.Request.URL.RawQuery = values.Encode()
			} else {
				c.Request = httptest.NewRequest("GET", "/", nil)
			}

			params := parsePaginationParams(c)
			assert.Equal(t, tt.expected, params)
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	cursor := Cursor{
		ID:        "123",
		Timestamp: testTime,
	}

	encoded := encodeCursor(cursor)
	assert.NotEmpty(t, encoded)

	decoded, err := decodeCursor(encoded)
	assert.NoError(t, err)
	assert.Equal(t, cursor.ID, decoded.ID)
	assert.True(t, cursor.Timestamp.Equal(decoded.Timestamp))
}

func TestDecodeInvalidCursor(t *testing.T) {
	_, err := decodeCursor("invalid-cursor")
	assert.Error(t, err)

	cursor, err := decodeCursor("")
	assert.NoError(t, err)
	assert.Nil(t, cursor)
}

func TestBuildPaginatedResponse(t *testing.T) {
	data := []interface{}{"item1", "item2"}
	params := PaginationParams{
		Page:   2,
		Limit:  10,
		Sort:   "name",
		Order:  "asc",
		Search: "test",
		Status: "active",
	}

	response := buildPaginatedResponse(data, 25, params)

	assert.Equal(t, data, response.Data)
	assert.Equal(t, 2, response.Pagination.Page)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, 25, response.Pagination.Total)
	assert.Equal(t, 3, response.Pagination.TotalPages)
	assert.True(t, response.Pagination.HasNext)
	assert.True(t, response.Pagination.HasPrev)
	assert.Equal(t, "name", response.Meta.Sort)
	assert.Equal(t, "asc", response.Meta.Order)
	assert.Contains(t, response.Meta.FiltersApplied, "search")
	assert.Contains(t, response.Meta.FiltersApplied, "status")
}

func TestBuildCursorResponse(t *testing.T) {
	data := []interface{}{"item1", "item2"}
	params := PaginationParams{
		Limit:  10,
		Sort:   "name",
		Order:  "desc",
		Search: "test",
	}

	response := buildCursorResponse(data, true, "next_cursor", "prev_cursor", params)

	assert.Equal(t, data, response.Data)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.True(t, response.Pagination.HasNext)
	assert.True(t, response.Pagination.HasPrev)
	assert.Equal(t, "next_cursor", response.Pagination.NextCursor)
	assert.Equal(t, "prev_cursor", response.Pagination.PrevCursor)
	assert.Equal(t, "name", response.Meta.Sort)
	assert.Equal(t, "desc", response.Meta.Order)
	assert.Contains(t, response.Meta.FiltersApplied, "search")
}