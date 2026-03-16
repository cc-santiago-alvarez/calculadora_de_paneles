import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
} from 'recharts';
import { MONTHS } from '../../utils/format';
import { Scenario } from '../../types';
import ChartContainer from '../../components/common/ChartContainer';

export default function IrradiationChart({ scenario }: { scenario: Scenario }) {
  const data = MONTHS.map((month, i) => ({
    month,
    GHI: parseFloat(scenario.irradiation.monthlyGHI[i]?.toFixed(2)),
    POA: parseFloat(scenario.irradiation.monthlyPOA[i]?.toFixed(2)),
  }));

  return (
    <ChartContainer title="Irradiacion Mensual (kWh/m²/dia)">
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid stroke="var(--color-border-subtle)" strokeDasharray="3 3" />
          <XAxis
            dataKey="month"
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
          />
          <YAxis
            tick={{ fontSize: 11, fill: 'var(--color-fg-tertiary)' }}
            style={{ fontFamily: 'var(--font-mono)' }}
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
          />
          <Legend />
          <Line type="monotone" dataKey="GHI" stroke="var(--color-danger)" name="GHI Horizontal" strokeWidth={2} />
          <Line type="monotone" dataKey="POA" stroke="var(--color-brand)" name="POA Inclinado" strokeWidth={2} />
        </LineChart>
      </ResponsiveContainer>
    </ChartContainer>
  );
}
