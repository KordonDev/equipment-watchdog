CREATE TABLE glove_ids_new (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  glove_id   TEXT NOT NULL UNIQUE
);
INSERT INTO glove_ids_new (id, created_at, glove_id)
  SELECT id, created_at, glove_id FROM glove_ids WHERE used = 1;
DROP TABLE glove_ids;
ALTER TABLE glove_ids_new RENAME TO glove_ids;
