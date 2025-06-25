package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/yourusername/pixel-game/internal/domain"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

// Card master data operations

func (r *CardRepository) GetAll(filter domain.CardFilter) ([]*domain.Card, error) {
	query := `
		SELECT id, name, type, rarity, cost, description, code_snippet, effects, visual_effects, created_at
		FROM cards
		WHERE 1=1`
	
	args := []interface{}{}
	argCounter := 1

	if filter.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argCounter)
		args = append(args, *filter.Type)
		argCounter++
	}

	if filter.Rarity != nil {
		query += fmt.Sprintf(" AND rarity = $%d", argCounter)
		args = append(args, *filter.Rarity)
		argCounter++
	}

	if filter.MinCost != nil {
		query += fmt.Sprintf(" AND cost >= $%d", argCounter)
		args = append(args, *filter.MinCost)
		argCounter++
	}

	if filter.MaxCost != nil {
		query += fmt.Sprintf(" AND cost <= $%d", argCounter)
		args = append(args, *filter.MaxCost)
		argCounter++
	}

	if filter.SearchTerm != nil && *filter.SearchTerm != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCounter, argCounter)
		searchPattern := "%" + *filter.SearchTerm + "%"
		args = append(args, searchPattern)
		argCounter++
	}

	query += " ORDER BY cost ASC, name ASC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]*domain.Card, 0)
	for rows.Next() {
		card := &domain.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Rarity,
			&card.Cost,
			&card.Description,
			&card.CodeSnippet,
			&card.Effects,
			&card.VisualEffects,
			&card.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *CardRepository) GetByID(id string) (*domain.Card, error) {
	query := `
		SELECT id, name, type, rarity, cost, description, code_snippet, effects, visual_effects, created_at
		FROM cards
		WHERE id = $1`

	card := &domain.Card{}
	err := r.db.QueryRow(query, id).Scan(
		&card.ID,
		&card.Name,
		&card.Type,
		&card.Rarity,
		&card.Cost,
		&card.Description,
		&card.CodeSnippet,
		&card.Effects,
		&card.VisualEffects,
		&card.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return card, nil
}

func (r *CardRepository) GetByIDs(ids []string) ([]*domain.Card, error) {
	if len(ids) == 0 {
		return []*domain.Card{}, nil
	}

	query := `
		SELECT id, name, type, rarity, cost, description, code_snippet, effects, visual_effects, created_at
		FROM cards
		WHERE id = ANY($1)
		ORDER BY cost ASC, name ASC`

	rows, err := r.db.Query(query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]*domain.Card, 0)
	for rows.Next() {
		card := &domain.Card{}
		err := rows.Scan(
			&card.ID,
			&card.Name,
			&card.Type,
			&card.Rarity,
			&card.Cost,
			&card.Description,
			&card.CodeSnippet,
			&card.Effects,
			&card.VisualEffects,
			&card.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *CardRepository) Create(card *domain.Card) error {
	query := `
		INSERT INTO cards (id, name, type, rarity, cost, description, code_snippet, effects, visual_effects, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at`

	err := r.db.QueryRow(
		query,
		card.ID,
		card.Name,
		card.Type,
		card.Rarity,
		card.Cost,
		card.Description,
		card.CodeSnippet,
		card.Effects,
		card.VisualEffects,
		time.Now(),
	).Scan(&card.CreatedAt)

	return err
}

func (r *CardRepository) Update(card *domain.Card) error {
	query := `
		UPDATE cards
		SET name = $2, type = $3, rarity = $4, cost = $5, description = $6, 
		    code_snippet = $7, effects = $8, visual_effects = $9
		WHERE id = $1`

	_, err := r.db.Exec(
		query,
		card.ID,
		card.Name,
		card.Type,
		card.Rarity,
		card.Cost,
		card.Description,
		card.CodeSnippet,
		card.Effects,
		card.VisualEffects,
	)

	return err
}

func (r *CardRepository) Delete(id string) error {
	query := `DELETE FROM cards WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// User card collection operations

func (r *CardRepository) GetUserCards(userID int) ([]*domain.UserCard, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.card_id, uc.acquired_at, uc.is_upgraded, uc.upgrade_path, uc.level,
		       c.id, c.name, c.type, c.rarity, c.cost, c.description, c.code_snippet, c.effects, c.visual_effects, c.created_at
		FROM user_cards uc
		INNER JOIN cards c ON uc.card_id = c.id
		WHERE uc.user_id = $1
		ORDER BY c.cost ASC, c.name ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userCards := make([]*domain.UserCard, 0)
	for rows.Next() {
		uc := &domain.UserCard{Card: &domain.Card{}}
		err := rows.Scan(
			&uc.ID,
			&uc.UserID,
			&uc.CardID,
			&uc.AcquiredAt,
			&uc.IsUpgraded,
			&uc.UpgradePath,
			&uc.Level,
			&uc.Card.ID,
			&uc.Card.Name,
			&uc.Card.Type,
			&uc.Card.Rarity,
			&uc.Card.Cost,
			&uc.Card.Description,
			&uc.Card.CodeSnippet,
			&uc.Card.Effects,
			&uc.Card.VisualEffects,
			&uc.Card.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		userCards = append(userCards, uc)
	}

	return userCards, nil
}

func (r *CardRepository) GetUserCard(userID int, cardID string) (*domain.UserCard, error) {
	query := `
		SELECT uc.id, uc.user_id, uc.card_id, uc.acquired_at, uc.is_upgraded, uc.upgrade_path, uc.level,
		       c.id, c.name, c.type, c.rarity, c.cost, c.description, c.code_snippet, c.effects, c.visual_effects, c.created_at
		FROM user_cards uc
		INNER JOIN cards c ON uc.card_id = c.id
		WHERE uc.user_id = $1 AND uc.card_id = $2`

	uc := &domain.UserCard{Card: &domain.Card{}}
	err := r.db.QueryRow(query, userID, cardID).Scan(
		&uc.ID,
		&uc.UserID,
		&uc.CardID,
		&uc.AcquiredAt,
		&uc.IsUpgraded,
		&uc.UpgradePath,
		&uc.Level,
		&uc.Card.ID,
		&uc.Card.Name,
		&uc.Card.Type,
		&uc.Card.Rarity,
		&uc.Card.Cost,
		&uc.Card.Description,
		&uc.Card.CodeSnippet,
		&uc.Card.Effects,
		&uc.Card.VisualEffects,
		&uc.Card.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return uc, nil
}

func (r *CardRepository) AddCardToUser(userCard *domain.UserCard) error {
	query := `
		INSERT INTO user_cards (user_id, card_id, acquired_at, is_upgraded, upgrade_path, level)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, acquired_at`

	err := r.db.QueryRow(
		query,
		userCard.UserID,
		userCard.CardID,
		time.Now(),
		false,
		"",
		1,
	).Scan(&userCard.ID, &userCard.AcquiredAt)

	if err != nil {
		return err
	}

	userCard.IsUpgraded = false
	userCard.Level = 1
	return nil
}

func (r *CardRepository) UpdateUserCard(userCard *domain.UserCard) error {
	query := `
		UPDATE user_cards
		SET is_upgraded = $3, upgrade_path = $4, level = $5
		WHERE user_id = $1 AND card_id = $2`

	_, err := r.db.Exec(
		query,
		userCard.UserID,
		userCard.CardID,
		userCard.IsUpgraded,
		userCard.UpgradePath,
		userCard.Level,
	)

	return err
}

func (r *CardRepository) RemoveCardFromUser(userID int, cardID string) error {
	query := `DELETE FROM user_cards WHERE user_id = $1 AND card_id = $2`
	_, err := r.db.Exec(query, userID, cardID)
	return err
}

// Deck operations

func (r *CardRepository) CreateDeck(deck *domain.Deck) error {
	query := `
		INSERT INTO decks (user_id, name, card_ids, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		deck.UserID,
		deck.Name,
		pq.Array(deck.CardIDs),
		deck.IsActive,
		now,
		now,
	).Scan(&deck.ID, &deck.CreatedAt, &deck.UpdatedAt)

	return err
}

func (r *CardRepository) GetUserDecks(userID int) ([]*domain.Deck, error) {
	query := `
		SELECT id, user_id, name, card_ids, is_active, created_at, updated_at
		FROM decks
		WHERE user_id = $1
		ORDER BY is_active DESC, created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	decks := make([]*domain.Deck, 0)
	for rows.Next() {
		deck := &domain.Deck{}
		err := rows.Scan(
			&deck.ID,
			&deck.UserID,
			&deck.Name,
			pq.Array(&deck.CardIDs),
			&deck.IsActive,
			&deck.CreatedAt,
			&deck.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}

	return decks, nil
}

func (r *CardRepository) GetDeck(deckID int) (*domain.Deck, error) {
	query := `
		SELECT id, user_id, name, card_ids, is_active, created_at, updated_at
		FROM decks
		WHERE id = $1`

	deck := &domain.Deck{}
	err := r.db.QueryRow(query, deckID).Scan(
		&deck.ID,
		&deck.UserID,
		&deck.Name,
		pq.Array(&deck.CardIDs),
		&deck.IsActive,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return deck, nil
}

func (r *CardRepository) UpdateDeck(deck *domain.Deck) error {
	query := `
		UPDATE decks
		SET name = $2, card_ids = $3, updated_at = $4
		WHERE id = $1`

	deck.UpdatedAt = time.Now()
	_, err := r.db.Exec(
		query,
		deck.ID,
		deck.Name,
		pq.Array(deck.CardIDs),
		deck.UpdatedAt,
	)

	return err
}

func (r *CardRepository) DeleteDeck(deckID int) error {
	query := `DELETE FROM decks WHERE id = $1`
	_, err := r.db.Exec(query, deckID)
	return err
}

func (r *CardRepository) SetActiveDeck(userID int, deckID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First, deactivate all user's decks
	_, err = tx.Exec(`UPDATE decks SET is_active = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Then activate the selected deck
	_, err = tx.Exec(`UPDATE decks SET is_active = true WHERE id = $1 AND user_id = $2`, deckID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *CardRepository) GetActiveDeck(userID int) (*domain.Deck, error) {
	query := `
		SELECT id, user_id, name, card_ids, is_active, created_at, updated_at
		FROM decks
		WHERE user_id = $1 AND is_active = true
		LIMIT 1`

	deck := &domain.Deck{}
	err := r.db.QueryRow(query, userID).Scan(
		&deck.ID,
		&deck.UserID,
		&deck.Name,
		pq.Array(&deck.CardIDs),
		&deck.IsActive,
		&deck.CreatedAt,
		&deck.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return deck, nil
}