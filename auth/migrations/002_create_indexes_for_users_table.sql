-- +goose Up
CREATE INDEX idx_users_login ON users(login);

CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_users_notification_methods ON users USING GIN (notification_methods);