-- Drop the CASCADE foreign key constraint
ALTER TABLE material_versions DROP FOREIGN KEY fk_material_id;
-- Recreate it without CASCADE (original state)
ALTER TABLE material_versions
ADD CONSTRAINT fk_material_id FOREIGN KEY (material_id) REFERENCES materials(id);