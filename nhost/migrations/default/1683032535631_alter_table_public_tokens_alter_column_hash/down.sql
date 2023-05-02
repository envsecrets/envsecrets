alter table "public"."tokens" alter column "hash" set default md5((random())::text);
