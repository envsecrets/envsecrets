alter table "public"."paths"
  add constraint "paths_parent_id_fkey"
  foreign key ("parent_id")
  references "public"."paths"
  ("id") on update restrict on delete cascade;
