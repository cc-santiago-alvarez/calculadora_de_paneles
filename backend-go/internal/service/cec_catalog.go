package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
)

const (
	cecModulesURL   = "https://raw.githubusercontent.com/NREL/SAM/develop/deploy/libraries/CEC%20Modules.csv"
	cecInvertersURL = "https://raw.githubusercontent.com/NREL/SAM/develop/deploy/libraries/CEC%20Inverters.csv"

	// Cost estimation: COP per watt for panels
	copPerWattPanel = 1500.0
	// Cost estimation: COP per kW for inverters
	copPerKwInverter = 800000.0

	// Only import panels in this power range
	minPanelPower = 300.0
	maxPanelPower = 700.0
)

type CECCatalogService struct {
	catalogRepo *repository.CatalogRepo
	httpClient  *http.Client
}

func NewCECCatalogService(catalogRepo *repository.CatalogRepo) *CECCatalogService {
	return &CECCatalogService{
		catalogRepo: catalogRepo,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

type SyncResult struct {
	PanelsUpserted   int `json:"panelsUpserted"`
	InvertersUpserted int `json:"invertersUpserted"`
	PanelErrors      int `json:"panelErrors"`
	InverterErrors   int `json:"inverterErrors"`
}

// SyncAll downloads and imports both panels and inverters from CEC/SAM.
func (s *CECCatalogService) SyncAll(ctx context.Context) (*SyncResult, error) {
	result := &SyncResult{}

	panels, panelErrs, err := s.fetchPanels(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching CEC panels: %w", err)
	}
	result.PanelErrors = panelErrs

	if len(panels) > 0 {
		upserted, err := s.catalogRepo.UpsertPanels(ctx, panels)
		if err != nil {
			return nil, fmt.Errorf("upserting panels: %w", err)
		}
		result.PanelsUpserted = upserted
	}

	inverters, invErrs, err := s.fetchInverters(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching CEC inverters: %w", err)
	}
	result.InverterErrors = invErrs

	if len(inverters) > 0 {
		upserted, err := s.catalogRepo.UpsertInverters(ctx, inverters)
		if err != nil {
			return nil, fmt.Errorf("upserting inverters: %w", err)
		}
		result.InvertersUpserted = upserted
	}

	return result, nil
}

func (s *CECCatalogService) fetchPanels(ctx context.Context) ([]model.PanelCatalog, int, error) {
	resp, err := s.httpClient.Get(cecModulesURL)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP GET modules: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("CEC modules HTTP %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // variable fields

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, 0, fmt.Errorf("parsing CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, 0, fmt.Errorf("CSV has no data rows")
	}

	// Build header index from first row
	header := records[0]
	idx := buildIndex(header)

	var panels []model.PanelCatalog
	parseErrs := 0

	for i := 1; i < len(records); i++ {
		row := records[i]
		panel, err := s.parsePanelRow(row, idx)
		if err != nil {
			parseErrs++
			continue
		}
		if panel.PowerWp < minPanelPower || panel.PowerWp > maxPanelPower {
			continue
		}
		panels = append(panels, *panel)
	}

	log.Printf("CEC panels: parsed %d, skipped %d errors", len(panels), parseErrs)
	return panels, parseErrs, nil
}

func (s *CECCatalogService) parsePanelRow(row []string, idx map[string]int) (*model.PanelCatalog, error) {
	getName := func(key string) string {
		if i, ok := idx[key]; ok && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}
	getFloat := func(key string) float64 {
		if i, ok := idx[key]; ok && i < len(row) {
			v, _ := strconv.ParseFloat(strings.TrimSpace(row[i]), 64)
			return v
		}
		return 0
	}

	name := getName("Name")
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}

	// Parse manufacturer from name: typically "Manufacturer Model"
	manufacturer := getName("Manufacturer")
	modelName := name
	if manufacturer != "" && strings.HasPrefix(name, manufacturer) {
		modelName = strings.TrimSpace(strings.TrimPrefix(name, manufacturer))
	}

	stc := getFloat("STC")
	if stc <= 0 {
		return nil, fmt.Errorf("invalid STC power")
	}

	area := getFloat("A_c")
	voc := getFloat("V_oc_ref")
	isc := getFloat("I_sc_ref")
	vmp := getFloat("V_mp_ref")
	imp := getFloat("I_mp_ref")
	noct := getFloat("T_NOCT")
	gammaPmp := getFloat("gamma_pmp")   // %/C already (negative)
	betaOc := getFloat("beta_oc")       // V/C
	length := getFloat("Length")
	width := getFloat("Width")
	technology := getName("Technology")

	// Calculate efficiency
	efficiency := 0.0
	if area > 0 {
		efficiency = stc / (area * 1000) // STC W / (area m² * 1000 W/m²)
	}

	// Map technology
	panelType := mapCECTechnology(technology)

	// Temperature coefficients: gamma_pmp is already in %/C
	tempCoeffPmax := gammaPmp
	// beta_oc is V/C, convert to %/C relative to Voc
	tempCoeffVoc := 0.0
	if voc > 0 {
		tempCoeffVoc = (betaOc / voc) * 100
	}

	// Estimate cost
	costCOP := stc * copPerWattPanel

	panel := &model.PanelCatalog{
		Manufacturer:  manufacturer,
		Model:         modelName,
		Type:          panelType,
		PowerWp:       math.Round(stc*100) / 100,
		Efficiency:    math.Round(efficiency*10000) / 10000,
		Area:          math.Round(area*100) / 100,
		Voc:           math.Round(voc*100) / 100,
		Isc:           math.Round(isc*100) / 100,
		Vmp:           math.Round(vmp*100) / 100,
		Imp:           math.Round(imp*100) / 100,
		TempCoeffPmax: math.Round(tempCoeffPmax*1000) / 1000,
		TempCoeffVoc:  math.Round(tempCoeffVoc*1000) / 1000,
		NOCT:          noct,
		PanelDimensions: model.Dimensions{
			Length: math.Round(length*1000) / 1000,
			Width:  math.Round(width*1000) / 1000,
		},
		Warranty: 25,
		CostCOP:  math.Round(costCOP),
		IsActive: true,
		Source:   "cec",
		CecID:    name,
	}

	return panel, nil
}

func (s *CECCatalogService) fetchInverters(ctx context.Context) ([]model.InverterCatalog, int, error) {
	resp, err := s.httpClient.Get(cecInvertersURL)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP GET inverters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("CEC inverters HTTP %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, 0, fmt.Errorf("reading header: %w", err)
	}
	idx := buildIndex(header)

	var inverters []model.InverterCatalog
	parseErrs := 0

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			parseErrs++
			continue
		}
		inv, err := s.parseInverterRow(row, idx)
		if err != nil {
			parseErrs++
			continue
		}
		// Only import inverters >= 1 kW
		if inv.RatedPowerKw < 1 {
			continue
		}
		inverters = append(inverters, *inv)
	}

	log.Printf("CEC inverters: parsed %d, skipped %d errors", len(inverters), parseErrs)
	return inverters, parseErrs, nil
}

func (s *CECCatalogService) parseInverterRow(row []string, idx map[string]int) (*model.InverterCatalog, error) {
	getName := func(key string) string {
		if i, ok := idx[key]; ok && i < len(row) {
			return strings.TrimSpace(row[i])
		}
		return ""
	}
	getFloat := func(key string) float64 {
		if i, ok := idx[key]; ok && i < len(row) {
			v, _ := strconv.ParseFloat(strings.TrimSpace(row[i]), 64)
			return v
		}
		return 0
	}

	name := getName("Name")
	if name == "" {
		return nil, fmt.Errorf("empty name")
	}

	manufacturer := getName("Manufacturer")
	modelName := name
	if manufacturer != "" && strings.HasPrefix(name, manufacturer) {
		modelName = strings.TrimSpace(strings.TrimPrefix(name, manufacturer))
	}

	// Paco = AC power output rating (W)
	paco := getFloat("Paco")
	if paco <= 0 {
		return nil, fmt.Errorf("invalid Paco")
	}

	pdco := getFloat("Pdco")     // DC power input rating (W)
	vdcmax := getFloat("Vdcmax") // Max DC voltage
	idcmax := getFloat("Idcmax") // Max DC current
	mpptLow := getFloat("Mppt_low")
	mpptHigh := getFloat("Mppt_high")
	vac := getFloat("Vac")       // AC output voltage

	// Efficiency
	efficiency := 0.0
	if pdco > 0 {
		efficiency = paco / pdco
	}

	ratedPowerKw := math.Round(paco/10) / 100 // W to kW with 2 decimals
	maxDCPowerKw := math.Round(pdco/10) / 100

	// Estimate cost
	costCOP := ratedPowerKw * copPerKwInverter

	inv := &model.InverterCatalog{
		Manufacturer:    manufacturer,
		Model:           modelName,
		Type:            "grid-tie",
		RatedPowerKw:    ratedPowerKw,
		MaxDCPowerKw:    maxDCPowerKw,
		Efficiency:      math.Round(efficiency*10000) / 10000,
		MPPTCount:       1, // CEC CSV doesn't specify MPPT count
		MPPTVoltageMin:  math.Round(mpptLow*10) / 10,
		MPPTVoltageMax:  math.Round(mpptHigh*10) / 10,
		MaxInputVoltage: math.Round(vdcmax*10) / 10,
		MaxInputCurrent: math.Round(idcmax*100) / 100,
		OutputVoltage:   vac,
		OutputPhases:    1,
		HasBatteryPort:  false,
		Warranty:        10,
		CostCOP:         math.Round(costCOP),
		IsActive:        true,
		Source:          "cec",
		CecID:          name,
	}

	return inv, nil
}

func mapCECTechnology(tech string) string {
	switch {
	case strings.Contains(strings.ToLower(tech), "mono"):
		return "monocrystalline"
	case strings.Contains(strings.ToLower(tech), "multi"), strings.Contains(strings.ToLower(tech), "poly"):
		return "polycrystalline"
	case strings.Contains(strings.ToLower(tech), "thin"), strings.Contains(strings.ToLower(tech), "cigs"),
		strings.Contains(strings.ToLower(tech), "cdte"), strings.Contains(strings.ToLower(tech), "a-si"):
		return "thin-film"
	default:
		return "monocrystalline"
	}
}

func buildIndex(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}
	return idx
}
