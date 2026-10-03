ALTER TABLE api_keys
ADD COLUMN name VARCHAR(100);

UPDATE api_keys
SET name = 'default';

ALTER TABLE api_keys
ALTER COLUMN name SET NOT NULL;
