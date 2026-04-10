create table if not exists devices (
    id            serial       primary key,
    name          text         not null,
    api_key_hash  text         not null unique,
    location      text         not null default '',
    active        boolean      not null default true,
    created_at    timestamp    not null default current_timestamp,
    updated_at    timestamp    not null default current_timestamp
);

create index idx_devices_api_key_hash on devices (api_key_hash);
create index idx_devices_active on devices (active);
