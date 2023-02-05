CREATE TABLE "public"."permission" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "membership_id" uuid NOT NULL, "scope" text NOT NULL DEFAULT '*/*', PRIMARY KEY ("id") , FOREIGN KEY ("membership_id") REFERENCES "public"."memberships"("id") ON UPDATE restrict ON DELETE restrict, UNIQUE ("membership_id", "scope"));
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
CREATE TRIGGER "set_public_permission_updated_at"
BEFORE UPDATE ON "public"."permission"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_permission_updated_at" ON "public"."permission" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
