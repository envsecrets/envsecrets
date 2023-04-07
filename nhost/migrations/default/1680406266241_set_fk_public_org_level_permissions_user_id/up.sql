alter table "public"."org_level_permissions"
  add constraint "org_level_permissions_user_id_fkey"
  foreign key ("user_id")
  references "auth"."users"
  ("id") on update restrict on delete cascade;
