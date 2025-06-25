package domain

import (
	"time"
)

type User struct {
	ID                int                `json:"id" db:"id"`
	Username          string             `json:"username" db:"username"`
	Email             string             `json:"email" db:"email"`
	PasswordHash      string             `json:"-" db:"password_hash"`
	Platform          Platform           `json:"platform" db:"platform"`
	IsActive          bool               `json:"is_active" db:"is_active"`
	LastLoginAt       *time.Time         `json:"last_login_at" db:"last_login_at"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`
	Profile           *UserProfile       `json:"profile,omitempty"`
	Stats             *UserStats         `json:"stats,omitempty"`
}

type Platform string

const (
	PlatformWeb     Platform = "web"
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

type UserProfile struct {
	UserID      int    `json:"user_id" db:"user_id"`
	DisplayName string `json:"display_name" db:"display_name"`
	Avatar      string `json:"avatar" db:"avatar"`
	Bio         string `json:"bio" db:"bio"`
	Level       int    `json:"level" db:"level"`
	Experience  int    `json:"experience" db:"experience"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserStats struct {
	UserID           int `json:"user_id" db:"user_id"`
	GamesPlayed      int `json:"games_played" db:"games_played"`
	GamesWon         int `json:"games_won" db:"games_won"`
	TotalPlayTime    int `json:"total_play_time" db:"total_play_time"`
	HighestLevel     int `json:"highest_level" db:"highest_level"`
	CardsCollected   int `json:"cards_collected" db:"cards_collected"`
	AchievementsCount int `json:"achievements_count" db:"achievements_count"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=20"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Platform Platform `json:"platform" binding:"required"`
}

type UpdateUserProfileRequest struct {
	DisplayName string `json:"display_name" binding:"max=50"`
	Avatar      string `json:"avatar" binding:"max=255"`
	Bio         string `json:"bio" binding:"max=500"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	UpdateLastLogin(userID int) error
	Delete(id int) error
	
	CreateProfile(profile *UserProfile) error
	GetProfile(userID int) (*UserProfile, error)
	UpdateProfile(profile *UserProfile) error
	
	GetStats(userID int) (*UserStats, error)
	UpdateStats(stats *UserStats) error
	IncrementGamesPlayed(userID int) error
	IncrementGamesWon(userID int) error
	AddPlayTime(userID int, seconds int) error
}