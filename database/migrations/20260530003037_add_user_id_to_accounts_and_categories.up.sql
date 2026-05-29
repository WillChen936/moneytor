ALTER TABLE accounts ADD COLUMN user_id bigint NOT NULL REFERENCES users (id);
ALTER TABLE categories ADD COLUMN user_id bigint NOT NULL REFERENCES users (id);
