-- +goose Up
create table if not exists user_troubleshooting_journey
(
    id                              bigserial primary key,
    user_troubleshooting_session_id bigint not null,
    from_troubleshooting_step_id    bigint not null,
    to_troubleshooting_step_id      bigint not null,
    description                     text        default null,
    created_at                      timestamptz default now(),

    constraint fk_user_troubleshooting_journey_to_user_ts_sess foreign key (user_troubleshooting_session_id) references users (id),
    constraint fk_user_troubleshooting_journey_to_to_ts_step foreign key (to_troubleshooting_step_id) references devices (id),
    constraint fk_user_troubleshooting_journey_to_from_user_ts_step foreign key (from_troubleshooting_step_id) references device_errors (id)
);

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
