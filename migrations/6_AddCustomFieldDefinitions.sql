create type "custom_field_type" as enum ('text', 'int', 'float', 'bool');

create table "custom_field_definitions"
(
    "id" uuid not null default gen_random_uuid(),
    "asset_type_id" uuid not null,
    "name" text not null,
    "field_type" custom_field_type not null,
    primary key ("id")
);

create index "custom_field_definitions_unique_name" on "custom_field_definitions" ("asset_type_id", "name");
create index "custom_field_definitions_name_search" on "custom_field_definitions" using gin ("name" gin_trgm_ops);

alter table "custom_field_definitions"
    add constraint "custom_field_definitions_asset_type" foreign key ("asset_type_id") references "asset_types";
