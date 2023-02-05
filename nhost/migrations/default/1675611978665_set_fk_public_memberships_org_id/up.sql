alter table "public"."memberships" drop constraint "memberships_workspace_id_fkey",
  add constraint "memberships_org_id_fkey"
  foreign key ("org_id")
  references "public"."organisations"
  ("id") on update restrict on delete cascade;
