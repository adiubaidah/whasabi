-- Drop the trigger
DROP TRIGGER IF EXISTS update_updated_at ON services;

-- Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the table
DROP TABLE IF EXISTS services;