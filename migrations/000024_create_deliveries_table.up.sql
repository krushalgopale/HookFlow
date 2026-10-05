CREATE TABLE deliveries (
  id VARCHAR(100) PRIMARY KEY,
  event_id VARCHAR(100) NOT NULL,
  destination_id VARCHAR(100) NOT NULL,
  status VARCHAR(50) NOT NULL,
  response_status INTEGER,
  response_body TEXT,
  error TEXT,
  attempt_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_deliveries_event
    FOREIGN KEY (event_id)
    REFERENCES events(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_deliveries_destination
    FOREIGN KEY (destination_id)
    REFERENCES destinations(id)
    ON DELETE CASCADE
);
