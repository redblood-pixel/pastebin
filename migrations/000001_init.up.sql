BEGIN;

create table users(
    id serial primary key,
    name varchar(64) not null unique,
    email varchar(256) not null unique,
    password_hashed varchar(256) not null,
    created_at timestamptz not null default now(),
    last_login timestamptz not null default now()
);

DROP TYPE IF EXISTS access_type;
create type access_type as enum ('public', 'private');

create table pastes(
    id serial primary key,
    title varchar(256) not null,
    content_location varchar(256) not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz not null,
    visibility access_type not null, 
    last_visited timestamptz default now() not null,
    user_id integer references users(id)
);

COMMIT;