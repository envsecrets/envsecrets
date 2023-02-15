alter table "public"."org_level_permissions" alter column "permissions" set default '{"projects_create": true}'::jsonb;
