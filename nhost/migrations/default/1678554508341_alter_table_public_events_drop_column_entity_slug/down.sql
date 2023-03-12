alter table "public"."events" alter column "entity_slug" drop not null;
alter table "public"."events" add column "entity_slug" text;
