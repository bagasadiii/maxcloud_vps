CREATE TABLE IF NOT EXISTS clients (
  client_id UUID PRIMARY KEY,
  email VARCHAR(50) UNIQUE NOT NULL,
  balance INT DEFAULT 0,
  suspended BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS billings ( 
  billing_id UUID PRIMARY KEY,
  client_id UUID NOT NULL,
  cpu INT NOT NULL,
  ram INT NOT NULL,
  storage INT NOT NULL,
  monthly_fee INT NOT NULL,
  cost_per_hour INT NOT NULL,
  total_fee INT NOT NULL,
  uptime INT NOT NULL,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  CONSTRAINT fk_billing_client FOREIGN KEY (client_id) REFERENCES clients(client_id) ON DELETE CASCADE
);

