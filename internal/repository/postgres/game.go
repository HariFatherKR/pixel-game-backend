package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourusername/pixel-game/internal/domain"
)

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

// Session management

func (r *GameRepository) CreateSession(session *domain.GameSession) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	session.StartedAt = time.Now()
	session.LastActionAt = time.Now()

	query := `
		INSERT INTO game_sessions (
			id, user_id, status, game_mode, current_floor, current_turn,
			turn_phase, player_state, enemy_state, game_state, deck_snapshot,
			score, cards_played, damage_dealt, damage_taken,
			started_at, last_action_at, turn_time_limit, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)`

	_, err := r.db.Exec(query,
		session.ID,
		session.UserID,
		session.Status,
		session.GameMode,
		session.CurrentFloor,
		session.CurrentTurn,
		session.TurnPhase,
		session.PlayerState,
		session.EnemyState,
		session.GameState,
		pq.Array(session.DeckSnapshot),
		session.Score,
		session.CardsPlayed,
		session.DamageDealt,
		session.DamageTaken,
		session.StartedAt,
		session.LastActionAt,
		session.TurnTimeLimit,
		session.CreatedAt,
		session.UpdatedAt,
	)

	return err
}

func (r *GameRepository) GetSession(sessionID uuid.UUID) (*domain.GameSession, error) {
	query := `
		SELECT 
			id, user_id, status, game_mode, current_floor, current_turn,
			turn_phase, player_state, enemy_state, game_state, deck_snapshot,
			score, cards_played, damage_dealt, damage_taken,
			started_at, completed_at, last_action_at, turn_time_limit,
			created_at, updated_at
		FROM game_sessions
		WHERE id = $1`

	session := &domain.GameSession{}
	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.Status,
		&session.GameMode,
		&session.CurrentFloor,
		&session.CurrentTurn,
		&session.TurnPhase,
		&session.PlayerState,
		&session.EnemyState,
		&session.GameState,
		pq.Array(&session.DeckSnapshot),
		&session.Score,
		&session.CardsPlayed,
		&session.DamageDealt,
		&session.DamageTaken,
		&session.StartedAt,
		&session.CompletedAt,
		&session.LastActionAt,
		&session.TurnTimeLimit,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}

func (r *GameRepository) GetActiveSession(userID int) (*domain.GameSession, error) {
	query := `
		SELECT 
			id, user_id, status, game_mode, current_floor, current_turn,
			turn_phase, player_state, enemy_state, game_state, deck_snapshot,
			score, cards_played, damage_dealt, damage_taken,
			started_at, completed_at, last_action_at, turn_time_limit,
			created_at, updated_at
		FROM game_sessions
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1`

	session := &domain.GameSession{}
	err := r.db.QueryRow(query, userID, domain.GameStatusActive).Scan(
		&session.ID,
		&session.UserID,
		&session.Status,
		&session.GameMode,
		&session.CurrentFloor,
		&session.CurrentTurn,
		&session.TurnPhase,
		&session.PlayerState,
		&session.EnemyState,
		&session.GameState,
		pq.Array(&session.DeckSnapshot),
		&session.Score,
		&session.CardsPlayed,
		&session.DamageDealt,
		&session.DamageTaken,
		&session.StartedAt,
		&session.CompletedAt,
		&session.LastActionAt,
		&session.TurnTimeLimit,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}

func (r *GameRepository) UpdateSession(session *domain.GameSession) error {
	session.UpdatedAt = time.Now()
	session.LastActionAt = time.Now()

	query := `
		UPDATE game_sessions SET
			status = $2,
			current_floor = $3,
			current_turn = $4,
			turn_phase = $5,
			player_state = $6,
			enemy_state = $7,
			game_state = $8,
			score = $9,
			cards_played = $10,
			damage_dealt = $11,
			damage_taken = $12,
			last_action_at = $13,
			updated_at = $14
		WHERE id = $1`

	_, err := r.db.Exec(query,
		session.ID,
		session.Status,
		session.CurrentFloor,
		session.CurrentTurn,
		session.TurnPhase,
		session.PlayerState,
		session.EnemyState,
		session.GameState,
		session.Score,
		session.CardsPlayed,
		session.DamageDealt,
		session.DamageTaken,
		session.LastActionAt,
		session.UpdatedAt,
	)

	return err
}

func (r *GameRepository) EndSession(sessionID uuid.UUID, status domain.GameStatus) error {
	now := time.Now()
	query := `
		UPDATE game_sessions SET
			status = $2,
			completed_at = $3,
			updated_at = $4
		WHERE id = $1`

	_, err := r.db.Exec(query, sessionID, status, now, now)
	return err
}

// Game state

func (r *GameRepository) SaveGameState(sessionID uuid.UUID, playerState *domain.PlayerState, enemyState *domain.EnemyState, gameState *domain.GameState) error {
	playerJSON, err := json.Marshal(playerState)
	if err != nil {
		return fmt.Errorf("failed to marshal player state: %w", err)
	}

	enemyJSON, err := json.Marshal(enemyState)
	if err != nil {
		return fmt.Errorf("failed to marshal enemy state: %w", err)
	}

	gameJSON, err := json.Marshal(gameState)
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	query := `
		UPDATE game_sessions SET
			player_state = $2,
			enemy_state = $3,
			game_state = $4,
			last_action_at = $5,
			updated_at = $6
		WHERE id = $1`

	now := time.Now()
	_, err = r.db.Exec(query, sessionID, playerJSON, enemyJSON, gameJSON, now, now)
	return err
}

func (r *GameRepository) LoadGameState(sessionID uuid.UUID) (*domain.PlayerState, *domain.EnemyState, *domain.GameState, error) {
	query := `
		SELECT player_state, enemy_state, game_state
		FROM game_sessions
		WHERE id = $1`

	var playerJSON, enemyJSON, gameJSON json.RawMessage
	err := r.db.QueryRow(query, sessionID).Scan(&playerJSON, &enemyJSON, &gameJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}

	var playerState domain.PlayerState
	if err := json.Unmarshal(playerJSON, &playerState); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal player state: %w", err)
	}

	var enemyState domain.EnemyState
	if err := json.Unmarshal(enemyJSON, &enemyState); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal enemy state: %w", err)
	}

	var gameState domain.GameState
	if err := json.Unmarshal(gameJSON, &gameState); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	return &playerState, &enemyState, &gameState, nil
}

// Actions

func (r *GameRepository) RecordAction(action *domain.GameAction) error {
	action.ID = uuid.New()
	action.Timestamp = time.Now()

	query := `
		INSERT INTO game_actions (id, session_id, action_type, card_id, target_id, action_data, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(query,
		action.ID,
		action.SessionID,
		action.ActionType,
		action.CardID,
		action.TargetID,
		action.ActionData,
		action.Timestamp,
	)

	return err
}

func (r *GameRepository) GetSessionActions(sessionID uuid.UUID) ([]*domain.GameAction, error) {
	query := `
		SELECT id, session_id, action_type, card_id, target_id, action_data, timestamp
		FROM game_actions
		WHERE session_id = $1
		ORDER BY timestamp ASC`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actions := make([]*domain.GameAction, 0)
	for rows.Next() {
		action := &domain.GameAction{}
		err := rows.Scan(
			&action.ID,
			&action.SessionID,
			&action.ActionType,
			&action.CardID,
			&action.TargetID,
			&action.ActionData,
			&action.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// Statistics

func (r *GameRepository) GetUserGameStats(userID int) (*domain.UserGameStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_games,
			COUNT(CASE WHEN status = 'COMPLETED' THEN 1 END) as games_won,
			COUNT(CASE WHEN status = 'FAILED' THEN 1 END) as games_lost,
			COALESCE(MAX(current_floor), 0) as highest_floor,
			COALESCE(SUM(score), 0) as total_score,
			COALESCE(MAX(score), 0) as highest_score,
			COALESCE(SUM(EXTRACT(EPOCH FROM (completed_at - started_at))::INT), 0) as total_play_time
		FROM game_sessions
		WHERE user_id = $1 AND status IN ('COMPLETED', 'FAILED')`

	stats := &domain.UserGameStats{}
	var totalPlayTime int
	
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalGames,
		&stats.GamesWon,
		&stats.GamesLost,
		&stats.HighestFloor,
		&stats.TotalScore,
		&stats.HighestScore,
		&totalPlayTime,
	)
	
	if err != nil {
		return nil, err
	}

	stats.TotalPlayTime = totalPlayTime
	if stats.TotalGames > 0 {
		stats.WinRate = float64(stats.GamesWon) / float64(stats.TotalGames)
		stats.AverageGameTime = totalPlayTime / stats.TotalGames
	}

	// Get favorite cards
	cardQuery := `
		SELECT card_id, COUNT(*) as usage_count
		FROM game_actions
		WHERE session_id IN (
			SELECT id FROM game_sessions WHERE user_id = $1
		) AND action_type = 'PLAY_CARD' AND card_id IS NOT NULL
		GROUP BY card_id
		ORDER BY usage_count DESC
		LIMIT 5`

	rows, err := r.db.Query(cardQuery, userID)
	if err != nil {
		return stats, nil // Return stats even if favorite cards query fails
	}
	defer rows.Close()

	favoriteCards := make([]string, 0)
	for rows.Next() {
		var cardID string
		var count int
		if err := rows.Scan(&cardID, &count); err == nil {
			favoriteCards = append(favoriteCards, cardID)
		}
	}
	stats.FavoriteCards = favoriteCards

	return stats, nil
}

func (r *GameRepository) UpdateGameStats(sessionID uuid.UUID) error {
	// This is called when a game ends to update user statistics
	// The actual statistics are calculated on-demand in GetUserGameStats
	// This method could be used to update cached statistics if needed
	
	// For now, we just verify the session exists
	query := `SELECT 1 FROM game_sessions WHERE id = $1`
	var exists int
	err := r.db.QueryRow(query, sessionID).Scan(&exists)
	
	if err == sql.ErrNoRows {
		return fmt.Errorf("game session not found: %s", sessionID)
	}
	
	return err
}