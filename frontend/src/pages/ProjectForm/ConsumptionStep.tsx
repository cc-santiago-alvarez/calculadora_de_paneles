import { useState } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { MONTHS_FULL } from '../../utils/format';
import Select from '../../components/common/Select';
import Button from '../../components/common/Button';

const ESTRATOS = [
  { value: 1, label: 'Estrato 1 - Bajo-Bajo', tariff: 450 },
  { value: 2, label: 'Estrato 2 - Bajo', tariff: 540 },
  { value: 3, label: 'Estrato 3 - Medio-Bajo', tariff: 680 },
  { value: 4, label: 'Estrato 4 - Medio', tariff: 800 },
  { value: 5, label: 'Estrato 5 - Medio-Alto', tariff: 960 },
  { value: 6, label: 'Estrato 6 - Alto', tariff: 960 },
];

const ESTRATO_OPTIONS = ESTRATOS.map((e) => ({
  value: String(e.value),
  label: e.label,
}));

const CONNECTION_OPTIONS = [
  { value: 'monofasica', label: 'Monofasica (120V)' },
  { value: 'bifasica', label: 'Bifasica (240V)' },
  { value: 'trifasica', label: 'Trifasica (208/480V)' },
];

export default function ConsumptionStep() {
  const { formData, setFormData } = useProjectStore();
  const [quickFillValue, setQuickFillValue] = useState('');

  const updateMonthly = (index: number, value: number) => {
    const newMonthly = [...formData.consumption.monthly];
    newMonthly[index] = value;
    setFormData({
      consumption: { ...formData.consumption, monthly: newMonthly },
    });
  };

  const setUniformConsumption = () => {
    const val = parseFloat(quickFillValue);
    if (!isNaN(val) && val >= 0) {
      setFormData({
        consumption: { ...formData.consumption, monthly: new Array(12).fill(val) },
      });
    }
  };

  const handleEstratoChange = (estrato: number) => {
    const estratoData = ESTRATOS.find((e) => e.value === estrato);
    setFormData({
      consumption: {
        ...formData.consumption,
        estrato,
        tariffPerKwh: estratoData?.tariff || 800,
      },
    });
  };

  const totalAnnual = formData.consumption.monthly.reduce((a, b) => a + b, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h3 className="text-xl font-bold text-fg-primary mb-1">Consumo Electrico Mensual</h3>
        <p className="text-sm text-fg-muted">
          Este es el dato mas importante para dimensionar tu sistema solar.
          Ingresa el consumo de cada mes en kWh — lo encuentras en tu factura de energia
          en la seccion "Consumo del periodo".
        </p>
      </div>

      {/* Project name */}
      <div>
        <label className="block text-sm font-medium text-fg-secondary mb-1">Nombre del Proyecto</label>
        <input
          type="text"
          value={formData.name}
          onChange={(e) => setFormData({ name: e.target.value })}
          className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
          placeholder="Ej: Casa de Juan - Bogota"
        />
      </div>

      {/* Estrato, tariff, connection */}
      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Estrato</label>
          <Select
            value={String(formData.consumption.estrato)}
            onChange={(val) => handleEstratoChange(parseInt(val))}
            options={ESTRATO_OPTIONS}
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Tarifa (COP/kWh)</label>
          <input
            type="number"
            value={formData.consumption.tariffPerKwh}
            onChange={(e) =>
              setFormData({ consumption: { ...formData.consumption, tariffPerKwh: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
          />
          <p className="text-xs text-fg-muted mt-1">Se ajusta automaticamente segun el estrato</p>
        </div>
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">Tipo de Conexion</label>
          <Select
            value={formData.consumption.connectionType}
            onChange={(val) =>
              setFormData({
                consumption: { ...formData.consumption, connectionType: val as any },
              })
            }
            options={CONNECTION_OPTIONS}
          />
        </div>
      </div>

      {/* Quick fill */}
      <div className="bg-inset p-4 rounded-lg">
        <p className="text-sm font-medium text-fg-secondary mb-2">Llenado rapido</p>
        <div className="flex items-center gap-3">
          <input
            type="number"
            value={quickFillValue}
            onChange={(e) => setQuickFillValue(e.target.value)}
            className="w-32 bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
            placeholder="kWh"
            min={0}
          />
          <Button
            variant="primary"
            size="sm"
            onClick={setUniformConsumption}
            disabled={!quickFillValue}
          >
            Aplicar a todos los meses
          </Button>
          <span className="text-xs text-fg-muted">
            Usa esto si tu consumo es similar cada mes
          </span>
        </div>
      </div>

      {/* Monthly consumption table - the main input */}
      <div>
        <p className="text-sm font-medium text-fg-secondary mb-3">
          Consumo mensual (kWh) — ingresa el valor de cada mes
        </p>
        <div className="grid grid-cols-1 gap-2">
          {MONTHS_FULL.map((month, i) => (
            <div
              key={month}
              className={`flex items-center gap-4 px-4 py-2.5 rounded-lg border transition-colors ${
                formData.consumption.monthly[i] > 0
                  ? 'border-[var(--color-success)] border-opacity-30 bg-success-soft dark:border-brand dark:border-opacity-30 dark:bg-brand-soft'
                  : 'border-[var(--color-border-default)] bg-surface'
              }`}
            >
              <span className="text-sm font-medium text-fg-secondary w-28">{month}</span>
              <input
                type="number"
                value={formData.consumption.monthly[i] || ''}
                onChange={(e) => updateMonthly(i, parseFloat(e.target.value) || 0)}
                className="flex-1 max-w-[200px] bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 font-mono tabular-nums"
                placeholder="0"
                min={0}
              />
              <span className="text-sm text-fg-muted">kWh</span>
              {/* Visual bar */}
              <div className="flex-1 bg-[var(--color-border-default)] rounded-full h-2 max-w-[200px]">
                <div
                  className="bg-brand dark:bg-info h-2 rounded-full transition-all"
                  style={{
                    width: `${totalAnnual > 0 ? (formData.consumption.monthly[i] / Math.max(...formData.consumption.monthly)) * 100 : 0}%`,
                  }}
                />
              </div>
              {formData.consumption.monthly[i] > 0 && (
                <span className="text-xs text-fg-muted w-24 text-right font-mono tabular-nums">
                  {new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(
                    formData.consumption.monthly[i] * formData.consumption.tariffPerKwh
                  )}
                  /mes
                </span>
              )}
            </div>
          ))}
        </div>
      </div>

    </div>
  );
}
