package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sg3t41/api/internal/application/usecase"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/infrastructure/persistence"
	"github.com/sg3t41/api/internal/interfaces/middleware"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *UserHandler) {
	gin.SetMode(gin.TestMode)
	
	userRepo := persistence.NewMemoryUserRepository()
	userService := service.NewUserService(userRepo)
	createUserUseCase := usecase.NewCreateUserUseCase(userService)
	getUserUseCase := usecase.NewGetUserUseCase(userService)
	getUsersUseCase := usecase.NewGetUsersUseCase(userService)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepo)
	updatePasswordUseCase := usecase.NewUpdatePasswordUseCase(userRepo)
	deleteUserUseCase := usecase.NewDeleteUserUseCase(userService)
	
	userHandler := NewUserHandler(
		createUserUseCase,
		getUserUseCase,
		getUsersUseCase,
		updateUserUseCase,
		updatePasswordUseCase,
		deleteUserUseCase,
	)
	
	router := gin.New()
	return router, userHandler
}

func TestUpdateUser_Success(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Create a test user first
	user, _ := entity.NewUser("test@example.com", "Test User")
	user.Password = "hashedpassword"
	
	// Setup route with auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: user.ID, IsAdmin: false})
		c.Next()
	})
	router.PATCH("/users/:id", handler.UpdateUser)
	
	// Test partial update
	updateReq := dto.UpdateUserRequest{
		Name: stringPtr("Updated Name"),
	}
	
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/"+user.ID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Name)
	assert.Equal(t, user.Email, response.Email)
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Setup route without auth middleware
	router.PATCH("/users/:id", handler.UpdateUser)
	
	updateReq := dto.UpdateUserRequest{
		Name: stringPtr("Updated Name"),
	}
	
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateUser_Forbidden(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Setup route with different user auth
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: "different-user", IsAdmin: false})
		c.Next()
	})
	router.PATCH("/users/:id", handler.UpdateUser)
	
	updateReq := dto.UpdateUserRequest{
		Name: stringPtr("Updated Name"),
	}
	
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdatePassword_Success(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Create a test user first
	user, _ := entity.NewUser("test@example.com", "Test User")
	user.Password = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8" // sha256 of "password"
	
	// Setup route with auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: user.ID, IsAdmin: false})
		c.Next()
	})
	router.PATCH("/users/:id/password", handler.UpdatePassword)
	
	updateReq := dto.UpdatePasswordRequest{
		CurrentPassword: "password",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/"+user.ID+"/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "password updated successfully", response["message"])
}

func TestUpdatePassword_PasswordMismatch(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Setup route with auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: "test-id", IsAdmin: false})
		c.Next()
	})
	router.PATCH("/users/:id/password", handler.UpdatePassword)
	
	updateReq := dto.UpdatePasswordRequest{
		CurrentPassword: "password",
		Password:        "newpassword123",
		ConfirmPassword: "different-password",
	}
	
	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/test-id/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func stringPtr(s string) *string {
	return &s
}

func TestDeleteUser_SoftDelete_Success(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Create a test user first
	user, _ := entity.NewUser("test@example.com", "Test User")
	user.Password = "hashedpassword"
	
	// Setup route with auth middleware (user can delete their own account)
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: user.ID, IsAdmin: false})
		c.Next()
	})
	router.DELETE("/users/:id", handler.DeleteUser)
	
	req := httptest.NewRequest("DELETE", "/users/"+user.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUser_HardDelete_Success(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Create a test user first
	user, _ := entity.NewUser("test@example.com", "Test User")
	user.Password = "hashedpassword"
	
	// Setup route with admin auth middleware (admin can hard delete)
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: "admin-id", IsAdmin: true})
		c.Next()
	})
	router.DELETE("/users/:id", handler.DeleteUser)
	
	req := httptest.NewRequest("DELETE", "/users/"+user.ID+"?hard=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUser_MissingID(t *testing.T) {
	router, handler := setupTestRouter()
	
	router.DELETE("/users/:id", handler.DeleteUser)
	
	req := httptest.NewRequest("DELETE", "/users/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUser_HardDeleteQueryParam(t *testing.T) {
	router, handler := setupTestRouter()
	
	// Setup route with admin auth
	router.Use(func(c *gin.Context) {
		c.Set("user", &middleware.User{ID: "admin-id", IsAdmin: true})
		c.Next()
	})
	router.DELETE("/users/:id", handler.DeleteUser)
	
	testCases := []struct {
		name  string
		query string
		code  int
	}{
		{
			name:  "hard=true should work",
			query: "hard=true",
			code:  http.StatusNoContent,
		},
		{
			name:  "hard=false should work",
			query: "hard=false",
			code:  http.StatusNoContent,
		},
		{
			name:  "no hard param should work",
			query: "",
			code:  http.StatusNoContent,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/users/test-user-id"
			if tc.query != "" {
				url += "?" + tc.query
			}
			
			req := httptest.NewRequest("DELETE", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tc.code, w.Code)
		})
	}
}