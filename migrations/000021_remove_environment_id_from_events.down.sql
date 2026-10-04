ALTER TABLE events
DROP CONSTRAINT fk_events_environment;

ALTER TABLE events
DROP COLUMN environment_id;
