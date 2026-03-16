import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
} from 'recharts';
import { Scenario } from '../../types';
import { formatNumber } from '../../utils/format';
import ChartContainer from '../../components/common/ChartContainer';

export default function DegradationChart({ scenario }: { scenario: Scenario }) {
  const data = scenario.production.yearly25.map((prod, i) => ({
    year: i + 1,
    produccion: Math.round(prod),
  }));

  return (
    <ChartContainer title="Produccion Anual con Degradacion (25 anos)">
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid stroke="var(--color-border-subtle)" strokeDasharray="3 3" />
          <XAxis
            dataKey="year"
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
          />
          <YAxis
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
            tickFormatter={(v) => `${formatNumber(v)} kWh`}
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
            formatter={(value: number) => [`${formatNumber(value)} kWh`, 'Produccion']}
            labelFormatter={(label) => `Ano ${label}`}
          />
          <Line type="monotone" dataKey="produccion" stroke="var(--color-brand)" strokeWidth={2} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </ChartContainer>
  );
}
