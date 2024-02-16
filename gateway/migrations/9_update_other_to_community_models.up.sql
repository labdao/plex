UPDATE tools
SET task_category = 'community-models'
WHERE task_category = 'other-models';

ALTER TABLE tools
ALTER COLUMN task_category SET DEFAULT 'community-models';