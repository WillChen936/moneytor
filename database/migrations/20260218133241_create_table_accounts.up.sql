CREATE TABLE accounts (
  id            bigint          GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name          text            NOT NULL,
  currency_id   smallint        NOT NULL,
  balance       bigint          NOT NULL DEFAULT 0,
  created_at    timestamptz     NOT NULL DEFAULT NOW(),
  updated_at    timestamptz     NULL
);
ALTER TABLE accounts ADD FOREIGN KEY (currency_id) REFERENCES currencies (id);