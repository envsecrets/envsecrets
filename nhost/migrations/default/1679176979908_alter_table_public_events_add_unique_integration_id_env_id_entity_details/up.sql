alter table "public"."events" add constraint "events_integration_id_env_id_entity_details_key" unique ("integration_id", "env_id", "entity_details");
