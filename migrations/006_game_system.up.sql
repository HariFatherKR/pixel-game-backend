-- Create game sessions table
CREATE TABLE game_sessions (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'COMPLETED', 'FAILED', 'PAUSED')),
    game_mode VARCHAR(30) NOT NULL CHECK (game_mode IN ('STORY', 'DAILY_CHALLENGE', 'EVENT')),
    current_floor INTEGER NOT NULL DEFAULT 1,
    current_turn INTEGER NOT NULL DEFAULT 1,
    turn_phase VARCHAR(20) NOT NULL DEFAULT 'START' CHECK (turn_phase IN ('START', 'DRAW', 'MAIN', 'END', 'ENEMY')),
    player_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    enemy_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    game_state JSONB NOT NULL DEFAULT '{}'::jsonb,
    deck_snapshot TEXT[] NOT NULL DEFAULT '{}',
    score INTEGER NOT NULL DEFAULT 0,
    cards_played INTEGER NOT NULL DEFAULT 0,
    damage_dealt INTEGER NOT NULL DEFAULT 0,
    damage_taken INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    last_action_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    turn_time_limit INTEGER NOT NULL DEFAULT 120, -- seconds
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create game actions table for recording all player actions
CREATE TABLE game_actions (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES game_sessions(id) ON DELETE CASCADE,
    action_type VARCHAR(30) NOT NULL,
    card_id VARCHAR(50),
    target_id VARCHAR(50),
    action_data JSONB,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_game_sessions_user_id ON game_sessions(user_id);
CREATE INDEX idx_game_sessions_status ON game_sessions(status);
CREATE INDEX idx_game_sessions_user_status ON game_sessions(user_id, status);
CREATE INDEX idx_game_sessions_game_mode ON game_sessions(game_mode);
CREATE INDEX idx_game_sessions_created_at ON game_sessions(created_at);

CREATE INDEX idx_game_actions_session_id ON game_actions(session_id);
CREATE INDEX idx_game_actions_action_type ON game_actions(action_type);
CREATE INDEX idx_game_actions_timestamp ON game_actions(timestamp);

-- Add trigger for updating updated_at
CREATE TRIGGER update_game_sessions_updated_at BEFORE UPDATE ON game_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();