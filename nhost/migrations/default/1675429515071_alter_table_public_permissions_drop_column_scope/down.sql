alter table "public"."permissions" alter column "scope" set default ''*/*'::text';
alter table "public"."permissions" alter column "scope" drop not null;
alter table "public"."permissions" add column "scope" text;
