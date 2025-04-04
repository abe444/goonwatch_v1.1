ALTER TABLE users ADD COLUMN IF NOT EXISTS current_rank TEXT; 
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_log_time TIMESTAMP; 

CREATE TABLE IF NOT EXISTS hall_of_gooners (
    user_id TEXT PRIMARY KEY REFERENCES users(id),
    last_streak INT NOT NULL,
    best_streak INT NOT NULL,
    relapse_count INT DEFAULT 1,
    last_relapse_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 
