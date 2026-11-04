-- 000004_add_event_meeting_location.up.sql
ALTER TABLE events
ADD COLUMN IF NOT EXISTS same_meeting_location BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS meeting_address TEXT,
ADD COLUMN IF NOT EXISTS meeting_latitude NUMERIC(10,7),
ADD COLUMN IF NOT EXISTS meeting_longitude NUMERIC(10,7),
ADD COLUMN IF NOT EXISTS location_notes TEXT;
