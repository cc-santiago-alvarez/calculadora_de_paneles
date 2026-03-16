package service

import (
	"math"
	"strconv"

	"github.com/dev13/calculadora-paneles-backend/data"
	"github.com/dev13/calculadora-paneles-backend/internal/calc"
	"github.com/dev13/calculadora-paneles-backend/internal/config"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
)

type FinancialInput struct {
	InstalledKwp        float64   `json:"installedKwp"`
	AnnualProductionKwh float64   `json:"annualProductionKwh"`
	MonthlyProductionKwh []float64 `json:"monthlyProductionKwh"`
	TariffPerKwh        float64   `json:"tariffPerKwh"`
	Estrato             int       `json:"estrato"`
	PanelCostCOP        float64   `json:"panelCostCOP"`
	NumberOfPanels      int       `json:"numberOfPanels"`
	InverterCostCOP     float64   `json:"inverterCostCOP"`
	BatteryCostCOP      float64   `json:"batteryCostCOP,omitempty"`
	SystemLifeYears     int       `json:"systemLifeYears,omitempty"`
	DiscountRate        float64   `json:"discountRate,omitempty"`
	DegradationRate     float64   `json:"degradationRate,omitempty"`
	TariffEscalation    float64   `json:"tariffEscalation,omitempty"`
}

type CostBreakdown struct {
	Panels    float64 `json:"panels"`
	Inverter  float64 `json:"inverter"`
	Batteries float64 `json:"batteries"`
	BOS       float64 `json:"bos"`
	Labor     float64 `json:"labor"`
}

type FinancialResult struct {
	InstallationCostCOP float64            `json:"installationCostCOP"`
	CostBreakdown       CostBreakdown      `json:"costBreakdown"`
	MonthlySavingsCOP   []float64          `json:"monthlySavingsCOP"`
	AnnualSavingsCOP    float64            `json:"annualSavingsCOP"`
	PaybackYears        *float64           `json:"paybackYears"`
	IRRPercent          float64            `json:"irrPercent"`
	NPVCOP              float64            `json:"npvCOP"`
	CO2AvoidedTonsYear  float64            `json:"co2AvoidedTonsYear"`
	CumulativeSavings25 []float64          `json:"cumulativeSavings25"`
	LCOE                model.SafeFloat64  `json:"lcoe"`
	Yearly25Savings     []float64          `json:"yearly25Savings"`
}

type FinancialModel struct{}

func NewFinancialModel() *FinancialModel {
	return &FinancialModel{}
}

func (f *FinancialModel) Analyze(input FinancialInput) FinancialResult {
	lifeYears := input.SystemLifeYears
	if lifeYears <= 0 {
		lifeYears = config.SystemDefaults.SystemLifeYears
	}
	discountRate := input.DiscountRate
	if discountRate <= 0 {
		discountRate = config.FinancialDefaults.DiscountRate
	}
	degradationRate := input.DegradationRate
	if degradationRate <= 0 {
		degradationRate = config.SystemDefaults.DegradationRatePerYear
	}
	tariffEscalation := input.TariffEscalation
	if tariffEscalation <= 0 {
		tariffEscalation = config.FinancialDefaults.TariffEscalationRate
	}

	// Cost breakdown
	panelsCost := input.PanelCostCOP * float64(input.NumberOfPanels)
	inverterCost := input.InverterCostCOP
	batteriesCost := input.BatteryCostCOP
	equipmentCost := panelsCost + inverterCost + batteriesCost
	bosCost := equipmentCost * config.FinancialDefaults.BOSCostPercent
	laborCost := equipmentCost * config.FinancialDefaults.LaborCostPercent
	totalInstallationCost := equipmentCost + bosCost + laborCost

	// Monthly savings (first year)
	tariff := input.TariffPerKwh
	if tariff <= 0 {
		tariff = getTariffByEstrato(input.Estrato)
	}
	monthlySavingsCOP := make([]float64, len(input.MonthlyProductionKwh))
	annualSavingsCOP := 0.0
	for i, kwh := range input.MonthlyProductionKwh {
		monthlySavingsCOP[i] = kwh * tariff
		annualSavingsCOP += monthlySavingsCOP[i]
	}

	// 25-year projections
	yearly25Savings := calc.GenerateAnnualSavings(
		input.AnnualProductionKwh, tariff, degradationRate, tariffEscalation, lifeYears,
	)
	cumSavings := calc.CumulativeSavings(yearly25Savings)

	// Cash flows for IRR
	cashFlows := make([]float64, lifeYears+1)
	cashFlows[0] = -totalInstallationCost
	maintenanceCost := totalInstallationCost * config.FinancialDefaults.MaintenanceCostPerYear
	for y := 0; y < lifeYears; y++ {
		cashFlows[y+1] = yearly25Savings[y] - maintenanceCost
	}

	// Financial metrics
	payback := calc.PaybackPeriod(totalInstallationCost, yearly25Savings)
	irrValue := calc.IRR(cashFlows, 100, 1e-7)
	npvValue := calc.NPV(cashFlows, discountRate)

	// CO2 avoided
	co2AvoidedTonsYear := (input.AnnualProductionKwh / 1000) * config.ColombiaCO2Factor

	// LCOE
	yearlyProduction := make([]float64, lifeYears)
	for y := 0; y < lifeYears; y++ {
		yearlyProduction[y] = input.AnnualProductionKwh * math.Pow(1-degradationRate, float64(y))
	}
	lcoeValue := calc.LCOE(totalInstallationCost, maintenanceCost, yearlyProduction, discountRate)

	return FinancialResult{
		InstallationCostCOP: totalInstallationCost,
		CostBreakdown: CostBreakdown{
			Panels:    panelsCost,
			Inverter:  inverterCost,
			Batteries: batteriesCost,
			BOS:       bosCost,
			Labor:     laborCost,
		},
		MonthlySavingsCOP:   monthlySavingsCOP,
		AnnualSavingsCOP:    annualSavingsCOP,
		PaybackYears:        payback,
		IRRPercent:          irrValue * 100,
		NPVCOP:              npvValue,
		CO2AvoidedTonsYear:  co2AvoidedTonsYear,
		CumulativeSavings25: cumSavings,
		LCOE:                model.SafeFloat64(lcoeValue),
		Yearly25Savings:     yearly25Savings,
	}
}

func getTariffByEstrato(estrato int) float64 {
	key := strconv.Itoa(estrato)
	if entry, ok := data.ColombianTariffs[key]; ok {
		return entry.TariffPerKwh
	}
	return 800
}
