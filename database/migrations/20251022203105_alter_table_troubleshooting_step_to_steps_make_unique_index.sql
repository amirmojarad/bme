-- +goose Up
create unique index uidx_from_step_id_to_step_id_troubleshooting_step_to_steps on troubleshooting_steps_to_steps (from_step_id, to_step_id);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
