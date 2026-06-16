CREATE TABLE sessions (
  id BIGSERIAL PRIMARY KEY,
  parent_id BIGINT,
  auth_id BIGINT NOT NULL,

  refresh_token VARCHAR UNIQUE NOT NULL,
  ip_address TEXT NOT NULL,
  user_agent TEXT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,

  FOREIGN KEY (parent_id) REFERENCES sessions (id) ON DELETE CASCADE,
  FOREIGN KEY (auth_id) REFERENCES auth (id) ON DELETE CASCADE
);

-- Index records by auth if active (not revoked)
CREATE INDEX idx_sessions_auth ON sessions (auth_id) WHERE revoked_at IS NULL;

-- Index records by refresh token if active (not revoked)
CREATE INDEX idx_sessions_refresh_token ON sessions (refresh_token) WHERE revoked_at IS NULL;
