-- +goose Up
create table if not exists troubleshooting_steps_to_steps
(
    id           bigserial primary key,
    from_step_id bigint not null,
    to_step_id   bigint not null,
    priority     int         default 1,
    created_by   bigint not null,
    updated_by   bigint not null,
    deleted_by      bigint       default null,
    created_at   timestamptz default now(),
    updated_at   timestamptz default now(),
    deleted_at   timestamptz default null,

    constraint fk_troubleshooting_steps_to_steps_from_step_id foreign key (from_step_id) references troubleshooting_steps (id),
    constraint fk_troubleshooting_steps_to_steps_to_step_id foreign key (to_step_id) references troubleshooting_steps (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
