import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
} from 'recharts';
import { MONTHS } from '../../utils/format';
import { Scenario } from '../../types';
import ChartContainer from '../../components/common/ChartContainer';

interface Props {
  scenario: Scenario;
  monthlyConsumption: number[];
}

export default function GenerationChart({ scenario, monthlyConsumption }: Props) {
  const data = MONTHS.map((month, i) => ({
    month,
    produccion: Math.round(scenario.production.monthlyKwh[i]),
    consumo: Math.round(monthlyConsumption[i]),
  }));

  return (
    <ChartContainer title="Generacion vs Consumo Mensual (kWh)">
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data}>
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
          <Bar dataKey="produccion" fill="var(--color-brand)" name="Produccion" />
          <Bar dataKey="consumo" fill="var(--color-info)" name="Consumo" />
        </BarChart>
      </ResponsiveContainer>
    </ChartContainer>
  );
}
