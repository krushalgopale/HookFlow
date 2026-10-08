CREATE TABLE delivery_attempts (
  id VARCHAR(100) PRIMARY KEY,
  delivery_id VARCHAR(100) NOT NULL,
  attempt_number INTEGER NOT NULL,
  status VARCHAR(50) NOT NULL,
  response_status INTEGER,
  response_body TEXT,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_delivery_attempts_delivery
    FOREIGN KEY (delivery_id)
    REFERENCES deliveries(id)
    ON DELETE CASCADE
);
