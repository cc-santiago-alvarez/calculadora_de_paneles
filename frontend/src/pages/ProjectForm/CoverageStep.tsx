import { useProjectStore } from '../../store/useProjectStore';

const COVERAGE_OPTIONS = [25, 50, 75, 100];

export default function CoverageStep() {
  const { formData, setFormData } = useProjectStore();

  const isCustom = !COVERAGE_OPTIONS.includes(formData.coveragePercentage);

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
    </div>
  );
}
