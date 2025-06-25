-- Drop indexes
DROP INDEX IF EXISTS idx_decks_is_active;
DROP INDEX IF EXISTS idx_decks_user_id;
DROP INDEX IF EXISTS idx_user_cards_card_id;
DROP INDEX IF EXISTS idx_user_cards_user_id;
DROP INDEX IF EXISTS idx_cards_name;
DROP INDEX IF EXISTS idx_cards_cost;
DROP INDEX IF EXISTS idx_cards_rarity;
DROP INDEX IF EXISTS idx_cards_type;

-- Drop tables
DROP TABLE IF EXISTS decks CASCADE;
DROP TABLE IF EXISTS user_cards CASCADE;
DROP TABLE IF EXISTS cards CASCADE;