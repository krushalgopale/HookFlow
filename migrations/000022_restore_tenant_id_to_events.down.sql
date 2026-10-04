ALTER TABLE events
ADD COLUMN tenant_id VARCHAR(100);

ALTER TABLE events
ADD CONSTRAINT fk_events_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id)
    ON DELETE CASCADE;
