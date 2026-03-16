package handler

import (
	"fmt"
	"net/http"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportsHandler struct {
	projectRepo  *repository.ProjectRepo
	scenarioRepo *repository.ScenarioRepo
	pdfGen       *service.ReportPDFGenerator
	excelGen     *service.ReportExcelGenerator
}

func NewReportsHandler(
	projectRepo *repository.ProjectRepo,
	scenarioRepo *repository.ScenarioRepo,
	pdfGen *service.ReportPDFGenerator,
	excelGen *service.ReportExcelGenerator,
) *ReportsHandler {
	return &ReportsHandler{
		projectRepo:  projectRepo,
		scenarioRepo: scenarioRepo,
		pdfGen:       pdfGen,
		excelGen:     excelGen,
	}
}

func (h *ReportsHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProjectID  string `json:"projectId" validate:"required"`
		ScenarioID string `json:"scenarioId" validate:"required"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	ctx := r.Context()

	projectID, _ := bson.ObjectIDFromHex(input.ProjectID)
	project, err := h.projectRepo.FindByID(ctx, projectID)
	if err != nil || project == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}

	scenarioID, _ := bson.ObjectIDFromHex(input.ScenarioID)
	scenario, err := h.scenarioRepo.FindByID(ctx, scenarioID)
	if err != nil || scenario == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Escenario no encontrado"))
		return
	}

	pdfBytes, err := h.pdfGen.Generate(*project, *scenario)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="reporte-%s.pdf"`, project.Name))
	w.Write(pdfBytes)
}

func (h *ReportsHandler) GenerateExcel(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProjectID  string `json:"projectId" validate:"required"`
		ScenarioID string `json:"scenarioId" validate:"required"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	ctx := r.Context()

	projectID, _ := bson.ObjectIDFromHex(input.ProjectID)
	project, err := h.projectRepo.FindByID(ctx, projectID)
	if err != nil || project == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}

	scenarioID, _ := bson.ObjectIDFromHex(input.ScenarioID)
	scenario, err := h.scenarioRepo.FindByID(ctx, scenarioID)
	if err != nil || scenario == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Escenario no encontrado"))
		return
	}

	excelBytes, err := h.excelGen.Generate(*project, *scenario)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="reporte-%s.xlsx"`, project.Name))
	w.Write(excelBytes)
}
