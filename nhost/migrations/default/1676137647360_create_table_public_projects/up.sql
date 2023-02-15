CREATE TABLE "public"."projects" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "name" text NOT NULL, "org_id" uuid NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("org_id") REFERENCES "public"."organisations"("id") ON UPDATE restrict ON DELETE cascade);
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
CREATE TRIGGER "set_public_projects_updated_at"
BEFORE UPDATE ON "public"."projects"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_projects_updated_at" ON "public"."projects" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
