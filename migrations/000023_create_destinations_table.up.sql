CREATE TABLE destinations (
  id VARCHAR(100) PRIMARY KEY,
  environment_id VARCHAR(100) NOT NULL,
  name VARCHAR(100) NOT NULL,
  url VARCHAR(500) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT fk_destinations_environment
      FOREIGN KEY (environment_id)
      REFERENCES environments(id)
      ON DELETE CASCADE
);
