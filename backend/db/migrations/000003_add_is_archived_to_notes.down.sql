-- Remove index
DROP INDEX IF EXISTS idx_notes_archived;

-- Remove is_archived column from notes table
ALTER TABLE notes DROP COLUMN IF EXISTS is_archived;
