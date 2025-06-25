package postgres

import (
	"database/sql"
	"time"

	"github.com/yourusername/pixel-game/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, platform, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	
	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Platform,
		true,
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return err
	}
	
	user.IsActive = true
	
	profile := &domain.UserProfile{
		UserID:      user.ID,
		DisplayName: user.Username,
		Level:       1,
		Experience:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	if err := r.CreateProfile(profile); err != nil {
		return err
	}
	
	stats := &domain.UserStats{
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	
	return r.createStats(stats)
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, platform, is_active, 
			   last_login_at, created_at, updated_at
		FROM users 
		WHERE id = $1 AND is_active = true`
	
	user := &domain.User{}
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Platform,
		&user.IsActive,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	
	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, platform, is_active, 
			   last_login_at, created_at, updated_at
		FROM users 
		WHERE username = $1 AND is_active = true`
	
	user := &domain.User{}
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Platform,
		&user.IsActive,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, password_hash, platform, is_active, 
			   last_login_at, created_at, updated_at
		FROM users 
		WHERE email = $1 AND is_active = true`
	
	user := &domain.User{}
	var lastLoginAt sql.NullTime
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Platform,
		&user.IsActive,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	
	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, platform = $4, updated_at = $5
		WHERE id = $1 AND is_active = true`
	
	user.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.Platform,
		user.UpdatedAt,
	)
	
	return err
}

func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `
		UPDATE users 
		SET last_login_at = $2, updated_at = $2
		WHERE id = $1 AND is_active = true`
	
	now := time.Now()
	_, err := r.db.Exec(query, userID, now)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `
		UPDATE users 
		SET is_active = false, updated_at = $2
		WHERE id = $1`
	
	_, err := r.db.Exec(query, id, time.Now())
	return err
}

func (r *UserRepository) CreateProfile(profile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, display_name, avatar, bio, level, experience, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`
	
	now := time.Now()
	err := r.db.QueryRow(
		query,
		profile.UserID,
		profile.DisplayName,
		profile.Avatar,
		profile.Bio,
		profile.Level,
		profile.Experience,
		now,
		now,
	).Scan(&profile.CreatedAt, &profile.UpdatedAt)
	
	return err
}

func (r *UserRepository) GetProfile(userID int) (*domain.UserProfile, error) {
	query := `
		SELECT user_id, display_name, avatar, bio, level, experience, created_at, updated_at
		FROM user_profiles 
		WHERE user_id = $1`
	
	profile := &domain.UserProfile{}
	
	err := r.db.QueryRow(query, userID).Scan(
		&profile.UserID,
		&profile.DisplayName,
		&profile.Avatar,
		&profile.Bio,
		&profile.Level,
		&profile.Experience,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return profile, nil
}

func (r *UserRepository) UpdateProfile(profile *domain.UserProfile) error {
	query := `
		UPDATE user_profiles 
		SET display_name = $2, avatar = $3, bio = $4, level = $5, experience = $6, updated_at = $7
		WHERE user_id = $1`
	
	profile.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(
		query,
		profile.UserID,
		profile.DisplayName,
		profile.Avatar,
		profile.Bio,
		profile.Level,
		profile.Experience,
		profile.UpdatedAt,
	)
	
	return err
}

func (r *UserRepository) GetStats(userID int) (*domain.UserStats, error) {
	query := `
		SELECT user_id, games_played, games_won, total_play_time, highest_level, 
			   cards_collected, achievements_count, created_at, updated_at
		FROM user_stats 
		WHERE user_id = $1`
	
	stats := &domain.UserStats{}
	
	err := r.db.QueryRow(query, userID).Scan(
		&stats.UserID,
		&stats.GamesPlayed,
		&stats.GamesWon,
		&stats.TotalPlayTime,
		&stats.HighestLevel,
		&stats.CardsCollected,
		&stats.AchievementsCount,
		&stats.CreatedAt,
		&stats.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return stats, nil
}

func (r *UserRepository) UpdateStats(stats *domain.UserStats) error {
	query := `
		UPDATE user_stats 
		SET games_played = $2, games_won = $3, total_play_time = $4, highest_level = $5,
			cards_collected = $6, achievements_count = $7, updated_at = $8
		WHERE user_id = $1`
	
	stats.UpdatedAt = time.Now()
	
	_, err := r.db.Exec(
		query,
		stats.UserID,
		stats.GamesPlayed,
		stats.GamesWon,
		stats.TotalPlayTime,
		stats.HighestLevel,
		stats.CardsCollected,
		stats.AchievementsCount,
		stats.UpdatedAt,
	)
	
	return err
}

func (r *UserRepository) IncrementGamesPlayed(userID int) error {
	query := `
		UPDATE user_stats 
		SET games_played = games_played + 1, updated_at = $2
		WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID, time.Now())
	return err
}

func (r *UserRepository) IncrementGamesWon(userID int) error {
	query := `
		UPDATE user_stats 
		SET games_won = games_won + 1, updated_at = $2
		WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID, time.Now())
	return err
}

func (r *UserRepository) AddPlayTime(userID int, seconds int) error {
	query := `
		UPDATE user_stats 
		SET total_play_time = total_play_time + $2, updated_at = $3
		WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID, seconds, time.Now())
	return err
}

func (r *UserRepository) createStats(stats *domain.UserStats) error {
	query := `
		INSERT INTO user_stats (user_id, games_played, games_won, total_play_time, 
								highest_level, cards_collected, achievements_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`
	
	now := time.Now()
	err := r.db.QueryRow(
		query,
		stats.UserID,
		0, // games_played
		0, // games_won
		0, // total_play_time
		0, // highest_level
		0, // cards_collected
		0, // achievements_count
		now,
		now,
	).Scan(&stats.CreatedAt, &stats.UpdatedAt)
	
	return err
}