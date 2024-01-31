-- This migration adds the display, task_category, and default_tool flags to the tools table.
ALTER TABLE tools ADD COLUMN display BOOLEAN DEFAULT true;
ALTER TABLE tools ADD COLUMN task_category VARCHAR(255);
ALTER TABLE tools ADD COLUMN default_tool BOOLEAN DEFAULT false;