CREATE TABLE IF NOT EXISTS material_proposals (
	id BIGINT PRIMARY KEY AUTO_INCREMENT,
	material_id BIGINT NOT NULL,
	material_version_id BIGINT NOT NULL,
	owner_teacher_id BIGINT NOT NULL,
	author_teacher_id BIGINT NOT NULL,
	title VARCHAR(255) NOT NULL,
	summary VARCHAR(255) NULL,
	description TEXT NULL,
	content TEXT NOT NULL,
	status ENUM("PENDING", "APPROVED", "REJECTED") NOT NULL DEFAULT "PENDING",
	decided_by_teacher_id BIGINT NULL,
	decided_at TIMESTAMP NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT fk_prop_material_id FOREIGN KEY (material_id) REFERENCES materials(id),
	CONSTRAINT fk_prop_material_version_id FOREIGN KEY (material_version_id) REFERENCES material_versions(id),
	CONSTRAINT fk_prop_owner_teacher_id FOREIGN KEY (owner_teacher_id) REFERENCES teachers(id),
	CONSTRAINT fk_prop_author_teacher_id FOREIGN KEY (author_teacher_id) REFERENCES teachers(id),
	CONSTRAINT fk_prop_decided_by_teacher_id FOREIGN KEY (decided_by_teacher_id) REFERENCES teachers(id),
	CONSTRAINT not_pending_when_decided CHECK (
		status != "PENDING"
		OR decided_at IS NOT NULL
	),
	CONSTRAINT not_rejected_when_decided CHECK (
		status != "REJECTED"
		OR decided_at IS NOT NULL
	)
)