alter table "public"."invites"
  add constraint "invites_receiver_email_fkey"
  foreign key ("receiver_email")
  references "auth"."users"
  ("email") on update restrict on delete restrict;
