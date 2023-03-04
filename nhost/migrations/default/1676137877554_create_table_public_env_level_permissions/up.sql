CREATE TABLE "public"."env_level_permissions" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "env_id" uuid NOT NULL, "user_id" uuid NOT NULL, "permissions" jsonb NOT NULL DEFAULT '{"secrets_write": false}', PRIMARY KEY ("id") , FOREIGN KEY ("env_id") REFERENCES "public"."environments"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON UPDATE restrict ON DELETE cascade);
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
CREATE TRIGGER "set_public_env_level_permissions_updated_at"
BEFORE UPDATE ON "public"."env_level_permissions"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_env_level_permissions_updated_at" ON "public"."env_level_permissions" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
