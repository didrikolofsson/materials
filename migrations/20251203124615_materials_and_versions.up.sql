CREATE TABLE IF NOT EXISTS subjects (
	id CHAR(36) PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS materials (
	id CHAR(36) PRIMARY KEY NOT NULL,
	teacher_id CHAR(36) NOT NULL,
	subject_id CHAR(36) NULL,
	original_material_id CHAR(36) NOT NULL,
	current_version_id CHAR(36) NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_teacher_id FOREIGN KEY (teacher_id) REFERENCES teachers(id),
	CONSTRAINT fk_subject_id FOREIGN KEY (subject_id) REFERENCES subjects(id),
	CONSTRAINT fk_original_material_id FOREIGN KEY (original_material_id) REFERENCES materials(id)
);
CREATE TABLE IF NOT EXISTS material_versions (
	id CHAR(36) PRIMARY KEY NOT NULL,
	title VARCHAR(255) NOT NULL,
	summary TEXT NULL,
	description TEXT NULL,
	content TEXT NOT NULL,
	version_number INT NOT NULL,
	is_main BOOLEAN NOT NULL DEFAULT FALSE,
	material_id CHAR(36) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_material_id FOREIGN KEY (material_id) REFERENCES materials(id),
	CONSTRAINT unique_version_number_per_material UNIQUE (material_id, version_number)
);