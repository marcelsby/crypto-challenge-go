CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) NOT NULL UNIQUE,
    user_document VARCHAR(500) NOT NULL,
    credit_card_token VARCHAR(500) NOT NULL,
    `value` DECIMAL(6, 2) NOT NULL
);