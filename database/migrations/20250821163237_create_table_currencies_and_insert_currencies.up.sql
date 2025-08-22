CREATE TABLE currencies (
  id              smallint  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  currency_code   text      NOT NULL
);

INSERT INTO currencies (
    currency_code
)
VALUES
    ('TWD'),
    ('CNY'),
    ('USD'),
    ('EUR'),
    ('GBP');