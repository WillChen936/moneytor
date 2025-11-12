CREATE TABLE accounts (
  id            bigint          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name          text            NOT NULL,
  currency_id   smallint        NOT NULL,
  balance       numeric(19,6)   NOT NULL DEFAULT 0,
  created_at    timestamptz     NOT NULL DEFAULT NOW(),
  updated_at    timestamptz     NULL
);
ALTER TABLE accounts ADD FOREIGN KEY (currency_id) REFERENCES currencies (id);


CREATE TABLE entries (
  id            bigint          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name          text            NOT NULL,
  note          text            NOT NULL DEFAULT '',
  account_id    bigint          NOT NULL,
  category_id   bigint          NOT NULL,
  amount        numeric(19,6)   NOT NULL,
  created_at    timestamptz     NOT NULL DEFAULT NOW()
);
CREATE INDEX ON entries(account_id);
ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts (id);
ALTER TABLE entries ADD FOREIGN KEY (category_id) REFERENCES categories (id);