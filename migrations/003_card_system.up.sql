-- Drop existing tables if they exist from previous schema
DROP TABLE IF EXISTS decks CASCADE;
DROP TABLE IF EXISTS user_cards CASCADE;
DROP TABLE IF EXISTS cards CASCADE;

-- Create cards table for master card data
CREATE TABLE cards (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('ACTION', 'EVENT', 'POWER')),
    rarity VARCHAR(20) NOT NULL CHECK (rarity IN ('COMMON', 'RARE', 'EPIC', 'LEGENDARY')),
    cost INTEGER NOT NULL CHECK (cost >= 0),
    description TEXT NOT NULL,
    code_snippet TEXT NOT NULL,
    effects JSONB NOT NULL DEFAULT '[]'::jsonb,
    visual_effects JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create user_cards table for user's card collection
CREATE TABLE user_cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    card_id VARCHAR(50) NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    acquired_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_upgraded BOOLEAN NOT NULL DEFAULT false,
    upgrade_path VARCHAR(50) DEFAULT '',
    level INTEGER NOT NULL DEFAULT 1 CHECK (level >= 1 AND level <= 5),
    UNIQUE(user_id, card_id)
);

-- Create decks table for user's deck management
CREATE TABLE decks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    card_ids VARCHAR(50)[] NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_cards_type ON cards(type);
CREATE INDEX idx_cards_rarity ON cards(rarity);
CREATE INDEX idx_cards_cost ON cards(cost);
CREATE INDEX idx_cards_name ON cards(name);
CREATE INDEX idx_user_cards_user_id ON user_cards(user_id);
CREATE INDEX idx_user_cards_card_id ON user_cards(card_id);
CREATE INDEX idx_decks_user_id ON decks(user_id);
CREATE INDEX idx_decks_is_active ON decks(user_id, is_active);

-- Add update trigger for decks
CREATE TRIGGER update_decks_updated_at BEFORE UPDATE ON decks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();