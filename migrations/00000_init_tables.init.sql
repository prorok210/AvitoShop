CREATE OR REPLACE FUNCTION update_updatedAt_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updatedAt = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance INT NOT NULL DEFAULT 1000,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens(
    token_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    access_token VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS merch(
    merch_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    price INT NOT NULL,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS  orders(
    order_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    merch_id INTEGER REFERENCES merch(merch_id) ON DELETE CASCADE,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS  transactions(
    transaction_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    recipient_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    amount_coins INT NOT NULL,
    createdAt TIMESTAMP DEFAULT NOW() NOT NULL,
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_users_updatedAt
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updatedAt_column();

CREATE TRIGGER update_tokens_updatedAt
BEFORE UPDATE ON tokens
FOR EACH ROW EXECUTE FUNCTION update_updatedAt_column();

CREATE TRIGGER update_merch_updatedAt
BEFORE UPDATE ON merch
FOR EACH ROW EXECUTE FUNCTION update_updatedAt_column();

CREATE TRIGGER update_orders_updatedAt
BEFORE UPDATE ON orders
FOR EACH ROW EXECUTE FUNCTION update_updatedAt_column();

CREATE TRIGGER update_transactions_updatedAt
BEFORE UPDATE ON transactions
FOR EACH ROW EXECUTE FUNCTION update_updatedAt_column();

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