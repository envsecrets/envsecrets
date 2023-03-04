alter table "public"."org_level_permissions" alter column "permissions" set default jsonb_build_object();
