alter table "public"."projects" drop constraint "projects_workspace_id_fkey",
  add constraint "projects_workspace_id_fkey"
  foreign key ("workspace_id")
  references "public"."workspaces"
  ("id") on update restrict on delete restrict;
