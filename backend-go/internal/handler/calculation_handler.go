package handler

import (
	"fmt"
	"math"
	"net/http"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CalculationHandler struct {
	projectRepo        *repository.ProjectRepo
	scenarioRepo       *repository.ScenarioRepo
	catalogRepo        *repository.CatalogRepo
	irradiationService *service.IrradiationService
	pvCalculator       *service.PVSystemCalculator
	batteryCalculator  *service.BatteryCalculator
	financialModel     *service.FinancialModel
}

func NewCalculationHandler(
	projectRepo *repository.ProjectRepo,
	scenarioRepo *repository.ScenarioRepo,
	catalogRepo *repository.CatalogRepo,
	irradiationService *service.IrradiationService,
	pvCalculator *service.PVSystemCalculator,
	batteryCalculator *service.BatteryCalculator,
	financialModel *service.FinancialModel,
) *CalculationHandler {
	return &CalculationHandler{
		projectRepo:        projectRepo,
		scenarioRepo:       scenarioRepo,
		catalogRepo:        catalogRepo,
		irradiationService: irradiationService,
		pvCalculator:       pvCalculator,
		batteryCalculator:  batteryCalculator,
		financialModel:     financialModel,
	}
}

func (h *CalculationHandler) FullCalculation(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProjectID    string `json:"projectId" validate:"required"`
		ScenarioName string `json:"scenarioName"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	if input.ScenarioName == "" {
		input.ScenarioName = "Escenario 1"
	}

	ctx := r.Context()

	projectID, err := bson.ObjectIDFromHex(input.ProjectID)
	if err != nil {
		middleware.WriteError(w, middleware.NewAppError(400, "ID de proyecto inválido"))
		return
	}

	project, err := h.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}
	if project == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Proyecto no encontrado"))
		return
	}

	panel, err := h.catalogRepo.FindPanelByID(ctx, project.Equipment.PanelID)
	if err != nil || panel == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Panel no encontrado en el catálogo"))
		return
	}

	inverter, err := h.catalogRepo.FindInverterByID(ctx, project.Equipment.InverterID)
	if err != nil || inverter == nil {
		middleware.WriteError(w, middleware.NewAppError(404, "Inversor no encontrado en el catálogo"))
		return
	}

	// 1. Get irradiation per slope
	slopes := project.Roof.EffectiveSlopes()
	totalArea := project.Roof.TotalArea()

	type slopeIrrad struct {
		slope model.Slope
		irrad *service.IrradiationResult
	}
	slopeData := make([]slopeIrrad, len(slopes))
	irradCache := map[string]*service.IrradiationResult{}

	for i, s := range slopes {
		key := fmt.Sprintf("%.1f_%.1f", s.Tilt, s.Azimuth)
		if cached, ok := irradCache[key]; ok {
			slopeData[i] = slopeIrrad{s, cached}
		} else {
			result, err := h.irradiationService.GetIrradiation(
				ctx, project.Location.Latitude, project.Location.Longitude,
				s.Tilt, s.Azimuth, 0.2,
			)
			if err != nil {
				middleware.WriteError(w, err)
				return
			}
			irradCache[key] = result
			slopeData[i] = slopeIrrad{s, result}
		}
	}

	// Compute area-weighted average HSP and POA for system sizing
	weightedHSP := 0.0
	weightedPOA := make([]float64, 12)
	weightedGHI := make([]float64, 12)
	for _, sd := range slopeData {
		areaFraction := sd.slope.Area / totalArea
		weightedHSP += sd.irrad.AnnualAvgHSP * areaFraction
		for m := 0; m < 12; m++ {
			if m < len(sd.irrad.MonthlyPOA) {
				weightedPOA[m] += sd.irrad.MonthlyPOA[m] * areaFraction
			}
			if m < len(sd.irrad.MonthlyGHI) {
				weightedGHI[m] += sd.irrad.MonthlyGHI[m] * areaFraction
			}
		}
	}

	// 2. Calculate system design
	annualConsumption := 0.0
	for _, v := range project.Consumption.Monthly {
		annualConsumption += v
	}
	dailyConsumption := annualConsumption / 365

	var shadingLoss []float64
	if project.Roof.ShadingProfile.HasShading {
		shadingLoss = project.Roof.ShadingProfile.MonthlyLoss[:]
	}

	coveragePercentage := project.CoveragePercentage
	if coveragePercentage <= 0 {
		coveragePercentage = 100
	}

	designInput := service.SystemDesignInput{
		DailyConsumptionKwh:   dailyConsumption,
		MonthlyConsumptionKwh: project.Consumption.Monthly,
		AvgHSP:                weightedHSP,
		MonthlyHSP:            weightedPOA,
		Panel:                 *panel,
		Inverter:              *inverter,
		RoofArea:              totalArea,
		UsablePercentage:      project.Roof.UsablePercentage,
		ShadingLoss:           shadingLoss,
		SystemType:            project.SystemType,
		CoveragePercentage:    coveragePercentage,
	}

	var systemDesign service.SystemDesignResult
	if len(slopes) <= 1 {
		systemDesign = h.pvCalculator.Calculate(designInput)
	} else {
		// Distribute panels proportionally across slopes
		singleResult := h.pvCalculator.Calculate(designInput)
		totalPanels := singleResult.NumberOfPanels

		slopeAllocations := make([]service.SlopeAllocation, len(slopeData))
		usableArea := totalArea * (project.Roof.UsablePercentage / 100)
		assignedPanels := 0
		largestIdx := 0
		largestArea := 0.0

		for i, sd := range slopeData {
			slopeUsable := sd.slope.Area * (project.Roof.UsablePercentage / 100)
			panels := int(math.Floor(slopeUsable / usableArea * float64(totalPanels)))
			slopeAllocations[i] = service.SlopeAllocation{
				PanelCount: panels,
				MonthlyHSP: sd.irrad.MonthlyPOA,
				AvgHSP:     sd.irrad.AnnualAvgHSP,
			}
			assignedPanels += panels
			if sd.slope.Area > largestArea {
				largestArea = sd.slope.Area
				largestIdx = i
			}
		}
		// Assign remainder panels to the largest slope
		remainder := totalPanels - assignedPanels
		slopeAllocations[largestIdx].PanelCount += remainder

		systemDesign = h.pvCalculator.CalculateMultiSlope(designInput, slopeAllocations)
	}

	// Use first slope's irradiation as primary for scenario storage
	irradiation := slopeData[0].irrad

	// 3. Battery calculation (if off-grid or hybrid)
	var batteryBank *model.BatteryBank
	var batteryCost float64
	if project.SystemType != "on-grid" {
		batteryResult := h.batteryCalculator.Calculate(service.BatteryInput{
			DailyConsumptionKwh: dailyConsumption,
			AutonomyDays:        2,
			SystemVoltage:       48,
			BatteryType:         "lithium",
			BatteryCapacityAh:   100,
			BatteryVoltage:      12,
		})
		batteryCost = batteryResult.EstimatedCostCOP
		batteryBank = &model.BatteryBank{
			CapacityKwh:       batteryResult.CapacityKwh,
			AutonomyDays:      batteryResult.AutonomyDays,
			NumberOfBatteries: batteryResult.NumberOfBatteries,
			BankVoltage:       batteryResult.BankVoltage,
		}
	}

	// 4. Financial analysis
	financial := h.financialModel.Analyze(service.FinancialInput{
		InstalledKwp:        systemDesign.ActualPowerKwp,
		AnnualProductionKwh: systemDesign.AnnualProductionKwh,
		MonthlyProductionKwh: systemDesign.MonthlyProductionKwh,
		TariffPerKwh:        project.Consumption.TariffPerKwh,
		Estrato:             project.Consumption.Estrato,
		PanelCostCOP:        panel.CostCOP,
		NumberOfPanels:      systemDesign.NumberOfPanels,
		InverterCostCOP:     inverter.CostCOP,
		BatteryCostCOP:      batteryCost,
	})

	// 5. Save scenario
	scenario := &model.Scenario{
		ProjectID: projectID,
		Name:      input.ScenarioName,
		InputSnapshot: map[string]interface{}{
			"location":           project.Location,
			"consumption":        project.Consumption,
			"roof":               project.Roof,
			"systemType":         project.SystemType,
			"coveragePercentage": coveragePercentage,
			"panel": map[string]interface{}{
				"manufacturer": panel.Manufacturer,
				"model":        panel.Model,
				"powerWp":      panel.PowerWp,
			},
			"inverter": map[string]interface{}{
				"manufacturer": inverter.Manufacturer,
				"model":        inverter.Model,
				"ratedPowerKw": inverter.RatedPowerKw,
			},
		},
		Irradiation: model.IrradiationData{
			Source:       irradiation.Source,
			MonthlyGHI:  weightedGHI,
			MonthlyPOA:  weightedPOA,
			AnnualAvgHSP: weightedHSP,
		},
		SystemDesign: model.SystemDesign{
			RequiredPowerKwp:   systemDesign.RequiredPowerKwp,
			NumberOfPanels:     systemDesign.NumberOfPanels,
			ActualPowerKwp:     systemDesign.ActualPowerKwp,
			RoofUtilization:    systemDesign.RoofUtilization,
			InverterCapacityKw: systemDesign.InverterCapacityKw,
			StringConfiguration: model.StringConfiguration{
				PanelsPerString: systemDesign.StringConfiguration.PanelsPerString,
				NumberOfStrings: systemDesign.StringConfiguration.NumberOfStrings,
				StringVoltage:   systemDesign.StringConfiguration.StringVoltage,
				StringCurrent:   systemDesign.StringConfiguration.StringCurrent,
			},
			BatteryBank: batteryBank,
		},
		Production: model.Production{
			MonthlyKwh:      systemDesign.MonthlyProductionKwh,
			AnnualKwh:       systemDesign.AnnualProductionKwh,
			DegradationRate: 0.005,
			Yearly25:        systemDesign.Yearly25Production,
		},
		Financial: model.Financial{
			InstallationCostCOP: financial.InstallationCostCOP,
			MonthlySavingsCOP:   financial.MonthlySavingsCOP,
			AnnualSavingsCOP:    financial.AnnualSavingsCOP,
			PaybackYears:        financial.PaybackYears,
			IRRPercent:          financial.IRRPercent,
			NPVCOP:              financial.NPVCOP,
			CO2AvoidedTonsYear:  financial.CO2AvoidedTonsYear,
			CumulativeSavings25: financial.CumulativeSavings25,
			LCOE:                financial.LCOE,
		},
		Losses: systemDesign.Losses,
	}

	if err := h.scenarioRepo.Create(ctx, scenario); err != nil {
		middleware.WriteError(w, err)
		return
	}

	h.projectRepo.PushScenario(ctx, projectID, scenario.ID)

	middleware.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"scenario": scenario,
		"warnings": systemDesign.Warnings,
	})
}
