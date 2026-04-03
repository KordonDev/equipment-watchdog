INSERT OR IGNORE INTO glove_ids (glove_id, used, created_at)
SELECT registration_code, 1, created_at
FROM equipments
WHERE type = 'gloves';
