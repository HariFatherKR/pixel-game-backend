-- Drop triggers
DROP TRIGGER IF EXISTS update_user_progression_updated_at ON user_progression;
DROP TRIGGER IF EXISTS update_decks_updated_at ON decks;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_user_cards_user;
DROP INDEX IF EXISTS idx_daily_challenges_date;
DROP INDEX IF EXISTS idx_leaderboard_user;
DROP INDEX IF EXISTS idx_leaderboard_mode_score;
DROP INDEX IF EXISTS idx_game_sessions_game_mode;
DROP INDEX IF EXISTS idx_game_sessions_user_status;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_platform;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS user_progression;
DROP TABLE IF EXISTS daily_challenges;
DROP TABLE IF EXISTS leaderboard;
DROP TABLE IF EXISTS game_sessions;
DROP TABLE IF EXISTS decks;
DROP TABLE IF EXISTS user_cards;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";