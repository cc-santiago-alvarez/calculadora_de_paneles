package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/dev13/calculadora-paneles-backend/data"
	"github.com/dev13/calculadora-paneles-backend/internal/calc"
	"github.com/dev13/calculadora-paneles-backend/internal/httpclient"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/dev13/calculadora-paneles-backend/internal/repository"
)

const cacheTTLDays = 30

type IrradiationResult struct {
	Source       string    `json:"source"`
	MonthlyGHI   []float64 `json:"monthlyGHI"`
	MonthlyPOA   []float64 `json:"monthlyPOA"`
	AnnualAvgHSP float64   `json:"annualAvgHSP"`
	Elevation    float64   `json:"elevation"`
}

type IrradiationService struct {
	cacheRepo       *repository.IrradiationCacheRepo
	pvgisClient     *httpclient.Client
	nasaClient      *httpclient.Client
	elevationClient *httpclient.Client
}

func NewIrradiationService(cacheRepo *repository.IrradiationCacheRepo, pvgisBaseURL, nasaBaseURL string) *IrradiationService {
	timeout := 10 * time.Second
	return &IrradiationService{
		cacheRepo:       cacheRepo,
		pvgisClient:     httpclient.New(pvgisBaseURL, timeout, 3),
		nasaClient:      httpclient.New(nasaBaseURL, timeout, 3),
		elevationClient: httpclient.New("https://api.open-meteo.com", 5*time.Second, 1),
	}
}

func (s *IrradiationService) GetIrradiation(ctx context.Context, latitude, longitude, tilt, azimuth, albedo float64) (*IrradiationResult, error) {
	locationKey := makeLocationKey(latitude, longitude)

	var source string
	var monthlyGHI []float64
	var elevation float64

	// 1. Check cache
	if cached, err := s.getFromCache(ctx, locationKey); err == nil && cached != nil {
		source = "cache"
		monthlyGHI = cached.MonthlyGHI
		elevation = cached.Elevation
	}

	// 2. Try IDEAM exact match
	if source == "" {
		if ideamData := getFromIDEAM(latitude, longitude); ideamData != nil {
			s.saveToCache(ctx, locationKey, "ideam", *ideamData)
			source = "ideam"
			monthlyGHI = ideamData.MonthlyGHI
			elevation = ideamData.Elevation
		}
	}

	// 3. Try PVGIS
	if source == "" {
		if pvgisData, err := s.fetchFromPVGIS(ctx, latitude, longitude); err == nil && pvgisData != nil {
			s.saveToCache(ctx, locationKey, "pvgis", *pvgisData)
			source = "pvgis"
			monthlyGHI = pvgisData.MonthlyGHI
			elevation = pvgisData.Elevation
		} else if err != nil {
			log.Printf("PVGIS fetch failed: %v", err)
		}
	}

	// 4. Try NASA POWER
	if source == "" {
		if nasaData, err := s.fetchFromNASA(ctx, latitude, longitude); err == nil && nasaData != nil {
			s.saveToCache(ctx, locationKey, "nasa_power", *nasaData)
			source = "nasa_power"
			monthlyGHI = nasaData.MonthlyGHI
			elevation = nasaData.Elevation
		} else if err != nil {
			log.Printf("NASA POWER fetch failed: %v", err)
		}
	}

	// 5. Fallback: nearest IDEAM zone
	if source == "" {
		fallback := getNearestIDEAM(latitude, longitude)
		source = "ideam_fallback"
		monthlyGHI = fallback.MonthlyGHI
		elevation = fallback.Elevation
	}

	// Resolve elevation from Open-Meteo if still 0
	if elevation == 0 {
		elevation = s.fetchElevation(ctx, latitude, longitude)
	}

	return s.buildResult(source, monthlyGHI, elevation, latitude, tilt, azimuth, albedo), nil
}

func (s *IrradiationService) buildResult(source string, monthlyGHI []float64, elevation, latitude, tilt, azimuth, albedo float64) *IrradiationResult {
	monthlyPOA := calc.MonthlyGHItoPOA(monthlyGHI, latitude, tilt, azimuth, albedo)
	sum := 0.0
	for _, v := range monthlyPOA {
		sum += v
	}
	return &IrradiationResult{
		Source:       source,
		MonthlyGHI:   monthlyGHI,
		MonthlyPOA:   monthlyPOA,
		AnnualAvgHSP: sum / 12,
		Elevation:    elevation,
	}
}

func makeLocationKey(lat, lon float64) string {
	return fmt.Sprintf("%.2f_%.2f", lat, lon)
}

func (s *IrradiationService) getFromCache(ctx context.Context, locationKey string) (*model.NormalizedGHI, error) {
	cached, err := s.cacheRepo.FindByLocationKey(ctx, locationKey)
	if err != nil || cached == nil {
		return nil, err
	}
	return &cached.Normalized, nil
}

func (s *IrradiationService) saveToCache(ctx context.Context, locationKey, source string, ghiData model.NormalizedGHI) {
	if err := s.cacheRepo.Upsert(ctx, locationKey, source, ghiData, cacheTTLDays); err != nil {
		log.Printf("Failed to save to cache: %v", err)
	}
}

func getFromIDEAM(latitude, longitude float64) *model.NormalizedGHI {
	for _, zone := range data.IDEAMZones {
		latDiff := math.Abs(zone.Latitude - latitude)
		lonDiff := math.Abs(zone.Longitude - longitude)
		if latDiff < 0.5 && lonDiff < 0.5 {
			ghi := make([]float64, 12)
			copy(ghi, zone.MonthlyGHI[:])
			return &model.NormalizedGHI{
				MonthlyGHI: ghi,
				AnnualGHI:  zone.AnnualAvgGHI * 365,
				Elevation:  0,
			}
		}
	}
	return nil
}

func getNearestIDEAM(latitude, longitude float64) model.NormalizedGHI {
	var nearest *data.IDEAMZone
	minDist := math.Inf(1)

	for _, zone := range data.IDEAMZones {
		z := zone // capture
		dist := math.Sqrt(math.Pow(z.Latitude-latitude, 2) + math.Pow(z.Longitude-longitude, 2))
		if dist < minDist {
			minDist = dist
			nearest = &z
		}
	}

	ghi := make([]float64, 12)
	copy(ghi, nearest.MonthlyGHI[:])
	return model.NormalizedGHI{
		MonthlyGHI: ghi,
		AnnualGHI:  nearest.AnnualAvgGHI * 365,
		Elevation:  0,
	}
}

func (s *IrradiationService) fetchFromPVGIS(ctx context.Context, latitude, longitude float64) (*model.NormalizedGHI, error) {
	var result map[string]interface{}
	err := s.pvgisClient.GetJSON(ctx, "/MRcalc", map[string]string{
		"lat":          strconv.FormatFloat(latitude, 'f', 4, 64),
		"lon":          strconv.FormatFloat(longitude, 'f', 4, 64),
		"horirrad":     "1",
		"outputformat": "json",
	}, &result)
	if err != nil {
		return nil, err
	}

	outputs, ok := result["outputs"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no outputs in PVGIS response")
	}
	monthly, ok := outputs["monthly"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no monthly data in PVGIS response")
	}

	monthlyGHI := make([]float64, 12)
	for i, m := range monthly {
		if i >= 12 {
			break
		}
		mm, ok := m.(map[string]interface{})
		if !ok {
			continue
		}
		if val, ok := mm["H(h)_m"].(float64); ok {
			monthlyGHI[i] = val / 30 // Convert monthly total to daily avg
		}
	}

	var elevation float64
	if inputs, ok := result["inputs"].(map[string]interface{}); ok {
		if loc, ok := inputs["location"].(map[string]interface{}); ok {
			if elev, ok := loc["elevation"].(float64); ok {
				elevation = elev
			}
		}
	}

	sum := 0.0
	for _, v := range monthlyGHI {
		sum += v
	}

	return &model.NormalizedGHI{
		MonthlyGHI: monthlyGHI,
		AnnualGHI:  sum * 365 / 12,
		Elevation:  elevation,
	}, nil
}

func (s *IrradiationService) fetchFromNASA(ctx context.Context, latitude, longitude float64) (*model.NormalizedGHI, error) {
	var result map[string]interface{}
	err := s.nasaClient.GetJSON(ctx, "", map[string]string{
		"parameters": "ALLSKY_SFC_SW_DWN",
		"community":  "RE",
		"longitude":  strconv.FormatFloat(longitude, 'f', 4, 64),
		"latitude":   strconv.FormatFloat(latitude, 'f', 4, 64),
		"start":      "2001",
		"end":        "2020",
		"format":     "JSON",
	}, &result)
	if err != nil {
		return nil, err
	}

	properties, ok := result["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no properties in NASA response")
	}
	parameter, ok := properties["parameter"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no parameter in NASA response")
	}
	solarData, ok := parameter["ALLSKY_SFC_SW_DWN"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no ALLSKY_SFC_SW_DWN in NASA response")
	}

	// Average monthly values across all years
	monthlyTotals := make([]float64, 12)
	monthlyCounts := make([]int, 12)

	for key, value := range solarData {
		if key == "ANN" {
			continue
		}
		if len(key) >= 6 {
			monthStr := key[4:6]
			month, err := strconv.Atoi(monthStr)
			if err != nil {
				continue
			}
			month-- // 0-indexed
			if month >= 0 && month < 12 {
				if val, ok := value.(float64); ok && val > 0 {
					monthlyTotals[month] += val
					monthlyCounts[month]++
				}
			}
		}
	}

	monthlyGHI := make([]float64, 12)
	for i := 0; i < 12; i++ {
		if monthlyCounts[i] > 0 {
			monthlyGHI[i] = monthlyTotals[i] / float64(monthlyCounts[i])
		} else {
			monthlyGHI[i] = 4.0
		}
	}

	sum := 0.0
	for _, v := range monthlyGHI {
		sum += v
	}

	return &model.NormalizedGHI{
		MonthlyGHI: monthlyGHI,
		AnnualGHI:  sum * 365 / 12,
		Elevation:  0,
	}, nil
}

func (s *IrradiationService) fetchElevation(ctx context.Context, latitude, longitude float64) float64 {
	var result struct {
		Elevation []float64 `json:"elevation"`
	}
	err := s.elevationClient.GetJSON(ctx, "/v1/elevation", map[string]string{
		"latitude":  strconv.FormatFloat(latitude, 'f', 4, 64),
		"longitude": strconv.FormatFloat(longitude, 'f', 4, 64),
	}, &result)
	if err != nil {
		log.Printf("Open-Meteo elevation fetch failed: %v", err)
		return 0
	}
	if len(result.Elevation) > 0 {
		return result.Elevation[0]
	}
	return 0
}
