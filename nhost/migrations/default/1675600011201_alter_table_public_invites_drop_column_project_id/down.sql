alter table "public"."invites" alter column "project_id" drop not null;
alter table "public"."invites" add column "project_id" uuid;
