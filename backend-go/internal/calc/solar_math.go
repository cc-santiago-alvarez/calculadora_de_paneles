package calc

import "math"

const (
	deg2rad = math.Pi / 180
	rad2deg = 180 / math.Pi
)

// RepresentativeDays is the day of year for the 15th of each month.
var RepresentativeDays = [12]int{17, 47, 75, 105, 135, 162, 198, 228, 258, 288, 318, 344}

// SolarDeclination calculates solar declination angle using Spencer (1971).
// Returns declination in degrees.
func SolarDeclination(dayOfYear int) float64 {
	B := float64(dayOfYear-1) * 360.0 / 365.0
	Br := B * deg2rad
	return (0.006918 -
		0.399912*math.Cos(Br) +
		0.070257*math.Sin(Br) -
		0.006758*math.Cos(2*Br) +
		0.000907*math.Sin(2*Br) -
		0.002697*math.Cos(3*Br) +
		0.00148*math.Sin(3*Br)) * rad2deg
}

// EquationOfTime returns the equation of time in minutes (Spencer 1971).
func EquationOfTime(dayOfYear int) float64 {
	B := float64(dayOfYear-1) * 360.0 / 365.0
	Br := B * deg2rad
	return 229.18 * (0.000075 +
		0.001868*math.Cos(Br) -
		0.032077*math.Sin(Br) -
		0.014615*math.Cos(2*Br) -
		0.04089*math.Sin(2*Br))
}

// SunsetHourAngle returns the sunrise/sunset hour angle in degrees (positive).
func SunsetHourAngle(latitude, declination float64) float64 {
	latRad := latitude * deg2rad
	decRad := declination * deg2rad
	cosWs := -math.Tan(latRad) * math.Tan(decRad)
	if cosWs > 1 {
		return 0 // No sunrise (polar night)
	}
	if cosWs < -1 {
		return 180 // No sunset (midnight sun)
	}
	return math.Acos(cosWs) * rad2deg
}

// ExtraterrestrialIrradiation calculates H0 on horizontal surface in kWh/m2/day.
func ExtraterrestrialIrradiation(latitude float64, dayOfYear int) float64 {
	Gsc := 1.367 // kW/m2 (solar constant)
	latRad := latitude * deg2rad
	dec := SolarDeclination(dayOfYear)
	decRad := dec * deg2rad
	ws := SunsetHourAngle(latitude, dec) * deg2rad
	dr := 1 + 0.033*math.Cos(2*math.Pi*float64(dayOfYear)/365)

	return (24 / math.Pi) * Gsc * dr *
		(math.Cos(latRad)*math.Cos(decRad)*math.Sin(ws) +
			ws*math.Sin(latRad)*math.Sin(decRad))
}

// DiffuseFraction estimates diffuse fraction using Erbs correlation.
func DiffuseFraction(Kt float64) float64 {
	if Kt <= 0.22 {
		return 1.0 - 0.09*Kt
	}
	if Kt <= 0.80 {
		return 0.9511 - 0.1604*Kt + 4.388*math.Pow(Kt, 2) - 16.638*math.Pow(Kt, 3) + 12.336*math.Pow(Kt, 4)
	}
	return 0.165
}

// BeamTiltFactor calculates geometric factor Rb for beam radiation on tilted surface.
func BeamTiltFactor(latitude, tilt, azimuth float64, dayOfYear int) float64 {
	latRad := latitude * deg2rad
	tiltRad := tilt * deg2rad
	azRad := azimuth * deg2rad
	dec := SolarDeclination(dayOfYear)
	decRad := dec * deg2rad

	ws := SunsetHourAngle(latitude, dec) * deg2rad
	wsTilt := SunsetHourAngle(latitude-tilt, dec) * deg2rad
	wsMin := math.Min(ws, wsTilt)

	// Simplified Rb for south-facing (azimuth~0) and near-equatorial latitudes
	if math.Abs(azimuth) < 5 {
		numerator := math.Cos(latRad-tiltRad)*math.Cos(decRad)*math.Sin(wsMin) +
			wsMin*math.Sin(latRad-tiltRad)*math.Sin(decRad)
		denominator := math.Cos(latRad)*math.Cos(decRad)*math.Sin(ws) +
			ws*math.Sin(latRad)*math.Sin(decRad)

		if denominator > 0 {
			return math.Max(0, numerator/denominator)
		}
		return 1
	}

	// For non-south orientations, approximate with hourly calculation
	_ = azRad
	var sumTilted, sumHorizontal float64
	steps := 24
	for h := -12; h < 12; h++ {
		hf := float64(h) * 24.0 / float64(steps)
		wRad := hf * 15 * deg2rad
		if math.Abs(wRad) > ws {
			continue
		}

		cosTheta := math.Sin(decRad)*math.Sin(latRad)*math.Cos(tiltRad) -
			math.Sin(decRad)*math.Cos(latRad)*math.Sin(tiltRad)*math.Cos(azRad) +
			math.Cos(decRad)*math.Cos(latRad)*math.Cos(tiltRad)*math.Cos(wRad) +
			math.Cos(decRad)*math.Sin(latRad)*math.Sin(tiltRad)*math.Cos(azRad)*math.Cos(wRad) +
			math.Cos(decRad)*math.Sin(tiltRad)*math.Sin(azRad)*math.Sin(wRad)

		cosThetaZ := math.Sin(decRad)*math.Sin(latRad) +
			math.Cos(decRad)*math.Cos(latRad)*math.Cos(wRad)

		if cosThetaZ > 0 {
			sumHorizontal += cosThetaZ
			sumTilted += math.Max(0, cosTheta)
		}
	}

	if sumHorizontal > 0 {
		return sumTilted / sumHorizontal
	}
	return 1
}

// LiuJordanTransposition converts horizontal GHI to plane-of-array (POA) irradiation.
// GHI in kWh/m2/day, returns POA in kWh/m2/day.
func LiuJordanTransposition(GHI, latitude, tilt, azimuth float64, dayOfYear int, albedo float64) float64 {
	H0 := ExtraterrestrialIrradiation(latitude, dayOfYear)
	if H0 <= 0 || GHI <= 0 {
		return 0
	}

	Kt := math.Min(GHI/H0, 1)
	fd := DiffuseFraction(Kt)
	DHI := fd * GHI
	DNI := GHI - DHI

	Rb := BeamTiltFactor(latitude, tilt, azimuth, dayOfYear)
	tiltRad := tilt * deg2rad

	beamTilted := DNI * Rb
	diffuseTilted := DHI * (1 + math.Cos(tiltRad)) / 2
	groundReflected := GHI * albedo * (1 - math.Cos(tiltRad)) / 2

	return math.Max(0, beamTilted+diffuseTilted+groundReflected)
}

// MonthlyGHItoPOA converts 12 monthly GHI values to POA values.
func MonthlyGHItoPOA(monthlyGHI []float64, latitude, tilt, azimuth, albedo float64) []float64 {
	poa := make([]float64, 12)
	for i := 0; i < 12 && i < len(monthlyGHI); i++ {
		poa[i] = LiuJordanTransposition(monthlyGHI[i], latitude, tilt, azimuth, RepresentativeDays[i], albedo)
	}
	return poa
}

// CellTemperature estimates cell temperature.
// T_cell = T_amb + (NOCT - 20) / 800 * G
func CellTemperature(ambientTemp, NOCT, irradiance float64) float64 {
	return ambientTemp + (NOCT-20)/800*irradiance
}

// TemperatureLoss calculates temperature loss factor.
// Returns loss as decimal (e.g. 0.05 for 5% loss).
func TemperatureLoss(cellTemp, tempCoeffPmax float64) float64 {
	delta := cellTemp - 25 // STC reference
	if delta <= 0 {
		return 0
	}
	return math.Abs(tempCoeffPmax/100) * delta
}

// SolarPosition calculates solar altitude and azimuth.
type SolarPos struct {
	Altitude float64
	Azimuth  float64
}

func SolarPosition(latitude, longitude float64, dayOfYear int, hourSolar float64) SolarPos {
	latRad := latitude * deg2rad
	dec := SolarDeclination(dayOfYear)
	decRad := dec * deg2rad
	hourAngle := (hourSolar - 12) * 15 * deg2rad

	sinAlt := math.Sin(latRad)*math.Sin(decRad) +
		math.Cos(latRad)*math.Cos(decRad)*math.Cos(hourAngle)
	altitude := math.Asin(sinAlt) * rad2deg

	cosAz := (math.Sin(decRad) - math.Sin(latRad)*sinAlt) /
		(math.Cos(latRad) * math.Cos(math.Asin(sinAlt)))
	azimuthVal := math.Acos(math.Max(-1, math.Min(1, cosAz))) * rad2deg

	if hourAngle > 0 {
		azimuthVal = 360 - azimuthVal
	}

	return SolarPos{Altitude: altitude, Azimuth: azimuthVal}
}
