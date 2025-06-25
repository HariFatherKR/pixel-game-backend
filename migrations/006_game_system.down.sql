-- Drop indexes
DROP INDEX IF EXISTS idx_game_actions_timestamp;
DROP INDEX IF EXISTS idx_game_actions_action_type;
DROP INDEX IF EXISTS idx_game_actions_session_id;

DROP INDEX IF EXISTS idx_game_sessions_created_at;
DROP INDEX IF EXISTS idx_game_sessions_game_mode;
DROP INDEX IF EXISTS idx_game_sessions_user_status;
DROP INDEX IF EXISTS idx_game_sessions_status;
DROP INDEX IF EXISTS idx_game_sessions_user_id;

-- Drop tables
DROP TABLE IF EXISTS game_actions CASCADE;
DROP TABLE IF EXISTS game_sessions CASCADE;