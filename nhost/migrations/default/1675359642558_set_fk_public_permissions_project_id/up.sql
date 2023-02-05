alter table "public"."permissions"
  add constraint "permissions_project_id_fkey"
  foreign key ("project_id")
  references "public"."projects"
  ("id") on update restrict on delete restrict;
