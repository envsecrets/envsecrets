alter table "public"."projects" alter column "user_id" drop not null;
alter table "public"."projects" add column "user_id" uuid;
