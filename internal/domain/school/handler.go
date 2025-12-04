package school

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SubjectsHandler struct {
	r SubjectsRepository
}

func NewSubjectsHandler(r SubjectsRepository) *SubjectsHandler {
	return &SubjectsHandler{r: r}
}

func (h *SubjectsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/subjects", h.handleListSubjects)
}

func (h *SubjectsHandler) handleListSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subjects, err := h.r.List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subjects)
}
