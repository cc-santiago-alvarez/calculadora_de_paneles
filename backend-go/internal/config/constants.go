package config

// SolarConstant is the solar irradiance at top of atmosphere in W/m².
const SolarConstant = 1367

// StandardTestConditions for PV panels.
var StandardTestConditions = struct {
	Irradiance      float64
	CellTemperature float64
	AirMass         float64
}{
	Irradiance:      1000,
	CellTemperature: 25,
	AirMass:         1.5,
}

// DefaultLosses for PV system calculations.
var DefaultLosses = struct {
	Soiling                  float64
	Mismatch                 float64
	Wiring                   float64
	InverterEfficiency       float64
	LightInducedDegradation  float64
}{
	Soiling:                  0.03,
	Mismatch:                 0.02,
	Wiring:                   0.02,
	InverterEfficiency:       0.96,
	LightInducedDegradation:  0.015,
}

// SystemDefaults for PV system design.
var SystemDefaults = struct {
	DegradationRatePerYear float64
	SystemLifeYears        int
	Albedo                 float64
	NOCT                   float64
	PanelSpacingFactor     float64
}{
	DegradationRatePerYear: 0.005,
	SystemLifeYears:        25,
	Albedo:                 0.2,
	NOCT:                   45,
	PanelSpacingFactor:     1.15,
}

// FinancialDefaults for financial analysis.
var FinancialDefaults = struct {
	DiscountRate            float64
	TariffEscalationRate    float64
	InstallationCostPerWp   float64
	BOSCostPercent          float64
	LaborCostPercent        float64
	MaintenanceCostPerYear  float64
}{
	DiscountRate:            0.08,
	TariffEscalationRate:    0.04,
	InstallationCostPerWp:   3500,
	BOSCostPercent:          0.15,
	LaborCostPercent:        0.10,
	MaintenanceCostPerYear:  0.01,
}

// ColombiaCO2Factor is the grid emission factor in tCO2/MWh.
const ColombiaCO2Factor = 0.126

// BatteryConfig holds battery type specifications.
type BatteryConfig struct {
	DOD        float64
	Efficiency float64
	CycleLife  int
	CostPerKwh float64
}

// BatteryDefaults for battery calculations.
var BatteryDefaults = struct {
	LeadAcid              BatteryConfig
	Lithium               BatteryConfig
	DefaultAutonomyDays   int
	ControllerEfficiency  float64
}{
	LeadAcid: BatteryConfig{
		DOD:        0.50,
		Efficiency: 0.85,
		CycleLife:  1500,
		CostPerKwh: 450000,
	},
	Lithium: BatteryConfig{
		DOD:        0.80,
		Efficiency: 0.95,
		CycleLife:  6000,
		CostPerKwh: 1200000,
	},
	DefaultAutonomyDays:  2,
	ControllerEfficiency: 0.95,
}

// Months in Spanish.
var Months = [12]string{
	"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
	"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
}

// DaysInMonth for a non-leap year.
var DaysInMonth = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
