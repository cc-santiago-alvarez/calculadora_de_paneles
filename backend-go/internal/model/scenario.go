package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type IrradiationData struct {
	Source       string    `json:"source" bson:"source"`
	MonthlyGHI   []float64 `json:"monthlyGHI" bson:"monthlyGHI"`
	MonthlyPOA   []float64 `json:"monthlyPOA" bson:"monthlyPOA"`
	AnnualAvgHSP float64   `json:"annualAvgHSP" bson:"annualAvgHSP"`
}

type StringConfiguration struct {
	PanelsPerString  int     `json:"panelsPerString" bson:"panelsPerString"`
	NumberOfStrings  int     `json:"numberOfStrings" bson:"numberOfStrings"`
	StringVoltage    float64 `json:"stringVoltage" bson:"stringVoltage"`
	StringCurrent    float64 `json:"stringCurrent" bson:"stringCurrent"`
}

type BatteryBank struct {
	CapacityKwh       float64 `json:"capacityKwh" bson:"capacityKwh"`
	AutonomyDays      float64 `json:"autonomyDays" bson:"autonomyDays"`
	NumberOfBatteries int     `json:"numberOfBatteries" bson:"numberOfBatteries"`
	BankVoltage       float64 `json:"bankVoltage" bson:"bankVoltage"`
}

type SystemDesign struct {
	RequiredPowerKwp    float64              `json:"requiredPowerKwp" bson:"requiredPowerKwp"`
	NumberOfPanels      int                  `json:"numberOfPanels" bson:"numberOfPanels"`
	ActualPowerKwp      float64              `json:"actualPowerKwp" bson:"actualPowerKwp"`
	RoofUtilization     float64              `json:"roofUtilization" bson:"roofUtilization"`
	InverterCapacityKw  float64              `json:"inverterCapacityKw" bson:"inverterCapacityKw"`
	StringConfiguration StringConfiguration  `json:"stringConfiguration" bson:"stringConfiguration"`
	BatteryBank         *BatteryBank         `json:"batteryBank,omitempty" bson:"batteryBank,omitempty"`
}

type Production struct {
	MonthlyKwh      []float64 `json:"monthlyKwh" bson:"monthlyKwh"`
	AnnualKwh       float64   `json:"annualKwh" bson:"annualKwh"`
	DegradationRate float64   `json:"degradationRate" bson:"degradationRate"`
	Yearly25        []float64 `json:"yearly25" bson:"yearly25"`
}

type Financial struct {
	InstallationCostCOP float64      `json:"installationCostCOP" bson:"installationCostCOP"`
	MonthlySavingsCOP   []float64    `json:"monthlySavingsCOP" bson:"monthlySavingsCOP"`
	AnnualSavingsCOP    float64      `json:"annualSavingsCOP" bson:"annualSavingsCOP"`
	PaybackYears        *float64     `json:"paybackYears" bson:"paybackYears"`
	IRRPercent          float64      `json:"irrPercent" bson:"irrPercent"`
	NPVCOP              float64      `json:"npvCOP" bson:"npvCOP"`
	CO2AvoidedTonsYear  float64      `json:"co2AvoidedTonsYear" bson:"co2AvoidedTonsYear"`
	CumulativeSavings25 []float64    `json:"cumulativeSavings25" bson:"cumulativeSavings25"`
	LCOE                SafeFloat64  `json:"lcoe" bson:"lcoe"`
}

type Losses struct {
	ShadingPercent     float64 `json:"shadingPercent" bson:"shadingPercent"`
	TemperaturePercent float64 `json:"temperaturePercent" bson:"temperaturePercent"`
	WiringPercent      float64 `json:"wiringPercent" bson:"wiringPercent"`
	InverterPercent    float64 `json:"inverterPercent" bson:"inverterPercent"`
	SoilingPercent     float64 `json:"soilingPercent" bson:"soilingPercent"`
	TotalSystemLoss    float64 `json:"totalSystemLoss" bson:"totalSystemLoss"`
}

type Scenario struct {
	ID            bson.ObjectID          `json:"_id,omitempty" bson:"_id,omitempty"`
	ProjectID     bson.ObjectID          `json:"projectId" bson:"projectId"`
	Name          string                 `json:"name" bson:"name"`
	InputSnapshot map[string]interface{} `json:"inputSnapshot" bson:"inputSnapshot"`
	Irradiation   IrradiationData        `json:"irradiation" bson:"irradiation"`
	SystemDesign  SystemDesign           `json:"systemDesign" bson:"systemDesign"`
	Production    Production             `json:"production" bson:"production"`
	Financial     Financial              `json:"financial" bson:"financial"`
	Losses        Losses                 `json:"losses" bson:"losses"`
	CreatedAt     time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt" bson:"updatedAt"`
}
