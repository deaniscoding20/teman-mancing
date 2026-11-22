-- 000005_add_comment_parent_and_photo.up.sql
ALTER TABLE spot_comments ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES spot_comments(id) ON DELETE CASCADE;
ALTER TABLE spot_comments ADD COLUMN IF NOT EXISTS photo_url TEXT;

CREATE INDEX IF NOT EXISTS idx_spot_comments_parent_id ON spot_comments(parent_id);
