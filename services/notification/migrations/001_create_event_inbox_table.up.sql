CREATE TABLE event_inbox (
  id UUID PRIMARY KEY,

  topic TEXT NOT NULL,
  payload JSONB NOT NULL,

  received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ,

  retry_count INTEGER NOT NULL DEFAULT 0,
  next_retry_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  last_attempted_at TIMESTAMPTZ,
  last_error TEXT
);

-- Index records by next retry time and received time if not completed
CREATE INDEX idx_event_inbox_pending ON event_inbox (next_retry_at, received_at) WHERE completed_at IS NULL;
