CREATE TABLE IF NOT EXISTS note_links (
    id SERIAL PRIMARY KEY,
    source_note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    target_note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_note_id, target_note_id)
);

CREATE INDEX IF NOT EXISTS idx_note_links_source ON note_links (source_note_id);
CREATE INDEX IF NOT EXISTS idx_note_links_target ON note_links (target_note_id);
