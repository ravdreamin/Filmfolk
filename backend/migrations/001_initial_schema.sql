-- FilmFolk Database Schema - Simplified
-- Essential tables for movie reviews platform

-- ============================================================================
-- ENUMS
-- ============================================================================

CREATE TYPE auth_provider AS ENUM ('email', 'google');
CREATE TYPE account_status AS ENUM ('active', 'suspended', 'banned');

-- ============================================================================
-- USERS TABLE
-- ============================================================================

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT,
    auth_provider auth_provider NOT NULL DEFAULT 'email',
    provider_id VARCHAR(255),
    status account_status NOT NULL DEFAULT 'active',
    avatar_url TEXT,
    bio TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,

    CONSTRAINT unique_provider_id UNIQUE(auth_provider, provider_id)
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

COMMENT ON TABLE users IS 'User accounts with multi-provider authentication';
COMMENT ON COLUMN users.password_hash IS 'NULL for OAuth users, required for email auth';

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

    -- Aggregated stats
    average_rating DECIMAL(3,2),
    total_reviews INTEGER DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_movie UNIQUE(title, release_year)
);

CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX idx_movies_average_rating ON movies(average_rating DESC);
CREATE INDEX idx_movies_release_year ON movies(release_year DESC);

COMMENT ON TABLE movies IS 'Movie catalog with TMDB integration';

-- ============================================================================
-- REVIEWS
-- ============================================================================

CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    review_text TEXT NOT NULL,

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
CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC);
CREATE INDEX idx_reviews_rating ON reviews(rating DESC);

COMMENT ON TABLE reviews IS 'User reviews with ratings and thread control';
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

COMMENT ON TABLE review_likes IS 'Tracks which users liked which reviews';

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

COMMENT ON TABLE comment_likes IS 'Tracks which users liked which comments';

-- ============================================================================
-- REFRESH TOKENS (JWT Authentication)
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

COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for secure authentication';

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_movies_updated_at BEFORE UPDATE ON movies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_review_comments_updated_at BEFORE UPDATE ON review_comments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
