BEGIN;

create table users(
    id serial primary key,
    name varchar(64) not null,
    email varchar(256) not null,
    password_hashed varchar(256) not null,
    created_at timestamp not null default now(),
    last_login timestamp not null default now()
);

create type access_type as enum ('public', 'private');

create table pastes(
    id serial primary key,
    title varchar(256) not null,
    content_location varchar(256) not null,
    created_at timestamp not null default now(),
    expires_at timestamp not null,
    visibility access_type not null, 
    last_visited timestamp default now(),
    user_id integer references users(id)
);

COMMIT;