-- +goose Up
alter table user_troubleshooting_sessions
    add column if not exists current_troubleshooting_step_id bigint;

alter table user_troubleshooting_sessions
    add constraint fk_user_troubleshooting_sessions_step foreign key (current_troubleshooting_step_id) references troubleshooting_steps (id);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
