USE application;


DROP TABLE IF EXISTS account;
CREATE TABLE account (
                                           id CHAR(36) NOT NULL PRIMARY KEY,                 -- store UUID as text like '550e8400-e29b-41d4-a716-446655440000'
                                        username VARCHAR(255) NOT NULL,

    -- SourceTracking fields
                                           sources VARCHAR(50),
                                           event VARCHAR(50),
                                           created_at BIGINT,
                                           updated_at BIGINT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE INDEX idx_a_id ON circle_credit_card_change (card_id);


