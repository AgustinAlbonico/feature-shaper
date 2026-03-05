package db

const SchemaSQL = `
CREATE TABLE IF NOT EXISTS projects (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  slug      TEXT UNIQUE NOT NULL,
  name      TEXT NOT NULL,
  path      TEXT,
  createdAt TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS features (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  projectSlug    TEXT NOT NULL REFERENCES projects(slug),
  slug           TEXT NOT NULL,
  title          TEXT NOT NULL,
  type           TEXT NOT NULL,
  status         TEXT NOT NULL DEFAULT 'draft',
  content        TEXT NOT NULL,
  version        INTEGER NOT NULL DEFAULT 1,
  topicKey       TEXT UNIQUE,
  normalizedHash TEXT,
  createdAt      TEXT NOT NULL DEFAULT (datetime('now')),
  updatedAt      TEXT NOT NULL DEFAULT (datetime('now')),
  UNIQUE(projectSlug, slug)
);

CREATE TABLE IF NOT EXISTS featureVersions (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  featureId INTEGER NOT NULL REFERENCES features(id) ON DELETE CASCADE,
  version   INTEGER NOT NULL,
  content   TEXT NOT NULL,
  changelog TEXT,
  createdAt TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE VIRTUAL TABLE IF NOT EXISTS featuresFts USING fts5(
  title,
  content,
  type,
  status,
  content='features',
  content_rowid='id'
);

CREATE TRIGGER IF NOT EXISTS features_ai AFTER INSERT ON features BEGIN
  INSERT INTO featuresFts(rowid, title, content, type, status)
  VALUES (new.id, new.title, new.content, new.type, new.status);
END;

CREATE TRIGGER IF NOT EXISTS features_au AFTER UPDATE ON features BEGIN
  INSERT INTO featuresFts(featuresFts, rowid, title, content, type, status)
  VALUES ('delete', old.id, old.title, old.content, old.type, old.status);
  INSERT INTO featuresFts(rowid, title, content, type, status)
  VALUES (new.id, new.title, new.content, new.type, new.status);
END;

CREATE TRIGGER IF NOT EXISTS features_ad AFTER DELETE ON features BEGIN
  INSERT INTO featuresFts(featuresFts, rowid, title, content, type, status)
  VALUES ('delete', old.id, old.title, old.content, old.type, old.status);
END;
`
