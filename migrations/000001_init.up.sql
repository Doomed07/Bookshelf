CREATE SCHEMA bookshelfapp;

CREATE TABLE bookshelfapp.users (
    id              SERIAL PRIMARY KEY,
    version         BIGINT NOT NULL DEFAULT 1,
    username        VARCHAR(30) NOT NULL,
    password_hash   TEXT NOT NULL,
    email           VARCHAR(254) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_username_format CHECK (username ~ '^[A-Za-z0-9_]{3,30}$'),
    CONSTRAINT users_password_hash_len CHECK (length(password_hash) = 60),
    CONSTRAINT users_email_lowercase CHECK (email = lower(email)),
    CONSTRAINT users_email_format CHECK (email ~ '^[a-z0-9._+-]+@([a-z0-9-]+\.)+[a-z]{2,}$')
);

CREATE UNIQUE INDEX users_username_lower_uidx ON bookshelfapp.users (lower(username));

CREATE TABLE bookshelfapp.sessions (
    token_hash BYTEA PRIMARY KEY CHECK (octet_length(token_hash) = 32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    CHECK (expires_at > created_at),

    user_id    INT NOT NULL REFERENCES bookshelfapp.users(id) ON DELETE CASCADE
);

CREATE INDEX sessions_user_id_idx    ON bookshelfapp.sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON bookshelfapp.sessions (expires_at);

CREATE TABLE bookshelfapp.books (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    author      VARCHAR(100) NOT NULL,
    year        SMALLINT NOT NULL,
    pages       SMALLINT NOT NULL,
    genres      VARCHAR(200)[] NOT NULL,
    description VARCHAR(2000) NOT NULL,
    rating_sum   INT NOT NULL DEFAULT 0,
    rating_count INT NOT NULL DEFAULT 0,
    reads_count  INT NOT NULL DEFAULT 0,
    score        SMALLINT GENERATED ALWAYS AS (
        CASE WHEN rating_count > 0
             THEN round(rating_sum::numeric / rating_count)::smallint
        END
    ) STORED,

    CONSTRAINT books_title_format  CHECK (char_length(title) >= 1 AND title = btrim(title)),
    CONSTRAINT books_author_format CHECK (char_length(author) >= 1 AND author = btrim(author)),
    CONSTRAINT books_year_format CHECK (year BETWEEN 1 AND EXTRACT(YEAR FROM now())),
    CONSTRAINT books_pages_format CHECK (pages BETWEEN 1 AND 32000),
    CONSTRAINT books_genres_not_empty CHECK (cardinality(genres) >= 1),
    CONSTRAINT books_description_format CHECK (char_length(description) >= 1 AND description = btrim(description)),
    CONSTRAINT books_score_range CHECK (score BETWEEN 1 AND 100),
    CONSTRAINT books_reads_count_range CHECK (reads_count >= 0),
    CONSTRAINT books_rating_sum_range   CHECK (rating_sum >= 0),
    CONSTRAINT books_rating_count_range CHECK (rating_count >= 0)
);

CREATE UNIQUE INDEX books_title_author_uniq
    ON bookshelfapp.books (lower(title), lower(author));

CREATE TABLE bookshelfapp.bookshelf (
    version     BIGINT NOT NULL DEFAULT 1,
    read        BOOLEAN NOT NULL DEFAULT FALSE,
    rating      SMALLINT,
    review      VARCHAR(5000),
    added_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at     TIMESTAMPTZ, 
    reviewed_at TIMESTAMPTZ,

    user_id     INT NOT NULL REFERENCES bookshelfapp.users(id) ON DELETE CASCADE,
    book_id     INT NOT NULL REFERENCES bookshelfapp.books(id),

    CONSTRAINT bookshelf_read_state CHECK (
        (read=TRUE AND read_at IS NOT NULL AND read_at >= added_at)
        OR
        (read=FALSE AND read_at IS NULL AND rating IS NULL AND review IS NULL) 
    ), 
    CONSTRAINT bookshelf_rating_range CHECK (rating BETWEEN 1 AND 100),
    CONSTRAINT bookshelf_review_range CHECK (char_length(review) >= 1 AND review = btrim(review)),
    CONSTRAINT bookshelf_pkey PRIMARY KEY (user_id, book_id),
    CONSTRAINT bookshelf_reviewed_at_state CHECK (
    (review IS NULL AND reviewed_at IS NULL)
    OR
    (review IS NOT NULL AND reviewed_at IS NOT NULL AND reviewed_at >= read_at))
);

CREATE INDEX bookshelf_book_id_idx ON bookshelfapp.bookshelf (book_id);
CREATE INDEX books_score_idx       ON bookshelfapp.books (score DESC NULLS LAST, id);
CREATE INDEX books_reads_count_idx ON bookshelfapp.books (reads_count DESC, id);
CREATE INDEX bookshelf_reviewed_at_idx
    ON bookshelfapp.bookshelf (reviewed_at DESC) WHERE review IS NOT NULL;

CREATE FUNCTION bookshelfapp.update_book_stats() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'UPDATE'
        AND OLD.read = NEW.read
        AND OLD.rating IS NOT DISTINCT FROM NEW.rating THEN
        RETURN NULL;
    END IF;

    IF TG_OP IN ('UPDATE', 'DELETE') THEN
        UPDATE bookshelfapp.books SET
            rating_sum   = rating_sum   - COALESCE(OLD.rating, 0),
            rating_count = rating_count - (OLD.rating IS NOT NULL)::int,
            reads_count  = reads_count  - OLD.read::int
        WHERE id = OLD.book_id;
    END IF;

    IF TG_OP IN ('INSERT', 'UPDATE') THEN
        UPDATE bookshelfapp.books SET
            rating_sum   = rating_sum   + COALESCE(NEW.rating, 0),
            rating_count = rating_count + (NEW.rating IS NOT NULL)::int,
            reads_count  = reads_count  + NEW.read::int
        WHERE id = NEW.book_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bookshelf_update_book_stats
AFTER INSERT OR DELETE OR UPDATE OF read, rating ON bookshelfapp.bookshelf
FOR EACH ROW EXECUTE FUNCTION bookshelfapp.update_book_stats();
