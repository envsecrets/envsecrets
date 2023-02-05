alter table "public"."environments" drop constraint "environments_project_id_fkey",
  add constraint "environments_project_id_fkey"
  foreign key ("project_id")
  references "public"."projects"
  ("id") on update restrict on delete cascade;
