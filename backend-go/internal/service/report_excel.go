package service

import (
	"fmt"

	"github.com/dev13/calculadora-paneles-backend/internal/config"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/xuri/excelize/v2"
)

type ReportExcelGenerator struct{}

func NewReportExcelGenerator() *ReportExcelGenerator {
	return &ReportExcelGenerator{}
}

func (g *ReportExcelGenerator) Generate(project model.Project, scenario model.Scenario) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Sheet 1: Resumen
	sheet1 := "Resumen"
	f.SetSheetName("Sheet1", sheet1)
	f.SetColWidth(sheet1, "A", "A", 30)
	f.SetColWidth(sheet1, "B", "B", 25)

	f.SetCellValue(sheet1, "A1", "Parámetro")
	f.SetCellValue(sheet1, "B1", "Valor")

	rows := []struct {
		param string
		value interface{}
	}{
		{"Proyecto", project.Name},
		{"Latitud", project.Location.Latitude},
		{"Longitud", project.Location.Longitude},
		{"Potencia Instalada (kWp)", scenario.SystemDesign.ActualPowerKwp},
		{"Número de Paneles", scenario.SystemDesign.NumberOfPanels},
		{"Producción Anual (kWh)", scenario.Production.AnnualKwh},
		{"Costo Instalación (COP)", scenario.Financial.InstallationCostCOP},
		{"Payback (años)", formatPayback(scenario.Financial.PaybackYears)},
		{"TIR (%)", scenario.Financial.IRRPercent},
		{"VPN (COP)", scenario.Financial.NPVCOP},
	}
	for i, r := range rows {
		cell := fmt.Sprintf("A%d", i+2)
		f.SetCellValue(sheet1, cell, r.param)
		cell = fmt.Sprintf("B%d", i+2)
		f.SetCellValue(sheet1, cell, r.value)
	}

	// Sheet 2: Producción Mensual
	sheet2 := "Producción Mensual"
	f.NewSheet(sheet2)
	f.SetColWidth(sheet2, "A", "A", 15)
	f.SetColWidth(sheet2, "B", "E", 18)

	headers2 := []string{"Mes", "GHI (kWh/m²/día)", "POA (kWh/m²/día)", "Producción (kWh)", "Ahorro (COP)"}
	for i, h := range headers2 {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet2, cell, h)
	}
	for i := 0; i < 12; i++ {
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), config.Months[i])
		if i < len(scenario.Irradiation.MonthlyGHI) {
			f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), scenario.Irradiation.MonthlyGHI[i])
		}
		if i < len(scenario.Irradiation.MonthlyPOA) {
			f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), scenario.Irradiation.MonthlyPOA[i])
		}
		if i < len(scenario.Production.MonthlyKwh) {
			f.SetCellValue(sheet2, fmt.Sprintf("D%d", row), scenario.Production.MonthlyKwh[i])
		}
		if i < len(scenario.Financial.MonthlySavingsCOP) {
			f.SetCellValue(sheet2, fmt.Sprintf("E%d", row), scenario.Financial.MonthlySavingsCOP[i])
		}
	}

	// Sheet 3: Proyección 25 Años
	sheet3 := "Proyección 25 Años"
	f.NewSheet(sheet3)
	f.SetColWidth(sheet3, "A", "A", 10)
	f.SetColWidth(sheet3, "B", "C", 25)

	headers3 := []string{"Año", "Producción (kWh)", "Ahorro Acumulado (COP)"}
	for i, h := range headers3 {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheet3, cell, h)
	}
	for y := 0; y < 25; y++ {
		row := y + 2
		f.SetCellValue(sheet3, fmt.Sprintf("A%d", row), y+1)
		if y < len(scenario.Production.Yearly25) {
			f.SetCellValue(sheet3, fmt.Sprintf("B%d", row), scenario.Production.Yearly25[y])
		}
		if y < len(scenario.Financial.CumulativeSavings25) {
			f.SetCellValue(sheet3, fmt.Sprintf("C%d", row), scenario.Financial.CumulativeSavings25[y])
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Excel: %w", err)
	}
	return buf.Bytes(), nil
}

func formatPayback(payback *float64) interface{} {
	if payback == nil {
		return "Nunca"
	}
	return *payback
}
