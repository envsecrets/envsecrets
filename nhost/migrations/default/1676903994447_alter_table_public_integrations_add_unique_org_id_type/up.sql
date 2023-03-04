alter table "public"."integrations" add constraint "integrations_org_id_type_key" unique ("org_id", "type");
