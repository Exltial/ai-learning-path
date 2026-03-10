-- Achievement System Database Schema
-- Created: 2026-03-10

-- Achievements table
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon_url VARCHAR(500),
    points INTEGER NOT NULL DEFAULT 0,
    achievement_type VARCHAR(50) NOT NULL,
    tier VARCHAR(50) NOT NULL DEFAULT 'bronze',
    criteria JSONB NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- User achievements (unlocked achievements)
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    earned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_notified BOOLEAN NOT NULL DEFAULT false,
    UNIQUE(user_id, achievement_id)
);

-- User levels and points
CREATE TABLE IF NOT EXISTS user_levels (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    level INTEGER NOT NULL DEFAULT 1,
    current_points INTEGER NOT NULL DEFAULT 0,
    total_points INTEGER NOT NULL DEFAULT 0,
    experience INTEGER NOT NULL DEFAULT 0,
    next_level_exp INTEGER NOT NULL DEFAULT 100,
    title VARCHAR(100) NOT NULL DEFAULT '初学者',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Points transactions
CREATE TABLE IF NOT EXISTS points_transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount INTEGER NOT NULL,
    balance_after INTEGER NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_id UUID,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- User learning streaks
CREATE TABLE IF NOT EXISTS user_streaks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_activity_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Level titles
CREATE TABLE IF NOT EXISTS level_titles (
    id UUID PRIMARY KEY,
    level INTEGER NOT NULL UNIQUE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500)
);

-- Daily challenges
CREATE TABLE IF NOT EXISTS daily_challenges (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    challenge_type VARCHAR(50) NOT NULL,
    target_count INTEGER NOT NULL,
    reward_points INTEGER NOT NULL,
    date DATE NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- User daily challenge progress
CREATE TABLE IF NOT EXISTS user_daily_challenge_progress (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    challenge_id UUID NOT NULL REFERENCES daily_challenges(id) ON DELETE CASCADE,
    current_count INTEGER NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, challenge_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_user_achievements_achievement_id ON user_achievements(achievement_id);
CREATE INDEX IF NOT EXISTS idx_user_levels_total_points ON user_levels(total_points DESC);
CREATE INDEX IF NOT EXISTS idx_points_transactions_user_id ON points_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_points_transactions_created_at ON points_transactions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_streaks_current_streak ON user_streaks(current_streak DESC);
CREATE INDEX IF NOT EXISTS idx_achievements_enabled ON achievements(is_enabled) WHERE is_enabled = true;
CREATE INDEX IF NOT EXISTS idx_achievements_type ON achievements(achievement_type);

-- Insert default level titles
INSERT INTO level_titles (id, level, title, description) VALUES
    (gen_random_uuid(), 1, '初学者', '刚入门的学习者'),
    (gen_random_uuid(), 2, '学习者', '开始系统学习'),
    (gen_random_uuid(), 3, '进阶者', '不断进步中'),
    (gen_random_uuid(), 4, '熟练者', '已经熟练掌握'),
    (gen_random_uuid(), 5, '专家', '领域专家'),
    (gen_random_uuid(), 6, '高手', '顶尖高手'),
    (gen_random_uuid(), 7, '大师', '一代大师'),
    (gen_random_uuid(), 8, '宗师', '开宗立派'),
    (gen_random_uuid(), 9, '传奇', '传奇人物'),
    (gen_random_uuid(), 10, '神话', '不朽神话')
ON CONFLICT (level) DO NOTHING;

-- Comments
COMMENT ON TABLE achievements IS '成就徽章定义表';
COMMENT ON TABLE user_achievements IS '用户已解锁成就表';
COMMENT ON TABLE user_levels IS '用户等级和积分表';
COMMENT ON TABLE points_transactions IS '积分交易记录表';
COMMENT ON TABLE user_streaks IS '用户学习连续天数表';
COMMENT ON TABLE level_titles IS '等级称号表';
COMMENT ON TABLE daily_challenges IS '每日挑战表';
COMMENT ON TABLE user_daily_challenge_progress IS '用户每日挑战进度表';
