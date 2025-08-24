CREATE TYPE transaction_type_enum AS ENUM ('income', 'expense', 'transfer');

CREATE TABLE categories (
    id smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    transaction_type_id transaction_type_enum NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE INDEX ON categories (name);