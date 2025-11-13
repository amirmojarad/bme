-- +goose Up
alter table user_troubleshooting_sessions
    add finished_at timestamptz default null;
alter table user_troubleshooting_journey
    add finished_at timestamptz default null;
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
