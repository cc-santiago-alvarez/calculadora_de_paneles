import { useMemo, useEffect } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { PANEL_FORMATS } from './panelFormats';

const COVERAGE_OPTIONS = [25, 50, 75, 100];

const SYSTEM_EFFICIENCY = (1 - 0.03) * (1 - 0.02) * (1 - 0.02) * 0.96;

export default function CoverageStep() {
  const { formData, setFormData, irradiationPreview } = useProjectStore();

  const avgHSP = irradiationPreview?.annualAvgHSP ?? 4.5;

  const annualConsumption = useMemo(
    () => formData.consumption.monthly.reduce((sum, v) => sum + v, 0),
    [formData.consumption.monthly]
  );
  const dailyConsumption = annualConsumption / 365;

  const coverageDecimal = formData.coveragePercentage / 100;
  const adjustedDaily = dailyConsumption * coverageDecimal;
  const requiredPowerKwp = adjustedDaily / avgHSP / SYSTEM_EFFICIENCY;
  const recommendedInverterKw = Math.ceil(requiredPowerKwp * 1.1);

  const isCustom = !COVERAGE_OPTIONS.includes(formData.coveragePercentage);

  // Save recommendedInverterKw to formData when it changes
  useEffect(() => {
    if (recommendedInverterKw !== formData.recommendedInverterKw) {
      setFormData({ recommendedInverterKw });
    }
  }, [recommendedInverterKw]);

  const handleFormatSelect = (formatKey: typeof PANEL_FORMATS[number]['key']) => {
    setFormData({
      panelFormat: formatKey,
      equipment: { panelId: '', inverterId: '' },
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Cobertura Energetica</h3>
        <p className="text-sm text-fg-secondary">
          Selecciona que porcentaje del consumo energetico deseas cubrir con el sistema solar.
        </p>
      </div>

      {/* Coverage options */}
      <div className="grid grid-cols-2 gap-4">
        {COVERAGE_OPTIONS.map((pct) => {
          const isSelected = formData.coveragePercentage === pct;
          return (
            <button
              key={pct}
              onClick={() => setFormData({ coveragePercentage: pct })}
              className={`p-4 rounded-lg border-2 text-left transition-colors duration-fast ease-decel ${
                isSelected
                  ? 'border-brand bg-brand-soft'
                  : 'border-[var(--color-border-default)] hover:border-[var(--color-border-strong)] bg-surface'
              }`}
            >
              <div className="flex items-center gap-3 mb-1">
                <div
                  className={`w-5 h-5 rounded-full border-2 flex items-center justify-center ${
                    isSelected ? 'border-brand' : 'border-[var(--color-border-default)]'
                  }`}
                >
                  {isSelected && (
                    <div className="w-3 h-3 rounded-full bg-brand" />
                  )}
                </div>
                <span className="text-2xl font-bold font-mono tabular-nums text-fg-primary">{pct}%</span>
              </div>
              <p className="text-sm text-fg-secondary ml-8">
                {pct === 100
                  ? 'Cubrir la totalidad del consumo'
                  : pct === 75
                  ? 'Cubrir tres cuartos del consumo'
                  : pct === 50
                  ? 'Cubrir la mitad del consumo'
                  : 'Cubrir un cuarto del consumo'}
              </p>
            </button>
          );
        })}
      </div>

      {/* Custom percentage */}
      <div className="flex items-center gap-3">
        <input
          type="checkbox"
          id="custom-coverage"
          checked={isCustom}
          onChange={(e) => {
            if (e.target.checked) {
              setFormData({ coveragePercentage: 60 });
            } else {
              setFormData({ coveragePercentage: 100 });
            }
          }}
          className="w-4 h-4 accent-[var(--color-brand)]"
        />
        <label htmlFor="custom-coverage" className="text-sm font-medium text-fg-secondary">
          Porcentaje personalizado
        </label>
        {isCustom && (
          <div className="flex items-center gap-2">
            <input
              type="range"
              min={10}
              max={100}
              step={5}
              value={formData.coveragePercentage}
              onChange={(e) => setFormData({ coveragePercentage: Number(e.target.value) })}
              className="w-40"
            />
            <span className="text-lg font-bold font-mono tabular-nums text-brand w-14 text-right">
              {formData.coveragePercentage}%
            </span>
          </div>
        )}
      </div>

      {/* Panel format selection */}
      <div>
        <h4 className="font-semibold text-fg-primary text-base mb-3">Formato de Panel Preferido</h4>
        <div className="grid grid-cols-3 gap-4">
          {PANEL_FORMATS.map((format) => {
            const isSelected = formData.panelFormat === format.key;
            const fmtPanels = Math.ceil((requiredPowerKwp * 1000) / format.watts);
            const fmtRoof = fmtPanels * format.area * 1.15;
            return (
              <button
                key={format.key}
                onClick={() => handleFormatSelect(format.key)}
                className={`p-4 rounded-lg border-2 text-left transition-colors duration-fast ease-decel ${
                  isSelected
                    ? 'border-brand bg-brand-soft'
                    : 'border-[var(--color-border-default)] hover:border-[var(--color-border-strong)] bg-surface'
                }`}
              >
                <div className="flex items-center gap-2 mb-2">
                  <div
                    className={`w-5 h-5 rounded-full border-2 flex items-center justify-center ${
                      isSelected ? 'border-brand' : 'border-[var(--color-border-default)]'
                    }`}
                  >
                    {isSelected && <div className="w-3 h-3 rounded-full bg-brand" />}
                  </div>
                  <span className="font-bold text-fg-primary">{format.label}</span>
                </div>
                <p className="text-xs text-fg-tertiary mb-3">{format.description}</p>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span className="text-fg-tertiary">Potencia:</span>
                    <span className="font-medium font-mono tabular-nums text-fg-primary">{format.watts}W</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-fg-tertiary">Paneles:</span>
                    <span className="font-medium font-mono tabular-nums text-fg-primary">{fmtPanels}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-fg-tertiary">Techo:</span>
                    <span className="font-medium font-mono tabular-nums text-fg-primary">~{fmtRoof.toFixed(0)} m2</span>
                  </div>
                </div>
              </button>
            );
          })}
        </div>
      </div>

    </div>
  );
}
