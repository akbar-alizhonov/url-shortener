create table if not exists url (
    id serial primary key,
    original_url text not null,
    alias text unique not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz,
    clicks bigint not null default 0
);
