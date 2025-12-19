package services

import "database/sql"

func nullStringToPointer(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func int64PtrToValue(ptr *int64) int64 {
	if ptr != nil {
		return *ptr
	}
	return 0
}
