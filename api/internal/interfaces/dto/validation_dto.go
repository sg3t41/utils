package dto

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email           string  `json:"email" validate:"required,email,max=255"`
	Name            string  `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Age             *int    `json:"age" validate:"omitempty,min=13,max=120"`
	Phone           *string `json:"phone" validate:"omitempty,phone"`
	Website         *string `json:"website" validate:"omitempty,url"`
	Address         *Address `json:"address" validate:"omitempty"`
}

// CreateUserWithPasswordRequest represents the request payload for creating a user with password
type CreateUserWithPasswordRequest struct {
	Email           string  `json:"email" validate:"required,email,max=255"`
	Password        string  `json:"password" validate:"required,password"`
	ConfirmPassword string  `json:"confirm_password" validate:"required,eqfield=Password"`
	Name            string  `json:"name" validate:"required,min=2,max=100,alpha_space"`
	Age             *int    `json:"age" validate:"omitempty,min=13,max=120"`
	Phone           *string `json:"phone" validate:"omitempty,phone"`
	Website         *string `json:"website" validate:"omitempty,url"`
	Address         *Address `json:"address" validate:"omitempty"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name     *string  `json:"name" validate:"omitempty,min=2,max=100,alpha_space"`
	Age      *int     `json:"age" validate:"omitempty,min=13,max=120"`
	Phone    *string  `json:"phone" validate:"omitempty,phone"`
	Website  *string  `json:"website" validate:"omitempty,url"`
	Address  *Address `json:"address" validate:"omitempty"`
}

// UpdatePasswordRequest represents the request payload for updating user password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

// Address represents a user's address with validation
type Address struct {
	Street     string `json:"street" validate:"required,max=200"`
	City       string `json:"city" validate:"required,max=100"`
	State      string `json:"state" validate:"required,len=2"`
	PostalCode string `json:"postal_code" validate:"required,postal_code"`
	Country    string `json:"country" validate:"required,iso3166_1_alpha2"`
}

// ListUsersQuery represents query parameters for listing users
type ListUsersQuery struct {
	Page   int    `form:"page" validate:"min=1"`
	Limit  int    `form:"limit" validate:"min=1,max=100"`
	Sort   string `form:"sort" validate:"omitempty,oneof=name email created_at updated_at"`
	Order  string `form:"order" validate:"omitempty,oneof=asc desc"`
	Status string `form:"status" validate:"omitempty,oneof=active inactive pending"`
	Search string `form:"search" validate:"omitempty,max=100"`
}

// GetUserQuery represents query parameters for getting a user
type GetUserQuery struct {
	Include string `form:"include" validate:"omitempty,oneof=address profile all"`
}