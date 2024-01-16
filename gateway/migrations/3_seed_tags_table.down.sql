-- Remove seeded entries from tags table
DELETE FROM tags WHERE name IN ('uploaded', 'generated');
