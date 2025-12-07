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
	v *validator.Validate
}

func NewHandlerDomainSchool(
	service *ServiceDomainSchool, validate *validator.Validate,
) *HandlerDomainSchool {
	return &HandlerDomainSchool{s: service, v: validate}
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

func (h *HandlerDomainSchool) handleListAllMaterialsByTeacherID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	teacherID, err := strconv.ParseInt(chi.URLParam(r, "teacher_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	materials, err := h.s.ListAllMaterialsByTeacherID(ctx, models.GenericID(teacherID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(materials); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerDomainSchool) handleListMaterialVersionsByMaterialID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	materialID, err := strconv.ParseInt(chi.URLParam(r, "material_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	materialVersions, err := h.s.ListMaterialVersionsByMaterialID(ctx, models.GenericID(materialID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(materialVersions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerDomainSchool) handleUpdateMaterialMainVersionByVersionID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	materialID, err := strconv.ParseInt(chi.URLParam(r, "material_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := models.UpdateMainVersionForMaterialParams{
		MaterialID: models.GenericID(materialID),
	}
	var body models.UpdateMainVersionForMaterialBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	request := models.UpdateMainVersionForMaterialRequest{
		Params: params,
		Body:   body,
	}
	if err = h.v.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.s.UpdateMainVersionForMaterialByVersionID(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.UpdateMainVersionForMaterialResponse{
		MaterialID:        request.Params.MaterialID,
		MaterialVersionID: request.Body.MaterialVersionID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HandlerDomainSchool) handleCreateMaterial(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
		TeacherID: models.GenericID(teacherID),
		SubjectID: models.GenericID(subjectID),
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

	if err = h.v.Struct(request); err != nil {
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

func (h *HandlerDomainSchool) handleGetMaterialVersionByVersionID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	materialID, err := strconv.ParseInt(chi.URLParam(r, "material_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	versionID, err := strconv.ParseInt(chi.URLParam(r, "version_id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	materialVersion, err := h.s.GetMaterialVersionByVersionID(ctx, models.GenericID(materialID), models.GenericID(versionID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(materialVersion); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
