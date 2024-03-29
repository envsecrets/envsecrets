CREATE TABLE "public"."invites" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "org_id" uuid NOT NULL, "user_id" uuid NOT NULL, "email" text NOT NULL, "role_id" uuid NOT NULL, "accepted" boolean NOT NULL DEFAULT False, PRIMARY KEY ("id") , FOREIGN KEY ("org_id") REFERENCES "public"."organisations"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("user_id") REFERENCES "auth"."users"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id") ON UPDATE restrict ON DELETE cascade, UNIQUE ("org_id", "email"));
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
CREATE TRIGGER "set_public_invites_updated_at"
BEFORE UPDATE ON "public"."invites"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_invites_updated_at" ON "public"."invites" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
