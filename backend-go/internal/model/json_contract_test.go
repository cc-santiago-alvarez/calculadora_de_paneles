package model

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TestProjectJSONContract verifies JSON output matches Mongoose's format.
func TestProjectJSONContract(t *testing.T) {
	id, _ := bson.ObjectIDFromHex("507f1f77bcf86cd799439011")
	panelID, _ := bson.ObjectIDFromHex("507f1f77bcf86cd799439012")
	inverterID, _ := bson.ObjectIDFromHex("507f1f77bcf86cd799439013")

	project := Project{
		ID:   id,
		Name: "Test Project",
		Location: Location{
			Latitude:  4.61,
			Longitude: -74.08,
			Altitude:  0,
		},
		Consumption: Consumption{
			Monthly:        [12]float64{200, 210, 205, 195, 190, 200, 210, 205, 195, 190, 200, 210},
			TariffPerKwh:   800,
			Estrato:        4,
			ConnectionType: "monofasica",
		},
		Roof: Roof{
			Area:             50,
			Azimuth:          0,
			Tilt:             10,
			UsablePercentage: 80,
			ShadingProfile: ShadingProfile{
				HasShading:  false,
				MonthlyLoss: [12]float64{},
			},
		},
		SystemType: "on-grid",
		Equipment: Equipment{
			PanelID:    panelID,
			InverterID: inverterID,
		},
		Scenarios: []bson.ObjectID{},
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
	}

	data, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("failed to marshal project: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify _id is a hex string (not {"$oid": "..."})
	idVal, ok := result["_id"]
	if !ok {
		t.Error("missing _id field")
	}
	idStr, ok := idVal.(string)
	if !ok {
		t.Errorf("_id should be string, got %T: %v", idVal, idVal)
	}
	if idStr != "507f1f77bcf86cd799439011" {
		t.Errorf("_id = %q, want %q", idStr, "507f1f77bcf86cd799439011")
	}

	// Verify createdAt is ISO format
	if _, ok := result["createdAt"]; !ok {
		t.Error("missing createdAt field")
	}

	// Verify nested equipment.panelId is a hex string
	equip := result["equipment"].(map[string]interface{})
	if pid, ok := equip["panelId"].(string); !ok || pid != "507f1f77bcf86cd799439012" {
		t.Errorf("equipment.panelId = %v, want string hex", equip["panelId"])
	}

	// Verify scenarios is empty array, not null
	scenarios := result["scenarios"]
	if scenarios == nil {
		t.Error("scenarios should be [] not null")
	}

	// Verify consumption.monthly is array of 12
	consumption := result["consumption"].(map[string]interface{})
	monthly := consumption["monthly"].([]interface{})
	if len(monthly) != 12 {
		t.Errorf("monthly should have 12 values, got %d", len(monthly))
	}

	// Verify altitude is included even when 0 (no omitempty)
	location := result["location"].(map[string]interface{})
	if _, ok := location["altitude"]; !ok {
		t.Error("altitude should be included even when 0")
	}
}

// TestFinancialPaybackNull verifies payback serializes as null when nil.
func TestFinancialPaybackNull(t *testing.T) {
	financial := Financial{
		InstallationCostCOP: 15000000,
		MonthlySavingsCOP:   []float64{100000},
		AnnualSavingsCOP:    1200000,
		PaybackYears:        nil, // Never pays back
		IRRPercent:          5.0,
		NPVCOP:              -5000000,
		CO2AvoidedTonsYear:  0.5,
		CumulativeSavings25: []float64{1200000},
		LCOE:                SafeFloat64(350),
	}

	data, err := json.Marshal(financial)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if result["paybackYears"] != nil {
		t.Errorf("paybackYears should be null, got %v", result["paybackYears"])
	}
}

// TestFinancialPaybackValue verifies payback serializes as number when set.
func TestFinancialPaybackValue(t *testing.T) {
	payback := 7.5
	financial := Financial{
		PaybackYears:        &payback,
		MonthlySavingsCOP:   []float64{},
		CumulativeSavings25: []float64{},
		LCOE:                SafeFloat64(350),
	}

	data, err := json.Marshal(financial)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	if val, ok := result["paybackYears"].(float64); !ok || val != 7.5 {
		t.Errorf("paybackYears should be 7.5, got %v", result["paybackYears"])
	}
}

// TestSafeFloat64InfSerializesAsNull verifies Infinity -> null.
func TestSafeFloat64InfSerializesAsNull(t *testing.T) {
	f := SafeFloat64(math.Inf(1))
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Infinity should serialize as null, got %s", string(data))
	}
}

// TestSafeFloat64NormalValue verifies normal floats serialize correctly.
func TestSafeFloat64NormalValue(t *testing.T) {
	f := SafeFloat64(350.25)
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(data) != "350.25" {
		t.Errorf("expected 350.25, got %s", string(data))
	}
}

// TestSafeFloat64NaN verifies NaN -> null.
func TestSafeFloat64NaN(t *testing.T) {
	f := SafeFloat64(math.NaN())
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(data) != "null" {
		t.Errorf("NaN should serialize as null, got %s", string(data))
	}
}

// TestScenarioJSONFields verifies all scenario fields are present.
func TestScenarioJSONFields(t *testing.T) {
	payback := 8.5
	scenario := Scenario{
		Irradiation: IrradiationData{
			Source:     "ideam",
			MonthlyGHI: make([]float64, 12),
			MonthlyPOA: make([]float64, 12),
		},
		SystemDesign: SystemDesign{
			StringConfiguration: StringConfiguration{},
		},
		Production: Production{
			MonthlyKwh: make([]float64, 12),
			Yearly25:   make([]float64, 25),
		},
		Financial: Financial{
			MonthlySavingsCOP:   make([]float64, 12),
			PaybackYears:        &payback,
			CumulativeSavings25: make([]float64, 25),
			LCOE:                SafeFloat64(350),
		},
		InputSnapshot: map[string]interface{}{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	data, err := json.Marshal(scenario)
	if err != nil {
		t.Fatalf("failed to marshal scenario: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	// Check all top-level fields exist
	requiredFields := []string{
		"projectId", "name", "inputSnapshot", "irradiation",
		"systemDesign", "production", "financial", "losses",
		"createdAt", "updatedAt",
	}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("missing field: %s", field)
		}
	}

	// Check financial sub-fields
	fin := result["financial"].(map[string]interface{})
	finFields := []string{
		"installationCostCOP", "monthlySavingsCOP", "annualSavingsCOP",
		"paybackYears", "irrPercent", "npvCOP", "co2AvoidedTonsYear",
		"cumulativeSavings25", "lcoe",
	}
	for _, field := range finFields {
		if _, ok := fin[field]; !ok {
			t.Errorf("missing financial field: %s", field)
		}
	}
}
