ALTER TABLE api_keys
ADD COLUMN tenant_id VARCHAR(100);

ALTER TABLE api_keys
ADD CONSTRAINT fk_api_keys_tenant
    FOREIGN KEY (tenant_id)
    REFERENCES tenants(id);
