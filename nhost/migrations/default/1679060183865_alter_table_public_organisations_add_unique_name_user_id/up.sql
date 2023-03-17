alter table "public"."organisations" add constraint "organisations_name_user_id_key" unique ("name", "user_id");
