create extension pgcrypto;

create table "assets"
(
    "id"              uuid        not null default gen_random_uuid(),
    "asset_type_id"   uuid        not null,
    "name"            text        not null,
    "serial_number"   text,
    "batch_number"    text,
    "manufacturer_id" uuid,
    "parent_id"       uuid,
    "created_at"      timestamptz not null default now(),
    "removed_at"      timestamptz,
    primary key ("id")
);

create index "assets_name_search_idx" on "assets" using gin ("name" gin_trgm_ops);
create unique index "assets_serial_number_unique" on "assets"
    using btree("asset_type_id", digest("serial_number", 'sha512'::text))
    where "serial_number" is not null;

alter table "assets"
    add constraint "fk_assets_asset_type" foreign key ("asset_type_id") references "asset_types";

alter table "assets"
    add constraint "fk_assets_assets" foreign key ("parent_id") references "assets";

alter table "assets"
    add constraint "fk_assets_manufacturers" foreign key ("manufacturer_id") references "manufacturers";