CREATE TABLE user_balance (
  id BIGSERIAL PRIMARY KEY,
  current NUMERIC(10,2) DEFAULT 0,
  withdrawn NUMERIC(10,2) DEFAULT 0,
  user_id BIGINT REFERENCES users(id)
  );