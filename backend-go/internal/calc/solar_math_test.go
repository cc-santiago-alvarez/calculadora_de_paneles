package calc

import (
	"math"
	"testing"
)

func TestSolarDeclination_Equinox(t *testing.T) {
	dec := SolarDeclination(80) // March 21
	if math.Abs(dec) >= 2 {
		t.Errorf("expected declination near 0 at equinox, got %f", dec)
	}
}

func TestSolarDeclination_SummerSolstice(t *testing.T) {
	dec := SolarDeclination(172) // June 21
	if dec <= 20 || dec >= 25 {
		t.Errorf("expected declination between 20 and 25, got %f", dec)
	}
}

func TestSolarDeclination_WinterSolstice(t *testing.T) {
	dec := SolarDeclination(355) // Dec 21
	if dec >= -20 || dec <= -25 {
		t.Errorf("expected declination between -25 and -20, got %f", dec)
	}
}

func TestSunsetHourAngle_EquatorEquinox(t *testing.T) {
	ws := SunsetHourAngle(0, 0)
	if math.Abs(ws-90) > 1 {
		t.Errorf("expected ~90 degrees, got %f", ws)
	}
}

func TestSunsetHourAngle_HighLatitudeSummer(t *testing.T) {
	summer := SunsetHourAngle(60, 23.45)
	equinox := SunsetHourAngle(60, 0)
	if summer <= equinox {
		t.Errorf("expected summer hour angle (%f) > equinox (%f)", summer, equinox)
	}
}

func TestExtraterrestrialIrradiation_Bogota(t *testing.T) {
	H0 := ExtraterrestrialIrradiation(4.61, 80)
	if H0 <= 8 || H0 >= 12 {
		t.Errorf("expected H0 between 8 and 12, got %f", H0)
	}
}

func TestLiuJordanTransposition_Colombian(t *testing.T) {
	poa := LiuJordanTransposition(4.3, 4.61, 10, 0, 80, 0.2)
	if poa <= 3 || poa >= 7 {
		t.Errorf("expected POA between 3 and 7, got %f", poa)
	}
}

func TestLiuJordanTransposition_TiltImprovement(t *testing.T) {
	poaFlat := LiuJordanTransposition(4.3, 4.61, 0, 0, 80, 0.2)
	poaTilted := LiuJordanTransposition(4.3, 4.61, 10, 0, 80, 0.2)
	if poaTilted < poaFlat*0.9 {
		t.Errorf("expected tilted POA (%f) >= 90%% of flat (%f)", poaTilted, poaFlat)
	}
}

func TestMonthlyGHItoPOA_Length(t *testing.T) {
	ghiBogota := []float64{4.4, 4.6, 4.5, 4.2, 4.0, 4.2, 4.4, 4.3, 4.1, 3.9, 4.1, 4.3}
	poa := MonthlyGHItoPOA(ghiBogota, 4.61, 10, 0, 0.2)
	if len(poa) != 12 {
		t.Errorf("expected 12 values, got %d", len(poa))
	}
	for i, v := range poa {
		if v <= 0 || v >= 8 {
			t.Errorf("month %d: expected POA between 0 and 8, got %f", i, v)
		}
	}
}

func TestCellTemperature_AboveAmbient(t *testing.T) {
	tCell := CellTemperature(25, 45, 800)
	if tCell <= 25 {
		t.Errorf("expected cell temp > 25, got %f", tCell)
	}
}

func TestCellTemperature_Formula(t *testing.T) {
	// T_cell = 25 + (45-20)/800 * 800 = 25 + 25 = 50
	tCell := CellTemperature(25, 45, 800)
	if math.Abs(tCell-50) > 0.1 {
		t.Errorf("expected ~50, got %f", tCell)
	}
}

func TestTemperatureLoss_AtSTC(t *testing.T) {
	loss := TemperatureLoss(25, -0.35)
	if loss != 0 {
		t.Errorf("expected 0 loss at STC, got %f", loss)
	}
}

func TestTemperatureLoss_Increase(t *testing.T) {
	loss := TemperatureLoss(50, -0.35)
	// 0.35% * 25°C = 8.75%
	expected := 0.0875
	if math.Abs(loss-expected) > 0.001 {
		t.Errorf("expected %f, got %f", expected, loss)
	}
}
