CREATE TABLE currencies (
  id                smallserial         PRIMARY KEY,
  currency_code     varchar(10)         NOT NULL
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