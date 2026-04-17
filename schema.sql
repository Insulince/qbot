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
    guild_id TEXT NOT NULL DEFAULT 'UNKNOWN',
    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    UNIQUE (tournament_id, user_id)
);

-- Migrated tournament_entries to support display_name field.
-- ALTER TABLE tournament_entries ADD COLUMN display_name TEXT NOT NULL DEFAULT 'MISSING';

-- Migrated tournament_entries to support guild_id field.
-- ALTER TABLE tournament_entries ADD COLUMN guild_id TEXT NOT NULL DEFAULT 'MISSING';

CREATE INDEX IF NOT EXISTS idx_tournament_entries_tournament_id ON tournament_entries(tournament_id);
CREATE INDEX IF NOT EXISTS idx_tournament_entries_user_id ON tournament_entries(user_id);
