alter table "public"."invites" alter column "scope" set default ''*/*'::text';
alter table "public"."invites" alter column "scope" drop not null;
alter table "public"."invites" add column "scope" text;
