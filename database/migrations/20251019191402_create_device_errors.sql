-- +goose Up
create table if not exists device_errors
(
    id          bigserial primary key,
    device_id   bigint       not null,
    title       varchar(255) not null,
    description text         default null,
    status      varchar(100) default 'active',
    created_by  bigint       not null,
    updated_by  bigint       not null,
    deleted_by      bigint       default null,
    created_at  timestamptz  default now(),
    updated_at  timestamptz  default now(),
    deleted_at  timestamptz  default null,

    constraint fk_device_errors_device foreign key (device_id) references devices (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
