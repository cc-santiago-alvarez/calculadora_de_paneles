package handler

import (
	"net/http"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectHandler struct {
	projectRepo  *repository.ProjectRepo
	scenarioRepo *repository.ScenarioRepo
}

func NewProjectHandler(projectRepo *repository.ProjectRepo, scenarioRepo *repository.ScenarioRepo) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo, scenarioRepo: scenarioRepo}
}

func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectRepo.FindAll(r.Context())
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "ID inválido"))
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if project == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name" validate:"required,min=1,max=200"`
		Location    struct {
			Latitude    float64 `json:"latitude" validate:"required,min=-90,max=90"`
			Longitude   float64 `json:"longitude" validate:"required,min=-180,max=180"`
			Altitude    float64 `json:"altitude"`
			ClimateZone string  `json:"climateZone"`
			Department  string  `json:"department"`
			City        string  `json:"city"`
		} `json:"location" validate:"required"`
		Consumption struct {
			Monthly        [12]float64 `json:"monthly" validate:"required"`
			TariffPerKwh   float64     `json:"tariffPerKwh" validate:"required,gt=0"`
			Estrato        int         `json:"estrato" validate:"required,min=1,max=6"`
			ConnectionType string      `json:"connectionType"`
		} `json:"consumption" validate:"required"`
		Roof struct {
			Area             float64 `json:"area" validate:"required,gt=0"`
			Azimuth          float64 `json:"azimuth"`
			Tilt             float64 `json:"tilt"`
			UsablePercentage float64 `json:"usablePercentage"`
			ShadingProfile   struct {
				HasShading  bool        `json:"hasShading"`
				MonthlyLoss [12]float64 `json:"monthlyLoss"`
			} `json:"shadingProfile"`
		} `json:"roof" validate:"required"`
		SystemType         string  `json:"systemType"`
		CoveragePercentage float64 `json:"coveragePercentage"`
		Equipment  struct {
			PanelID       string `json:"panelId" validate:"required"`
			InverterID    string `json:"inverterId" validate:"required"`
			PanelOverride *struct {
				Watts float64 `json:"watts"`
				Area  float64 `json:"area"`
			} `json:"panelOverride,omitempty"`
		} `json:"equipment" validate:"required"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	// Set defaults
	if input.Consumption.ConnectionType == "" {
		input.Consumption.ConnectionType = "monofasica"
	}
	if input.SystemType == "" {
		input.SystemType = "on-grid"
	}
	if input.CoveragePercentage <= 0 || input.CoveragePercentage > 100 {
		input.CoveragePercentage = 100
	}
	if input.Roof.Tilt == 0 {
		input.Roof.Tilt = 10
	}
	if input.Roof.UsablePercentage == 0 {
		input.Roof.UsablePercentage = 80
	}

	panelID, _ := bson.ObjectIDFromHex(input.Equipment.PanelID)
	inverterID, _ := bson.ObjectIDFromHex(input.Equipment.InverterID)

	project := &model.Project{
		Name: input.Name,
		Location: model.Location{
			Latitude:    input.Location.Latitude,
			Longitude:   input.Location.Longitude,
			Altitude:    input.Location.Altitude,
			ClimateZone: input.Location.ClimateZone,
			Department:  input.Location.Department,
			City:        input.Location.City,
		},
		Consumption: model.Consumption{
			Monthly:        input.Consumption.Monthly,
			TariffPerKwh:   input.Consumption.TariffPerKwh,
			Estrato:        input.Consumption.Estrato,
			ConnectionType: input.Consumption.ConnectionType,
		},
		Roof: model.Roof{
			Area:             input.Roof.Area,
			Azimuth:          input.Roof.Azimuth,
			Tilt:             input.Roof.Tilt,
			UsablePercentage: input.Roof.UsablePercentage,
			ShadingProfile: model.ShadingProfile{
				HasShading:  input.Roof.ShadingProfile.HasShading,
				MonthlyLoss: input.Roof.ShadingProfile.MonthlyLoss,
			},
		},
		SystemType:         input.SystemType,
		CoveragePercentage: input.CoveragePercentage,
		Equipment: model.Equipment{
			PanelID:    panelID,
			InverterID: inverterID,
		},
	}

	if input.Equipment.PanelOverride != nil {
		project.Equipment.PanelOverride = &model.PanelOverride{
			Watts: input.Equipment.PanelOverride.Watts,
			Area:  input.Equipment.PanelOverride.Area,
		}
	}

	if err := h.projectRepo.Create(r.Context(), project); err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "ID inválido"))
		return
	}

	var body map[string]interface{}
	if err := middleware.ReadJSON(r, &body); err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "JSON inválido"))
		return
	}

	project, err := h.projectRepo.Update(r.Context(), id, body)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if project == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}
	middleware.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "ID inválido"))
		return
	}

	if err := h.projectRepo.Delete(r.Context(), id); err != nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}
	h.scenarioRepo.DeleteByProjectID(r.Context(), id)
	middleware.WriteJSON(w, http.StatusOK, map[string]string{"message": "Proyecto eliminado"})
}

func (h *ProjectHandler) GetProjectScenarios(w http.ResponseWriter, r *http.Request) {
	id, err := bson.ObjectIDFromHex(chi.URLParam(r, "id"))
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "ID inválido"))
		return
	}

	scenarios, err := h.scenarioRepo.FindByProjectID(r.Context(), id)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusOK, scenarios)
}
