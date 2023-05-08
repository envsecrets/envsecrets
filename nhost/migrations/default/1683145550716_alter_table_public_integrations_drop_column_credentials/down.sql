alter table "public"."integrations" alter column "credentials" drop not null;
alter table "public"."integrations" add column "credentials" jsonb;
