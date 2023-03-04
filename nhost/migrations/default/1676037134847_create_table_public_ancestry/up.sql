CREATE TABLE "public"."ancestry" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "path_id" uuid NOT NULL, "parent_id" uuid NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("path_id") REFERENCES "public"."paths"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("parent_id") REFERENCES "public"."paths"("id") ON UPDATE restrict ON DELETE cascade);
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
CREATE TRIGGER "set_public_ancestry_updated_at"
BEFORE UPDATE ON "public"."ancestry"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_ancestry_updated_at" ON "public"."ancestry" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
