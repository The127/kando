create table if not exists "users"
(
    "id"              uuid        not null default gen_random_uuid(),
    "display_name"    text        not null,
    "username"        text unique not null,
    "hashed_password" bytea       not null,
    primary key ("id")
);

alter table "users"
    add constraint "display_name_length_check"
        check (length("display_name") < 100);

alter table "users"
    add constraint "email_username_check"
        check (length("username") <= 100);

create table if not exists "sessions"
(
    "id"                  uuid        not null default gen_random_uuid(),
    "user_id"             uuid        not null,
    "created_timestamp"   timestamptz not null default now(),
    "refreshed_timestamp" timestamptz not null default now(),
    primary key ("id")
);

alter table "sessions"
    add constraint "fk_sessions_users" foreign key ("user_id") references "users";

alter table "sessions"
    add constraint "created_before_refreshed"
        check ("created_timestamp" <= "refreshed_timestamp");

