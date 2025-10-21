CREATE TYPE user_role AS ENUM (
    'user',
    'moderator',
    'admin'
);

CREATE TYPE movie_status AS ENUM (
    'pending_approval',
    'approved',
    'rejected'
);

CREATE TYPE review_status AS ENUM (
    'pending_moderation',
    'published',
    'rejected'
);



-- USERS TABLE: Manages user accounts and roles

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE users IS 'Stores user account information, including authentication details and role.';
COMMENT ON COLUMN users.role IS 'Defines the access level of the user (user, moderator, admin).';



-- MOVIES TABLE: The central catalog of movies with moderation


CREATE TABLE movies (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    release_year INTEGER NOT NULL,
    genre VARCHAR(100),
    summary TEXT,
    -- This can store the ID from an external API like TMDb
    external_api_id VARCHAR(50) UNIQUE,
    status movie_status NOT NULL DEFAULT 'pending_approval',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- A movie is uniquely identified by its title and year of release.
    UNIQUE(title, release_year)
);

COMMENT ON TABLE movies IS 'Stores the central catalog of movies.';
COMMENT ON COLUMN movies.summary IS 'Can be populated by an AI to provide a unique movie summary.';
COMMENT ON COLUMN movies.external_api_id IS 'Stores the unique ID from an external data source like TMDb for easy lookups.';
COMMENT ON COLUMN movies.status IS 'Tracks the moderation status of a movie entry.';


-- REVIEWS TABLE: Connects users and movies, with content moderation

CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    -- Foreign keys with ON DELETE CASCADE to maintain data integrity.
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id BIGINT NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL,
    review_text TEXT NOT NULL,
    sentiment VARCHAR(50),
    status review_status NOT NULL DEFAULT 'published',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Ensure a user can only review a specific movie once.
    UNIQUE(user_id, movie_id),

    -- Add a check constraint to ensure the rating is within a valid range.
    CONSTRAINT rating_check CHECK (rating >= 1 AND rating <= 10)
);

COMMENT ON TABLE reviews IS 'Stores user-submitted reviews, linking users and movies.';
COMMENT ON COLUMN reviews.sentiment IS 'Populated by an AI to classify the tone of the review (e.g., positive, negative).';
COMMENT ON COLUMN reviews.status IS 'Tracks the moderation status of a review, managed by the AI filter and human moderators.';
COMMENT ON CONSTRAINT rating_check ON reviews IS 'Ensures review ratings are between 1 and 10.';


-- INDEXES: Improve query performance on frequently searched columns

-- Create indexes on foreign key columns for faster joins.
CREATE INDEX ON reviews (user_id);
CREATE INDEX ON reviews (movie_id);

-- Create an index on movie status for quickly finding movies needing approval.
CREATE INDEX ON movies (status);

-- Create an index on review status for quickly finding reviews needing moderation.
CREATE INDEX ON reviews (status);
