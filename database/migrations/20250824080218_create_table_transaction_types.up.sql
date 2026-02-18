CREATE TABLE transaction_types (
    id      smallint    PRIMARY KEY,
    name    text        NOT NULL
);

INSERT INTO transaction_types (
    id,
    name
)
VALUES
    (1, 'expense'),
    (2, 'income'),
    (3, 'transfer');