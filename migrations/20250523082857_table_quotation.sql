-- +goose Up
CREATE TABLE IF NOT EXISTS quotation_s (
    id serial PRIMARY KEY,
    quotation VARCHAR(10) NOT NULL UNIQUE,
    price_24h NUMERIC(20, 8) NOT NULL CHECK (price_24h >= 0),
    volume_24h NUMERIC(20, 8) NOT NULL CHECK (volume_24h >= 0),
    last_trade_price NUMERIC(20, 8) NOT NULL CHECK (last_trade_price >= 0)
);

-- +goose Down
DROP TABLE IF EXISTS quotation_s;