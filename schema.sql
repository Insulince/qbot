-- example table for testing purposes
CREATE TABLE IF NOT EXISTS tbl (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    data TEXT,
    created_at TEXT DEFAULT (datetime('now', 'utc'))
);
