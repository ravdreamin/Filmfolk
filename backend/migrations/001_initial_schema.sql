-- FilmFolk Database Schema
-- Complete schema for the movie review and social platform

-- ============================================================================
-- ENUMS
-- ============================================================================

CREATE TYPE user_role AS ENUM ('user', 'moderator', 'admin');
CREATE TYPE auth_provider AS ENUM ('email', 'google', 'facebook', 'instagram', 'twitter', 'guest');
CREATE TYPE account_status AS ENUM ('active', 'suspended', 'banned');
CREATE TYPE movie_status AS ENUM ('pending_approval', 'approved', 'rejected');
CREATE TYPE review_status AS ENUM ('pending_moderation', 'published', 'rejected');
CREATE TYPE list_type AS ENUM ('watched', 'dropped', 'plan_to_watch', 'custom');
CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read');
CREATE TYPE moderation_action AS ENUM ('review_flagged', 'review_removed', 'user_warned', 'user_suspended', 'user_banned');
CREATE TYPE community_type AS ENUM ('public', 'private', 'restricted');

-- ============================================================================
-- CORE TABLES
-- ============================================================================

-- USERS TABLE
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT,
    role user_role NOT NULL DEFAULT 'user',
    auth_provider auth_provider NOT NULL DEFAULT 'email',
    provider_id VARCHAR(255),
    status account_status NOT NULL DEFAULT 'active',
    avatar_url TEXT,
    bio TEXT,

    -- Gamification
    total_reviews INTEGER DEFAULT 0,
    total_comments INTEGER DEFAULT 0,
    total_likes_received INTEGER DEFAULT 0,
    engagement_score INTEGER DEFAULT 0,
    current_title_id BIGINT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,

    CONSTRAINT unique_provider_id UNIQUE(auth_provider, provider_id)
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_engagement_score ON users(engagement_score DESC);

COMMENT ON TABLE users IS 'Stores user accounts with multi-provider authentication support';
COMMENT ON COLUMN users.password_hash IS 'NULL for OAuth users, required for email auth';
COMMENT ON COLUMN users.engagement_score IS 'Calculated score for gamification and friend recommendations';

-- ============================================================================
-- USER TITLES (GAMIFICATION)
-- ============================================================================

CREATE TABLE user_titles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    required_reviews INTEGER DEFAULT 0,
    required_comments INTEGER DEFAULT 0,
    required_engagement_score INTEGER DEFAULT 0,
    icon_url TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_titles_sort_order ON user_titles(sort_order);

COMMENT ON TABLE user_titles IS 'Defines user title tiers based on engagement';

-- Add foreign key for users.current_title_id
ALTER TABLE users ADD CONSTRAINT fk_users_current_title
    FOREIGN KEY (current_title_id) REFERENCES user_titles(id) ON DELETE SET NULL;

-- ============================================================================
-- MOVIES
-- ============================================================================

CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    release_year INTEGER NOT NULL,
    genres VARCHAR(255)[],
    summary TEXT,
    poster_url TEXT,
    backdrop_url TEXT,
    runtime_minutes INTEGER,
    language VARCHAR(50),

    -- External API integration
    tmdb_id INTEGER UNIQUE,
    imdb_id VARCHAR(20),

    status movie_status NOT NULL DEFAULT 'pending_approval',
    submitted_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    -- Aggregated stats
    average_rating DECIMAL(3,2),
    total_reviews INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_movie UNIQUE(title, release_year)
);

CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_status ON movies(status);
CREATE INDEX idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX idx_movies_average_rating ON movies(average_rating DESC);
CREATE INDEX idx_movies_release_year ON movies(release_year DESC);

COMMENT ON TABLE movies IS 'Central catalog of movies with moderation workflow';

-- ============================================================================
-- CASTS
-- ============================================================================

CREATE TABLE casts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    profile_url TEXT,
    tmdb_person_id INTEGER UNIQUE,
    bio TEXT,
    birth_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_casts_name ON casts(name);
CREATE INDEX idx_casts_tmdb_person_id ON casts(tmdb_person_id);

COMMENT ON TABLE casts IS 'Stores information about actors, directors, and crew';

-- ============================================================================
-- MOVIE_CAST (Junction Table)
-- ============================================================================

CREATE TABLE movie_casts (
    id BIGSERIAL PRIMARY KEY,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    cast_id BIGINT NOT NULL REFERENCES casts(id) ON DELETE CASCADE,
    role VARCHAR(100) NOT NULL, -- e.g., 'actor', 'director', 'producer'
    character_name VARCHAR(255), -- for actors
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_movie_cast_role UNIQUE(movie_id, cast_id, role, character_name)
);

CREATE INDEX idx_movie_casts_movie_id ON movie_casts(movie_id);
CREATE INDEX idx_movie_casts_cast_id ON movie_casts(cast_id);

COMMENT ON TABLE movie_casts IS 'Links movies with cast members and their roles';

-- ============================================================================
-- REVIEWS
-- ============================================================================

CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    review_text TEXT NOT NULL,

    -- AI Analysis
    sentiment VARCHAR(50),
    ai_flagged BOOLEAN DEFAULT FALSE,
    ai_flag_reason TEXT,

    status review_status NOT NULL DEFAULT 'published',

    -- Thread control
    is_thread_locked BOOLEAN DEFAULT FALSE,

    -- Engagement
    likes_count INTEGER DEFAULT 0,
    comments_count INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_user_movie_review UNIQUE(user_id, movie_id)
);

CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_movie_id ON reviews(movie_id);
CREATE INDEX idx_reviews_status ON reviews(status);
CREATE INDEX idx_reviews_ai_flagged ON reviews(ai_flagged) WHERE ai_flagged = TRUE;
CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC);

COMMENT ON TABLE reviews IS 'User reviews with AI moderation and thread control';
COMMENT ON COLUMN reviews.is_thread_locked IS 'Review author can lock their review thread to prevent further comments';

-- ============================================================================
-- REVIEW COMMENTS (THREADED)
-- ============================================================================

CREATE TABLE review_comments (
    id BIGSERIAL PRIMARY KEY,
    review_id BIGINT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id BIGINT REFERENCES review_comments(id) ON DELETE CASCADE,
    comment_text TEXT NOT NULL,

    -- Moderation
    ai_flagged BOOLEAN DEFAULT FALSE,
    is_removed BOOLEAN DEFAULT FALSE,
    removed_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,

    likes_count INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_review_comments_review_id ON review_comments(review_id);
CREATE INDEX idx_review_comments_user_id ON review_comments(user_id);
CREATE INDEX idx_review_comments_parent_comment_id ON review_comments(parent_comment_id);
CREATE INDEX idx_review_comments_created_at ON review_comments(created_at);

COMMENT ON TABLE review_comments IS 'Threaded comments on reviews';

-- ============================================================================
-- REVIEW LIKES
-- ============================================================================

CREATE TABLE review_likes (
    id BIGSERIAL PRIMARY KEY,
    review_id BIGINT NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_review_like UNIQUE(review_id, user_id)
);

CREATE INDEX idx_review_likes_review_id ON review_likes(review_id);
CREATE INDEX idx_review_likes_user_id ON review_likes(user_id);

-- ============================================================================
-- COMMENT LIKES
-- ============================================================================

CREATE TABLE comment_likes (
    id BIGSERIAL PRIMARY KEY,
    comment_id BIGINT NOT NULL REFERENCES review_comments(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_comment_like UNIQUE(comment_id, user_id)
);

CREATE INDEX idx_comment_likes_comment_id ON comment_likes(comment_id);
CREATE INDEX idx_comment_likes_user_id ON comment_likes(user_id);

-- ============================================================================
-- USER LISTS
-- ============================================================================

CREATE TABLE user_lists (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    list_type list_type NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_user_list_name UNIQUE(user_id, name)
);

CREATE INDEX idx_user_lists_user_id ON user_lists(user_id);
CREATE INDEX idx_user_lists_type ON user_lists(list_type);

COMMENT ON TABLE user_lists IS 'User-created movie lists (watched, dropped, plan to watch, custom)';

-- ============================================================================
-- USER LIST ITEMS
-- ============================================================================

CREATE TABLE user_list_items (
    id BIGSERIAL PRIMARY KEY,
    list_id BIGINT NOT NULL REFERENCES user_lists(id) ON DELETE CASCADE,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    notes TEXT,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_list_movie UNIQUE(list_id, movie_id)
);

CREATE INDEX idx_user_list_items_list_id ON user_list_items(list_id);
CREATE INDEX idx_user_list_items_movie_id ON user_list_items(movie_id);

-- ============================================================================
-- FRIENDSHIPS
-- ============================================================================

CREATE TABLE friendships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'accepted', 'rejected', 'blocked'

    -- Recommendation score
    taste_similarity_score DECIMAL(5,2),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_friendship UNIQUE(user_id, friend_id),
    CONSTRAINT no_self_friendship CHECK (user_id != friend_id)
);

CREATE INDEX idx_friendships_user_id ON friendships(user_id);
CREATE INDEX idx_friendships_friend_id ON friendships(friend_id);
CREATE INDEX idx_friendships_status ON friendships(status);

COMMENT ON TABLE friendships IS 'Manages friend relationships with taste-based recommendation scores';

-- ============================================================================
-- DIRECT MESSAGES
-- ============================================================================

CREATE TABLE direct_messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_text TEXT NOT NULL,
    status message_status NOT NULL DEFAULT 'sent',
    is_deleted_by_sender BOOLEAN DEFAULT FALSE,
    is_deleted_by_receiver BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,

    CONSTRAINT no_self_message CHECK (sender_id != receiver_id)
);

CREATE INDEX idx_direct_messages_sender_id ON direct_messages(sender_id);
CREATE INDEX idx_direct_messages_receiver_id ON direct_messages(receiver_id);
CREATE INDEX idx_direct_messages_created_at ON direct_messages(created_at DESC);
CREATE INDEX idx_direct_messages_conversation ON direct_messages(sender_id, receiver_id, created_at DESC);

COMMENT ON TABLE direct_messages IS 'One-on-one direct messaging between users';

-- ============================================================================
-- COMMUNITIES
-- ============================================================================

CREATE TABLE communities (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    type community_type NOT NULL DEFAULT 'public',
    created_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    avatar_url TEXT,
    member_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_communities_name ON communities(name);
CREATE INDEX idx_communities_type ON communities(type);
CREATE INDEX idx_communities_created_by ON communities(created_by_user_id);

COMMENT ON TABLE communities IS 'Topic-based community chat rooms';

-- ============================================================================
-- COMMUNITY MEMBERS
-- ============================================================================

CREATE TABLE community_members (
    id BIGSERIAL PRIMARY KEY,
    community_id BIGINT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_moderator BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_community_member UNIQUE(community_id, user_id)
);

CREATE INDEX idx_community_members_community_id ON community_members(community_id);
CREATE INDEX idx_community_members_user_id ON community_members(user_id);

-- ============================================================================
-- COMMUNITY MESSAGES
-- ============================================================================

CREATE TABLE community_messages (
    id BIGSERIAL PRIMARY KEY,
    community_id BIGINT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_text TEXT NOT NULL,
    is_removed BOOLEAN DEFAULT FALSE,
    removed_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_community_messages_community_id ON community_messages(community_id);
CREATE INDEX idx_community_messages_user_id ON community_messages(user_id);
CREATE INDEX idx_community_messages_created_at ON community_messages(created_at DESC);

COMMENT ON TABLE community_messages IS 'Messages in community chat rooms';

-- ============================================================================
-- WORLD CHAT MESSAGES
-- ============================================================================

CREATE TABLE world_chat_messages (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_text TEXT NOT NULL,
    is_removed BOOLEAN DEFAULT FALSE,
    removed_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_world_chat_messages_created_at ON world_chat_messages(created_at DESC);
CREATE INDEX idx_world_chat_messages_user_id ON world_chat_messages(user_id);

COMMENT ON TABLE world_chat_messages IS 'Global chat accessible to all users';

-- ============================================================================
-- MODERATION LOGS
-- ============================================================================

CREATE TABLE moderation_logs (
    id BIGSERIAL PRIMARY KEY,
    moderator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action moderation_action NOT NULL,
    reason TEXT NOT NULL,

    -- References to moderated content
    review_id BIGINT REFERENCES reviews(id) ON DELETE SET NULL,
    comment_id BIGINT REFERENCES review_comments(id) ON DELETE SET NULL,

    -- Suspension/ban details
    duration_days INTEGER,
    expires_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_moderation_logs_moderator_id ON moderation_logs(moderator_id);
CREATE INDEX idx_moderation_logs_target_user_id ON moderation_logs(target_user_id);
CREATE INDEX idx_moderation_logs_action ON moderation_logs(action);
CREATE INDEX idx_moderation_logs_created_at ON moderation_logs(created_at DESC);

COMMENT ON TABLE moderation_logs IS 'Audit log for all moderation actions';

-- ============================================================================
-- USER WARNINGS
-- ============================================================================

CREATE TABLE user_warnings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issued_by_moderator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    moderation_log_id BIGINT REFERENCES moderation_logs(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_warnings_user_id ON user_warnings(user_id);
CREATE INDEX idx_user_warnings_is_active ON user_warnings(is_active) WHERE is_active = TRUE;

COMMENT ON TABLE user_warnings IS 'Tracks warnings issued to users by moderators';

-- ============================================================================
-- NOTIFICATIONS
-- ============================================================================

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    link_url TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read) WHERE is_read = FALSE;
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);

COMMENT ON TABLE notifications IS 'User notifications for various activities';

-- ============================================================================
-- REFRESH TOKENS (for JWT)
-- ============================================================================

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for authentication';
