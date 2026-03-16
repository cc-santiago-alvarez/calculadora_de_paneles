package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/handler"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func New(
	projectHandler *handler.ProjectHandler,
	irradiationHandler *handler.IrradiationHandler,
	calculationHandler *handler.CalculationHandler,
	financialHandler *handler.FinancialHandler,
	catalogHandler *handler.CatalogHandler,
	reportsHandler *handler.ReportsHandler,
	geocodeHandler *handler.GeocodingHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	// Health check
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Projects
	r.Route("/api/v1/projects", func(r chi.Router) {
		r.Get("/", projectHandler.GetProjects)
		r.Get("/{id}", projectHandler.GetProjectByID)
		r.Post("/", projectHandler.CreateProject)
		r.Put("/{id}", projectHandler.UpdateProject)
		r.Delete("/{id}", projectHandler.DeleteProject)
		r.Get("/{id}/scenarios", projectHandler.GetProjectScenarios)
	})

	// Irradiation
	r.Post("/api/v1/irradiation/fetch", irradiationHandler.FetchIrradiation)

	// Calculation
	r.Post("/api/v1/calculation/full", calculationHandler.FullCalculation)

	// Financial
	r.Post("/api/v1/financial/analyze", financialHandler.AnalyzeFinancial)

	// Catalog
	r.Get("/api/v1/catalog/panels", catalogHandler.GetPanels)
	r.Get("/api/v1/catalog/inverters", catalogHandler.GetInverters)
	r.Post("/api/v1/catalog/sync", catalogHandler.SyncCEC)

	// Reports
	r.Post("/api/v1/reports/pdf", reportsHandler.GeneratePDF)
	r.Post("/api/v1/reports/excel", reportsHandler.GenerateExcel)

	// Geocoding
	r.Get("/api/v1/geocode/reverse", geocodeHandler.ReverseGeocode)
	r.Get("/api/v1/geocode/search", geocodeHandler.SearchAddress)

	return r
}
