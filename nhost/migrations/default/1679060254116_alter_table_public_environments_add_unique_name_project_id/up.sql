alter table "public"."environments" add constraint "environments_name_project_id_key" unique ("name", "project_id");
