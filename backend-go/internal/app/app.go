package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dev13/calculadora-paneles-backend/data"
	"github.com/dev13/calculadora-paneles-backend/internal/config"
	"github.com/dev13/calculadora-paneles-backend/internal/handler"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/dev13/calculadora-paneles-backend/internal/router"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
)

// App contains all initialized infrastructure.
type App struct {
	Router http.Handler
	DB     *repository.MongoDB
	Config *config.Config
}

// New initializes the full application (DB, repos, services, handlers, router).
func New() (*App, error) {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := repository.ConnectMongoDB(ctx, cfg.MongoDBURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := db.EnsureIndexes(ctx); err != nil {
		log.Printf("Warning: failed to ensure indexes: %v", err)
	}

	projectRepo := repository.NewProjectRepo(db.Database)
	scenarioRepo := repository.NewScenarioRepo(db.Database)
	catalogRepo := repository.NewCatalogRepo(db.Database)
	cacheRepo := repository.NewIrradiationCacheRepo(db.Database)

	seedCatalog(ctx, catalogRepo)

	irradiationService := service.NewIrradiationService(cacheRepo, cfg.PVGISBaseURL, cfg.NASAPowerURL)
	pvCalculator := service.NewPVSystemCalculator()
	batteryCalculator := service.NewBatteryCalculator()
	financialModel := service.NewFinancialModel()
	pdfGenerator := service.NewReportPDFGenerator()
	excelGenerator := service.NewReportExcelGenerator()
	cecCatalogService := service.NewCECCatalogService(catalogRepo)

	projectHandler := handler.NewProjectHandler(projectRepo, scenarioRepo)
	irradiationHandler := handler.NewIrradiationHandler(irradiationService)
	calculationHandler := handler.NewCalculationHandler(
		projectRepo, scenarioRepo, catalogRepo,
		irradiationService, pvCalculator, batteryCalculator, financialModel,
	)
	financialHandler := handler.NewFinancialHandler(financialModel)
	catalogHandler := handler.NewCatalogHandler(catalogRepo, cecCatalogService)
	reportsHandler := handler.NewReportsHandler(projectRepo, scenarioRepo, pdfGenerator, excelGenerator)
	geocodeHandler := handler.NewGeocodingHandler()

	r := router.New(
		projectHandler,
		irradiationHandler,
		calculationHandler,
		financialHandler,
		catalogHandler,
		reportsHandler,
		geocodeHandler,
	)

	return &App{
		Router: r,
		DB:     db,
		Config: cfg,
	}, nil
}

// StartHTTPServer starts the HTTP server in a goroutine and returns it.
func (a *App) StartHTTPServer() *http.Server {
	addr := fmt.Sprintf(":%d", a.Config.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("✓ Backend running on http://localhost%s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	return server
}

// Close closes database connections.
func (a *App) Close() {
	a.DB.Close(context.Background())
}

func seedCatalog(ctx context.Context, catalogRepo *repository.CatalogRepo) {
	panels := make([]model.PanelCatalog, len(data.DefaultPanels))
	for i, p := range data.DefaultPanels {
		panels[i] = model.PanelCatalog{
			Manufacturer:  p.Manufacturer,
			Model:         p.Model,
			Type:          p.Type,
			PowerWp:       p.PowerWp,
			Efficiency:    p.Efficiency,
			Area:          p.Area,
			Voc:           p.Voc,
			Isc:           p.Isc,
			Vmp:           p.Vmp,
			Imp:           p.Imp,
			TempCoeffPmax: p.TempCoeffPmax,
			TempCoeffVoc:  p.TempCoeffVoc,
			NOCT:          p.NOCT,
			Weight:        p.Weight,
			PanelDimensions: model.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
			},
			Warranty: p.Warranty,
			CostCOP:  p.CostCOP,
			IsActive: true,
		}
	}
	if n, err := catalogRepo.UpsertPanels(ctx, panels); err != nil {
		log.Printf("Warning: failed to upsert panels: %v", err)
	} else if n > 0 {
		log.Printf("Catalog: upserted %d panels", n)
	}

	inverters := make([]model.InverterCatalog, len(data.DefaultInverters))
	for i, inv := range data.DefaultInverters {
		inverters[i] = model.InverterCatalog{
			Manufacturer:    inv.Manufacturer,
			Model:           inv.Model,
			Type:            inv.Type,
			RatedPowerKw:    inv.RatedPowerKw,
			MaxDCPowerKw:    inv.MaxDCPowerKw,
			Efficiency:      inv.Efficiency,
			MPPTCount:       inv.MPPTCount,
			MPPTVoltageMin:  inv.MPPTVoltageMin,
			MPPTVoltageMax:  inv.MPPTVoltageMax,
			MaxInputVoltage: inv.MaxInputVoltage,
			MaxInputCurrent: inv.MaxInputCurrent,
			OutputVoltage:   inv.OutputVoltage,
			OutputPhases:    inv.OutputPhases,
			HasBatteryPort:  inv.HasBatteryPort,
			Weight:          inv.Weight,
			Warranty:        inv.Warranty,
			CostCOP:         inv.CostCOP,
			IsActive:        true,
		}
	}
	if n, err := catalogRepo.UpsertInverters(ctx, inverters); err != nil {
		log.Printf("Warning: failed to upsert inverters: %v", err)
	} else if n > 0 {
		log.Printf("Catalog: upserted %d inverters", n)
	}
}
