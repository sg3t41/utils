package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/interfaces/dto"
	"github.com/sg3t41/api/internal/interfaces/middleware"
)

type UserHandler struct {
	createUserUseCase     *usecase.CreateUserUseCase
	getUserUseCase        *usecase.GetUserUseCase
	getUsersUseCase       *usecase.GetUsersUseCase
	updateUserUseCase     *usecase.UpdateUserUseCase
	updatePasswordUseCase *usecase.UpdatePasswordUseCase
	deleteUserUseCase     *usecase.DeleteUserUseCase
}

func NewUserHandler(
	createUserUseCase *usecase.CreateUserUseCase,
	getUserUseCase *usecase.GetUserUseCase,
	getUsersUseCase *usecase.GetUsersUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
	updatePasswordUseCase *usecase.UpdatePasswordUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase:     createUserUseCase,
		getUserUseCase:        getUserUseCase,
		getUsersUseCase:       getUsersUseCase,
		updateUserUseCase:     updateUserUseCase,
		updatePasswordUseCase: updatePasswordUseCase,
		deleteUserUseCase:     deleteUserUseCase,
	}
}



func (h *UserHandler) CreateUser(c *gin.Context) {
	createUserReq, ok := GetValidatedBody[dto.CreateUserRequest](c)
	if !ok {
		return
	}

	input := usecase.CreateUserInput{
		Email: createUserReq.Email,
		Name:  createUserReq.Name,
	}

	output, err := h.createUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := convertUserToResponse(output.User)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	input := usecase.GetUserInput{
		ID: id,
	}

	output, err := h.getUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        output.User.ID,
		Email:     output.User.Email,
		Name:      output.User.Name,
		Version:   output.User.Version,
		CreatedAt: output.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: output.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	params := parsePaginationParams(c)
	useCursor := params.Cursor != ""

	input := usecase.GetUsersInput{
		Limit:       params.Limit,
		Page:        params.Page,
		Cursor:      params.Cursor,
		Sort:        params.Sort,
		Order:       params.Order,
		Search:      params.Search,
		Status:      params.Status,
		CreatedFrom: params.CreatedFrom,
		CreatedTo:   params.CreatedTo,
		UseCursor:   useCursor,
	}

	output, err := h.getUsersUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users := make([]UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Version:   user.Version,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	if useCursor {
		var nextCursor, prevCursor string
		hasNext := len(output.Users) > params.Limit
		
		if hasNext {
			users = users[:params.Limit]
			lastUser := output.Users[params.Limit-1]
			nextCursor = encodeCursor(Cursor{
				ID:        lastUser.ID,
				Timestamp: lastUser.CreatedAt,
			})
		}

		if params.Cursor != "" {
			prevCursor = params.Cursor
		}

		response := buildCursorResponse(users, hasNext, nextCursor, prevCursor, params)
		c.JSON(http.StatusOK, response)
	} else {
		response := buildPaginatedResponse(users, output.Total, params)
		c.JSON(http.StatusOK, response)
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	hardDelete := c.Query("hard") == "true"

	input := usecase.DeleteUserInput{
		ID:   id,
		Hard: hardDelete,
	}

	_, err := h.deleteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	currentUser := middleware.GetUserFromContext(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if currentUser.ID != id && !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	req, exists := c.Get("validated_body")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "バリデーションが実行されていません"})
		return
	}

	updateUserReq, ok := req.(*dto.UpdateUserRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	input := usecase.UpdateUserInput{
		ID:   id,
		Name: updateUserReq.Name,
	}

	output, err := h.updateUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case usecase.ErrEmailAlreadyTaken:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case usecase.ErrVersionConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        output.User.ID,
		Email:     output.User.Email,
		Name:      output.User.Name,
		Version:   output.User.Version,
		CreatedAt: output.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: output.User.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	currentUser := middleware.GetUserFromContext(c)
	if currentUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if currentUser.ID != id && !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	req, exists := c.Get("validated_body")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "バリデーションが実行されていません"})
		return
	}

	updatePasswordReq, ok := req.(*dto.UpdatePasswordRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエスト形式です"})
		return
	}

	input := usecase.UpdatePasswordInput{
		UserID:          id,
		OldPassword:     updatePasswordReq.CurrentPassword,
		NewPassword:     updatePasswordReq.NewPassword,
		ConfirmPassword: updatePasswordReq.ConfirmPassword,
	}

	_, err := h.updatePasswordUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		switch err {
		case usecase.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case usecase.ErrInvalidOldPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrPasswordMismatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}