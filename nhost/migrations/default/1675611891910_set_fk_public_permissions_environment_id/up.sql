alter table "public"."permissions" drop constraint "permissions_environment_id_fkey",
  add constraint "permissions_environment_id_fkey"
  foreign key ("environment_id")
  references "public"."environments"
  ("id") on update restrict on delete cascade;
