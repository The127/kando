create table "locations"
(
    "id"   uuid not null default gen_random_uuid(),
    "name" text not null,
    primary key ("id")
);

create unique index "locations_name_unique" on "locations" ("name");
create index "locations_name_search" on "locations" using gin ("name" gin_trgm_ops);

create table "asset_locations"
(
    "id"             uuid        not null default gen_random_uuid(),
    "asset_id"       uuid        not null,
    "location_id"    uuid        not null,
    "checked_in_at"  timestamptz not null default now(),
    "checked_out_at" timestamptz,
    primary key ("id")
);

alter table "asset_locations"
    add constraint "fk_asset_locations_asset" foreign key ("asset_id") references "assets";

alter table "asset_locations"
    add constraint "fk_asset_locations_location" foreign key ("location_id") references "locations";