alter table "public"."permissions" drop constraint "permissions_environment_id_fkey",
  add constraint "permissions_environment_id_fkey"
  foreign key ("membership_id")
  references "public"."memberships"
  ("id") on update restrict on delete cascade;
