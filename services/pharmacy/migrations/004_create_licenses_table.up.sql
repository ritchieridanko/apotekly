CREATE TABLE licenses (
  id UUID PRIMARY KEY,
  pharmacy_id UUID NOT NULL,

  no VARCHAR NOT NULL,
  authority VARCHAR NOT NULL,
  issued_at DATE NOT NULL,
  expires_at DATE NOT NULL,
  attachment VARCHAR,
  status VARCHAR NOT NULL DEFAULT 'pending', -- opts: "pending", "active", "rejected", "expired", "revoked", etc.

  verified_at TIMESTAMPTZ,
  verified_by UUID,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  FOREIGN KEY (pharmacy_id) REFERENCES pharmacies (id) ON DELETE CASCADE
);

-- Index records by pharmacy if active (not deleted)
CREATE INDEX idx_licenses_pharmacy ON licenses (pharmacy_id) WHERE deleted_at IS NULL;

-- Set license (by no and authority) unique for active (not deleted) records
CREATE UNIQUE INDEX idx_licenses_unique ON licenses (no, authority) WHERE deleted_at IS NULL;
