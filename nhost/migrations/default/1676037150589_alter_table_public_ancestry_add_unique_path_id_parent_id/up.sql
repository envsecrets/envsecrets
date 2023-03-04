alter table "public"."ancestry" add constraint "ancestry_path_id_parent_id_key" unique ("path_id", "parent_id");
