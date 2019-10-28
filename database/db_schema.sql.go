package database

// Schema contains the SQL DDL used in creating the database.
// Its contents are expected to be valid SQLite 3.
const Schema = `
-- The following is a SQLite schema.
-- It will be created in a database file called bigboofer_data.sqlite3 if it doesn't exist.
-- (See database/db_interface.go for implementation details)

CREATE TABLE IF NOT EXISTS whitelist (
    id INTEGER PRIMARY KEY,
    group_id INTEGER,
    user_id INTEGER,
    added_on DATETIME
);

CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY,
    group_id INTEGER,
    channel_id INTEGER,
    passphrase STRING
)
`
