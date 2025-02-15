CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS tokens(
    token_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    access_token VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS merch(
    merch_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS  orders(
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    merch_id INTEGER REFERENCES merch(merch_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS  transactions(
    transaction_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    recipient_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    amount_coins INT NOT NULL
);

CREATE INDEX idx_users_name ON users(name);

CREATE INDEX idx_tokens_access_token ON tokens(access_token);

CREATE INDEX idx_orders_user_merch ON orders(user_id, merch_id);

CREATE INDEX idx_merch_name ON merch(name);

CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_recipient ON transactions(recipient_id);


INSERT INTO merch(name, price)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);