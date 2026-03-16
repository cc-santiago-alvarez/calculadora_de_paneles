import { useMemo } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { formatCOP, formatNumber } from '../../utils/format';
import { PANEL_FORMATS } from './panelFormats';

const SYSTEM_EFFICIENCY = (1 - 0.03) * (1 - 0.02) * (1 - 0.02) * 0.96;

function SidebarCard({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-surface border border-[var(--color-border-default)] rounded-lg p-4 animate-fade-in-up">
      <h4 className="text-xs font-semibold text-fg-tertiary uppercase tracking-wider mb-3">{title}</h4>
      {children}
    </div>
  );
}

function SidebarItem({ label, value, unit }: { label: string; value: string; unit?: string }) {
  return (
    <div className="py-2 border-b border-[var(--color-border-subtle)] last:border-0">
      <span className="text-xs text-fg-tertiary block">{label}</span>
      <div className="flex items-baseline gap-1 mt-0.5">
        <span className="text-lg font-bold font-mono tabular-nums text-fg-primary">{value}</span>
        {unit && <span className="text-xs text-fg-muted">{unit}</span>}
      </div>
    </div>
  );
}

function ConsumptionSidebar() {
  const { formData } = useProjectStore();
  const totalAnnual = formData.consumption.monthly.reduce((a, b) => a + b, 0);
  const avgMonthly = totalAnnual / 12;
  const dailyAvg = totalAnnual / 365;
  const allFilled = formData.consumption.monthly.every((v) => v > 0);

  return (
    <SidebarCard title="Resumen de Consumo">
      <SidebarItem label="Consumo Anual Total" value={formatNumber(totalAnnual)} unit="kWh/año" />
      <SidebarItem label="Promedio Mensual" value={formatNumber(avgMonthly)} unit="kWh/mes" />
      <SidebarItem label="Consumo Diario" value={dailyAvg.toFixed(1)} unit="kWh/día" />
      <SidebarItem
        label="Costo Anual Estimado"
        value={formatCOP(totalAnnual * formData.consumption.tariffPerKwh)}
        unit="COP/año"
      />
      {!allFilled && (
        <p className="text-xs text-danger mt-2">
          Completa los 12 meses para mayor precision.
        </p>
      )}
    </SidebarCard>
  );
}

function LocationSidebar() {
  const { irradiationPreview, formData } = useProjectStore();

  return (
    <SidebarCard title="Datos de Ubicacion">
      {irradiationPreview ? (
        <>
          <SidebarItem label="HSP Promedio" value={irradiationPreview.annualAvgHSP.toFixed(1)} unit="kWh/m²/día" />
          <SidebarItem label="Fuente" value={irradiationPreview.source} />
          <SidebarItem label="Altitud" value={String(Math.round(irradiationPreview.elevation))} unit="m" />
        </>
      ) : (
        <>
          <SidebarItem label="Altitud" value={String(formData.location.altitude)} unit="m" />
          <p className="text-xs text-fg-muted mt-2">
            Haz clic en el mapa para obtener datos de irradiacion.
          </p>
        </>
      )}
    </SidebarCard>
  );
}

function CoverageSidebar() {
  const { formData, irradiationPreview } = useProjectStore();
  const avgHSP = irradiationPreview?.annualAvgHSP ?? 4.5;

  const annualConsumption = useMemo(
    () => formData.consumption.monthly.reduce((sum, v) => sum + v, 0),
    [formData.consumption.monthly]
  );

  const coverageDecimal = formData.coveragePercentage / 100;
  const dailyConsumption = annualConsumption / 365;
  const adjustedDaily = dailyConsumption * coverageDecimal;
  const requiredPowerKwp = adjustedDaily / avgHSP / SYSTEM_EFFICIENCY;
  const recommendedInverterKw = Math.ceil(requiredPowerKwp * 1.1);

  const selectedFormat = PANEL_FORMATS.find((f) => f.key === formData.panelFormat) ?? PANEL_FORMATS[0];
  const numPanels = Math.ceil((requiredPowerKwp * 1000) / selectedFormat.watts);
  const roofNeeded = numPanels * selectedFormat.area * 1.15;

  return (
    <SidebarCard title="Dimensionamiento">
      <SidebarItem label={`Consumo a cubrir (${formData.coveragePercentage}%)`} value={formatNumber(annualConsumption * coverageDecimal)} unit="kWh/año" />
      <SidebarItem label="HSP Promedio" value={avgHSP.toFixed(2)} unit="h/día" />
      <SidebarItem label="Potencia Requerida" value={requiredPowerKwp.toFixed(2)} unit="kWp" />
      <div className="my-2 border-t border-[var(--color-border-default)]" />
      <SidebarItem label={`Paneles (${selectedFormat.label} ${selectedFormat.watts}W)`} value={String(numPanels)} unit="unidades" />
      <SidebarItem label="Area de Techo" value={`~${roofNeeded.toFixed(0)}`} unit="m²" />
      <SidebarItem label="Inversor Recomendado" value={`≥${recommendedInverterKw}`} unit="kW" />
    </SidebarCard>
  );
}

function EquipmentSidebar() {
  const { formData, panels, inverters, irradiationPreview } = useProjectStore();
  const avgHSP = irradiationPreview?.annualAvgHSP ?? 4.5;

  const annualConsumption = formData.consumption.monthly.reduce((sum, v) => sum + v, 0);
  const coverageDecimal = formData.coveragePercentage / 100;
  const dailyConsumption = annualConsumption / 365;
  const adjustedDaily = dailyConsumption * coverageDecimal;
  const requiredPowerKwp = adjustedDaily / avgHSP / SYSTEM_EFFICIENCY;

  const selectedPanel = panels.find((p) => p._id === formData.equipment.panelId);
  const selectedInverter = inverters.find((i) => i._id === formData.equipment.inverterId);

  const selectedFormat = PANEL_FORMATS.find((f) => f.key === formData.panelFormat) ?? PANEL_FORMATS[0];
  const panelWatts = selectedPanel?.powerWp ?? selectedFormat.watts;
  const numPanels = Math.ceil((requiredPowerKwp * 1000) / panelWatts);

  const panelTotal = selectedPanel ? numPanels * selectedPanel.costCOP : 0;
  const inverterTotal = selectedInverter ? selectedInverter.costCOP : 0;
  const grandTotal = panelTotal + inverterTotal;

  return (
    <SidebarCard title="Resumen de Equipo">
      {selectedPanel ? (
        <>
          <SidebarItem label={`${numPanels}× ${selectedPanel.manufacturer} ${selectedPanel.powerWp}W`} value={formatCOP(panelTotal)} unit="" />
        </>
      ) : (
        <p className="text-xs text-fg-muted py-2">Selecciona un panel solar</p>
      )}
      {selectedInverter ? (
        <>
          <SidebarItem label={`1× ${selectedInverter.manufacturer} ${selectedInverter.ratedPowerKw}kW`} value={formatCOP(inverterTotal)} unit="" />
        </>
      ) : (
        <p className="text-xs text-fg-muted py-2">Selecciona un inversor</p>
      )}
      {(selectedPanel || selectedInverter) && (
        <div className="mt-2 pt-2 border-t border-[var(--color-border-strong)]">
          <div className="flex justify-between items-baseline">
            <span className="text-sm font-semibold text-fg-primary">Total Estimado</span>
            <span className="text-xl font-bold font-mono tabular-nums text-brand">{formatCOP(grandTotal)}</span>
          </div>
        </div>
      )}
    </SidebarCard>
  );
}

const SIDEBAR_MAP: Record<number, React.FC> = {
  0: ConsumptionSidebar,
  1: LocationSidebar,
  4: CoverageSidebar,
  5: EquipmentSidebar,
};

export default function StepSidebar({ step }: { step: number }) {
  const SidebarContent = SIDEBAR_MAP[step];
  if (!SidebarContent) return null;

  return (
    <aside className="w-72 shrink-0">
      <div className="sticky top-8 space-y-4">
        <SidebarContent />
      </div>
    </aside>
  );
}
