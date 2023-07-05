-- Drop foreign key constraint
ALTER TABLE "accounts" DROP CONSTRAINT accounts_owner_fkey;

-- Drop unique index
DROP INDEX IF EXISTS accounts_owner_currency_idx;

-- Drop table
DROP TABLE IF EXISTS users;
