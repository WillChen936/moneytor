CREATE TABLE currencies (
  id    smallint  PRIMARY KEY,
  code  text      NOT NULL
);

INSERT INTO currencies (
    id,
    code
)
VALUES
    (1, 'TWD'),
    (2, 'CNY'),
    (3, 'USD'),
    (4, 'EUR'),
    (5, 'GBP');