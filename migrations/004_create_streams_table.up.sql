CREATE TABLE streams (
    id SERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
    thirdparty_status VARCHAR(50),
    streamed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
