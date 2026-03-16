package calc

import "math"

// NPV calculates Net Present Value.
func NPV(cashFlows []float64, discountRate float64) float64 {
	sum := 0.0
	for y, cf := range cashFlows {
		sum += cf / math.Pow(1+discountRate, float64(y))
	}
	return sum
}

// IRR calculates Internal Rate of Return using Newton-Raphson.
func IRR(cashFlows []float64, maxIterations int, tolerance float64) float64 {
	if maxIterations <= 0 {
		maxIterations = 100
	}
	if tolerance <= 0 {
		tolerance = 1e-7
	}

	r := 0.1 // Initial guess: 10%

	for i := 0; i < maxIterations; i++ {
		f := 0.0
		fPrime := 0.0

		for t, cf := range cashFlows {
			tf := float64(t)
			denominator := math.Pow(1+r, tf)
			f += cf / denominator
			fPrime += (-tf * cf) / math.Pow(1+r, tf+1)
		}

		if math.Abs(fPrime) < 1e-12 {
			break
		}

		rNew := r - f/fPrime
		if math.Abs(rNew-r) < tolerance {
			return rNew
		}
		r = rNew

		// Guard against divergence
		if r < -0.99 {
			r = -0.5
		}
		if r > 10 {
			r = 5
		}
	}
	return r
}

// PaybackPeriod calculates simple payback period in years.
// Returns nil if investment never pays back (equivalent to Infinity in TS).
func PaybackPeriod(initialInvestment float64, annualSavings []float64) *float64 {
	cumulative := -initialInvestment
	for y, saving := range annualSavings {
		cumulative += saving
		if cumulative >= 0 {
			prevCumulative := cumulative - saving
			result := float64(y) + (-prevCumulative)/saving
			return &result
		}
	}
	return nil // Never pays back
}

// LCOE calculates Levelized Cost of Energy.
func LCOE(initialCost, annualMaintenanceCost float64, annualProductionKwh []float64, discountRate float64) float64 {
	totalCost := initialCost
	totalProduction := 0.0

	for y, prod := range annualProductionKwh {
		discount := math.Pow(1+discountRate, float64(y+1))
		totalCost += annualMaintenanceCost / discount
		totalProduction += prod / discount
	}

	if totalProduction > 0 {
		return totalCost / totalProduction
	}
	return math.Inf(1)
}

// GenerateAnnualSavings generates 25-year annual savings with tariff escalation and degradation.
func GenerateAnnualSavings(baseAnnualKwh, baseTariff, degradationRate, tariffEscalation float64, years int) []float64 {
	if years <= 0 {
		years = 25
	}
	savings := make([]float64, years)
	for y := 0; y < years; y++ {
		production := baseAnnualKwh * math.Pow(1-degradationRate, float64(y))
		tariff := baseTariff * math.Pow(1+tariffEscalation, float64(y))
		savings[y] = production * tariff
	}
	return savings
}

// CumulativeSavings calculates cumulative savings over time.
func CumulativeSavings(annualSavings []float64) []float64 {
	cumulative := make([]float64, len(annualSavings))
	total := 0.0
	for i, saving := range annualSavings {
		total += saving
		cumulative[i] = total
	}
	return cumulative
}
