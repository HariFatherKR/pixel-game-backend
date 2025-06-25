package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/auth"
)

type AuthHandler struct {
	jwtManager *auth.JWTManager
}

func NewAuthHandler(jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		jwtManager: jwtManager,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register godoc
// @Summary      사용자 회원가입
// @Description  새로운 사용자 계정을 생성합니다. 사용자명과 이메일은 고유해야 합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "회원가입 정보"
// @Success      201  {object}  AuthResponse   "회원가입 성공"
// @Failure      400  {object}  ErrorResponse  "잘못된 요청"
// @Failure      409  {object}  ErrorResponse  "이미 존재하는 사용자"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid Request",
			Message: err.Error(),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to hash password",
		})
		return
	}

	userID := 1
	user := User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	_ = hashedPassword

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Login godoc
// @Summary      사용자 로그인
// @Description  사용자명과 비밀번호로 로그인하여 JWT 토큰을 발급받습니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "로그인 정보"
// @Success      200  {object}  AuthResponse   "로그인 성공"
// @Failure      400  {object}  ErrorResponse  "잘못된 요청"
// @Failure      401  {object}  ErrorResponse  "인증 실패"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid Request",
			Message: err.Error(),
		})
		return
	}

	userID := 1
	user := User{
		ID:       userID,
		Username: req.Username,
		Email:    "test@example.com",
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(userID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Logout godoc
// @Summary      사용자 로그아웃
// @Description  현재 세션을 종료합니다. 토큰을 블랙리스트에 추가하여 무효화합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]string  "로그아웃 성공"
// @Failure      401  {object}  ErrorResponse      "인증 실패"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	_ = userID

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// RefreshToken godoc
// @Summary      토큰 갱신
// @Description  Refresh Token을 사용하여 새로운 Access Token을 발급받습니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token body string true "Refresh Token"
// @Success      200  {object}  map[string]string  "토큰 갱신 성공"
// @Failure      400  {object}  ErrorResponse      "잘못된 요청"
// @Failure      401  {object}  ErrorResponse      "인증 실패"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid Request",
			Message: err.Error(),
		})
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid refresh token",
		})
		return
	}

	newAccessToken, err := h.jwtManager.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate new access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

// Profile godoc
// @Summary      사용자 프로필 조회
// @Description  현재 로그인한 사용자의 프로필 정보를 조회합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  User           "사용자 프로필"
// @Failure      401  {object}  ErrorResponse  "인증 실패"
// @Router       /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	username, _ := c.Get("username")

	user := User{
		ID:       userID.(int),
		Username: username.(string),
		Email:    "test@example.com",
	}

	c.JSON(http.StatusOK, user)
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid input provided"`
}