// Package handlers provides HTTP handlers for the application.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/didrikolofsson/materials/internal/models"
	"github.com/didrikolofsson/materials/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	svc      *services.Services
	validate *validator.Validate
}

func New(svc *services.Services, validate *validator.Validate) *Handlers {
	return &Handlers{
		svc:      svc,
		validate: validate,
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "pong"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) ListTeachers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teachers, err := h.svc.ListTeachers(ctx)
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

func (h *Handlers) ListSubjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subjects, err := h.svc.ListSubjects(ctx)
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

func (h *Handlers) ListAllMaterials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	materials, err := h.svc.ListAllMaterials(ctx)
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

func (h *Handlers) ListMaterialVersionsByMaterialID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	materialID := chi.URLParam(r, "id")
	materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(materialIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialVersions, err := h.svc.ListMaterialVersionsByMaterialID(ctx, materialIDInt)
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

func (h *Handlers) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teacher, err := h.svc.GetTeacherByID(ctx, teacherIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teacher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetTeacherMaterials(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materials, err := h.svc.GetTeacherMaterials(ctx, teacherIDInt)
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

func (h *Handlers) CreateInitialTeacherMaterial(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req models.CreateMaterialRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialID, err := h.svc.CreateInitialTeacherMaterial(ctx, teacherIDInt, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]int64{"id": materialID}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetTeacherMaterialByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialID := chi.URLParam(r, "material_id")
	materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(materialIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	material, err := h.svc.GetTeacherMaterialByID(ctx, teacherIDInt, materialIDInt)
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

func (h *Handlers) UpdateTeacherMaterialByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialID := chi.URLParam(r, "material_id")
	materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(materialIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req models.UpdateMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	material, err := h.svc.UpdateTeacherMaterialByID(ctx, teacherIDInt, materialIDInt, req)
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

func (h *Handlers) DeleteTeacherMaterialByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teacherID := chi.URLParam(r, "id")
	teacherIDInt, err := strconv.ParseInt(teacherID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	materialID := chi.URLParam(r, "material_id")
	materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(teacherIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(materialIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteTeacherMaterialByID(ctx, teacherIDInt, materialIDInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) UpdateMaterialVersionMain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	materialID := chi.URLParam(r, "id")
	materialIDInt, err := strconv.ParseInt(materialID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	versionID := chi.URLParam(r, "version_id")
	versionIDInt, err := strconv.ParseInt(versionID, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(materialIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Var(versionIDInt, "required,min=1"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateMaterialVersionMain(ctx, materialIDInt, versionIDInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"id": versionID}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
