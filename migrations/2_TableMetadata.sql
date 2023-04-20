create extension if not exists "hstore";

create table if not exists "db_object_metadata"
(
    "table_oid" oid    not null,
    "metadata"  hstore not null,
    primary key ("table_oid")
);

create or replace function set_table_metadata("table_name" text, "key" text, "value" text,
                                              "table_schema" text default "current_schema"())
    returns void
    language plpgsql
as
$$
begin
    insert into "db_object_metadata" ("table_oid", "metadata")
    values (("table_schema" || '.' || "table_name")::regclass::oid, hstore("key", "value"))
    on conflict ("table_oid") do update
        set "metadata" = "db_object_metadata"."metadata" || hstore("key", "value");
end;
$$;

create or replace function unset_table_metadata("table_name" text, "key" text, "table_schema" text default "current_schema"())
    returns void
    language plpgsql
as
$$
begin
    update "db_object_metadata"
    set "metadata" = delete("metadata", "key")
    where "table_oid" = ("table_schema" || '.' || "table_name")::regclass::oid;
end;
$$;

create or replace function get_table_metadata("table_name" text, "key" text, "table_schema" text default "current_schema"())
    returns text
    language plpgsql
as
$$
begin
    return (select "metadata" -> "key"
            from "db_object_metadata"
            where "table_oid" = ("table_schema" || '.' || "table_name")::regclass::oid);
end;
$$;

create or replace function update_table_metadata_entries()
    returns event_trigger
    language plpgsql
as
$$
begin
    if tg_tag = 'DROP TABLE' then
        delete from "db_object_metadata"
        where "table_oid" in(select objid from pg_event_trigger_dropped_objects());
    end if;
end ;
$$;

create event trigger "update_table_metadata_trg" on sql_drop
    execute function update_table_metadata_entries();