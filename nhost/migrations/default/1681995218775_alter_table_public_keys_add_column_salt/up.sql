alter table "public"."keys" add column "salt" text
 not null unique;
