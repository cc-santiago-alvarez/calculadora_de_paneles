package service

import (
	"fmt"
	"math"

	"github.com/dev13/calculadora-paneles-backend/internal/calc"
	"github.com/dev13/calculadora-paneles-backend/internal/config"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
)

type SystemDesignInput struct {
	DailyConsumptionKwh   float64
	MonthlyConsumptionKwh [12]float64
	AvgHSP                float64
	MonthlyHSP            []float64
	Panel                 model.PanelCatalog
	Inverter              model.InverterCatalog
	RoofArea              float64
	UsablePercentage      float64
	MonthlyAmbientTemp    []float64 // nil = default 25C
	ShadingLoss           []float64 // nil = no shading
	SystemType            string
	CoveragePercentage    float64 // 0-100, defaults to 100 if <= 0
	PanelConfig           *PanelConfigInput // nil = automatic string configuration
}

// PanelConfigInput carries a user-defined panel wiring that overrides the
// automatically derived string configuration. PanelsPerString (S) are wired in
// series; NumberOfStrings (P) strings are wired in parallel.
type PanelConfigInput struct {
	ConnectionType  string
	PanelsPerString int
	NumberOfStrings int
}

type StringConfig struct {
	PanelsPerString int     `json:"panelsPerString"`
	NumberOfStrings int     `json:"numberOfStrings"`
	StringVoltage   float64 `json:"stringVoltage"`
	StringCurrent   float64 `json:"stringCurrent"`
}

type SystemDesignResult struct {
	RequiredPowerKwp    float64      `json:"requiredPowerKwp"`
	NumberOfPanels      int          `json:"numberOfPanels"`
	ActualPowerKwp      float64      `json:"actualPowerKwp"`
	RoofUtilization     float64      `json:"roofUtilization"`
	InverterCapacityKw  float64      `json:"inverterCapacityKw"`
	StringConfiguration StringConfig `json:"stringConfiguration"`
	MonthlyProductionKwh []float64   `json:"monthlyProductionKwh"`
	AnnualProductionKwh  float64     `json:"annualProductionKwh"`
	Yearly25Production   []float64   `json:"yearly25Production"`
	Losses              model.Losses `json:"losses"`
	Warnings            []string     `json:"warnings"`
}

type PVSystemCalculator struct{}

func NewPVSystemCalculator() *PVSystemCalculator {
	return &PVSystemCalculator{}
}

func (c *PVSystemCalculator) Calculate(input SystemDesignInput) SystemDesignResult {
	var warnings []string

	baseEfficiency := (1 - config.DefaultLosses.Soiling) *
		(1 - config.DefaultLosses.Mismatch) *
		(1 - config.DefaultLosses.Wiring) *
		config.DefaultLosses.InverterEfficiency

	coverageFactor := input.CoveragePercentage / 100
	if coverageFactor <= 0 || coverageFactor > 1 {
		coverageFactor = 1
	}

	requiredPowerKwp := (input.DailyConsumptionKwh * coverageFactor) / input.AvgHSP / baseEfficiency
	panelWp := input.Panel.PowerWp
	panelAreaM2 := input.Panel.Area

	usableArea := input.RoofArea * (input.UsablePercentage / 100)
	maxPanelsByRoof := int(math.Floor(usableArea / (panelAreaM2 * config.SystemDefaults.PanelSpacingFactor)))

	manualConfig := input.PanelConfig != nil
	var numberOfPanels int

	if manualConfig {
		// The user manually defined the panel wiring; it drives the calculation.
		numberOfPanels = input.PanelConfig.PanelsPerString * input.PanelConfig.NumberOfStrings
		if maxPanelsByRoof > 0 && numberOfPanels > maxPanelsByRoof {
			warnings = append(warnings, fmt.Sprintf(
				"La configuración manual usa %d paneles pero solo caben %d en el techo disponible.",
				numberOfPanels, maxPanelsByRoof,
			))
		}
	} else {
		numberOfPanels = int(math.Ceil(requiredPowerKwp * 1000 / panelWp))
		if numberOfPanels > maxPanelsByRoof {
			warnings = append(warnings, fmt.Sprintf(
				"Se requieren %d paneles pero solo caben %d en el techo disponible. Se ajusta al máximo.",
				numberOfPanels, maxPanelsByRoof,
			))
			numberOfPanels = maxPanelsByRoof
		}
	}

	if numberOfPanels < 1 {
		numberOfPanels = 1
		warnings = append(warnings, "El área del techo es muy pequeña. Se configuró el mínimo de 1 panel.")
	}

	actualPowerKwp := float64(numberOfPanels) * panelWp / 1000
	roofUtilization := float64(numberOfPanels) * panelAreaM2 * config.SystemDefaults.PanelSpacingFactor / input.RoofArea * 100

	var stringConfig StringConfig
	if manualConfig {
		stringConfig = manualStringConfiguration(*input.PanelConfig, input.Panel)
		validateStringConfigAgainstInverter(*input.PanelConfig, input.Panel, input.Inverter, &warnings)
	} else {
		stringConfig = calculateStringConfiguration(numberOfPanels, input.Panel, input.Inverter, &warnings)
	}

	inverterRatio := input.Inverter.RatedPowerKw / actualPowerKwp
	if inverterRatio < 0.8 {
		warnings = append(warnings, fmt.Sprintf(
			"El inversor (%.1fkW) está subdimensionado para el array (%.1fkWp). Ratio: %.0f%%",
			input.Inverter.RatedPowerKw, actualPowerKwp, inverterRatio*100,
		))
	} else if inverterRatio > 1.2 {
		warnings = append(warnings, fmt.Sprintf(
			"El inversor (%.1fkW) está sobredimensionado para el array (%.1fkWp). Ratio: %.0f%%",
			input.Inverter.RatedPowerKw, actualPowerKwp, inverterRatio*100,
		))
	}

	ambientTemps := input.MonthlyAmbientTemp
	if ambientTemps == nil {
		ambientTemps = make([]float64, 12)
		for i := range ambientTemps {
			ambientTemps[i] = 25
		}
	}
	shadingLoss := input.ShadingLoss
	if shadingLoss == nil {
		shadingLoss = make([]float64, 12)
	}

	var totalTempLoss float64
	monthlyProductionKwh := make([]float64, 12)
	for m := 0; m < 12 && m < len(input.MonthlyHSP); m++ {
		tCell := calc.CellTemperature(ambientTemps[m], input.Panel.NOCT, 800)
		tLoss := calc.TemperatureLoss(tCell, input.Panel.TempCoeffPmax)
		totalTempLoss += tLoss

		monthLossFactor := (1 - tLoss) *
			(1 - config.DefaultLosses.Soiling) *
			(1 - config.DefaultLosses.Mismatch) *
			(1 - config.DefaultLosses.Wiring) *
			config.DefaultLosses.InverterEfficiency *
			(1 - shadingLoss[m])

		monthlyProductionKwh[m] = actualPowerKwp * input.MonthlyHSP[m] * float64(config.DaysInMonth[m]) * monthLossFactor
	}

	annualProductionKwh := 0.0
	for _, v := range monthlyProductionKwh {
		annualProductionKwh += v
	}
	avgTempLoss := totalTempLoss / 12
	avgShadingLoss := 0.0
	for _, v := range shadingLoss {
		avgShadingLoss += v
	}
	avgShadingLoss /= 12

	yearly25Production := make([]float64, config.SystemDefaults.SystemLifeYears)
	for y := 0; y < config.SystemDefaults.SystemLifeYears; y++ {
		degradationFactor := math.Pow(1-config.SystemDefaults.DegradationRatePerYear, float64(y))
		yearly25Production[y] = annualProductionKwh * degradationFactor
	}

	totalSystemLoss := 1 - (1-avgTempLoss)*
		(1-config.DefaultLosses.Soiling)*
		(1-config.DefaultLosses.Mismatch)*
		(1-config.DefaultLosses.Wiring)*
		config.DefaultLosses.InverterEfficiency*
		(1-avgShadingLoss)

	if warnings == nil {
		warnings = []string{}
	}

	return SystemDesignResult{
		RequiredPowerKwp:     requiredPowerKwp,
		NumberOfPanels:       numberOfPanels,
		ActualPowerKwp:       actualPowerKwp,
		RoofUtilization:      roofUtilization,
		InverterCapacityKw:   input.Inverter.RatedPowerKw,
		StringConfiguration:  stringConfig,
		MonthlyProductionKwh: monthlyProductionKwh,
		AnnualProductionKwh:  annualProductionKwh,
		Yearly25Production:   yearly25Production,
		Losses: model.Losses{
			ShadingPercent:     avgShadingLoss * 100,
			TemperaturePercent: avgTempLoss * 100,
			WiringPercent:      config.DefaultLosses.Wiring * 100,
			InverterPercent:    (1 - config.DefaultLosses.InverterEfficiency) * 100,
			SoilingPercent:     config.DefaultLosses.Soiling * 100,
			TotalSystemLoss:    totalSystemLoss * 100,
		},
		Warnings: warnings,
	}
}

// SlopeAllocation describes how panels and irradiation are distributed per slope.
type SlopeAllocation struct {
	PanelCount int
	MonthlyHSP []float64
	AvgHSP     float64
}

// CalculateMultiSlope performs system design with per-slope production calculation.
// It uses Calculate() for sizing (panels, strings, warnings) with weighted-average HSP,
// then computes production per slope and aggregates.
func (c *PVSystemCalculator) CalculateMultiSlope(input SystemDesignInput, slopes []SlopeAllocation) SystemDesignResult {
	// Use standard Calculate for system sizing (panel count, string config, warnings)
	result := c.Calculate(input)

	// If only one slope, the standard calculation is already correct
	if len(slopes) <= 1 {
		return result
	}

	// Distribute panels proportionally to each slope's panel count
	// (caller is responsible for setting PanelCount on each slope)
	totalPanels := 0
	for _, s := range slopes {
		totalPanels += s.PanelCount
	}
	if totalPanels == 0 {
		return result
	}

	// Recalculate production per slope
	panelWp := input.Panel.PowerWp
	actualPowerPerPanel := panelWp / 1000

	ambientTemps := input.MonthlyAmbientTemp
	if ambientTemps == nil {
		ambientTemps = make([]float64, 12)
		for i := range ambientTemps {
			ambientTemps[i] = 25
		}
	}
	shadingLoss := input.ShadingLoss
	if shadingLoss == nil {
		shadingLoss = make([]float64, 12)
	}

	monthlyProductionKwh := make([]float64, 12)
	var totalTempLoss float64

	for _, slope := range slopes {
		if slope.PanelCount == 0 {
			continue
		}
		slopePowerKwp := float64(slope.PanelCount) * actualPowerPerPanel

		for m := 0; m < 12 && m < len(slope.MonthlyHSP); m++ {
			tCell := calc.CellTemperature(ambientTemps[m], input.Panel.NOCT, 800)
			tLoss := calc.TemperatureLoss(tCell, input.Panel.TempCoeffPmax)
			totalTempLoss += tLoss

			monthLossFactor := (1 - tLoss) *
				(1 - config.DefaultLosses.Soiling) *
				(1 - config.DefaultLosses.Mismatch) *
				(1 - config.DefaultLosses.Wiring) *
				config.DefaultLosses.InverterEfficiency *
				(1 - shadingLoss[m])

			monthlyProductionKwh[m] += slopePowerKwp * slope.MonthlyHSP[m] * float64(config.DaysInMonth[m]) * monthLossFactor
		}
	}

	annualProductionKwh := 0.0
	for _, v := range monthlyProductionKwh {
		annualProductionKwh += v
	}

	avgTempLoss := totalTempLoss / float64(12*len(slopes))
	avgShadingLoss := 0.0
	for _, v := range shadingLoss {
		avgShadingLoss += v
	}
	avgShadingLoss /= 12

	yearly25Production := make([]float64, config.SystemDefaults.SystemLifeYears)
	for y := 0; y < config.SystemDefaults.SystemLifeYears; y++ {
		degradationFactor := math.Pow(1-config.SystemDefaults.DegradationRatePerYear, float64(y))
		yearly25Production[y] = annualProductionKwh * degradationFactor
	}

	totalSystemLoss := 1 - (1-avgTempLoss)*
		(1-config.DefaultLosses.Soiling)*
		(1-config.DefaultLosses.Mismatch)*
		(1-config.DefaultLosses.Wiring)*
		config.DefaultLosses.InverterEfficiency*
		(1-avgShadingLoss)

	// Update result with multi-slope production
	result.MonthlyProductionKwh = monthlyProductionKwh
	result.AnnualProductionKwh = annualProductionKwh
	result.Yearly25Production = yearly25Production
	result.Losses = model.Losses{
		ShadingPercent:     avgShadingLoss * 100,
		TemperaturePercent: avgTempLoss * 100,
		WiringPercent:      config.DefaultLosses.Wiring * 100,
		InverterPercent:    (1 - config.DefaultLosses.InverterEfficiency) * 100,
		SoilingPercent:     config.DefaultLosses.Soiling * 100,
		TotalSystemLoss:    totalSystemLoss * 100,
	}

	return result
}

func calculateStringConfiguration(totalPanels int, panel model.PanelCatalog, inverter model.InverterCatalog, warnings *[]string) StringConfig {
	vocAdjusted := panel.Voc * (1 + (panel.TempCoeffVoc/100)*(-10-25))
	vmpAdjusted := panel.Vmp * (1 + (panel.TempCoeffVoc/100)*(50-25))

	maxPanelsPerString := int(math.Floor(inverter.MaxInputVoltage / vocAdjusted))
	minPanelsPerString := int(math.Ceil(inverter.MPPTVoltageMin / vmpAdjusted))

	panelsPerString := maxPanelsPerString
	if panelsPerString > totalPanels {
		panelsPerString = totalPanels
	}
	if panelsPerString < minPanelsPerString {
		panelsPerString = minPanelsPerString
	}

	if panelsPerString > totalPanels {
		panelsPerString = totalPanels
		*warnings = append(*warnings, "Número de paneles insuficiente para el rango MPPT del inversor.")
	}

	numberOfStrings := int(math.Ceil(float64(totalPanels) / float64(panelsPerString)))

	maxStrings := int(math.Floor((inverter.MaxInputCurrent * float64(inverter.MPPTCount)) / panel.Imp))
	if numberOfStrings > maxStrings {
		*warnings = append(*warnings, fmt.Sprintf(
			"Se necesitan %d strings pero el inversor soporta máximo %d. Se ajusta.",
			numberOfStrings, maxStrings,
		))
		numberOfStrings = maxStrings
	}

	stringVoltage := float64(panelsPerString) * panel.Vmp
	stringCurrent := panel.Imp * float64(numberOfStrings)

	return StringConfig{
		PanelsPerString: panelsPerString,
		NumberOfStrings: numberOfStrings,
		StringVoltage:   math.Round(stringVoltage*100) / 100,
		StringCurrent:   math.Round(stringCurrent*100) / 100,
	}
}

// manualStringConfiguration builds the array electrical values from a user-defined
// panel wiring: panels in series add voltage, strings in parallel add current.
func manualStringConfiguration(cfg PanelConfigInput, panel model.PanelCatalog) StringConfig {
	s := cfg.PanelsPerString
	if s < 1 {
		s = 1
	}
	p := cfg.NumberOfStrings
	if p < 1 {
		p = 1
	}
	stringVoltage := float64(s) * panel.Vmp
	stringCurrent := float64(p) * panel.Imp

	return StringConfig{
		PanelsPerString: s,
		NumberOfStrings: p,
		StringVoltage:   math.Round(stringVoltage*100) / 100,
		StringCurrent:   math.Round(stringCurrent*100) / 100,
	}
}

// validateStringConfigAgainstInverter emits warnings when a manual panel
// configuration exceeds the inverter's voltage or current limits.
func validateStringConfigAgainstInverter(cfg PanelConfigInput, panel model.PanelCatalog, inverter model.InverterCatalog, warnings *[]string) {
	s := float64(cfg.PanelsPerString)
	p := float64(cfg.NumberOfStrings)

	// Voc at -10°C (worst cold case) must stay below the inverter max input voltage.
	vocCold := panel.Voc * (1 + (panel.TempCoeffVoc/100)*(-10-25))
	arrayVocCold := s * vocCold
	if inverter.MaxInputVoltage > 0 && arrayVocCold > inverter.MaxInputVoltage {
		*warnings = append(*warnings, fmt.Sprintf(
			"El voltaje del string (%.0fV a -10°C) supera el máximo del inversor (%.0fV). Reduce los paneles en serie.",
			arrayVocCold, inverter.MaxInputVoltage,
		))
	}

	// Vmp at 50°C (worst hot case) should stay within the MPPT range.
	vmpHot := panel.Vmp * (1 + (panel.TempCoeffVoc/100)*(50-25))
	arrayVmpHot := s * vmpHot
	if inverter.MPPTVoltageMin > 0 && arrayVmpHot < inverter.MPPTVoltageMin {
		*warnings = append(*warnings, fmt.Sprintf(
			"El voltaje del string (%.0fV a 50°C) está por debajo del rango MPPT del inversor (%.0fV). Aumenta los paneles en serie.",
			arrayVmpHot, inverter.MPPTVoltageMin,
		))
	}

	// Total array current = parallel strings × Isc must stay below the inverter input current.
	maxCurrent := inverter.MaxInputCurrent * float64(inverter.MPPTCount)
	arrayCurrent := p * panel.Isc
	if maxCurrent > 0 && arrayCurrent > maxCurrent {
		*warnings = append(*warnings, fmt.Sprintf(
			"La corriente del arreglo (%.1fA) supera la entrada máxima del inversor (%.1fA). Reduce las cadenas en paralelo.",
			arrayCurrent, maxCurrent,
		))
	}
}
