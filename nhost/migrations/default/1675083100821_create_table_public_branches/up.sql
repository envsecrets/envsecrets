CREATE TABLE "public"."branches" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "environment_id" uuid NOT NULL, "name" text NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("environment_id") REFERENCES "public"."environments"("id") ON UPDATE restrict ON DELETE restrict, UNIQUE ("environment_id", "name"));
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
CREATE TRIGGER "set_public_branches_updated_at"
BEFORE UPDATE ON "public"."branches"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_branches_updated_at" ON "public"."branches" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
