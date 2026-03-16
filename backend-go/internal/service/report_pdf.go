package service

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/dev13/calculadora-paneles-backend/internal/config"
	"github.com/dev13/calculadora-paneles-backend/internal/model"
	"github.com/go-pdf/fpdf"
)

// Color constants matching pdfmake's style
var (
	colorPrimary   = [3]int{37, 99, 235}   // #2563eb blue
	colorHeaderBg  = [3]int{37, 99, 235}   // Blue header background
	colorHeaderTxt = [3]int{255, 255, 255}  // White header text
	colorAltRow    = [3]int{245, 247, 250}  // Light gray alternating row
	colorWhite     = [3]int{255, 255, 255}
	colorBlack     = [3]int{0, 0, 0}
	colorGrayLine  = [3]int{220, 220, 220}
)

type ReportPDFGenerator struct{}

func NewReportPDFGenerator() *ReportPDFGenerator {
	return &ReportPDFGenerator{}
}

func (g *ReportPDFGenerator) Generate(project model.Project, scenario model.Scenario) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 20)

	// Use cp1252 translation for Latin-1 chars (covers Spanish: ñ, á, é, í, ó, ú, ü)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1252")

	pdf.SetHeaderFunc(func() {
		pdf.SetFont("Helvetica", "I", 8)
		pdf.SetTextColor(150, 150, 150)
		pdf.CellFormat(0, 5, tr("Calculadora de Paneles Solares - Colombia"), "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 5, time.Now().Format("02/01/2006"), "", 1, "R", false, 0, "")
		pdf.Ln(2)
	})

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Helvetica", "I", 8)
		pdf.SetTextColor(150, 150, 150)
		pdf.CellFormat(0, 10, fmt.Sprintf("Página %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("{nb}")

	pdf.AddPage()

	// ═══════════ TITLE ═══════════
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.CellFormat(0, 14, tr("Reporte de Dimensionamiento Solar"), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Subtitle: project name
	pdf.SetFont("Helvetica", "B", 14)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.CellFormat(0, 10, tr(fmt.Sprintf("Proyecto: %s", project.Name)), "", 1, "L", false, 0, "")

	// Date and scenario name
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("Fecha: %s", time.Now().Format("02/01/2006"))), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, tr(fmt.Sprintf("Escenario: %s", scenario.Name)), "", 1, "L", false, 0, "")
	pdf.Ln(6)

	// Horizontal rule
	g.drawHRule(pdf)
	pdf.Ln(6)

	// ═══════════ UBICACIÓN ═══════════
	g.sectionTitle(pdf, tr, "Ubicación")
	dept := project.Location.Department
	if dept == "" {
		dept = "-"
	}
	city := project.Location.City
	if city == "" {
		city = "-"
	}
	g.kvTable(pdf, tr, [][2]string{
		{"Latitud", fmt.Sprintf("%.4f°", project.Location.Latitude)},
		{"Longitud", fmt.Sprintf("%.4f°", project.Location.Longitude)},
		{"Altitud", fmt.Sprintf("%.0f m.s.n.m.", project.Location.Altitude)},
		{"Departamento", dept},
		{"Ciudad", city},
		{"Zona climática", project.Location.ClimateZone},
	})
	pdf.Ln(6)

	// ═══════════ CONSUMO ═══════════
	g.sectionTitle(pdf, tr, "Consumo Eléctrico")
	annualConsumption := 0.0
	for _, v := range project.Consumption.Monthly {
		annualConsumption += v
	}
	g.kvTable(pdf, tr, [][2]string{
		{"Consumo anual", fmt.Sprintf("%.0f kWh", annualConsumption)},
		{"Consumo diario promedio", fmt.Sprintf("%.1f kWh", annualConsumption/365)},
		{"Tarifa", fmt.Sprintf("%s/kWh", formatCOP(project.Consumption.TariffPerKwh))},
		{"Estrato", fmt.Sprintf("%d", project.Consumption.Estrato)},
		{"Conexión", project.Consumption.ConnectionType},
		{"Tipo de sistema", project.SystemType},
	})
	pdf.Ln(6)

	// ═══════════ DISEÑO DEL SISTEMA ═══════════
	g.sectionTitle(pdf, tr, "Diseño del Sistema")
	g.kvTable(pdf, tr, [][2]string{
		{"Potencia requerida", fmt.Sprintf("%.2f kWp", scenario.SystemDesign.RequiredPowerKwp)},
		{"Número de paneles", fmt.Sprintf("%d", scenario.SystemDesign.NumberOfPanels)},
		{"Potencia instalada", fmt.Sprintf("%.2f kWp", scenario.SystemDesign.ActualPowerKwp)},
		{"Capacidad inversor", fmt.Sprintf("%.1f kW", scenario.SystemDesign.InverterCapacityKw)},
		{"Utilización del techo", fmt.Sprintf("%.1f%%", scenario.SystemDesign.RoofUtilization)},
	})
	pdf.Ln(2)

	// String configuration sub-section
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(80, 80, 80)
	pdf.CellFormat(0, 6, tr("Configuración de strings:"), "", 1, "L", false, 0, "")
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	g.kvTable(pdf, tr, [][2]string{
		{"Paneles por string", fmt.Sprintf("%d", scenario.SystemDesign.StringConfiguration.PanelsPerString)},
		{"Número de strings", fmt.Sprintf("%d", scenario.SystemDesign.StringConfiguration.NumberOfStrings)},
		{"Voltaje de string", fmt.Sprintf("%.1f V", scenario.SystemDesign.StringConfiguration.StringVoltage)},
		{"Corriente de string", fmt.Sprintf("%.1f A", scenario.SystemDesign.StringConfiguration.StringCurrent)},
	})
	pdf.Ln(2)

	// Battery bank if present
	if scenario.SystemDesign.BatteryBank != nil {
		bb := scenario.SystemDesign.BatteryBank
		pdf.SetFont("Helvetica", "B", 9)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(0, 6, tr("Banco de baterías:"), "", 1, "L", false, 0, "")
		pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
		g.kvTable(pdf, tr, [][2]string{
			{"Capacidad", fmt.Sprintf("%.1f kWh", bb.CapacityKwh)},
			{"Autonomía", fmt.Sprintf("%.0f días", bb.AutonomyDays)},
			{"Número de baterías", fmt.Sprintf("%d", bb.NumberOfBatteries)},
			{"Voltaje del banco", fmt.Sprintf("%.0f V", bb.BankVoltage)},
		})
		pdf.Ln(2)
	}
	pdf.Ln(4)

	// ═══════════ PRODUCCIÓN ESTIMADA (TABLE) ═══════════
	g.sectionTitle(pdf, tr, "Producción Estimada")

	// Table header
	colWidths := []float64{35, 50, 50, 50}
	headers := []string{"Mes", "Irradiación (kWh/m²/día)", "Producción (kWh)", "Ahorro (COP)"}

	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetFillColor(colorHeaderBg[0], colorHeaderBg[1], colorHeaderBg[2])
	pdf.SetTextColor(colorHeaderTxt[0], colorHeaderTxt[1], colorHeaderTxt[2])
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 7, tr(h), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table body
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	for i := 0; i < 12; i++ {
		if i%2 == 0 {
			pdf.SetFillColor(colorAltRow[0], colorAltRow[1], colorAltRow[2])
		} else {
			pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
		}
		fill := true

		pdf.CellFormat(colWidths[0], 5.5, tr(config.Months[i]), "LR", 0, "L", fill, 0, "")

		poa := "-"
		if i < len(scenario.Irradiation.MonthlyPOA) {
			poa = fmt.Sprintf("%.2f", scenario.Irradiation.MonthlyPOA[i])
		}
		pdf.CellFormat(colWidths[1], 5.5, poa, "LR", 0, "R", fill, 0, "")

		prod := "-"
		if i < len(scenario.Production.MonthlyKwh) {
			prod = fmt.Sprintf("%.0f", scenario.Production.MonthlyKwh[i])
		}
		pdf.CellFormat(colWidths[2], 5.5, prod, "LR", 0, "R", fill, 0, "")

		savings := "-"
		if i < len(scenario.Financial.MonthlySavingsCOP) {
			savings = formatCOP(scenario.Financial.MonthlySavingsCOP[i])
		}
		pdf.CellFormat(colWidths[3], 5.5, savings, "LR", 1, "R", fill, 0, "")
	}

	// Total row
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetFillColor(colorHeaderBg[0], colorHeaderBg[1], colorHeaderBg[2])
	pdf.SetTextColor(colorHeaderTxt[0], colorHeaderTxt[1], colorHeaderTxt[2])
	pdf.CellFormat(colWidths[0], 6, tr("Total Anual"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(colWidths[1], 6, "", "1", 0, "C", true, 0, "")
	pdf.CellFormat(colWidths[2], 6, fmt.Sprintf("%.0f kWh", scenario.Production.AnnualKwh), "1", 0, "R", true, 0, "")
	pdf.CellFormat(colWidths[3], 6, formatCOP(scenario.Financial.AnnualSavingsCOP), "1", 1, "R", true, 0, "")
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.Ln(8)

	// ═══════════ PÉRDIDAS ═══════════
	g.sectionTitle(pdf, tr, "Desglose de Pérdidas")
	lossRows := [][2]string{
		{"Sombreado", fmt.Sprintf("%.1f%%", scenario.Losses.ShadingPercent)},
		{"Temperatura", fmt.Sprintf("%.1f%%", scenario.Losses.TemperaturePercent)},
		{"Cableado DC/AC", fmt.Sprintf("%.1f%%", scenario.Losses.WiringPercent)},
		{"Eficiencia del inversor", fmt.Sprintf("%.1f%%", scenario.Losses.InverterPercent)},
		{"Suciedad (soiling)", fmt.Sprintf("%.1f%%", scenario.Losses.SoilingPercent)},
	}
	g.kvTable(pdf, tr, lossRows)
	// Bold total
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetFillColor(colorAltRow[0], colorAltRow[1], colorAltRow[2])
	pdf.CellFormat(80, 7, tr("Pérdida total del sistema"), "B", 0, "L", true, 0, "")
	pdf.CellFormat(70, 7, fmt.Sprintf("%.1f%%", scenario.Losses.TotalSystemLoss), "B", 1, "R", true, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.Ln(8)

	// ═══════════ ANÁLISIS FINANCIERO ═══════════
	g.sectionTitle(pdf, tr, "Análisis Financiero")

	paybackStr := "No se recupera"
	if scenario.Financial.PaybackYears != nil {
		paybackStr = fmt.Sprintf("%.1f años", *scenario.Financial.PaybackYears)
	}
	lcoeVal := scenario.Financial.LCOE.Float64()
	lcoeStr := "N/A"
	if !math.IsInf(lcoeVal, 0) && !math.IsNaN(lcoeVal) {
		lcoeStr = fmt.Sprintf("%s/kWh", formatCOP(lcoeVal))
	}

	g.kvTable(pdf, tr, [][2]string{
		{"Costo de instalación", formatCOP(scenario.Financial.InstallationCostCOP)},
		{"Ahorro anual (Año 1)", formatCOP(scenario.Financial.AnnualSavingsCOP)},
		{"Período de recuperación", paybackStr},
		{"TIR (Tasa Interna de Retorno)", fmt.Sprintf("%.1f%%", scenario.Financial.IRRPercent)},
		{"VPN (Valor Presente Neto)", formatCOP(scenario.Financial.NPVCOP)},
		{"LCOE", lcoeStr},
		{"CO₂ evitado por año", fmt.Sprintf("%.2f toneladas", scenario.Financial.CO2AvoidedTonsYear)},
	})
	pdf.Ln(8)

	// ═══════════ FUENTE DE IRRADIACIÓN ═══════════
	g.sectionTitle(pdf, tr, "Fuente de Datos de Irradiación")
	sourceNames := map[string]string{
		"cache":          "Caché local",
		"ideam":          "IDEAM (Atlas Solar Colombia)",
		"pvgis":          "PVGIS (JRC European Commission)",
		"nasa_power":     "NASA POWER",
		"ideam_fallback": "IDEAM (zona más cercana)",
	}
	sourceName := scenario.Irradiation.Source
	if name, ok := sourceNames[sourceName]; ok {
		sourceName = name
	}
	g.kvTable(pdf, tr, [][2]string{
		{"Fuente", sourceName},
		{"HSP promedio anual", fmt.Sprintf("%.2f horas", scenario.Irradiation.AnnualAvgHSP)},
	})

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// sectionTitle draws a styled section title with blue color and underline.
func (g *ReportPDFGenerator) sectionTitle(pdf *fpdf.Fpdf, tr func(string) string, title string) {
	pdf.SetFont("Helvetica", "B", 12)
	pdf.SetTextColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.CellFormat(0, 8, tr(title), "", 1, "L", false, 0, "")
	// Underline
	pdf.SetDrawColor(colorPrimary[0], colorPrimary[1], colorPrimary[2])
	pdf.SetLineWidth(0.5)
	y := pdf.GetY()
	pdf.Line(pdf.GetX(), y, 195, y)
	pdf.Ln(3)
	pdf.SetDrawColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetLineWidth(0.2)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetFont("Helvetica", "", 10)
}

// kvTable draws a key-value table with alternating row backgrounds.
func (g *ReportPDFGenerator) kvTable(pdf *fpdf.Fpdf, tr func(string) string, rows [][2]string) {
	for i, row := range rows {
		if i%2 == 0 {
			pdf.SetFillColor(colorAltRow[0], colorAltRow[1], colorAltRow[2])
		} else {
			pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
		}
		pdf.CellFormat(80, 6, tr(row[0]), "B", 0, "L", true, 0, "")
		pdf.CellFormat(70, 6, tr(row[1]), "B", 1, "R", true, 0, "")
	}
}

// drawHRule draws a horizontal rule across the page.
func (g *ReportPDFGenerator) drawHRule(pdf *fpdf.Fpdf) {
	pdf.SetDrawColor(colorGrayLine[0], colorGrayLine[1], colorGrayLine[2])
	pdf.SetLineWidth(0.5)
	y := pdf.GetY()
	pdf.Line(10, y, 200, y)
	pdf.SetDrawColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetLineWidth(0.2)
}

// formatCOP formats a float64 as Colombian Pesos: $ 4.200.000
// Matches Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP',
//
//	minimumFractionDigits: 0, maximumFractionDigits: 0 })
func formatCOP(amount float64) string {
	if math.IsInf(amount, 0) || math.IsNaN(amount) {
		return "$ 0"
	}

	negative := amount < 0
	if negative {
		amount = -amount
	}

	// Round to integer (no decimals, matching TS config)
	whole := int64(math.Round(amount))
	result := ""

	if whole == 0 {
		result = "0"
	} else {
		for whole > 0 {
			if result != "" {
				result = "." + result
			}
			chunk := whole % 1000
			whole /= 1000
			if whole > 0 {
				result = fmt.Sprintf("%03d", chunk) + result
			} else {
				result = fmt.Sprintf("%d", chunk) + result
			}
		}
	}

	if negative {
		return "-$ " + result
	}
	return "$ " + result
}
