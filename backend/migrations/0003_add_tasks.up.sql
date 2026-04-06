CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    member_id INTEGER,
    `group` TEXT,
    equipment_type TEXT,
    `type` TEXT,
    CONSTRAINT fk_tasks_member FOREIGN KEY(member_id) REFERENCES members(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_tasks_member_id ON tasks(member_id);
