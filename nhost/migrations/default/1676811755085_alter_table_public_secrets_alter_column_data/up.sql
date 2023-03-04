alter table "public"."secrets" alter column "data" set default jsonb_build_object();
