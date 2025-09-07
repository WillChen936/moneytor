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

CREATE TABLE categories (
    id                      integer                 GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                    text                    NOT NULL,
    transaction_type_id     smallint                NOT NULL,
    created_at              timestamptz             NOT NULL DEFAULT NOW(),
    updated_at              timestamptz             NULL
);
CREATE INDEX ON categories (name);
ALTER TABLE categories ADD FOREIGN KEY (transaction_type_id) REFERENCES transaction_types (id);