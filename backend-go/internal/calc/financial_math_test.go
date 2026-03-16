package calc

import (
	"math"
	"testing"
)

func TestNPV_Simple(t *testing.T) {
	cashFlows := []float64{-10000, 3000, 3000, 3000, 3000, 3000}
	result := NPV(cashFlows, 0.1)
	if result <= 1000 || result >= 2000 {
		t.Errorf("expected NPV between 1000 and 2000, got %f", result)
	}
}

func TestNPV_InitialOnly(t *testing.T) {
	result := NPV([]float64{-5000}, 0.1)
	if result != -5000 {
		t.Errorf("expected -5000, got %f", result)
	}
}

func TestIRR_SimpleCashFlows(t *testing.T) {
	cashFlows := []float64{-100, 50, 50, 50}
	result := IRR(cashFlows, 100, 1e-7)
	if result <= 0.2 || result >= 0.3 {
		t.Errorf("expected IRR between 0.2 and 0.3, got %f", result)
	}
}

func TestIRR_BreakEven(t *testing.T) {
	cashFlows := []float64{-100, 100}
	result := IRR(cashFlows, 100, 1e-7)
	if math.Abs(result) > 0.01 {
		t.Errorf("expected IRR near 0, got %f", result)
	}
}

func TestPaybackPeriod_Normal(t *testing.T) {
	result := PaybackPeriod(10000, []float64{3000, 3000, 3000, 3000, 3000})
	if result == nil {
		t.Fatal("expected non-nil payback")
	}
	if *result <= 3 || *result >= 4 {
		t.Errorf("expected payback between 3 and 4, got %f", *result)
	}
}

func TestPaybackPeriod_NeverPaysBack(t *testing.T) {
	result := PaybackPeriod(100000, []float64{100, 100, 100})
	if result != nil {
		t.Errorf("expected nil (never pays back), got %f", *result)
	}
}

func TestLCOE_Simple(t *testing.T) {
	yearlyProd := make([]float64, 25)
	for i := range yearlyProd {
		yearlyProd[i] = 5000
	}
	result := LCOE(15000000, 150000, yearlyProd, 0.08)
	if result <= 200 || result >= 500 {
		t.Errorf("expected LCOE between 200 and 500, got %f", result)
	}
}

func TestGenerateAnnualSavings(t *testing.T) {
	savings := GenerateAnnualSavings(5000, 800, 0.005, 0.04, 3)
	if len(savings) != 3 {
		t.Fatalf("expected 3 values, got %d", len(savings))
	}
	// Year 1: 5000 * 800 = 4,000,000
	if math.Abs(savings[0]-4000000) > 100 {
		t.Errorf("expected ~4000000, got %f", savings[0])
	}
	// Year 2 should be > 99% of year 1 due to escalation > degradation
	if savings[1] <= savings[0]*0.99 {
		t.Errorf("expected year 2 (%f) > 99%% of year 1 (%f)", savings[1], savings[0])
	}
}

func TestCumulativeSavings(t *testing.T) {
	result := CumulativeSavings([]float64{100, 200, 300})
	expected := []float64{100, 300, 600}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("index %d: expected %f, got %f", i, expected[i], v)
		}
	}
}
