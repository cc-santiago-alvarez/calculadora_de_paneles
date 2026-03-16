package data

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed ideam-zones.json
var ideamZonesJSON []byte

//go:embed colombian-tariffs.json
var colombianTariffsJSON []byte

//go:embed default-panels.json
var defaultPanelsJSON []byte

//go:embed default-inverters.json
var defaultInvertersJSON []byte

// IDEAMZone represents a zone entry from ideam-zones.json.
type IDEAMZone struct {
	AnnualAvgGHI float64    `json:"annualAvgGHI"`
	MonthlyGHI   [12]float64 `json:"monthlyGHI"`
	Capital      string     `json:"capital"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
}

type ideamData struct {
	Zones map[string]IDEAMZone `json:"zones"`
}

// TariffEntry represents a tariff entry from colombian-tariffs.json.
type TariffEntry struct {
	Name           string  `json:"name"`
	TariffPerKwh   float64 `json:"tariffPerKwh"`
	SubsidyPercent float64 `json:"subsidyPercent"`
	BaseRate       float64 `json:"baseRate"`
	Description    string  `json:"description"`
}

type tariffData struct {
	Tariffs map[string]TariffEntry `json:"tariffs"`
}

// PanelEntry for default-panels.json.
type PanelEntry struct {
	Manufacturer  string  `json:"manufacturer"`
	Model         string  `json:"model"`
	Type          string  `json:"type"`
	PowerWp       float64 `json:"powerWp"`
	Efficiency    float64 `json:"efficiency"`
	Area          float64 `json:"area"`
	Voc           float64 `json:"voc"`
	Isc           float64 `json:"isc"`
	Vmp           float64 `json:"vmp"`
	Imp           float64 `json:"imp"`
	TempCoeffPmax float64 `json:"tempCoeffPmax"`
	TempCoeffVoc  float64 `json:"tempCoeffVoc"`
	NOCT          float64 `json:"NOCT"`
	Weight        float64 `json:"weight"`
	Dimensions    struct {
		Length float64 `json:"length"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Warranty int     `json:"warranty"`
	CostCOP  float64 `json:"costCOP"`
}

// InverterEntry for default-inverters.json.
type InverterEntry struct {
	Manufacturer    string  `json:"manufacturer"`
	Model           string  `json:"model"`
	Type            string  `json:"type"`
	RatedPowerKw    float64 `json:"ratedPowerKw"`
	MaxDCPowerKw    float64 `json:"maxDCPowerKw"`
	Efficiency      float64 `json:"efficiency"`
	MPPTCount       int     `json:"mpptCount"`
	MPPTVoltageMin  float64 `json:"mpptVoltageMin"`
	MPPTVoltageMax  float64 `json:"mpptVoltageMax"`
	MaxInputVoltage float64 `json:"maxInputVoltage"`
	MaxInputCurrent float64 `json:"maxInputCurrent"`
	OutputVoltage   float64 `json:"outputVoltage"`
	OutputPhases    int     `json:"outputPhases"`
	HasBatteryPort  bool    `json:"hasBatteryPort"`
	Weight          float64 `json:"weight"`
	Warranty        int     `json:"warranty"`
	CostCOP         float64 `json:"costCOP"`
}

var (
	IDEAMZones       map[string]IDEAMZone
	ColombianTariffs map[string]TariffEntry
	DefaultPanels    []PanelEntry
	DefaultInverters []InverterEntry
)

func init() {
	var iz ideamData
	if err := json.Unmarshal(ideamZonesJSON, &iz); err != nil {
		log.Fatalf("failed to parse ideam-zones.json: %v", err)
	}
	IDEAMZones = iz.Zones

	var td tariffData
	if err := json.Unmarshal(colombianTariffsJSON, &td); err != nil {
		log.Fatalf("failed to parse colombian-tariffs.json: %v", err)
	}
	ColombianTariffs = td.Tariffs

	if err := json.Unmarshal(defaultPanelsJSON, &DefaultPanels); err != nil {
		log.Fatalf("failed to parse default-panels.json: %v", err)
	}

	if err := json.Unmarshal(defaultInvertersJSON, &DefaultInverters); err != nil {
		log.Fatalf("failed to parse default-inverters.json: %v", err)
	}
}
