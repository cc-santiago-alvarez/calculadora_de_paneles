package service

import (
	"math"

	"github.com/dev13/calculadora-paneles-backend/internal/calc"
	"github.com/dev13/calculadora-paneles-backend/internal/config"
)

type Obstacle struct {
	AzimuthStart   float64 `json:"azimuthStart"`
	AzimuthEnd     float64 `json:"azimuthEnd"`
	ElevationAngle float64 `json:"elevationAngle"`
}

type ShadingResult struct {
	MonthlyLoss []float64 `json:"monthlyLoss"`
	AnnualLoss  float64   `json:"annualLoss"`
}

type ShadingAnalyzer struct{}

func NewShadingAnalyzer() *ShadingAnalyzer {
	return &ShadingAnalyzer{}
}

func (a *ShadingAnalyzer) Analyze(latitude, longitude float64, obstacles []Obstacle) ShadingResult {
	if len(obstacles) == 0 {
		return ShadingResult{
			MonthlyLoss: make([]float64, 12),
			AnnualLoss:  0,
		}
	}

	monthlyLoss := make([]float64, 12)

	for m := 0; m < 12; m++ {
		dayOfYear := calc.RepresentativeDays[m]
		var totalIrradiance, shadedIrradiance float64

		for hour := 6.0; hour <= 18.0; hour += 0.5 {
			pos := calc.SolarPosition(latitude, longitude, dayOfYear, hour)
			if pos.Altitude <= 0 {
				continue
			}

			weight := math.Sin(pos.Altitude * math.Pi / 180)
			totalIrradiance += weight

			isShaded := false
			for _, obs := range obstacles {
				if isInAzimuthRange(pos.Azimuth, obs.AzimuthStart, obs.AzimuthEnd) && pos.Altitude < obs.ElevationAngle {
					isShaded = true
					break
				}
			}

			if isShaded {
				shadedIrradiance += weight
			}
		}

		if totalIrradiance > 0 {
			monthlyLoss[m] = shadedIrradiance / totalIrradiance
		}
	}

	totalDays := 0
	for _, d := range config.DaysInMonth {
		totalDays += d
	}
	annualLoss := 0.0
	for m, loss := range monthlyLoss {
		annualLoss += loss * float64(config.DaysInMonth[m])
	}
	annualLoss /= float64(totalDays)

	return ShadingResult{
		MonthlyLoss: monthlyLoss,
		AnnualLoss:  annualLoss,
	}
}

func isInAzimuthRange(azimuth, start, end float64) bool {
	azimuth = math.Mod(math.Mod(azimuth, 360)+360, 360)
	start = math.Mod(math.Mod(start, 360)+360, 360)
	end = math.Mod(math.Mod(end, 360)+360, 360)

	if start <= end {
		return azimuth >= start && azimuth <= end
	}
	return azimuth >= start || azimuth <= end
}
