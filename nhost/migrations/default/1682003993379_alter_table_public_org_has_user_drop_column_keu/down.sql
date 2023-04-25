comment on column "public"."org_has_user"."keu" is E'Membership table.';
alter table "public"."org_has_user" alter column "keu" drop not null;
alter table "public"."org_has_user" add column "keu" text;
