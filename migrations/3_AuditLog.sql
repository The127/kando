create or replace function current_user_oid()
    returns oid
    language plpgsql
as
$$
begin
    return (select oid from "pg_roles" where rolname = current_user);
end;
$$;

create or replace function get_application_user_id()
    returns uuid
    language plpgsql
as
$$
begin
    return current_setting('application_user_id')::uuid;
exception
    when sqlstate '42704' then
        return null;
end;
$$;

create or replace function jsonb_diff_val("val1" jsonb, "val2" jsonb)
    returns jsonb
as
$$
declare
    result jsonb;
    v record;
begin
    result = "val1";
    for v in select * from jsonb_each("val2") loop
            if result @> jsonb_build_object(v.key,v.value)
            then result = result - v.key;
            elsif result ? v.key then continue;
            else
                result = result || jsonb_build_object(v.key,'null');
            end if ;
        end loop ;
    return result;
end
$$ language plpgsql;

create type "audit_log_event_type" as enum ('insert', 'update', 'delete');

create table if not exists "audit_logs"
(
    "id" uuid not null default gen_random_uuid(),
    "table_oid" oid not null,
    "user_oid" oid not null,
    "application_user_id" uuid,
    "row_id" uuid not null,
    "changes" jsonb not null,
    "timestamp" timestamptz not null default now(),
    "event_type" audit_log_event_type not null,
    primary key ("id")
);

select set_table_metadata('db_object_metadata', 'no_audit', 'true');
select set_table_metadata('audit_logs', 'no_audit', 'true');

create or replace function create_insert_audit_log_entries()
    returns trigger
    language plpgsql
as
$$
begin
    insert into "audit_logs"
    (
        "table_oid",
        "user_oid",
        "application_user_id",
        "row_id",
        "changes",
        "event_type"
    )
    values
        (
            tg_relid,
            current_user_oid(),
            get_application_user_id(),
            new.id,
            row_to_json(new)::jsonb - 'id',
            'insert'
        );

    return null;
end;
$$;

create or replace function create_update_audit_log_entries()
    returns trigger
    language plpgsql
as
$$
begin
    insert into "audit_logs"
    (
        "table_oid",
        "user_oid",
        "application_user_id",
        "row_id",
        "changes",
        "event_type"
    )
    values
        (
            tg_relid,
            current_user_oid(),
            get_application_user_id(),
            new.id,
            jsonb_diff_val(row_to_json(new)::jsonb, row_to_json(old)::jsonb),
            'update'
        );

    return null;
end;
$$;

create or replace function create_delete_audit_log_entries()
    returns trigger
    language plpgsql
as
$$
begin
    insert into "audit_logs"
    (
        "table_oid",
        "user_oid",
        "application_user_id",
        "row_id",
        "changes",
        "event_type"
    )
    values
        (
            tg_relid,
            (select current_user_oid()),
            get_application_user_id(),
            old.id,
            '{}'::jsonb,
            'delete'
        );

    return null;
end;
$$;

create or replace function create_audit_log_triggers_for_table("full_table_name" text)
    returns void
    language plpgsql
as
$$
declare
    "table_name" text;
    "table_schema" text;
begin
    "table_schema" = split_part("full_table_name", '.', 1);
    "table_name" = split_part("full_table_name", '.', 2);

    if get_table_metadata("table_name", 'no_audit', "table_schema") is not null then
        return;
    end if;

    execute format('create trigger "%s_%s_insert_audit_trg" after insert on "%s"."%s" ' ||
                   'for each row execute function create_insert_audit_log_entries()',
                   "table_schema", "table_name", "table_schema", "table_name");

    execute format('create trigger "%s_%s_update_audit_trg" after update on "%s"."%s" ' ||
                   'for each row execute function create_update_audit_log_entries()',
                   "table_schema", "table_name", "table_schema", "table_name");

    execute format('create trigger "%s_%s_delete_audit_trg" after delete on "%s"."%s" ' ||
                   'for each row execute function create_delete_audit_log_entries()',
                   "table_schema", "table_name", "table_schema", "table_name");
end;
$$;

create or replace function create_audit_log_triggers_on_create_table()
    returns event_trigger
    language plpgsql
as
$$
declare
    "row" record;
begin
    if tg_tag <> 'CREATE TABLE' then
        return;
    end if;

    for "row" in
        select * from pg_event_trigger_ddl_commands()
        loop
            select create_audit_log_triggers_for_table("row".object_identity);
        end loop;
end;
$$;

create event trigger "create_audit_log_triggers_on_create_table_trg" on ddl_command_end
execute function create_audit_log_triggers_on_create_table();

do
$$
    declare
        "row" record;
    begin
        for "row" in
            select "schemaname", "tablename" from pg_tables where "schemaname" = 'public'
            loop
                perform create_audit_log_triggers_for_table("row"."schemaname" || '.'|| "row"."tablename");
            end loop;
    end;
$$;