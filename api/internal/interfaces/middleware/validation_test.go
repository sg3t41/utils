package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ValidateJSON - Valid Request", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,password"`
			Name     string `json:"name" validate:"required,alpha_space"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			req, exists := c.Get("validated_body")
			assert.True(t, exists)
			
			testReq, ok := req.(*TestRequest)
			assert.True(t, ok)
			assert.Equal(t, "test@example.com", testReq.Email)
			
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		body := map[string]interface{}{
			"email":    "test@example.com",
			"password": "Password123!",
			"name":     "Test User",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ValidateJSON - Invalid Email", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Email string `json:"email" validate:"required,email"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		body := map[string]interface{}{
			"email": "invalid-email",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, "バリデーションエラー", response["error"])
		assert.Equal(t, "VALIDATION_ERROR", response["code"])
		assert.NotNil(t, response["errors"])
	})

	t.Run("ValidateQuery - Valid Query Parameters", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestQuery struct {
			Page  int    `form:"page" validate:"min=1"`
			Limit int    `form:"limit" validate:"min=1,max=100"`
			Sort  string `form:"sort" validate:"omitempty,oneof=name email"`
		}

		router := gin.New()
		router.GET("/test", middleware.ValidateQuery(&TestQuery{}), func(c *gin.Context) {
			query, exists := c.Get("validated_query")
			assert.True(t, exists)
			
			testQuery, ok := query.(*TestQuery)
			assert.True(t, ok)
			assert.Equal(t, 1, testQuery.Page)
			assert.Equal(t, 10, testQuery.Limit)
			assert.Equal(t, "name", testQuery.Sort)
			
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		req := httptest.NewRequest("GET", "/test?page=1&limit=10&sort=name", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ValidateQuery - Invalid Query Parameters", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestQuery struct {
			Page  int `form:"page" validate:"min=1"`
			Limit int `form:"limit" validate:"min=1,max=100"`
		}

		router := gin.New()
		router.GET("/test", middleware.ValidateQuery(&TestQuery{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		req := httptest.NewRequest("GET", "/test?page=0&limit=101", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.Equal(t, "バリデーションエラー", response["error"])
		assert.Equal(t, "VALIDATION_ERROR", response["code"])
		assert.NotNil(t, response["errors"])
	})
}

func TestCustomValidators(t *testing.T) {
	t.Run("validatePassword - Valid Password", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Password string `json:"password" validate:"password"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		body := map[string]interface{}{
			"password": "Password123!",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validatePassword - Invalid Password", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Password string `json:"password" validate:"password"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		testCases := []string{
			"password",        // no uppercase, number, special char
			"PASSWORD123!",    // no lowercase
			"Password!",       // no number
			"Password123",     // no special char
			"Pass1!",          // too short
		}

		for _, password := range testCases {
			body := map[string]interface{}{
				"password": password,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Password %s should be invalid", password)
		}
	})

	t.Run("validateUsername - Valid Username", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Username string `json:"username" validate:"username"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		validUsernames := []string{
			"user123",
			"test_user",
			"user-name",
			"User123",
			"a123",
		}

		for _, username := range validUsernames {
			body := map[string]interface{}{
				"username": username,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Username %s should be valid", username)
		}
	})

	t.Run("validateUsername - Invalid Username", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Username string `json:"username" validate:"username"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		invalidUsernames := []string{
			"us",                                   // too short
			"this-is-a-very-long-username-123456", // too long
			"user name",                           // space not allowed
			"user@name",                          // @ not allowed
			"user.name",                          // . not allowed
		}

		for _, username := range invalidUsernames {
			body := map[string]interface{}{
				"username": username,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Username %s should be invalid", username)
		}
	})

	t.Run("validatePhone - Valid Phone Numbers", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Phone string `json:"phone" validate:"phone"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		validPhones := []string{
			"+1234567890",
			"1234567890",
			"+819012345678",
			"9012345678",
		}

		for _, phone := range validPhones {
			body := map[string]interface{}{
				"phone": phone,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Phone %s should be valid", phone)
		}
	})

	t.Run("validatePhone - Invalid Phone Numbers", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Phone string `json:"phone" validate:"phone"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		invalidPhones := []string{
			"123",            // too short
			"+0123456789",    // starts with 0
			"abc123456789",   // contains letters
			"123-456-7890",   // contains hyphens
			"",               // empty
		}

		for _, phone := range invalidPhones {
			body := map[string]interface{}{
				"phone": phone,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Phone %s should be invalid", phone)
		}
	})

	t.Run("validateAlphaSpace - Valid Names", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Name string `json:"name" validate:"alpha_space"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		validNames := []string{
			"John Doe",
			"Mary Jane",
			"Alice",
			"Bob Smith Jr",
		}

		for _, name := range validNames {
			body := map[string]interface{}{
				"name": name,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Name %s should be valid", name)
		}
	})

	t.Run("validateAlphaSpace - Invalid Names", func(t *testing.T) {
		middleware := NewValidationMiddleware()
		
		type TestRequest struct {
			Name string `json:"name" validate:"alpha_space"`
		}

		router := gin.New()
		router.POST("/test", middleware.ValidateJSON(&TestRequest{}), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		invalidNames := []string{
			"John123",        // contains numbers
			"John@Doe",       // contains special chars
			"John_Doe",       // contains underscore
			"John-Doe",       // contains hyphen
		}

		for _, name := range invalidNames {
			body := map[string]interface{}{
				"name": name,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "Name %s should be invalid", name)
		}
	})
}