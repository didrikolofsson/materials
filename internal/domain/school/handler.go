package school

import (
	"encoding/json"
	"net/http"
)

type SchoolHandler struct {
	subjectService *SubjectService
}

func NewSchoolHandler(subjectService *SubjectService) *SchoolHandler {
	return &SchoolHandler{subjectService: subjectService}
}

func (h *SchoolHandler) handleListSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.subjectService.ListSubjects(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(subjects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
