alter table "public"."invites"
  add constraint "invites_project_id_fkey"
  foreign key ("project_id")
  references "public"."projects"
  ("id") on update restrict on delete restrict;
