-- Schema for git-resume cache database

-- Processed commits cache
CREATE TABLE IF NOT EXISTS processed_commits (
    hash TEXT PRIMARY KEY,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analysis results
CREATE TABLE IF NOT EXISTS analysis_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    commit_hash TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    project TEXT NOT NULL,
    category TEXT NOT NULL,
    impact_summary TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (commit_hash) REFERENCES processed_commits(hash)
);

-- Token usage tracking
CREATE TABLE IF NOT EXISTS token_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    batch_id TEXT NOT NULL,
    input_tokens INTEGER NOT NULL,
    output_tokens INTEGER NOT NULL,
    cost_estimate REAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_results_date ON analysis_results(date);
CREATE INDEX IF NOT EXISTS idx_results_project ON analysis_results(project);
CREATE INDEX IF NOT EXISTS idx_results_category ON analysis_results(category);
