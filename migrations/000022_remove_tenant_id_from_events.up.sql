ALTER TABLE events
DROP CONSTRAINT IF EXISTS fk_events_tenant;

ALTER TABLE events
DROP COLUMN tenant_id;
