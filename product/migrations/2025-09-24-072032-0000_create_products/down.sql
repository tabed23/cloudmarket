-- This file should undo anything in `up.sql`
DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE products;