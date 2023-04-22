create table "tags"
(
    "id"   uuid not null default gen_random_uuid(),
    "name" text not null,
    primary key ("id")
);

create unique index "tags_name_unique" on "tags" ("name");
create index "tags_name_search" on "tags" using gin ("name" gin_trgm_ops);

create table asset_tags
(
    "id"       uuid not null default gen_random_uuid(),
    "asset_id" uuid not null,
    "tag_id"   uuid not null,
    primary key ("id")
);

alter table "asset_tags"
    add constraint "fk_asset_tags_asset" foreign key ("asset_id") references "assets";

alter table "asset_tags"
    add constraint "fk_asset_tags_tag" foreign key ("tag_id") references "tags";

create table "manufacturer_tags"
(
    "id"              uuid not null default gen_random_uuid(),
    "manufacturer_id" uuid not null,
    "tag_id"          uuid not null,
    primary key ("id")
);

alter table "manufacturer_tags"
    add constraint "fk_manufacturer_tags_manufacturers" foreign key ("manufacturer_id") references "manufacturers";

alter table "manufacturer_tags"
    add constraint "fk_manufacturer_tags_tags" foreign key ("tag_id") references "tags";

create table "location_tags"
(
    "id"          uuid not null default gen_random_uuid(),
    "location_id" uuid not null,
    "tag_id"      uuid not null,
    primary key ("id")
);

alter table "location_tags"
    add constraint "fk_location_tags_locations" foreign key ("location_id") references "locations";

alter table "location_tags"
    add constraint "fk_location_tags_tags" foreign key ("tag_id") references "tags";