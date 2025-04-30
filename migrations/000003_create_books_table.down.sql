-- Drop indexes first
DROP INDEX IF EXISTS idx_books_category_id;
DROP INDEX IF EXISTS idx_books_isbn;
DROP INDEX IF EXISTS idx_books_author;
DROP INDEX IF EXISTS idx_books_title;

-- Drop the books table
DROP TABLE IF EXISTS books;
