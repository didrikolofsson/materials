CREATE TABLE IF NOT EXISTS subjects (
	id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS materials (
	id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
	title VARCHAR(255) NOT NULL,
	description TEXT NULL,
	owner_teacher_id CHAR(36) NOT NULL,
	subject_id CHAR(36) NOT NULL,
	current_version_id CHAR(36) NULL,
	original_material_id CHAR(36) NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_materials_subject FOREIGN KEY (subject_id) REFERENCES subjects(id),
	CONSTRAINT fk_material_owner_teacher FOREIGN KEY (owner_teacher_id) REFERENCES teachers(id),
	CONSTRAINT fk_material_original_id FOREIGN KEY (original_material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS material_versions (
	id CHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
	material_id CHAR(36) NOT NULL,
	teacher_id CHAR(36) NOT NULL,
	-- Parent version id, if null this is the first version
	parent_id CHAR(36) NULL,
	version_number INT NOT NULL,
	summary TEXT NULL,
	content TEXT NOT NULL,
	is_main BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_version_material_id FOREIGN KEY (material_id) REFERENCES materials(id),
	CONSTRAINT fk_version_author_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id),
	CONSTRAINT fk_version_parent_id FOREIGN KEY (parent_id) REFERENCES material_versions(id) ON DELETE CASCADE
);