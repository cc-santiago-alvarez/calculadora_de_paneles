package service

import (
	"math"

	"github.com/dev13/calculadora-paneles-backend/internal/config"
)

type BatteryInput struct {
	DailyConsumptionKwh float64 `json:"dailyConsumptionKwh"`
	AutonomyDays        float64 `json:"autonomyDays"`
	SystemVoltage       float64 `json:"systemVoltage"`
	BatteryType         string  `json:"batteryType"` // "leadAcid" or "lithium"
	BatteryCapacityAh   float64 `json:"batteryCapacityAh"`
	BatteryVoltage      float64 `json:"batteryVoltage"`
}

type BatteryResult struct {
	CapacityKwh       float64 `json:"capacityKwh"`
	RequiredAh        float64 `json:"requiredAh"`
	NumberOfBatteries int     `json:"numberOfBatteries"`
	SeriesBatteries   int     `json:"seriesBatteries"`
	ParallelBatteries int     `json:"parallelBatteries"`
	BankVoltage       float64 `json:"bankVoltage"`
	TotalAh           float64 `json:"totalAh"`
	AutonomyDays      float64 `json:"autonomyDays"`
	EstimatedCostCOP  float64 `json:"estimatedCostCOP"`
}

type BatteryCalculator struct{}

func NewBatteryCalculator() *BatteryCalculator {
	return &BatteryCalculator{}
}

func (c *BatteryCalculator) Calculate(input BatteryInput) BatteryResult {
	var batteryConfig config.BatteryConfig
	if input.BatteryType == "lithium" {
		batteryConfig = config.BatteryDefaults.Lithium
	} else {
		batteryConfig = config.BatteryDefaults.LeadAcid
	}

	DOD := batteryConfig.DOD
	inverterEfficiency := 0.95
	controllerEfficiency := config.BatteryDefaults.ControllerEfficiency

	requiredAh := (input.DailyConsumptionKwh * input.AutonomyDays * 1000) /
		(input.SystemVoltage * DOD * inverterEfficiency * controllerEfficiency)

	seriesBatteries := int(math.Ceil(input.SystemVoltage / input.BatteryVoltage))
	parallelBatteries := int(math.Ceil(requiredAh / input.BatteryCapacityAh))

	numberOfBatteries := seriesBatteries * parallelBatteries
	totalAh := float64(parallelBatteries) * input.BatteryCapacityAh
	bankVoltage := float64(seriesBatteries) * input.BatteryVoltage
	capacityKwh := (totalAh * bankVoltage) / 1000

	estimatedCostCOP := capacityKwh * batteryConfig.CostPerKwh

	return BatteryResult{
		CapacityKwh:       capacityKwh,
		RequiredAh:        requiredAh,
		NumberOfBatteries: numberOfBatteries,
		SeriesBatteries:   seriesBatteries,
		ParallelBatteries: parallelBatteries,
		BankVoltage:       bankVoltage,
		TotalAh:           totalAh,
		AutonomyDays:      input.AutonomyDays,
		EstimatedCostCOP:  estimatedCostCOP,
	}
}
