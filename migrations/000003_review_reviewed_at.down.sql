DROP INDEX IF EXISTS bookshelfapp.bookshelf_reviewed_at_idx;
ALTER TABLE bookshelfapp.bookshelf DROP CONSTRAINT IF EXISTS bookshelf_reviewed_at_state;
ALTER TABLE bookshelfapp.bookshelf DROP COLUMN IF EXISTS reviewed_at;