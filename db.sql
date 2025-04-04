-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    streak INTEGER DEFAULT 0,
    current_rank TEXT,
    last_log_time TIMESTAMP
);

-- Create rank history table
CREATE TABLE IF NOT EXISTS rank_history (
    user_id TEXT REFERENCES users(id),
    rank TEXT NOT NULL,
    achieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create streak history table
CREATE TABLE IF NOT EXISTS streak_history (
    user_id TEXT REFERENCES users(id),
    streak INTEGER NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create hall of gooners table
CREATE TABLE IF NOT EXISTS hall_of_gooners (
    user_id TEXT PRIMARY KEY REFERENCES users(id),
    last_streak INTEGER NOT NULL,
    best_streak INTEGER NOT NULL,
    relapse_count INTEGER DEFAULT 1,
    last_relapse_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_rank_history_user_id ON rank_history(user_id);
CREATE INDEX IF NOT EXISTS idx_streak_history_user_id ON streak_history(user_id);
CREATE INDEX IF NOT EXISTS idx_users_streak ON users(streak DESC);
CREATE INDEX IF NOT EXISTS idx_users_last_log_time ON users(last_log_time);

-- Add helpful comments
COMMENT ON TABLE users IS 'Stores user information and current streaks';
COMMENT ON TABLE rank_history IS 'Tracks rank progression history for users';
COMMENT ON TABLE streak_history IS 'Tracks streak history for users';
COMMENT ON TABLE hall_of_gooners IS 'Tracks users who have relapsed';

-- Grant necessary permissions (adjust as needed)
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO your_bot_user;