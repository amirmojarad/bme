-- +goose Up
create table if not exists users
(
    id              bigserial primary key,
    username        varchar(255) not null,
    first_name      varchar(255) default '',
    last_name       varchar(255) default '',
    phone_number    varchar(255) default '',
    hashed_password varchar(255) not null,
    created_at      timestamptz  default now(),
    updated_at      timestamptz  default now(),
    deleted_at      timestamptz  default null
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
