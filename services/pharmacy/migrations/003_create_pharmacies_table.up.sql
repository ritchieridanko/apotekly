CREATE TABLE pharmacies (
  id UUID PRIMARY KEY,
  auth_id BIGINT NOT NULL,

  -- Identity
  name VARCHAR NOT NULL,
  legal_name VARCHAR, -- registered business/legal name
  description TEXT,

  -- Status and Activity
  status VARCHAR NOT NULL DEFAULT 'pending', -- opts: "pending", "accepted", "rejected", "suspended", etc.
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  online_hours JSONB, -- e.g. { "mon": ["08:00-20:00"], "fri": ["08:00-12:00", "14:00-18:00"] }

  -- Address
  country VARCHAR NOT NULL,
  subdivision_1 TEXT,
  subdivision_2 TEXT,
  subdivision_3 TEXT,
  subdivision_4 TEXT,
  street TEXT NOT NULL,
  postal_code VARCHAR NOT NULL,

  -- Geolocation
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  location GEOGRAPHY(Point, 4326),

  -- Contact
  email VARCHAR,
  phone VARCHAR,
  website VARCHAR,
  whatsapp VARCHAR,

  -- Media
  profile_picture VARCHAR,
  profile_banner VARCHAR,

  -- Verification
  verified_at TIMESTAMPTZ,
  verified_by UUID,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Set auth unique for active (not deleted) records
CREATE UNIQUE INDEX idx_pharmacies_unique_auth ON pharmacies (auth_id) WHERE deleted_at IS NULL;

-- Index records by name if active (not deleted)
CREATE INDEX idx_pharmacies_name ON pharmacies USING GIN (name gin_trgm_ops) WHERE deleted_at IS NULL;

-- Index records by location if active (not deleted)
CREATE INDEX idx_pharmacies_location ON pharmacies USING GIST (location) WHERE is_active = TRUE AND deleted_at IS NULL;
