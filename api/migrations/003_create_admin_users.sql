create table if not exists admins (
    id            serial       primary key,
    email         text         not null unique,
    name          text         not null default '',
    active        boolean      not null default true,
    created_at    timestamp    not null default current_timestamp,
    updated_at    timestamp    not null default current_timestamp
);

create index idx_admins_email on admins (email);

-- seed the first admin (replace with your email before running)
insert into admins (email, name) values ('jaracahteague@gmail.com', 'jar');
