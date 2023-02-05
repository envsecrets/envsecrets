alter table "public"."permissions" add constraint "permissions_membership_id_scope_key" unique ("membership_id", "scope");
