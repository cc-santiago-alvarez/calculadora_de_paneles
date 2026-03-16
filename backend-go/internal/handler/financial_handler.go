package handler

import (
	"net/http"

	"github.com/dev13/calculadora-paneles-backend/internal/middleware"
	"github.com/dev13/calculadora-paneles-backend/internal/service"
)

type FinancialHandler struct {
	financialModel *service.FinancialModel
}

func NewFinancialHandler(financialModel *service.FinancialModel) *FinancialHandler {
	return &FinancialHandler{financialModel: financialModel}
}

func (h *FinancialHandler) AnalyzeFinancial(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AnnualProductionKwh float64 `json:"annualProductionKwh" validate:"required,gt=0"`
		TariffPerKwh        float64 `json:"tariffPerKwh" validate:"required,gt=0"`
		Estrato             int     `json:"estrato" validate:"required,min=1,max=6"`
		InstalledKwp        float64 `json:"installedKwp" validate:"required,gt=0"`
		SystemLifeYears     int     `json:"systemLifeYears"`
		DiscountRate        float64 `json:"discountRate"`
		DegradationRate     float64 `json:"degradationRate"`
		TariffEscalation    float64 `json:"tariffEscalation"`
		InstallationCostCOP float64 `json:"installationCostCOP"`
	}

	if !middleware.ValidateAndDecode(w, r, &input) {
		return
	}

	// Generate approximate monthly from annual
	monthlyProductionKwh := make([]float64, 12)
	for i := range monthlyProductionKwh {
		monthlyProductionKwh[i] = input.AnnualProductionKwh / 12
	}

	inverterCostCOP := input.InstallationCostCOP
	if inverterCostCOP <= 0 {
		inverterCostCOP = input.InstalledKwp * 3500000
	}

	result := h.financialModel.Analyze(service.FinancialInput{
		InstalledKwp:        input.InstalledKwp,
		AnnualProductionKwh: input.AnnualProductionKwh,
		MonthlyProductionKwh: monthlyProductionKwh,
		TariffPerKwh:        input.TariffPerKwh,
		Estrato:             input.Estrato,
		PanelCostCOP:        0,
		NumberOfPanels:      0,
		InverterCostCOP:     inverterCostCOP,
		SystemLifeYears:     input.SystemLifeYears,
		DiscountRate:        input.DiscountRate,
		DegradationRate:     input.DegradationRate,
		TariffEscalation:    input.TariffEscalation,
	})

	middleware.WriteJSON(w, http.StatusOK, result)
}
