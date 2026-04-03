DELETE FROM glove_ids
WHERE glove_id IN (
  SELECT registration_code FROM equipments WHERE type = 'gloves'
);
