-- +goose Up
create table if not exists troubleshooting_steps
(
    id              bigserial primary key,
    device_id       bigint       not null,
    device_error_id bigint       not null,
    title           varchar(255) not null,
    description     text         default null,
    hints           jsonb        default null,
    status          varchar(100) default 'active',
    created_by      bigint       not null,
    updated_by      bigint       not null,
    deleted_by      bigint       default null,
    created_at      timestamptz  default now(),
    updated_at      timestamptz  default now(),
    deleted_at      timestamptz  default null,

    constraint fk_troubleshooting_steps_device foreign key (device_id) references devices (id),
    constraint fk_troubleshooting_steps_device_errors foreign key (device_error_id) references device_errors (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
