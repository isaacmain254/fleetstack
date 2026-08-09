-- UP
CREATE TABLE IF NOT EXISTS projects (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  repo_url TEXT NOT NULL UNIQUE,
  branch TEXT NOT NULL DEFAULT 'main',
  root_directory TEXT NOT NULL DEFAULT '/home/cursor/.fleetstack/projects/',
  clone_path TEXT NOT NULL DEFAULT '/home/cursor/.fleetstack/projects/',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
