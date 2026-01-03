CREATE TABLE credential_ids (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    credential_id TEXT NOT NULL,
    PRIMARY KEY (user_id, credential_id)
);

INSERT INTO credential_ids (user_id, credential_id)
    SELECT id, credential_id FROM users WHERE credential_id IS NOT NULL;

ALTER TABLE users
    DROP credential_id;

