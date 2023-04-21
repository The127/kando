create extension if not exists "pg_trgm";

create table "manufacturers"
(
    "id"   uuid not null default gen_random_uuid(),
    "name" text not null,
    primary key ("id")
);

create unique index "manufacturers_name_unique_idx" on "manufacturers" ("name");

create index "manufacturers_name_search_idx" on "manufacturers" using gin ("name" gin_trgm_ops);