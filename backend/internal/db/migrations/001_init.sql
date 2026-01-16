-- 001_init.sql
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    username TEXT,
    first_name TEXT,
    language TEXT DEFAULT 'en',
    timezone TEXT DEFAULT 'UTC',
    theme_index INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    original_input TEXT,
    task_type TEXT NOT NULL,
    priority INTEGER DEFAULT 2,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    due_at TIMESTAMPTZ,
    reminder_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_user_type ON tasks(user_id, task_type, completed);
CREATE INDEX IF NOT EXISTS idx_tasks_reminder ON tasks(due_at) WHERE NOT reminder_sent AND NOT completed;
CREATE INDEX IF NOT EXISTS idx_events_user_date ON events(user_id, created_at);
