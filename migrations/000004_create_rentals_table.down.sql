-- Drop indexes first
DROP INDEX IF EXISTS idx_rentals_due_date;
DROP INDEX IF EXISTS idx_rentals_status;
DROP INDEX IF EXISTS idx_rentals_book_id;
DROP INDEX IF EXISTS idx_rentals_user_id;

-- Drop the rentals table
DROP TABLE IF EXISTS rentals;
