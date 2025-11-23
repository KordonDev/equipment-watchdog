-- Drop all tables in reverse order
PRAGMA foreign_keys = OFF;
DROP TABLE IF EXISTS webauthn_sessions;
DROP TABLE IF EXISTS glove_ids;
DROP TABLE IF EXISTS registration_codes;
DROP TABLE IF EXISTS changes;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS equipments;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS user_credentials;
DROP TABLE IF EXISTS users;
PRAGMA foreign_keys = ON;

