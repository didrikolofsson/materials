package school

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HandlerDomainSchool struct {
	s *ServiceDomainSchool
}

func NewHandlerDomainSchool(service *ServiceDomainSchool) *HandlerDomainSchool {
	return &HandlerDomainSchool{s: service}
}

func (h *HandlerDomainSchool) handleListSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.s.ListSubjects(ctx)
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

func (h *HandlerDomainSchool) handleListTeachers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teachers, err := h.s.ListTeachers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teachers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerDomainSchool) handleGetMaterialByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	material, err := h.s.GetMaterialByID(ctx, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(material); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
