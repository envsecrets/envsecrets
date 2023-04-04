alter table "public"."org_has_user"
  add constraint "org_has_user_role_id_fkey"
  foreign key ("role_id")
  references "public"."roles"
  ("id") on update restrict on delete cascade;
