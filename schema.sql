CREATE TABLE IF NOT EXISTS tournaments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    short_name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS tournament_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tournament_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    waves INTEGER NOT NULL,
    display_name TEXT NOT NULL DEFAULT 'UNKNOWN',
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    UNIQUE (tournament_id, user_id)
);

-- Migrated tournament_entries to support display_name field.
-- ALTER TABLE tournament_entries ADD COLUMN display_name TEXT NOT NULL DEFAULT 'MISSING';
