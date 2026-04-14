alter table devices add column if not exists last_status jsonb not null default '{}';
