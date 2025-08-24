CREATE TABLE currencies (
  id              smallint  PRIMARY KEY,
  currency_code   text      NOT NULL
);

INSERT INTO currencies (
    id,
    currency_code
)
VALUES
    (1, 'TWD'),
    (2, 'CNY'),
    (3, 'USD'),
    (4, 'EUR'),
    (5, 'GBP');