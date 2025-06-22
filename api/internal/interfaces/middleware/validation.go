package middleware

import (
	"net/http"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/ja"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ValidationMiddleware handles input validation for API requests
type ValidationMiddleware struct {
	validator  *validator.Validate
	translator ut.Translator
}

// NewValidationMiddleware creates a new validation middleware instance
func NewValidationMiddleware() *ValidationMiddleware {
	v := validator.New()

	// Register custom validators
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("alpha_space", validateAlphaSpace)
	v.RegisterValidation("postal_code", validatePostalCode)

	// Setup translator for internationalization
	translator := setupTranslator(v)

	return &ValidationMiddleware{
		validator:  v,
		translator: translator,
	}
}

// setupTranslator configures internationalization for validation error messages
func setupTranslator(v *validator.Validate) ut.Translator {
	ja := ja.New()
	en := en.New()
	uni := ut.New(en, ja, en)

	// Get Japanese translator (fallback to English if not available)
	translator, _ := uni.GetTranslator("ja")
	if translator == nil {
		translator, _ = uni.GetTranslator("en")
	}

	// Register default translations
	en_translations.RegisterDefaultTranslations(v, translator)
	ja_translations.RegisterDefaultTranslations(v, translator)

	// Register custom validation error messages
	registerCustomTranslations(v, translator)

	return translator
}

// registerCustomTranslations registers custom error messages for validators
func registerCustomTranslations(v *validator.Validate, trans ut.Translator) {
	v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "パスワードは8文字以上で、大文字、小文字、数字、特殊文字を含む必要があります", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	v.RegisterTranslation("username", trans, func(ut ut.Translator) error {
		return ut.Add("username", "ユーザー名は3-30文字の英数字、アンダースコア、ハイフンのみ使用可能です", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username", fe.Field())
		return t
	})

	v.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "有効な国際電話番号を入力してください", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	v.RegisterTranslation("alpha_space", trans, func(ut ut.Translator) error {
		return ut.Add("alpha_space", "国際的な文字、数字、スペース、記号、絵文字が使用可能です", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha_space", fe.Field())
		return t
	})

	v.RegisterTranslation("postal_code", trans, func(ut ut.Translator) error {
		return ut.Add("postal_code", "有効な郵便番号を入力してください", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("postal_code", fe.Field())
		return t
	})
}

// ValidateJSON validates JSON request body
func (v *ValidationMiddleware) ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "無効なJSON形式です",
				"code":    "INVALID_JSON",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if err := v.validator.Struct(obj); err != nil {
			c.JSON(http.StatusBadRequest, v.formatValidationErrors(err))
			c.Abort()
			return
		}

		c.Set("validated_body", obj)
		c.Next()
	}
}

// ValidateQuery validates query parameters
func (v *ValidationMiddleware) ValidateQuery(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "無効なクエリパラメータです",
				"code":    "INVALID_QUERY",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if err := v.validator.Struct(obj); err != nil {
			c.JSON(http.StatusBadRequest, v.formatValidationErrors(err))
			c.Abort()
			return
		}

		c.Set("validated_query", obj)
		c.Next()
	}
}

// formatValidationErrors formats validation errors into a consistent response
func (v *ValidationMiddleware) formatValidationErrors(err error) gin.H {
	var errors []gin.H

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			errors = append(errors, gin.H{
				"field":   fe.Field(),
				"value":   fe.Value(),
				"message": fe.Translate(v.translator),
				"code":    fe.Tag(),
			})
		}
	}

	return gin.H{
		"error":  "バリデーションエラー",
		"code":   "VALIDATION_ERROR",
		"errors": errors,
	}
}

// Custom validation functions

// validatePassword validates password complexity
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}

	// Must contain uppercase, lowercase, number, and special character
	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateUsername validates username format
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Alphanumeric, underscore, hyphen, 3-30 characters
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{3,30}$`, username)
	return matched
}

// validatePhone validates international phone number format
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Basic international phone number validation
	// Must be 4-15 digits, optionally starting with +, and first digit cannot be 0
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{3,14}$`, phone)
	return matched
}

// validateAlphaSpace validates Unicode characters for international names including emojis
func validateAlphaSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// Allow Unicode letters, numbers, spaces, marks, punctuation, and symbols (including emojis)
	// This supports all international languages, emojis, and special characters
	for _, char := range value {
		if !unicode.IsLetter(char) && 
		   !unicode.IsSpace(char) && 
		   !unicode.IsNumber(char) && 
		   !unicode.IsMark(char) &&      // Diacritics, accents (Arabic, etc.)
		   !unicode.IsPunct(char) &&     // Punctuation (apostrophes, hyphens in names)
		   !unicode.IsSymbol(char) &&    // Symbols and emojis (🌟, ❤️, etc.)
		   char != '\'' &&               // Allow apostrophes (O'Connor, etc.)
		   char != '-' &&                // Allow hyphens (Jean-Pierre, etc.)
		   char != '.' {                 // Allow dots (Jr., Sr., etc.)
			return false
		}
	}
	return len(value) > 0
}

// validatePostalCode validates postal code format (basic implementation)
func validatePostalCode(fl validator.FieldLevel) bool {
	postalCode := fl.Field().String()

	// Basic postal code validation (customize based on requirements)
	// This supports various formats: 12345, 12345-6789, A1A 1A1, etc.
	matched, _ := regexp.MatchString(`^[A-Za-z0-9\s\-]{3,10}$`, postalCode)
	return matched
}