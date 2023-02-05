alter table "public"."invites" drop constraint "invites_workspace_id_fkey",
  add constraint "invites_org_id_fkey"
  foreign key ("org_id")
  references "public"."organisations"
  ("id") on update restrict on delete cascade;
