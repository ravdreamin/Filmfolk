-- Add Followers/Following Feature
-- Self-referential many-to-many relationship for user following

-- ============================================================================
-- ADD FOLLOWER COUNTS TO USERS TABLE
-- ============================================================================

ALTER TABLE users ADD COLUMN followers_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN following_count INTEGER DEFAULT 0;

CREATE INDEX idx_users_followers_count ON users(followers_count DESC);
CREATE INDEX idx_users_following_count ON users(following_count DESC);

COMMENT ON COLUMN users.followers_count IS 'Cached count of users following this user';
COMMENT ON COLUMN users.following_count IS 'Cached count of users this user is following';

-- ============================================================================
-- FOLLOWERS TABLE (Self-referential many-to-many)
-- ============================================================================

CREATE TABLE followers (
    id BIGSERIAL PRIMARY KEY,
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT unique_follow UNIQUE(follower_id, following_id),
    CONSTRAINT no_self_follow CHECK (follower_id != following_id)
);

CREATE INDEX idx_followers_follower_id ON followers(follower_id);
CREATE INDEX idx_followers_following_id ON followers(following_id);
CREATE INDEX idx_followers_created_at ON followers(created_at DESC);

COMMENT ON TABLE followers IS 'User following relationships';
COMMENT ON COLUMN followers.follower_id IS 'User who is following';
COMMENT ON COLUMN followers.following_id IS 'User being followed';

-- ============================================================================
-- FUNCTIONS TO UPDATE FOLLOWER COUNTS
-- ============================================================================

-- Function to increment follower counts when a follow is created
CREATE OR REPLACE FUNCTION increment_follower_counts()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment followers_count for the user being followed
    UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.following_id;

    -- Increment following_count for the user who is following
    UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Function to decrement follower counts when a follow is deleted
CREATE OR REPLACE FUNCTION decrement_follower_counts()
RETURNS TRIGGER AS $$
BEGIN
    -- Decrement followers_count for the user being unfollowed
    UPDATE users SET followers_count = GREATEST(followers_count - 1, 0) WHERE id = OLD.following_id;

    -- Decrement following_count for the user who is unfollowing
    UPDATE users SET following_count = GREATEST(following_count - 1, 0) WHERE id = OLD.follower_id;

    RETURN OLD;
END;
$$ LANGUAGE 'plpgsql';

-- ============================================================================
-- TRIGGERS
-- ============================================================================

-- Trigger to update counts on follow
CREATE TRIGGER update_follower_counts_on_insert
AFTER INSERT ON followers
FOR EACH ROW
EXECUTE FUNCTION increment_follower_counts();

-- Trigger to update counts on unfollow
CREATE TRIGGER update_follower_counts_on_delete
AFTER DELETE ON followers
FOR EACH ROW
EXECUTE FUNCTION decrement_follower_counts();
