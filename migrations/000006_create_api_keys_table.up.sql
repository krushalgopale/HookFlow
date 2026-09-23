CREATE TABLE api_keys (
    id VARCHAR(100) PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_api_keys_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
);
