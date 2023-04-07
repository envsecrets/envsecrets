alter table "public"."invites" add column "key" text
 not null default md5(random()::text);
