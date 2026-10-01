ALTER TABLE tenants
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE  tenants
ADD CONSTRAINT fk_tenants_user
    FOREIGN KEY (user_id)
    REFERENCES users(id);
