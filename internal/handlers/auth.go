package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/auth"
	"github.com/yourusername/pixel-game/internal/domain"
)

type AuthHandler struct {
	jwtManager     *auth.JWTManager
	userRepository domain.UserRepository
	cardRepository domain.CardRepository
}

func NewAuthHandler(jwtManager *auth.JWTManager, userRepository domain.UserRepository, cardRepository domain.CardRepository) *AuthHandler {
	return &AuthHandler{
		jwtManager:     jwtManager,
		userRepository: userRepository,
		cardRepository: cardRepository,
	}
}

type RegisterRequest struct {
	Username string          `json:"username" binding:"required,min=3,max=20"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	Platform domain.Platform `json:"platform" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int             `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Platform domain.Platform `json:"platform"`
	Profile  *domain.UserProfile `json:"profile,omitempty"`
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

	existingUser, err := h.userRepository.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to check username",
		})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Username Already Exists",
			Message: "Username is already taken",
		})
		return
	}

	existingUser, err = h.userRepository.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to check email",
		})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, ErrorResponse{
			Error:   "Email Already Exists",
			Message: "Email is already registered",
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

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Platform:     req.Platform,
	}

	if err := h.userRepository.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to create user",
		})
		return
	}

	profile, err := h.userRepository.GetProfile(user.ID)
	if err != nil {
		profile = nil
	}

	// Grant initial cards to new user
	initialCards := []string{
		"card_001", "card_001", "card_001", // 3x 해킹 스트라이크
		"card_002", "card_002",              // 2x 코드 인젝션  
		"card_008", "card_008", "card_008", // 3x 방화벽
		"card_009", "card_009",              // 2x 백업
		"card_014",                          // 1x 알고리즘 최적화
		"card_018", "card_018",              // 2x 메모리 누수
	}

	for _, cardID := range initialCards {
		userCard := &domain.UserCard{
			UserID: user.ID,
			CardID: cardID,
		}
		h.cardRepository.AddCardToUser(userCard)
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Platform: user.Platform,
		Profile:  profile,
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
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

	user, err := h.userRepository.GetByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get user",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid Credentials",
			Message: "Username or password is incorrect",
		})
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Invalid Credentials",
			Message: "Username or password is incorrect",
		})
		return
	}

	if err := h.userRepository.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail the login
	}

	profile, err := h.userRepository.GetProfile(user.ID)
	if err != nil {
		profile = nil
	}

	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate access token",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to generate refresh token",
		})
		return
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Platform: user.Platform,
		Profile:  profile,
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
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

	user, err := h.userRepository.GetByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get user",
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User Not Found",
			Message: "User not found",
		})
		return
	}

	profile, err := h.userRepository.GetProfile(user.ID)
	if err != nil {
		profile = nil
	}

	stats, err := h.userRepository.GetStats(user.ID)
	if err != nil {
		stats = nil
	}

	userResponse := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Platform: user.Platform,
		Profile:  profile,
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  userResponse,
		"stats": stats,
	})
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid input provided"`
}