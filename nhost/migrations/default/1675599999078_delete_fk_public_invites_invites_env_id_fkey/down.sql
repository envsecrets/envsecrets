alter table "public"."invites"
  add constraint "invites_env_id_fkey"
  foreign key ("env_id")
  references "public"."environments"
  ("id") on update restrict on delete restrict;
