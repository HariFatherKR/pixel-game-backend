package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/pixel-game/internal/database"
	"github.com/yourusername/pixel-game/internal/game/rewards"
)

// RewardRepositoryImpl PostgreSQL 보상 리포지토리 구현
type RewardRepositoryImpl struct {
	db *database.DB
}

// NewRewardRepository 새로운 보상 리포지토리 생성
func NewRewardRepository(db *database.DB) *RewardRepositoryImpl {
	return &RewardRepositoryImpl{
		db: db,
	}
}

// SaveRewardBundle 보상 묶음 저장
func (r *RewardRepositoryImpl) SaveRewardBundle(sessionID string, bundle *rewards.RewardBundle) error {
	baseRewardsJSON, err := json.Marshal(bundle.BaseRewards)
	if err != nil {
		return fmt.Errorf("기본 보상 직렬화 실패: %w", err)
	}

	choiceRewardsJSON, err := json.Marshal(bundle.ChoiceRewards)
	if err != nil {
		return fmt.Errorf("선택 보상 직렬화 실패: %w", err)
	}

	query := `
		INSERT INTO reward_bundles (
			id, session_id, source_type, source_id, floor_number,
			base_rewards, choice_rewards, is_completed, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`

	now := time.Now()
	_, err = r.db.Exec(query,
		bundle.ID, sessionID, bundle.SourceType, bundle.SourceID, bundle.FloorNumber,
		baseRewardsJSON, choiceRewardsJSON, bundle.IsCompleted, now, now,
	)

	if err != nil {
		return fmt.Errorf("보상 묶음 저장 실패: %w", err)
	}

	return nil
}

// GetRewardBundle 보상 묶음 조회
func (r *RewardRepositoryImpl) GetRewardBundle(sessionID string, bundleID string) (*rewards.RewardBundle, error) {
	query := `
		SELECT id, session_id, source_type, source_id, floor_number,
			   base_rewards, choice_rewards, is_completed, created_at, updated_at
		FROM reward_bundles 
		WHERE session_id = $1 AND id = $2`

	var bundle rewards.RewardBundle
	var baseRewardsJSON, choiceRewardsJSON []byte
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(query, sessionID, bundleID).Scan(
		&bundle.ID, &sessionID, &bundle.SourceType, &bundle.SourceID, &bundle.FloorNumber,
		&baseRewardsJSON, &choiceRewardsJSON, &bundle.IsCompleted, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("보상 묶음을 찾을 수 없습니다")
		}
		return nil, fmt.Errorf("보상 묶음 조회 실패: %w", err)
	}

	// JSON 역직렬화
	err = json.Unmarshal(baseRewardsJSON, &bundle.BaseRewards)
	if err != nil {
		return nil, fmt.Errorf("기본 보상 역직렬화 실패: %w", err)
	}

	err = json.Unmarshal(choiceRewardsJSON, &bundle.ChoiceRewards)
	if err != nil {
		return nil, fmt.Errorf("선택 보상 역직렬화 실패: %w", err)
	}

	return &bundle, nil
}

// GetPendingRewards 대기 중인 보상 목록
func (r *RewardRepositoryImpl) GetPendingRewards(sessionID string) ([]*rewards.RewardBundle, error) {
	query := `
		SELECT id, session_id, source_type, source_id, floor_number,
			   base_rewards, choice_rewards, is_completed, created_at, updated_at
		FROM reward_bundles 
		WHERE session_id = $1 AND is_completed = false
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("대기 중인 보상 조회 실패: %w", err)
	}
	defer rows.Close()

	var bundles []*rewards.RewardBundle

	for rows.Next() {
		bundle := &rewards.RewardBundle{}
		var baseRewardsJSON, choiceRewardsJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&bundle.ID, &sessionID, &bundle.SourceType, &bundle.SourceID, &bundle.FloorNumber,
			&baseRewardsJSON, &choiceRewardsJSON, &bundle.IsCompleted, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("보상 데이터 스캔 실패: %w", err)
		}

		// JSON 역직렬화
		err = json.Unmarshal(baseRewardsJSON, &bundle.BaseRewards)
		if err != nil {
			return nil, fmt.Errorf("기본 보상 역직렬화 실패: %w", err)
		}

		err = json.Unmarshal(choiceRewardsJSON, &bundle.ChoiceRewards)
		if err != nil {
			return nil, fmt.Errorf("선택 보상 역직렬화 실패: %w", err)
		}

		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

// MarkRewardCompleted 보상 완료 처리
func (r *RewardRepositoryImpl) MarkRewardCompleted(sessionID string, bundleID string) error {
	query := `
		UPDATE reward_bundles 
		SET is_completed = true, updated_at = $1
		WHERE session_id = $2 AND id = $3`

	result, err := r.db.Exec(query, time.Now(), sessionID, bundleID)
	if err != nil {
		return fmt.Errorf("보상 완료 처리 실패: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("보상 완료 결과 확인 실패: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("해당 보상을 찾을 수 없습니다")
	}

	return nil
}

// GetRewardHistory 보상 히스토리
func (r *RewardRepositoryImpl) GetRewardHistory(sessionID string) ([]*rewards.RewardBundle, error) {
	query := `
		SELECT id, session_id, source_type, source_id, floor_number,
			   base_rewards, choice_rewards, is_completed, created_at, updated_at
		FROM reward_bundles 
		WHERE session_id = $1 AND is_completed = true
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("보상 히스토리 조회 실패: %w", err)
	}
	defer rows.Close()

	var bundles []*rewards.RewardBundle

	for rows.Next() {
		bundle := &rewards.RewardBundle{}
		var baseRewardsJSON, choiceRewardsJSON []byte
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&bundle.ID, &sessionID, &bundle.SourceType, &bundle.SourceID, &bundle.FloorNumber,
			&baseRewardsJSON, &choiceRewardsJSON, &bundle.IsCompleted, &createdAt, &updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("보상 히스토리 스캔 실패: %w", err)
		}

		// JSON 역직렬화
		err = json.Unmarshal(baseRewardsJSON, &bundle.BaseRewards)
		if err != nil {
			return nil, fmt.Errorf("기본 보상 역직렬화 실패: %w", err)
		}

		err = json.Unmarshal(choiceRewardsJSON, &bundle.ChoiceRewards)
		if err != nil {
			return nil, fmt.Errorf("선택 보상 역직렬화 실패: %w", err)
		}

		bundles = append(bundles, bundle)
	}

	return bundles, nil
}

// SaveRewardSelection 보상 선택 내역 저장
func (r *RewardRepositoryImpl) SaveRewardSelection(sessionID string, bundleID string, selectedRewardIDs []string) error {
	selectionJSON, err := json.Marshal(selectedRewardIDs)
	if err != nil {
		return fmt.Errorf("선택 보상 직렬화 실패: %w", err)
	}

	query := `
		INSERT INTO reward_selections (
			session_id, bundle_id, selected_reward_ids, created_at
		) VALUES ($1, $2, $3, $4)`

	_, err = r.db.Exec(query, sessionID, bundleID, selectionJSON, time.Now())
	if err != nil {
		return fmt.Errorf("보상 선택 저장 실패: %w", err)
	}

	return nil
}

// GetRewardSelection 보상 선택 내역 조회
func (r *RewardRepositoryImpl) GetRewardSelection(sessionID string, bundleID string) ([]string, error) {
	query := `
		SELECT selected_reward_ids 
		FROM reward_selections 
		WHERE session_id = $1 AND bundle_id = $2`

	var selectionJSON []byte
	err := r.db.QueryRow(query, sessionID, bundleID).Scan(&selectionJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("보상 선택 내역을 찾을 수 없습니다")
		}
		return nil, fmt.Errorf("보상 선택 조회 실패: %w", err)
	}

	var selectedRewardIDs []string
	err = json.Unmarshal(selectionJSON, &selectedRewardIDs)
	if err != nil {
		return nil, fmt.Errorf("선택 보상 역직렬화 실패: %w", err)
	}

	return selectedRewardIDs, nil
}