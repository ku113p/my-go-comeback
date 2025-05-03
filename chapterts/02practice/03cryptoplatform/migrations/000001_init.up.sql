-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    telegram_chat_id BIGINT
);

-- Create token_price table
CREATE TABLE IF NOT EXISTS token_price (
    symbol TEXT NOT NULL PRIMARY KEY,
    price DOUBLE PRECISION NOT NULL,
    name TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL
);

-- Create notification table
CREATE TABLE IF NOT EXISTS notification (
    id UUID NOT NULL PRIMARY KEY,
    symbol TEXT NOT NULL,
    sign TEXT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    user_id UUID,
    CONSTRAINT notification_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);
