alter table "public"."workspaces"
  add constraint "workspaces_user_id_fkey"
  foreign key ("user_id")
  references "auth"."users"
  ("id") on update restrict on delete restrict;
