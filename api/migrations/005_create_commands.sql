create table if not exists commands (
    id            serial       primary key,
    device_id     integer      not null references devices(id),
    action        text         not null,
    payload       jsonb        not null default '{}',
    status        text         not null default 'pending',
    created_at    timestamp    not null default current_timestamp,
    updated_at    timestamp    not null default current_timestamp
);

create index idx_commands_device_id on commands (device_id);
create index idx_commands_device_pending on commands (device_id, status) where status = 'pending';
