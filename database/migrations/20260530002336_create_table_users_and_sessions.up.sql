CREATE TABLE users (
  id               bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username         text        NOT NULL,
  email            text        NOT NULL,
  hashed_password  text        NOT NULL,
  created_at       timestamptz NOT NULL DEFAULT NOW(),
  updated_at       timestamptz NULL
);

CREATE UNIQUE INDEX ON users (email);

CREATE TABLE sessions (
  id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        bigint      NOT NULL,
  refresh_token  text        NOT NULL,
  expires_at     timestamptz NOT NULL,
  created_at     timestamptz NOT NULL DEFAULT NOW()
);

ALTER TABLE sessions ADD FOREIGN KEY (user_id) REFERENCES users (id);
