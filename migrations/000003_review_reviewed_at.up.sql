ALTER TABLE bookshelfapp.bookshelf ADD COLUMN reviewed_at TIMESTAMPTZ;

UPDATE bookshelfapp.bookshelf SET reviewed_at = read_at WHERE review IS NOT NULL;

ALTER TABLE bookshelfapp.bookshelf
    ADD CONSTRAINT bookshelf_reviewed_at_state CHECK (
        (review IS NULL AND reviewed_at IS NULL)
        OR
        (review IS NOT NULL AND reviewed_at IS NOT NULL AND reviewed_at >= read_at)
    );

CREATE INDEX bookshelf_reviewed_at_idx
    ON bookshelfapp.bookshelf (reviewed_at DESC) WHERE review IS NOT NULL;