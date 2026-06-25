-- +goose Up
ALTER TABLE projects ADD COLUMN IF NOT EXISTS access_level VARCHAR(16) NOT NULL DEFAULT 'private';
UPDATE projects SET access_level = 'public' WHERE is_public = true AND access_level = 'private';
CREATE INDEX IF NOT EXISTS idx_projects_access_level ON projects (access_level);

-- +goose Down
DROP INDEX IF EXISTS idx_projects_access_level;
ALTER TABLE projects DROP COLUMN IF EXISTS access_level;
