-- 000001_init_schema.down.sql
-- Teman Mancing Database Schema Reversion

DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS messages CASCADE;
DROP TABLE IF EXISTS conversation_participants CASCADE;
DROP TABLE IF EXISTS conversations CASCADE;
DROP TABLE IF EXISTS event_participants CASCADE;
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS spot_comments CASCADE;
DROP TABLE IF EXISTS recommended_baits CASCADE;
DROP TABLE IF EXISTS spot_techniques CASCADE;
DROP TABLE IF EXISTS fishing_techniques CASCADE;
DROP TABLE IF EXISTS fishing_spot_fish_types CASCADE;
DROP TABLE IF EXISTS fish_types CASCADE;
DROP TABLE IF EXISTS fishing_spot_photos CASCADE;
DROP TABLE IF EXISTS fishing_spots CASCADE;
DROP TABLE IF EXISTS users CASCADE;
