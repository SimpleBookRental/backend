-- Drop indexes first
DROP INDEX IF EXISTS idx_payments_transaction_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_rental_id;
DROP INDEX IF EXISTS idx_payments_user_id;

-- Drop the payments table
DROP TABLE IF EXISTS payments;
