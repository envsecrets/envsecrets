alter table "public"."projects" add constraint "projects_name_org_id_key" unique ("name", "org_id");
