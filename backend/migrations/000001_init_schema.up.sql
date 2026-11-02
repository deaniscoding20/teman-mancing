-- 000001_init_schema.up.sql
-- Teman Mancing Database Schema Initialization

-- Enable pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    avatar_key TEXT,
    bio TEXT,
    location VARCHAR(150),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. Fishing Spots
CREATE TABLE IF NOT EXISTS fishing_spots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fishing_spots_user_id ON fishing_spots(user_id);
CREATE INDEX IF NOT EXISTS idx_fishing_spots_lat_lng ON fishing_spots(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_fishing_spots_created_at ON fishing_spots(created_at DESC);

-- 3. Fishing Spot Photos
CREATE TABLE IF NOT EXISTS fishing_spot_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spot_id UUID NOT NULL REFERENCES fishing_spots(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fishing_spot_photos_spot_id ON fishing_spot_photos(spot_id);

-- 4. Fish Types Catalog
CREATE TABLE IF NOT EXISTS fish_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_fish_types_name ON fish_types(name);

-- 5. Fishing Spot Fish Types (Pivot Many-to-Many)
CREATE TABLE IF NOT EXISTS fishing_spot_fish_types (
    spot_id UUID NOT NULL REFERENCES fishing_spots(id) ON DELETE CASCADE,
    fish_type_id UUID NOT NULL REFERENCES fish_types(id) ON DELETE CASCADE,
    PRIMARY KEY (spot_id, fish_type_id)
);

-- 6. Fishing Techniques Catalog
CREATE TABLE IF NOT EXISTS fishing_techniques (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_fishing_techniques_name ON fishing_techniques(name);

-- 7. Spot Techniques (Pivot Many-to-Many)
CREATE TABLE IF NOT EXISTS spot_techniques (
    spot_id UUID NOT NULL REFERENCES fishing_spots(id) ON DELETE CASCADE,
    technique_id UUID NOT NULL REFERENCES fishing_techniques(id) ON DELETE CASCADE,
    PRIMARY KEY (spot_id, technique_id)
);

-- 8. Recommended Baits
CREATE TABLE IF NOT EXISTS recommended_baits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spot_id UUID NOT NULL REFERENCES fishing_spots(id) ON DELETE CASCADE,
    name VARCHAR(150) NOT NULL,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_recommended_baits_spot_id ON recommended_baits(spot_id);

-- 9. Spot Comments
CREATE TABLE IF NOT EXISTS spot_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spot_id UUID NOT NULL REFERENCES fishing_spots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_spot_comments_spot_id ON spot_comments(spot_id);
CREATE INDEX IF NOT EXISTS idx_spot_comments_user_id ON spot_comments(user_id);

-- 10. Events (Janjian Mancing)
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    cover_image_key TEXT,
    fishing_spot_id UUID REFERENCES fishing_spots(id) ON DELETE SET NULL,
    destination_name VARCHAR(200) NOT NULL,
    destination_address TEXT NOT NULL,
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    event_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME,
    meeting_point TEXT NOT NULL,
    difficulty VARCHAR(50) NOT NULL DEFAULT 'pemula',
    has_parking_fee BOOLEAN NOT NULL DEFAULT FALSE,
    parking_fee DECIMAL(12, 2),
    requires_boat BOOLEAN NOT NULL DEFAULT FALSE,
    boat_rental_fee DECIMAL(12, 2),
    additional_cost TEXT,
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    max_participants INTEGER,
    registration_deadline TIMESTAMPTZ,
    status VARCHAR(30) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_events_difficulty CHECK (difficulty IN ('pemula', 'menengah', 'ahli')),
    CONSTRAINT chk_events_visibility CHECK (visibility IN ('public', 'private')),
    CONSTRAINT chk_events_status CHECK (status IN ('open', 'closed', 'cancelled', 'completed'))
);

CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events(creator_id);
CREATE INDEX IF NOT EXISTS idx_events_date_status ON events(event_date, status);
CREATE INDEX IF NOT EXISTS idx_events_visibility ON events(visibility);
CREATE INDEX IF NOT EXISTS idx_events_spot_id ON events(fishing_spot_id);

-- 11. Event Participants
CREATE TABLE IF NOT EXISTS event_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(30) NOT NULL DEFAULT 'participant',
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_event_user UNIQUE (event_id, user_id),
    CONSTRAINT chk_participant_role CHECK (role IN ('creator', 'co-host', 'participant')),
    CONSTRAINT chk_participant_status CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_event_participants_event_status ON event_participants(event_id, status);
CREATE INDEX IF NOT EXISTS idx_event_participants_user_id ON event_participants(user_id);

-- 12. Conversations
CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL DEFAULT 'direct',
    event_id UUID UNIQUE REFERENCES events(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_conversation_type CHECK (type IN ('direct', 'group'))
);

CREATE INDEX IF NOT EXISTS idx_conversations_event_id ON conversations(event_id) WHERE event_id IS NOT NULL;

-- 13. Conversation Participants
CREATE TABLE IF NOT EXISTS conversation_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_conversation_user UNIQUE (conversation_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_conv_participants_user ON conversation_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_conv_participants_conv ON conversation_participants(conversation_id);

-- 14. Messages
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at ASC);

-- 15. Notifications
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    reference_id UUID,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read, created_at DESC);
