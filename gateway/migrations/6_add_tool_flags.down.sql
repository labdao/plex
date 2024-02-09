-- This migration removes the display, task_category, and default_tool columns from the tools table.
ALTER TABLE tools DROP COLUMN display;
ALTER TABLE tools DROP COLUMN task_category;
ALTER TABLE tools DROP COLUMN default_tool;
