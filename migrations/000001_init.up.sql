CREATE SCHEMA bookshelfapp;

CREATE TABLE bookshelfapp.users (
    id          SERIAL PRIMARY KEY,
    version     BIGINT NOT NULL DEFAULT 1,
    username    VARCHAR(30) NOT NULL UNIQUE,
    email       VARCHAR(254) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_username_format CHECK (username ~ '^[A-Za-z0-9_]{3,30}$'),
    CONSTRAINT users_email_lowercase CHECK (email = lower(email)),
    CONSTRAINT users_email_format CHECK (email ~ '^[a-z0-9._+-]+@([a-z0-9-]+\.)+[a-z]{2,}$')
);

CREATE TABLE bookshelfapp.books (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(200) NOT NULL,
    author      VARCHAR(100) NOT NULL,
    year        SMALLINT NOT NULL,
    pages       SMALLINT NOT NULL,
    genres      VARCHAR(200)[] NOT NULL,
    description VARCHAR(2000) NOT NULL,
    score       SMALLINT,
    reads_count INT NOT NULL DEFAULT 0, 

    CONSTRAINT books_title_format  CHECK (char_length(title) >= 1 AND title = btrim(title)),
    CONSTRAINT books_author_format CHECK (char_length(author) >= 1 AND author = btrim(author)),
    CONSTRAINT books_year_format CHECK (year BETWEEN 1 AND EXTRACT(YEAR FROM now())),
    CONSTRAINT books_pages_format CHECK (pages BETWEEN 1 AND 32000),
    CONSTRAINT books_genres_not_empty CHECK (cardinality(genres) >= 1),
    CONSTRAINT books_description_format CHECK (char_length(description) >= 1 AND description = btrim(description)),
    CONSTRAINT books_score_range CHECK (score BETWEEN 1 AND 100),
    CONSTRAINT books_reads_count_range CHECK (reads_count >= 0)
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

    user_id     INT NOT NULL REFERENCES bookshelfapp.users(id) ON DELETE CASCADE,
    book_id     INT NOT NULL REFERENCES bookshelfapp.books(id),

    CONSTRAINT bookshelf_read_state CHECK (
        (read=TRUE AND read_at IS NOT NULL AND read_at >= added_at)
        OR
        (read=FALSE AND read_at IS NULL AND rating IS NULL AND review IS NULL) 
    ), 
    CONSTRAINT bookshelf_rating_range CHECK (rating BETWEEN 1 AND 100),
    CONSTRAINT bookshelf_review_range CHECK (char_length(review) >= 1 AND review = btrim(review)),
    CONSTRAINT bookshelf_pkey PRIMARY KEY (user_id, book_id)
);

CREATE INDEX bookshelf_book_id_idx ON bookshelfapp.bookshelf (book_id);

