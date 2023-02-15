CREATE TABLE "public"."path_has_user" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "path_id" uuid NOT NULL, "permissions" jsonb NOT NULL DEFAULT '{ "subpath_create": false, "secrets_write": false }', "user_id" uuid NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("path_id") REFERENCES "public"."paths"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON UPDATE restrict ON DELETE restrict);
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
CREATE TRIGGER "set_public_path_has_user_updated_at"
BEFORE UPDATE ON "public"."path_has_user"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_path_has_user_updated_at" ON "public"."path_has_user" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
