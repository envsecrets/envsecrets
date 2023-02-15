alter table "public"."paths" add column "location" ltree
 not null unique;
