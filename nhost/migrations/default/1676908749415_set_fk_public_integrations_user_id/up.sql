alter table "public"."integrations"
  add constraint "integrations_user_id_fkey"
  foreign key ("user_id")
  references "auth"."users"
  ("id") on update restrict on delete cascade;
