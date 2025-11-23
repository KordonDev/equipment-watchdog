-- Initial schema
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  name TEXT NOT NULL UNIQUE,
  is_approved INTEGER NOT NULL DEFAULT 0,
  is_admin INTEGER NOT NULL DEFAULT 0,
  password TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS user_credentials (
  id BLOB PRIMARY KEY,
  user_id INTEGER NOT NULL,
  public_key BLOB NOT NULL,
  attestation_type TEXT,
  authenticator_aa_guid BLOB,
  authenticator_sign_count INTEGER NOT NULL DEFAULT 0,
  authenticator_clone_warning INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME,
  CONSTRAINT fk_user_credentials_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_credentials_user_id ON user_credentials(user_id);

CREATE TABLE IF NOT EXISTS members (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  name TEXT NOT NULL UNIQUE,
  "group" TEXT
);

CREATE TABLE IF NOT EXISTS equipments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  `type` TEXT,
  registration_code TEXT NOT NULL UNIQUE,
  member_id INTEGER,
  `size` TEXT,
  CONSTRAINT fk_equipments_member FOREIGN KEY(member_id) REFERENCES members(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_equipments_member_id ON equipments(member_id);
CREATE INDEX IF NOT EXISTS idx_equipments_type ON equipments(type);

CREATE TABLE IF NOT EXISTS orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  fulfilled_at DATETIME,
  `type` TEXT,
  member_id INTEGER,
  `size` TEXT,
  CONSTRAINT fk_orders_member FOREIGN KEY(member_id) REFERENCES members(id)
);
CREATE INDEX IF NOT EXISTS idx_orders_member_id ON orders(member_id);
CREATE INDEX IF NOT EXISTS idx_orders_fulfilled_at ON orders(fulfilled_at);

CREATE TABLE IF NOT EXISTS changes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  member_id INTEGER,
  equipment_id INTEGER,
  order_id INTEGER,
  `action` TEXT,
  user_id INTEGER,
  CONSTRAINT fk_changes_member FOREIGN KEY(member_id) REFERENCES members(id),
  CONSTRAINT fk_changes_equipment FOREIGN KEY(equipment_id) REFERENCES equipments(id),
  CONSTRAINT fk_changes_order FOREIGN KEY(order_id) REFERENCES orders(id),
  CONSTRAINT fk_changes_user FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_changes_member_id ON changes(member_id);
CREATE INDEX IF NOT EXISTS idx_changes_equipment_id ON changes(equipment_id);
CREATE INDEX IF NOT EXISTS idx_changes_order_id ON changes(order_id);
CREATE INDEX IF NOT EXISTS idx_changes_user_id ON changes(user_id);

CREATE TABLE IF NOT EXISTS registration_codes (
  id TEXT PRIMARY KEY,
  reserved_until DATETIME
);

CREATE TABLE IF NOT EXISTS glove_ids (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  glove_id TEXT NOT NULL UNIQUE,
  used INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS webauthn_sessions (
  username TEXT PRIMARY KEY,
  challenge TEXT,
  user_id BLOB,
  allowed_credential_ids BLOB,
  expires DATETIME,
  user_verification TEXT
);

