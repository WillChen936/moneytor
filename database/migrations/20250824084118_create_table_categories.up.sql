CREATE TABLE categories (
    id                      bigint                 GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                    text                    NOT NULL,
    transaction_type_id     smallint                NOT NULL,
    created_at              timestamptz             NOT NULL DEFAULT NOW(),
    updated_at              timestamptz             NULL
);
CREATE INDEX ON categories (name);
ALTER TABLE categories ADD FOREIGN KEY (transaction_type_id) REFERENCES transaction_types (id);