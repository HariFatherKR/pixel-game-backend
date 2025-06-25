-- 보상 묶음 테이블
CREATE TABLE reward_bundles (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    source_type VARCHAR(20) NOT NULL, -- COMBAT, EVENT, BOSS
    source_id VARCHAR(100) NOT NULL,
    floor_number INTEGER NOT NULL,
    base_rewards JSONB NOT NULL DEFAULT '[]',
    choice_rewards JSONB NOT NULL DEFAULT '[]',
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 보상 선택 내역 테이블
CREATE TABLE reward_selections (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    bundle_id VARCHAR(36) NOT NULL,
    selected_reward_ids JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (bundle_id) REFERENCES reward_bundles(id) ON DELETE CASCADE
);

-- 인덱스 생성
CREATE INDEX idx_reward_bundles_session_id ON reward_bundles(session_id);
CREATE INDEX idx_reward_bundles_completed ON reward_bundles(is_completed);
CREATE INDEX idx_reward_bundles_floor ON reward_bundles(floor_number);
CREATE INDEX idx_reward_selections_session_bundle ON reward_selections(session_id, bundle_id);

-- 보상 통계 뷰
CREATE VIEW reward_stats AS
SELECT 
    session_id,
    COUNT(*) as total_bundles,
    COUNT(CASE WHEN is_completed THEN 1 END) as completed_bundles,
    AVG(floor_number) as avg_floor,
    MAX(floor_number) as max_floor
FROM reward_bundles
GROUP BY session_id;