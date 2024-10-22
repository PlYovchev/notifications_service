-- Init schema and main notification table script

CREATE SCHEMA IF NOT EXISTS notifications_schema AUTHORIZATION postgres;

CREATE TABLE IF NOT EXISTS notifications_schema.notification (
    id SERIAL PRIMARY KEY,
    key TEXT,
    message TEXT NOT NULL,
    status TEXT NOT NULL,
    delivery_channel TEXT NOT NULL, 
    created_at TIMESTAMP default current_timestamp
)