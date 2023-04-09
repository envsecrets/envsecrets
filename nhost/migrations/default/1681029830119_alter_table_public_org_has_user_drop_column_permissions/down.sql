comment on column "public"."org_has_user"."permissions" is E'Membership table.';
alter table "public"."org_has_user" alter column "permissions" set default jsonb_build_object();
alter table "public"."org_has_user" alter column "permissions" drop not null;
alter table "public"."org_has_user" add column "permissions" jsonb;
