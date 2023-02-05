alter table "public"."permissions" drop constraint "permission_membership_id_fkey",
  add constraint "permissions_membership_id_fkey"
  foreign key ("membership_id")
  references "public"."memberships"
  ("id") on update restrict on delete cascade;
