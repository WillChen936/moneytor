CREATE TABLE accounts (
  id            bigint          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  owner         text            NOT NULL,
  currency_id   smallint        NOT NULL,
  balance       numeric(19,6)   NOT NULL DEFAULT 0,
  created_at    timestamptz     NOT NULL DEFAULT (now())
);
CREATE INDEX ON accounts(owner);
ALTER TABLE accounts ADD FOREIGN KEY (currency_id) REFERENCES currencies (id);


CREATE TABLE entries (
  id            bigint          PRIMARY KEY,
  account_id    bigint          NOT NULL,
  amount        numeric(19,6)   NOT NULL,
  created_at    timestamptz     NOT NULL DEFAULT (now())
);
CREATE INDEX ON entries(account_id);
ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts (id);


CREATE TABLE transfers (
  id                bigint          PRIMARY KEY,
  from_account_id   bigint          NOT NULL,
  to_account_id     bigint          NOT NULL,
  amount            numeric(19,6)   NOT NULL,
  created_at        timestamptz     NOT NULL DEFAULT (now())
);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);
ALTER TABLE transfers ADD CONSTRAINT transfers_amount_should_be_positive CHECK (amount > 0);
ALTER TABLE transfers ADD FOREIGN KEY (from_account_id) REFERENCES accounts (id);
ALTER TABLE transfers ADD FOREIGN KEY (to_account_id) REFERENCES accounts (id);