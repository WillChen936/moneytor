CREATE TABLE entries (
  id                bigint          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name              text            NOT NULL,
  note              text            NOT NULL DEFAULT '',
  from_account_id   bigint          NOT NULL,
  to_account_id     bigint          NULL,
  category_id       bigint          NOT NULL,
  amount            bigint          NOT NULL,
  created_at        timestamptz     NOT NULL DEFAULT NOW()
);

CREATE INDEX ON entries(from_account_id);
CREATE INDEX ON entries(to_account_id);
ALTER TABLE entries ADD FOREIGN KEY (from_account_id) REFERENCES accounts (id);
ALTER TABLE entries ADD FOREIGN KEY (to_account_id) REFERENCES accounts (id);
ALTER TABLE entries ADD FOREIGN KEY (category_id) REFERENCES categories (id);