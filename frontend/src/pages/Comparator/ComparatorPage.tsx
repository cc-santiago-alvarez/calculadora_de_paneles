import { useProjectStore } from '../../store/useProjectStore';
import { formatCOP, formatNumber } from '../../utils/format';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';
import {
  RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, Radar, Legend, ResponsiveContainer,
} from 'recharts';

const COLORS = ['#f59e0b', '#3b82f6', '#22c55e', '#a855f7'];

export default function ComparatorPage() {
  const { comparisonScenarios, removeFromComparison, clearComparison } = useProjectStore();

  if (comparisonScenarios.length === 0) {
    return (
      <Card className="text-center py-20">
        <p className="text-fg-secondary text-lg mb-2">No hay escenarios para comparar.</p>
        <p className="text-fg-tertiary text-sm">
          Calcula un proyecto y agrega escenarios desde la pagina de resultados.
        </p>
      </Card>
    );
  }

  // Normalize metrics for radar chart
  const maxValues = {
    production: Math.max(...comparisonScenarios.map((s) => s.production.annualKwh)),
    irr: Math.max(...comparisonScenarios.map((s) => s.financial.irrPercent)),
    npv: Math.max(...comparisonScenarios.map((s) => s.financial.npvCOP)),
    co2: Math.max(...comparisonScenarios.map((s) => s.financial.co2AvoidedTonsYear)),
    efficiency: 100,
  };

  const radarData = [
    { metric: 'Produccion', ...Object.fromEntries(comparisonScenarios.map((s, i) => [`s${i}`, (s.production.annualKwh / maxValues.production) * 100])) },
    { metric: 'TIR', ...Object.fromEntries(comparisonScenarios.map((s, i) => [`s${i}`, (s.financial.irrPercent / Math.max(maxValues.irr, 1)) * 100])) },
    { metric: 'VPN', ...Object.fromEntries(comparisonScenarios.map((s, i) => [`s${i}`, (s.financial.npvCOP / Math.max(maxValues.npv, 1)) * 100])) },
    { metric: 'CO2 Evitado', ...Object.fromEntries(comparisonScenarios.map((s, i) => [`s${i}`, (s.financial.co2AvoidedTonsYear / Math.max(maxValues.co2, 1)) * 100])) },
    { metric: 'Uso Techo', ...Object.fromEntries(comparisonScenarios.map((s, i) => [`s${i}`, s.systemDesign.roofUtilization])) },
  ];

  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-fg-primary">Comparador de Escenarios</h2>
        <Button variant="danger" size="sm" onClick={clearComparison}>
          Limpiar Comparacion
        </Button>
      </div>

      {/* Radar Chart */}
      {comparisonScenarios.length >= 2 && (
        <Card padding="none">
          <h4 className="font-semibold text-fg-primary p-4 pb-0">Comparacion Visual</h4>
          <div className="p-4">
            <ResponsiveContainer width="100%" height={400}>
              <RadarChart data={radarData}>
                <PolarGrid />
                <PolarAngleAxis dataKey="metric" tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }} />
                <PolarRadiusAxis tick={false} />
                {comparisonScenarios.map((s, i) => (
                  <Radar
                    key={s._id}
                    name={s.name}
                    dataKey={`s${i}`}
                    stroke={COLORS[i % COLORS.length]}
                    fill={COLORS[i % COLORS.length]}
                    fillOpacity={0.15}
                  />
                ))}
                <Legend />
              </RadarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      )}

      {/* Comparison Table */}
      <Card padding="none" className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-inset">
            <tr>
              <th className="px-4 py-3 text-left text-fg-secondary font-medium">Metrica</th>
              {comparisonScenarios.map((s, i) => (
                <th key={s._id} className="px-4 py-3 text-center text-fg-secondary font-medium" style={{ color: COLORS[i % COLORS.length] }}>
                  <div className="flex items-center justify-center gap-2">
                    {s.name}
                    <button
                      onClick={() => removeFromComparison(s._id)}
                      className="text-fg-muted hover:text-danger text-xs"
                    >
                      x
                    </button>
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            <CompRow label="Potencia Instalada" values={comparisonScenarios.map((s) => `${s.systemDesign.actualPowerKwp.toFixed(2)} kWp`)} />
            <CompRow label="Numero de Paneles" values={comparisonScenarios.map((s) => `${s.systemDesign.numberOfPanels}`)} />
            <CompRow label="Produccion Anual" values={comparisonScenarios.map((s) => `${formatNumber(s.production.annualKwh)} kWh`)} />
            <CompRow label="HSP Promedio" values={comparisonScenarios.map((s) => `${s.irradiation.annualAvgHSP.toFixed(2)}`)} />
            <CompRow label="Costo Instalacion" values={comparisonScenarios.map((s) => formatCOP(s.financial.installationCostCOP))} />
            <CompRow label="Ahorro Anual" values={comparisonScenarios.map((s) => formatCOP(s.financial.annualSavingsCOP))} />
            <CompRow label="Payback" values={comparisonScenarios.map((s) => `${s.financial.paybackYears != null ? s.financial.paybackYears.toFixed(1) : 'N/A'} anos`)} />
            <CompRow label="TIR" values={comparisonScenarios.map((s) => `${s.financial.irrPercent.toFixed(1)}%`)} />
            <CompRow label="VPN" values={comparisonScenarios.map((s) => formatCOP(s.financial.npvCOP))} />
            <CompRow label="LCOE" values={comparisonScenarios.map((s) => `${formatCOP(s.financial.lcoe)}/kWh`)} />
            <CompRow label="CO2 Evitado" values={comparisonScenarios.map((s) => `${s.financial.co2AvoidedTonsYear.toFixed(2)} ton/ano`)} />
            <CompRow label="Perdida Total" values={comparisonScenarios.map((s) => `${s.losses.totalSystemLoss.toFixed(1)}%`)} />
            <CompRow label="Uso Techo" values={comparisonScenarios.map((s) => `${s.systemDesign.roofUtilization.toFixed(1)}%`)} />
          </tbody>
        </table>
      </Card>
    </div>
  );
}

function CompRow({ label, values }: { label: string; values: string[] }) {
  return (
    <tr className="border-t border-[var(--color-border-subtle)] hover:bg-inset">
      <td className="px-4 py-2 text-fg-secondary font-medium">{label}</td>
      {values.map((v, i) => (
        <td key={i} className="px-4 py-2 text-center font-mono tabular-nums text-fg-primary">{v}</td>
      ))}
    </tr>
  );
}
