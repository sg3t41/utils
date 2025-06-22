package dto


// CreateArticleRequest represents the request payload for creating an article
type CreateArticleRequest struct {
	Title         string  `json:"title" validate:"required,min=1,max=500"`
	Content       string  `json:"content" validate:"required,min=1"`
	Summary       string  `json:"summary" validate:"max=1000"`
	Tags          []string `json:"tags" validate:"dive,min=1,max=50"`
	FeaturedImage *string `json:"featured_image" validate:"omitempty,max=500"`
}

// UpdateArticleRequest represents the request payload for updating an article
type UpdateArticleRequest struct {
	Title         *string  `json:"title" validate:"omitempty,min=1,max=500"`
	Content       *string  `json:"content" validate:"omitempty,min=1"`
	Summary       *string  `json:"summary" validate:"omitempty,max=1000"`
	Tags          []string `json:"tags" validate:"dive,min=1,max=50"`
	FeaturedImage *string  `json:"featured_image" validate:"omitempty,max=500"`
}

// ArticleResponse represents the response payload for an article
type ArticleResponse struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	Summary        string  `json:"summary"`
	Status         string  `json:"status"`
	AuthorID       string  `json:"author_id"`
	Tags           []string `json:"tags"`
	FeaturedImage  *string `json:"featured_image"`
	ThumbnailImage *string `json:"thumbnail_image"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	PublishedAt    *string `json:"published_at"`
}

// ArticleListResponse represents the response payload for article list
type ArticleListResponse struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Summary        string  `json:"summary"`
	Status         string  `json:"status"`
	Tags           []string `json:"tags"`
	ThumbnailImage *string `json:"thumbnail_image"`
	CreatedAt      string  `json:"created_at"`
	PublishedAt    *string `json:"published_at"`
}

// ArticlesResponse represents the paginated response for articles list
type ArticlesResponse struct {
	Data       []ArticleListResponse `json:"data"`
	Pagination PaginationMeta        `json:"pagination"`
	Meta       SortMeta              `json:"meta"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// SortMeta represents sorting metadata
type SortMeta struct {
	Sort  string `json:"sort"`
	Order string `json:"order"`
}

// ListArticlesQuery represents query parameters for listing articles
type ListArticlesQuery struct {
	Page     int    `form:"page" validate:"min=1"`
	Limit    int    `form:"limit" validate:"min=1,max=100"`
	Sort     string `form:"sort" validate:"omitempty,oneof=created_at updated_at published_at title"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
	Status   string `form:"status" validate:"omitempty,oneof=draft published archived"`
	Tag      string `form:"tag" validate:"omitempty,max=50"`
	Search   string `form:"search" validate:"omitempty,max=100"`
	DateFrom string `form:"date_from" validate:"omitempty,datetime=2006-01-02"`
	DateTo   string `form:"date_to" validate:"omitempty,datetime=2006-01-02"`
}

// GetArticleQuery represents query parameters for getting an article
type GetArticleQuery struct {
	Include string `form:"include" validate:"omitempty,oneof=author full"`
}

// PublishArticleRequest represents the request payload for publishing an article
type PublishArticleRequest struct {
	// Empty for now, but can be extended with publish options
}

// ImageUploadResponse represents the response payload for image upload
type ImageUploadResponse struct {
	ImagePath      string `json:"image_path"`
	ThumbnailPath  string `json:"thumbnail_path"`
	OriginalName   string `json:"original_name"`
	Size           int64  `json:"size"`
	ContentType    string `json:"content_type"`
}

// Helper functions to convert between DTO and domain objects

// ToCreateArticleInput converts CreateArticleRequest to usecase input
func (req *CreateArticleRequest) ToCreateArticleInput(authorID string) interface{} {
	// This will be implemented when we create the specific usecase input type
	return nil
}

// FromArticleEntity converts Article entity to ArticleResponse
func FromArticleEntity(article interface{}) *ArticleResponse {
	// This will be implemented when we have the Article entity imported
	return nil
}

// FromArticleEntityList converts Article entities to ArticleListResponse slice
func FromArticleEntityList(articles interface{}) []ArticleListResponse {
	// This will be implemented when we have the Article entity imported
	return nil
}