BEGIN;

create table users(
    id serial primary key,
    name varchar(63) not null unique check (
        char_length(name) between 6 and 63
    ),
    email varchar(255) not null unique check (
        char_length(email) between 6 and 255 and
        email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    password_hashed varchar(255) not null,
    created_at timestamptz not null default now(),
    last_login timestamptz not null default now()
);


CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE tokens(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER REFERENCES users(id) NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

DROP TYPE IF EXISTS access_type;
create type access_type as enum ('public', 'private');

create table pastes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title varchar(256) not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz not null,
    visibility access_type not null, 
    last_visited timestamptz default now() not null,
    burn_after_read BOOLEAN NOT NULL default FALSE,
    user_id integer references users(id)
);

create table pastes_passwords(
    id SERIAL PRIMARY KEY,
    paste_id UUID REFERENCES pastes(id),
    password_hashed VARCHAR(255) not null
);

create index idx_pastes_userid ON pastes(user_id);
create index idx_pastes_created_at on pastes(created_at DESC);

create index idx_pastes_passwords ON pastes_passwords(paste_id);

COMMIT;