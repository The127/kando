create table "asset_types"
(
    "id"   uuid not null default gen_random_uuid(),
    "name" text not null,
    primary key ("id")
);

create unique index "asset_types_name_unique_idx" on "asset_types" ("name");

create index "asset_types_name_search_idx" on "asset_types" using gin ("name" gin_trgm_ops);