ALTER TABLE api_keys
DROP CONSTRAINT fk_api_keys_environment;

ALTER TABLE api_keys 
DROP COLUMN environment_id;
