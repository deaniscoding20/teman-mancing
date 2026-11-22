DROP INDEX IF EXISTS idx_spot_comments_parent_id;
ALTER TABLE spot_comments DROP COLUMN IF EXISTS photo_url;
ALTER TABLE spot_comments DROP COLUMN IF EXISTS parent_id;
