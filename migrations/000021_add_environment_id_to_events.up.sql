ALTER TABLE events
ADD COLUMN environment_id VARCHAR(100) NOT NULL;

ALTER TABLE events
ADD CONSTRAINT fk_events_environment
    FOREIGN KEY (environment_id)
    REFERENCES environments(id)
    ON DELETE CASCADE;
