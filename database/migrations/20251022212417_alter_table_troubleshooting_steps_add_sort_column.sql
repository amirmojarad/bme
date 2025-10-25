-- +goose Up
alter table troubleshooting_steps
    add column if not exists sort bigint default 0;
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
