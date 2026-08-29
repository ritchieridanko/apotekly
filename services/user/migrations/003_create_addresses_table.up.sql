CREATE TABLE addresses (
  id BIGSERIAL PRIMARY KEY,
  auth_id BIGINT NOT NULL,
  
  label VARCHAR NOT NULL, -- e.g. "Home", "Work", etc.
  recipient VARCHAR NOT NULL,
  phone VARCHAR NOT NULL,
  notes TEXT,
  is_primary BOOLEAN NOT NULL,
  
  country VARCHAR NOT NULL,
  subdivision_1 TEXT,
  subdivision_2 TEXT,
  subdivision_3 TEXT,
  subdivision_4 TEXT,
  street TEXT NOT NULL,
  postal_code VARCHAR NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  location GEOGRAPHY(Point, 4326),
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index records by auth
CREATE INDEX idx_addresses_auth ON addresses (auth_id);
