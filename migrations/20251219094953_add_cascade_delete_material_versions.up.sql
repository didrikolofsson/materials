-- Drop the existing foreign key constraint
ALTER TABLE material_versions DROP FOREIGN KEY fk_material_id;
-- Recreate it with ON DELETE CASCADE
ALTER TABLE material_versions
ADD CONSTRAINT fk_material_id FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE;