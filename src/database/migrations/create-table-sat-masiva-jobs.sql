CREATE TABLE IF NOT EXISTS sat_masiva_jobs (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    year INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    month VARCHAR(40),
    message TEXT,
    error TEXT,
    logs JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP,
    CONSTRAINT fk_sat_masiva_jobs_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sat_masiva_jobs_user_id ON sat_masiva_jobs (user_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_sat_masiva_jobs_one_active
    ON sat_masiva_jobs (user_id)
    WHERE status IN ('queued', 'running');
