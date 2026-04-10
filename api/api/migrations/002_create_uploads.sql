create table if not exists uploads (
    id              serial       primary key,
    device_id       integer      references devices(id),
    resource_type   text,
    resource_id     text,
    status          text         not null default 'pending',
    filename        text         not null,
    original_name   text         not null,
    mime_type       text         not null,
    size            bigint       not null default 0,
    url             text         not null default '',
    sort_order      integer      not null default 0,
    expires_at      timestamp,
    created_at      timestamp    not null default current_timestamp,
    updated_at      timestamp    not null default current_timestamp
);

create index idx_uploads_device_id on uploads (device_id);
create index idx_uploads_filename on uploads (filename);
create index idx_uploads_status on uploads (status);
create index idx_uploads_resource on uploads (resource_type, resource_id);
