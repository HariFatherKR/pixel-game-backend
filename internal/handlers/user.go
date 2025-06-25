package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/pixel-game/internal/domain"
)

type UserHandler struct {
	userRepository domain.UserRepository
}

func NewUserHandler(userRepository domain.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

// UpdateProfile godoc
// @Summary      사용자 프로필 수정
// @Description  현재 로그인한 사용자의 프로필 정보를 수정합니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request body domain.UpdateUserProfileRequest true "프로필 수정 정보"
// @Success      200  {object}  domain.UserProfile "수정된 프로필"
// @Failure      400  {object}  ErrorResponse      "잘못된 요청"
// @Failure      401  {object}  ErrorResponse      "인증 실패"
// @Failure      404  {object}  ErrorResponse      "사용자를 찾을 수 없음"
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req domain.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid Request",
			Message: err.Error(),
		})
		return
	}

	profile, err := h.userRepository.GetProfile(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get profile",
		})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Profile Not Found",
			Message: "User profile not found",
		})
		return
	}

	if req.DisplayName != "" {
		profile.DisplayName = req.DisplayName
	}
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}

	if err := h.userRepository.UpdateProfile(profile); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetStats godoc
// @Summary      사용자 통계 조회
// @Description  현재 로그인한 사용자의 게임 통계를 조회합니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  domain.UserStats  "사용자 통계"
// @Failure      401  {object}  ErrorResponse     "인증 실패"
// @Failure      404  {object}  ErrorResponse     "통계를 찾을 수 없음"
// @Router       /users/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	stats, err := h.userRepository.GetStats(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get stats",
		})
		return
	}
	if stats == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Stats Not Found",
			Message: "User stats not found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetCollection godoc
// @Summary      사용자 카드 컬렉션 조회
// @Description  현재 로그인한 사용자의 카드 컬렉션을 조회합니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]interface{}  "카드 컬렉션 (추후 구현)"
// @Failure      401  {object}  ErrorResponse           "인증 실패"
// @Router       /users/collection [get]
func (h *UserHandler) GetCollection(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"collection": []interface{}{},
		"message":    "Card collection feature will be implemented in Phase 5",
	})
}

// IncrementGamesPlayed godoc
// @Summary      게임 플레이 횟수 증가
// @Description  사용자의 게임 플레이 횟수를 1 증가시킵니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]string "성공 메시지"
// @Failure      401  {object}  ErrorResponse     "인증 실패"
// @Router       /users/stats/games-played [post]
func (h *UserHandler) IncrementGamesPlayed(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	if err := h.userRepository.IncrementGamesPlayed(userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to increment games played",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Games played incremented successfully",
	})
}

// IncrementGamesWon godoc
// @Summary      게임 승리 횟수 증가
// @Description  사용자의 게임 승리 횟수를 1 증가시킵니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  map[string]string "성공 메시지"
// @Failure      401  {object}  ErrorResponse     "인증 실패"
// @Router       /users/stats/games-won [post]
func (h *UserHandler) IncrementGamesWon(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	if err := h.userRepository.IncrementGamesWon(userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to increment games won",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Games won incremented successfully",
	})
}

// AddPlayTime godoc
// @Summary      플레이 시간 추가
// @Description  사용자의 총 플레이 시간을 추가합니다.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        seconds path int true "추가할 시간(초)"
// @Success      200  {object}  map[string]string "성공 메시지"
// @Failure      400  {object}  ErrorResponse     "잘못된 요청"
// @Failure      401  {object}  ErrorResponse     "인증 실패"
// @Router       /users/stats/play-time/{seconds} [post]
func (h *UserHandler) AddPlayTime(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	secondsStr := c.Param("seconds")
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil || seconds < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid Parameter",
			Message: "Seconds must be a positive integer",
		})
		return
	}

	if err := h.userRepository.AddPlayTime(userID.(int), seconds); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to add play time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Play time added successfully",
		"seconds": seconds,
	})
}