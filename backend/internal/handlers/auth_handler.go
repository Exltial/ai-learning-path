package handlers

import (
	"net/http"
	"ai-learning-platform/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest represents a registration request
// @Description User registration with username, email, and password
type RegisterRequest struct {
	// Username for the account (3-50 characters)
	// Required: true
	// Example: john_doe
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	
	// Email address (must be valid and unique)
	// Required: true
	// Example: john@example.com
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	
	// Password (minimum 8 characters)
	// Required: true
	// Example: SecurePass123
	Password string `json:"password" binding:"required,min=8" example:"SecurePass123"`
	
	// Optional avatar URL
	// Example: https://example.com/avatar.jpg
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
}

// LoginRequest represents a login request
// @Description User login with email and password
type LoginRequest struct {
	// Email address
	// Required: true
	// Example: john@example.com
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	
	// Password
	// Required: true
	// Example: SecurePass123
	Password string `json:"password" binding:"required" example:"SecurePass123"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	// Current valid JWT token
	// Required: true
	Token string `json:"token" binding:"required"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Creates a new user account and returns JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "User created successfully with token"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Failure 409 {object} map[string]interface{} "Username or email already exists"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid input: " + err.Error(),
			},
		})
		return
	}

	// Call service to register user
	user, token, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		message := err.Error()
		
		// Handle specific errors
		switch err {
		case services.ErrUsernameExists:
			status = http.StatusConflict
			code = "CONFLICT"
			message = "Username already exists"
		case services.ErrEmailExists:
			status = http.StatusConflict
			code = "CONFLICT"
			message = "Email already registered"
		}

		c.JSON(status, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    code,
				"message": message,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"created_at": user.CreatedAt,
			},
			"token": token,
			"token_type": "Bearer",
			"expires_in": 86400, // 24 hours in seconds
		},
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with token"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid input: " + err.Error(),
			},
		})
		return
	}

	// Call service to authenticate user
	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "Invalid email or password",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"role":       user.Role,
				"avatar_url": user.AvatarURL,
			},
			"token":      token,
			"token_type": "Bearer",
			"expires_in": 86400, // 24 hours
		},
	})
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Refreshes an existing JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Current token"
// @Success 200 {object} map[string]interface{} "New token generated"
// @Failure 400 {object} map[string]interface{} "Invalid token"
// @Failure 401 {object} map[string]interface{} "Token expired or invalid"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Invalid input: " + err.Error(),
			},
		})
		return
	}

	// Call service to refresh token
	newToken, err := h.authService.RefreshToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "UNAUTHORIZED",
				"message": "Invalid or expired token",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"token":      newToken,
			"token_type": "Bearer",
			"expires_in": 86400,
		},
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidates the current JWT token
// @Tags Authentication
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "BAD_REQUEST",
				"message": "Missing authorization header",
			},
		})
		return
	}

	// Extract token from "Bearer <token>"
	token := authHeader[7:] // Remove "Bearer " prefix

	// Invalidate token (add to blacklist)
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": map[string]interface{}{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to logout",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}
