alter table "public"."tokens" alter column "hash" drop not null;
alter table "public"."tokens" add column "hash" text;
