package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.projectService.Create(r.Context(), &project); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	respondJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Project not found")
		return
	}

	respondJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project.ID = id
	if err := h.projectService.Update(r.Context(), &project); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	respondJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := h.projectService.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get("company_id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	projects, err := h.projectService.ListByCompany(r.Context(), companyID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	respondJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	projects, err := h.projectService.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	respondJSON(w, http.StatusOK, projects)
}
