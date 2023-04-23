create table "custom_fields"
(
    "id" uuid not null default gen_random_uuid(),
    "asset_id" uuid not null,
    "custom_field_definition_id" uuid not null,
    "field_type" custom_field_type not null,
    "value" jsonb not null,
    primary key ("id")
);

alter table "custom_fields"
    add constraint "fk_asset_fields_asset" foreign key ("asset_id") references "assets";

alter table "custom_fields"
    add constraint "fk_custom_fields_custom_field_definition" foreign key ("custom_field_definition_id") references "custom_field_definitions";
