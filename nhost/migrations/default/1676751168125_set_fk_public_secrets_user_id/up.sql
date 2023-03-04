alter table "public"."secrets"
  add constraint "secrets_user_id_fkey"
  foreign key ("user_id")
  references "auth"."users"
  ("id") on update restrict on delete cascade;
