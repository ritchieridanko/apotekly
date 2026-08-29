CREATE TABLE users (
  id UUID PRIMARY KEY,
  auth_id BIGINT NOT NULL,
  
  name VARCHAR NOT NULL,
  sex VARCHAR,
  birthdate DATE,
  phone VARCHAR,
  profile_picture VARCHAR,
  profile_banner VARCHAR,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

-- Set auth unique for active (not deleted) records
CREATE UNIQUE INDEX idx_users_unique_auth ON users (auth_id) WHERE deleted_at IS NULL;
