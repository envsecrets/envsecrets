alter table "public"."projects" drop constraint "projects_workspace_id_fkey",
  add constraint "projects_org_id_fkey"
  foreign key ("org_id")
  references "public"."organisations"
  ("id") on update restrict on delete cascade;
