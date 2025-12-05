package school

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/didrikolofsson/materials/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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

func (h *HandlerDomainSchool) handleCreateMaterial(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	validate := validator.New()

	teacherID, err := strconv.ParseInt(chi.URLParam(r, "teacher_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subjectID, err := strconv.ParseInt(chi.URLParam(r, "subject_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := models.CreateMaterialParams{
		TeacherID: models.TeacherID(teacherID),
		SubjectID: models.SubjectID(subjectID),
	}

	var body models.CreateMaterialBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	request := models.CreateMaterialRequest{
		Params: params,
		Body:   body,
	}

	if err = validate.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialID, materialVersionID, err := h.s.CreateMaterialWithInitialVersion(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.CreateMaterialResponse{
		MaterialID:        materialID,
		MaterialVersionID: materialVersionID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
