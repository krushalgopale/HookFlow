ALTER TABLE events
ADD CONSTRAINT fk_events_tenant
FOREIGN KEY (tenant_id)
REFERENCES tenants(id);
