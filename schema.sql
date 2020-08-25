CREATE TABLE users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  created DATETIME NOT NULL,
  updated DATETIME NOT NULL
);

CREATE TABLE emails (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  address TEXT UNIQUE NOT NULL,
  created DATETIME NOT NULL
);
