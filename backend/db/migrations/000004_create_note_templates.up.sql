CREATE TABLE IF NOT EXISTS note_templates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    default_title VARCHAR(200),
    content_template TEXT NOT NULL,
    default_tags VARCHAR(255),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_note_templates_user_id ON note_templates (user_id);
CREATE INDEX IF NOT EXISTS idx_note_templates_is_default ON note_templates (user_id, is_default);
