CREATE TABLE products (
  id UUID PRIMARY KEY,
  pharmacy_id UUID NOT NULL,
  
  -- Identity
  brand_name VARCHAR NOT NULL,
  generic_name VARCHAR NOT NULL,
  registration_number VARCHAR NOT NULL,
  description TEXT,
  requires_rx BOOLEAN NOT NULL,

  -- Characteristics
  dosage_form VARCHAR NOT NULL, -- e.g. "Tablet", "Capsule", "Gel", "Cream", etc.
  strength VARCHAR, -- e.g. "500mg", "250mg/5mL", etc.
  pack_unit VARCHAR NOT NULL, -- e.g. "Box", "Bottle", etc.
  pack_size INT NOT NULL, -- e.g. "10 tablets", "100 mL", etc.
  height_cm NUMERIC NOT NULL,
  length_cm NUMERIC NOT NULL,
  width_cm NUMERIC NOT NULL,
  weight_g NUMERIC NOT NULL,

  -- Pharmacy-specific Information
  price NUMERIC NOT NULL,
  currency VARCHAR NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL,

  -- Manufacturer Information
  manufactured_by VARCHAR NOT NULL,
  manufactured_in VARCHAR NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Index records by pharmacy if active (not deleted)
CREATE INDEX idx_products_pharmacy ON products (pharmacy_id) WHERE deleted_at IS NULL AND is_active = TRUE;

-- Index records by brand name if active (not deleted)
CREATE INDEX idx_products_brand_name ON products USING GIN (brand_name gin_trgm_ops) WHERE deleted_at IS NULL AND is_active = TRUE;

-- Index records by generic name if active (not deleted)
CREATE INDEX idx_products_generic_name ON products USING GIN (generic_name gin_trgm_ops) WHERE deleted_at IS NULL AND is_active = TRUE;
