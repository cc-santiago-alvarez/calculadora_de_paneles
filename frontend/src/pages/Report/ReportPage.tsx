import { useState } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { reportsApi } from '../../api/reports';
import { formatCOP, formatNumber, MONTHS_FULL } from '../../utils/format';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';

export default function ReportPage() {
  const { currentProject, currentScenario } = useProjectStore();
  const [generating, setGenerating] = useState(false);

  if (!currentProject || !currentScenario) {
    return (
      <div className="text-center py-20">
        <p className="text-fg-secondary text-lg">Calcula un proyecto primero para generar reportes.</p>
      </div>
    );
  }

  const s = currentScenario;

  const downloadFile = (blob: Blob, filename: string) => {
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handlePDF = async () => {
    setGenerating(true);
    try {
      const blob = await reportsApi.generatePDF(currentProject._id, s._id);
      downloadFile(new Blob([blob]), `reporte-${currentProject.name}.pdf`);
    } catch (err) {
      console.error('PDF generation failed:', err);
    } finally {
      setGenerating(false);
    }
  };

  const handleExcel = async () => {
    setGenerating(true);
    try {
      const blob = await reportsApi.generateExcel(currentProject._id, s._id);
      downloadFile(new Blob([blob]), `reporte-${currentProject.name}.xlsx`);
    } catch (err) {
      console.error('Excel generation failed:', err);
    } finally {
      setGenerating(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-fg-primary">Reporte del Proyecto</h2>
        <div className="flex gap-3">
          <Button variant="danger" onClick={handlePDF} disabled={generating}>
            {generating ? 'Generando...' : 'Descargar PDF'}
          </Button>
          <Button className="bg-success hover:bg-[hsl(145,63%,36%)] text-white" onClick={handleExcel} disabled={generating}>
            {generating ? 'Generando...' : 'Descargar Excel'}
          </Button>
        </div>
      </div>

      {/* Preview */}
      <Card padding="lg" className="space-y-8">
        <div className="text-center border-b border-[var(--color-border-subtle)] pb-6">
          <h1 className="text-xl font-bold text-fg-primary">Reporte de Dimensionamiento Solar</h1>
          <p className="text-fg-secondary mt-1">{currentProject.name}</p>
          <p className="text-sm text-fg-muted">{new Date().toLocaleDateString('es-CO')}</p>
        </div>

        {/* Location */}
        <section>
          <h3 className="text-fg-secondary font-semibold border-b border-[var(--color-border-subtle)] pb-1 mb-3">Ubicacion</h3>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <p><span className="text-fg-tertiary">Latitud:</span> <span className="font-mono tabular-nums text-fg-primary">{currentProject.location.latitude}</span></p>
            <p><span className="text-fg-tertiary">Longitud:</span> <span className="font-mono tabular-nums text-fg-primary">{currentProject.location.longitude}</span></p>
            <p><span className="text-fg-tertiary">Departamento:</span> <span className="font-mono tabular-nums text-fg-primary">{currentProject.location.department}</span></p>
            <p><span className="text-fg-tertiary">Ciudad:</span> <span className="font-mono tabular-nums text-fg-primary">{currentProject.location.city}</span></p>
          </div>
        </section>

        {/* System Design */}
        <section>
          <h3 className="text-fg-secondary font-semibold border-b border-[var(--color-border-subtle)] pb-1 mb-3">Diseno del Sistema</h3>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <p><span className="text-fg-tertiary">Potencia Requerida:</span> <span className="font-mono tabular-nums text-fg-primary">{s.systemDesign.requiredPowerKwp.toFixed(2)} kWp</span></p>
            <p><span className="text-fg-tertiary">Potencia Instalada:</span> <span className="font-mono tabular-nums text-fg-primary">{s.systemDesign.actualPowerKwp.toFixed(2)} kWp</span></p>
            <p><span className="text-fg-tertiary">Paneles:</span> <span className="font-mono tabular-nums text-fg-primary">{s.systemDesign.numberOfPanels}</span></p>
            <p><span className="text-fg-tertiary">Inversor:</span> <span className="font-mono tabular-nums text-fg-primary">{s.systemDesign.inverterCapacityKw} kW</span></p>
          </div>
        </section>

        {/* Monthly Production Table */}
        <section>
          <h3 className="text-fg-secondary font-semibold border-b border-[var(--color-border-subtle)] pb-1 mb-3">Produccion Mensual</h3>
          <table className="w-full text-sm">
            <thead className="bg-inset">
              <tr>
                <th className="px-3 py-2 text-left text-fg-secondary font-medium">Mes</th>
                <th className="px-3 py-2 text-right text-fg-secondary font-medium">Irradiacion (kWh/m2/dia)</th>
                <th className="px-3 py-2 text-right text-fg-secondary font-medium">Produccion (kWh)</th>
                <th className="px-3 py-2 text-right text-fg-secondary font-medium">Ahorro (COP)</th>
              </tr>
            </thead>
            <tbody>
              {MONTHS_FULL.map((month, i) => (
                <tr key={month} className="border-t border-[var(--color-border-subtle)]">
                  <td className="px-3 py-1 text-fg-primary">{month}</td>
                  <td className="px-3 py-1 text-right font-mono tabular-nums text-fg-primary">{s.irradiation.monthlyPOA[i]?.toFixed(2)}</td>
                  <td className="px-3 py-1 text-right font-mono tabular-nums text-fg-primary">{formatNumber(s.production.monthlyKwh[i])}</td>
                  <td className="px-3 py-1 text-right font-mono tabular-nums text-fg-primary">{formatCOP(s.financial.monthlySavingsCOP[i])}</td>
                </tr>
              ))}
              <tr className="border-t-2 border-[var(--color-border-strong)] font-bold">
                <td className="px-3 py-2 text-fg-primary">Total Anual</td>
                <td className="px-3 py-2 text-right text-fg-primary">-</td>
                <td className="px-3 py-2 text-right font-mono tabular-nums text-fg-primary">{formatNumber(s.production.annualKwh)}</td>
                <td className="px-3 py-2 text-right font-mono tabular-nums text-fg-primary">{formatCOP(s.financial.annualSavingsCOP)}</td>
              </tr>
            </tbody>
          </table>
        </section>

        {/* Financial */}
        <section>
          <h3 className="text-fg-secondary font-semibold border-b border-[var(--color-border-subtle)] pb-1 mb-3">Analisis Financiero</h3>
          <div className="grid grid-cols-2 gap-2 text-sm">
            <p><span className="text-fg-tertiary">Costo Total:</span> <span className="font-mono tabular-nums text-fg-primary">{formatCOP(s.financial.installationCostCOP)}</span></p>
            <p><span className="text-fg-tertiary">Payback:</span> <span className="font-mono tabular-nums text-fg-primary">{s.financial.paybackYears != null ? s.financial.paybackYears.toFixed(1) : 'N/A'} anos</span></p>
            <p><span className="text-fg-tertiary">TIR:</span> <span className="font-mono tabular-nums text-fg-primary">{s.financial.irrPercent.toFixed(1)}%</span></p>
            <p><span className="text-fg-tertiary">VPN:</span> <span className="font-mono tabular-nums text-fg-primary">{formatCOP(s.financial.npvCOP)}</span></p>
            <p><span className="text-fg-tertiary">LCOE:</span> <span className="font-mono tabular-nums text-fg-primary">{formatCOP(s.financial.lcoe)}/kWh</span></p>
            <p><span className="text-fg-tertiary">CO2 Evitado:</span> <span className="font-mono tabular-nums text-fg-primary">{s.financial.co2AvoidedTonsYear.toFixed(2)} ton/ano</span></p>
          </div>
        </section>
      </Card>
    </div>
  );
}
