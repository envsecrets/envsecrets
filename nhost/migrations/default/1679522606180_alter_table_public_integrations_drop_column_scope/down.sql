alter table "public"."integrations" alter column "scope" drop not null;
alter table "public"."integrations" add column "scope" jsonb;
