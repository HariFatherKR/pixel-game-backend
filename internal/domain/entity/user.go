package entity

import (
	"time"

	"github.com/google/uuid"
)

type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformWeb     Platform = "web"
)

type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Platform     Platform   `json:"platform" db:"platform"`
	DeviceID     *string    `json:"device_id,omitempty" db:"device_id"`
	LastLogin    *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type UserStats struct {
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	TotalRuns        int       `json:"total_runs" db:"total_runs"`
	TotalWins        int       `json:"total_wins" db:"total_wins"`
	HighestAscension int       `json:"highest_ascension" db:"highest_ascension"`
	AchievementPoints int      `json:"achievement_points" db:"achievement_points"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(username, email, passwordHash string, platform Platform) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Platform:     platform,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

func (u *User) SetDeviceID(deviceID string) {
	u.DeviceID = &deviceID
	u.UpdatedAt = time.Now()
}