alter table "public"."roles" alter column "permissions" set default jsonb_build_object();
