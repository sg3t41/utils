package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/domain/entity"
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
		HandleUseCaseError(c, err, "ユーザー")
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
		HandleUseCaseError(c, err, "ユーザー")
		return
	}

	response := convertUserToResponse(output.User)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	params := parsePaginationParams(c)
	
	input := h.buildGetUsersInput(params)
	output, err := h.getUsersUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "ユーザー")
		return
	}

	h.sendUsersResponse(c, output, params)
}

// buildGetUsersInput GetUsersInputの構築
func (h *UserHandler) buildGetUsersInput(params PaginationParams) usecase.GetUsersInput {
	return usecase.GetUsersInput{
		Limit:       params.Limit,
		Page:        params.Page,
		Cursor:      params.Cursor,
		Sort:        params.Sort,
		Order:       params.Order,
		Search:      params.Search,
		Status:      params.Status,
		CreatedFrom: params.CreatedFrom,
		CreatedTo:   params.CreatedTo,
		UseCursor:   params.Cursor != "",
	}
}

// sendUsersResponse ユーザー一覧のレスポンス送信
func (h *UserHandler) sendUsersResponse(c *gin.Context, output *usecase.GetUsersOutput, params PaginationParams) {
	users := make([]UserResponse, len(output.Users))
	for i, user := range output.Users {
		users[i] = convertUserToResponse(user)
	}

	if params.Cursor != "" {
		h.sendCursorPaginatedResponse(c, users, output.Users, params)
	} else {
		response := buildPaginatedResponse(users, output.Total, params)
		c.JSON(http.StatusOK, response)
	}
}

// sendCursorPaginatedResponse カーソルページネーションレスポンスの送信
func (h *UserHandler) sendCursorPaginatedResponse(c *gin.Context, users []UserResponse, allUsers []*entity.User, params PaginationParams) {
	var nextCursor, prevCursor string
	hasNext := len(allUsers) > params.Limit
	
	if hasNext {
		users = users[:params.Limit]
		lastUser := allUsers[params.Limit-1]
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
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	hardDelete := c.Query("hard") == "true"

	input := usecase.DeleteUserInput{
		ID:   id,
		Hard: hardDelete,
	}

	_, err := h.deleteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		HandleUseCaseError(c, err, "ユーザー")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	if !h.validateUserAccess(c, id) {
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
		HandleUseCaseError(c, err, "ユーザー")
		return
	}

	response := convertUserToResponse(output.User)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	id, ok := ParseIDParam(c)
	if !ok {
		return
	}

	if !h.validateUserAccess(c, id) {
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
		HandleUseCaseError(c, err, "パスワード")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}

// validateUserAccess ユーザーアクセス権限の検証
func (h *UserHandler) validateUserAccess(c *gin.Context, targetUserID string) bool {
	currentUser := middleware.GetUserFromContext(c)
	if currentUser == nil {
		SendErrorResponse(c, http.StatusUnauthorized, "認証されていません")
		return false
	}

	if currentUser.ID != targetUserID && !currentUser.IsAdmin {
		SendErrorResponse(c, http.StatusForbidden, "アクセス権限がありません")
		return false
	}

	return true
}