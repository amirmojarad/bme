-- +goose Up
create table if not exists user_troubleshooting_sessions
(
    id              bigserial primary key,
    user_id         bigint not null,
    device_id       bigint not null,
    device_error_id bigint not null,
    status          varchar(255) default 'active',
    created_at      timestamptz  default now(),
    deleted_at      timestamptz  default null,

    constraint fk_user_troubleshooting_sessions_to_users foreign key (user_id) references users (id),
    constraint fk_user_troubleshooting_sessions_to_devices foreign key (device_id) references devices (id),
    constraint fk_user_troubleshooting_sessions_to_device_errors foreign key (device_error_id) references device_errors (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_user_troubleshooting_sessions_active
    ON user_troubleshooting_sessions (user_id)
    WHERE status = 'active';
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
