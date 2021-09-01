ALTER TABLE management.projects ADD COLUMN register_on_project_resource_owner BOOLEAN;

ALTER TABLE authz.applications ADD COLUMN register_on_project_resource_owner BOOLEAN;
ALTER TABLE auth.applications ADD COLUMN register_on_project_resource_owner BOOLEAN;
ALTER TABLE management.applications ADD COLUMN register_on_project_resource_owner BOOLEAN;

ALTER TABLE management.projects ADD COLUMN login_policy_setting SMALLINT;

ALTER TABLE authz.applications ADD COLUMN login_policy_setting SMALLINT;
ALTER TABLE auth.applications ADD COLUMN login_policy_setting SMALLINT;
ALTER TABLE management.applications ADD COLUMN login_policy_setting SMALLINT;
