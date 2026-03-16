import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ReferenceLine,
} from 'recharts';
import { Scenario } from '../../types';
import { formatCOP } from '../../utils/format';
import ChartContainer from '../../components/common/ChartContainer';

export default function CashFlowChart({ scenario }: { scenario: Scenario }) {
  const data = scenario.financial.cumulativeSavings25.map((cumulative, i) => ({
    year: i + 1,
    ahorroAcumulado: Math.round(cumulative),
    inversionInicial: Math.round(scenario.financial.installationCostCOP),
  }));

  return (
    <ChartContainer title="Flujo de Caja Acumulado (25 anos)">
      <ResponsiveContainer width="100%" height={300}>
        <AreaChart data={data}>
          <CartesianGrid stroke="var(--color-border-subtle)" strokeDasharray="3 3" />
          <XAxis
            dataKey="year"
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
            label={{ value: 'Ano', position: 'bottom' }}
          />
          <YAxis
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
            tickFormatter={(v) => `${(v / 1000000).toFixed(0)}M`}
          />
          <Tooltip
            contentStyle={{
              borderRadius: 8,
              border: '1px solid var(--color-border-default)',
              backgroundColor: 'var(--color-bg-surface-raised)',
              color: 'var(--color-fg-primary)',
              boxShadow: 'var(--shadow-md)',
              fontSize: 12,
            }}
            formatter={(value: number) => formatCOP(value)}
            labelFormatter={(label) => `Ano ${label}`}
          />
          <ReferenceLine
            y={scenario.financial.installationCostCOP}
            stroke="var(--color-danger)"
            strokeDasharray="5 5"
            label="Inversion"
          />
          <Area
            type="monotone"
            dataKey="ahorroAcumulado"
            stroke="var(--color-success)"
            fill="var(--color-success-soft)"
            name="Ahorro Acumulado"
          />
        </AreaChart>
      </ResponsiveContainer>
    </ChartContainer>
  );
}
