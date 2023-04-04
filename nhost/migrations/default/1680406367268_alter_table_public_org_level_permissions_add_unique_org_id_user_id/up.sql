alter table "public"."org_level_permissions" add constraint "org_level_permissions_org_id_user_id_key" unique ("org_id", "user_id");
