ALTER TABLE users
    ADD credential_id TEXT DEFAULT "";

-- move data from table to old column
INSERT INTO users (credential_id)
    SELECT credential_id FROM credential_ids WHERE user_id = users.id;

DROP TABLE IF EXISTS credential_ids;

