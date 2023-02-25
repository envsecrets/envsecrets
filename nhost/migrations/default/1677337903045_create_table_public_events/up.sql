CREATE TABLE "public"."events" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "integration_id" uuid NOT NULL, "env_id" uuid NOT NULL, "entity_slug" text NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("integration_id") REFERENCES "public"."integrations"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("env_id") REFERENCES "public"."environments"("id") ON UPDATE restrict ON DELETE cascade);
CREATE OR REPLACE FUNCTION "public"."set_current_timestamp_updated_at"()
RETURNS TRIGGER AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER "set_public_events_updated_at"
BEFORE UPDATE ON "public"."events"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_events_updated_at" ON "public"."events" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
