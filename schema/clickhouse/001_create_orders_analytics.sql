CREATE TABLE IF NOT EXISTS orders_analytics
(
    order_id UInt64,
    user_id UInt64,
    total Float64,
    created_at DateTime
)
ENGINE = MergeTree
ORDER BY (created_at, user_id);
