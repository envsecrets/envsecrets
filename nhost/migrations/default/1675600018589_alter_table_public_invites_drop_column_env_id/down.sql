alter table "public"."invites" alter column "env_id" drop not null;
alter table "public"."invites" add column "env_id" uuid;
