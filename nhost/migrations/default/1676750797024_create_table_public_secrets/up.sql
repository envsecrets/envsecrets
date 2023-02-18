CREATE TABLE "public"."secrets" ("created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "env_id" uuid NOT NULL, "data" jsonb NOT NULL, "version" integer NOT NULL DEFAULT 0, "id" uuid NOT NULL DEFAULT gen_random_uuid(), PRIMARY KEY ("id") , FOREIGN KEY ("env_id") REFERENCES "public"."environments"("id") ON UPDATE restrict ON DELETE cascade, UNIQUE ("env_id", "version"));
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
CREATE TRIGGER "set_public_secrets_updated_at"
BEFORE UPDATE ON "public"."secrets"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_secrets_updated_at" ON "public"."secrets" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
