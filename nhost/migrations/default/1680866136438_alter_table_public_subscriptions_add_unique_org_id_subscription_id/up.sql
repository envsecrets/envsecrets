alter table "public"."subscriptions" add constraint "subscriptions_org_id_subscription_id_key" unique ("org_id", "subscription_id");
