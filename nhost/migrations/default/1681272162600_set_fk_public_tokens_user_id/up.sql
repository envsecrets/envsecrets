alter table "public"."tokens"
  add constraint "tokens_user_id_fkey"
  foreign key ("user_id")
  references "auth"."users"
  ("id") on update restrict on delete cascade;
