-- Indexes for performance and multi-tenancy isolation
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_locations_tenant_id ON locations(tenant_id);
CREATE INDEX idx_locations_user_id ON locations(user_id);
CREATE INDEX idx_streams_tenant_id ON streams(tenant_id);
