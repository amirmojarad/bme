-- +goose Up
create unique index uidx_user_username on users (username);
create unique index uidx_phone_number on users (phone_number);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
