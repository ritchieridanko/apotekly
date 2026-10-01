CREATE TABLE oauth (
  id BIGSERIAL PRIMARY KEY,
  auth_id BIGINT NOT NULL,
  
  provider VARCHAR NOT NULL,
  subject VARCHAR NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  
  FOREIGN KEY (auth_id) REFERENCES auth (id) ON DELETE CASCADE
);

-- Set oauth (by auth and provider) unique for active (not deleted) records
CREATE UNIQUE INDEX idx_oauth_unique_auth_provider ON oauth (auth_id, provider) WHERE deleted_at IS NULL;

-- Set oauth (by provider and subject) unique for active (not deleted) records
CREATE UNIQUE INDEX idx_oauth_unique_provider_subject ON oauth (provider, subject) WHERE deleted_at IS NULL;
