import { useProjectStore } from '../../store/useProjectStore';
import { RoofType, Slope } from '../../types';

const ROOF_TYPES: { key: RoofType; label: string; description: string; slopeCount: number }[] = [
  { key: 'plana', label: 'Plana', description: 'Superficie horizontal sin pendiente', slopeCount: 0 },
  { key: 'una_agua', label: 'Una agua', description: 'Una sola pendiente', slopeCount: 1 },
  { key: 'dos_aguas', label: 'Dos aguas', description: 'Dos pendientes opuestas', slopeCount: 2 },
  { key: 'cuatro_aguas', label: 'Cuatro aguas', description: 'Cuatro pendientes', slopeCount: 4 },
];

const DEFAULT_AZIMUTHS: Record<string, number[]> = {
  una_agua: [0],
  dos_aguas: [0, 180],
  cuatro_aguas: [0, 90, 180, -90],
};

function generateSlopes(roofType: RoofType, currentArea: number): Slope[] {
  const azimuths = DEFAULT_AZIMUTHS[roofType];
  if (!azimuths) return [];
  const areaPerSlope = Math.round((currentArea / azimuths.length) * 10) / 10;
  return azimuths.map((azimuth) => ({
    area: areaPerSlope,
    tilt: 15,
    azimuth,
  }));
}

const inputClass =
  'w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2';

export default function RoofStep() {
  const { formData, setFormData } = useProjectStore();
  const roof = formData.roof;
  const currentType = ROOF_TYPES.find((t) => t.key === roof.roofType) || ROOF_TYPES[0];
  const isFlat = roof.roofType === 'plana';

  const handleRoofTypeChange = (roofType: RoofType) => {
    const slopes = generateSlopes(roofType, roof.area);
    const totalArea = slopes.length > 0 ? slopes.reduce((sum, s) => sum + s.area, 0) : roof.area;
    setFormData({
      roof: {
        ...roof,
        roofType,
        slopes,
        area: totalArea,
        tilt: roofType === 'plana' ? 0 : roof.tilt,
      },
    });
  };

  const handleSlopeChange = (index: number, field: keyof Slope, value: number) => {
    const newSlopes = [...roof.slopes];
    newSlopes[index] = { ...newSlopes[index], [field]: value };
    const totalArea = newSlopes.reduce((sum, s) => sum + s.area, 0);
    setFormData({
      roof: { ...roof, slopes: newSlopes, area: totalArea },
    });
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Caracteristicas del Techo</h3>
        <p className="text-sm text-fg-secondary">
          Selecciona el tipo de techo y define las caracteristicas de cada caida.
        </p>
      </div>

      {/* Roof type selector */}
      <div>
        <label className="block text-sm font-medium text-fg-secondary mb-3">Tipo de Techo</label>
        <div className="grid grid-cols-4 gap-3">
          {ROOF_TYPES.map((type) => (
            <button
              key={type.key}
              type="button"
              onClick={() => handleRoofTypeChange(type.key)}
              className={`p-3 rounded-lg border-2 text-left transition-colors ${
                roof.roofType === type.key
                  ? 'border-[var(--color-brand)] bg-[var(--color-brand)]/10'
                  : 'border-[var(--color-border-default)] hover:border-[var(--color-border-emphasis)]'
              }`}
            >
              <div className="text-sm font-medium text-fg-primary">{type.label}</div>
              <div className="text-xs text-fg-muted mt-1">{type.description}</div>
            </button>
          ))}
        </div>
      </div>

      {/* Flat roof: single surface */}
      {isFlat && (
        <div className="grid grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-fg-secondary mb-1">
              Area Total del Techo (m2)
            </label>
            <input
              type="number"
              value={roof.area}
              onChange={(e) =>
                setFormData({ roof: { ...roof, area: parseFloat(e.target.value) || 0 } })
              }
              className={inputClass}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-fg-secondary mb-1">
              Porcentaje Utilizable (%)
            </label>
            <input
              type="number"
              min={0}
              max={100}
              value={roof.usablePercentage}
              onChange={(e) =>
                setFormData({
                  roof: { ...roof, usablePercentage: parseFloat(e.target.value) || 0 },
                })
              }
              className={inputClass}
            />
            <p className="text-xs text-fg-muted mt-1">
              Area efectiva:{' '}
              <span className="font-mono tabular-nums">
                {((roof.area * roof.usablePercentage) / 100).toFixed(1)}
              </span>{' '}
              m2
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-fg-secondary mb-1">
              Azimut (grados)
            </label>
            <input
              type="number"
              min={-180}
              max={180}
              value={roof.azimuth}
              onChange={(e) =>
                setFormData({ roof: { ...roof, azimuth: parseFloat(e.target.value) || 0 } })
              }
              className={inputClass}
            />
            <p className="text-xs text-fg-muted mt-1">
              0° = Sur, 90° = Oeste, -90° = Este, 180° = Norte
            </p>
          </div>
        </div>
      )}

      {/* Sloped roof: per-slope configuration */}
      {!isFlat && (
        <div className="space-y-4">
          <div className="grid grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-fg-secondary mb-1">
                Porcentaje Utilizable (%)
              </label>
              <input
                type="number"
                min={0}
                max={100}
                value={roof.usablePercentage}
                onChange={(e) =>
                  setFormData({
                    roof: { ...roof, usablePercentage: parseFloat(e.target.value) || 0 },
                  })
                }
                className={inputClass}
              />
            </div>
            <div className="flex items-end">
              <p className="text-sm text-fg-secondary">
                Area total:{' '}
                <span className="font-mono tabular-nums font-semibold text-fg-primary">
                  {roof.area.toFixed(1)}
                </span>{' '}
                m2 — Area efectiva:{' '}
                <span className="font-mono tabular-nums font-semibold text-fg-primary">
                  {((roof.area * roof.usablePercentage) / 100).toFixed(1)}
                </span>{' '}
                m2
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            {roof.slopes.map((slope, i) => (
              <div
                key={i}
                className="bg-inset border border-[var(--color-border-subtle)] rounded-lg p-4"
              >
                <h4 className="text-sm font-semibold text-fg-primary mb-3">Caida {i + 1}</h4>
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <label className="block text-xs font-medium text-fg-secondary mb-1">
                      Area (m2)
                    </label>
                    <input
                      type="number"
                      min={0}
                      value={slope.area}
                      onChange={(e) =>
                        handleSlopeChange(i, 'area', parseFloat(e.target.value) || 0)
                      }
                      className={inputClass}
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-fg-secondary mb-1">
                      Inclinacion (grados)
                    </label>
                    <input
                      type="number"
                      min={0}
                      max={90}
                      value={slope.tilt}
                      onChange={(e) =>
                        handleSlopeChange(i, 'tilt', parseFloat(e.target.value) || 0)
                      }
                      className={inputClass}
                    />
                    <p className="text-xs text-fg-muted mt-1">0° = horizontal, 90° = vertical</p>
                  </div>
                  <div>
                    <label className="block text-xs font-medium text-fg-secondary mb-1">
                      Azimut (grados)
                    </label>
                    <input
                      type="number"
                      min={-180}
                      max={180}
                      value={slope.azimuth}
                      onChange={(e) =>
                        handleSlopeChange(i, 'azimuth', parseFloat(e.target.value) || 0)
                      }
                      className={inputClass}
                    />
                    <p className="text-xs text-fg-muted mt-1">0°=Sur, 90°=O, -90°=E</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Shading */}
      <div className="border-t border-[var(--color-border-subtle)] pt-4">
        <div className="flex items-center gap-3 mb-4">
          <input
            type="checkbox"
            id="hasShading"
            checked={roof.shadingProfile.hasShading}
            onChange={(e) =>
              setFormData({
                roof: {
                  ...roof,
                  shadingProfile: {
                    ...roof.shadingProfile,
                    hasShading: e.target.checked,
                  },
                },
              })
            }
            className="w-4 h-4 rounded accent-[var(--color-brand)]"
          />
          <label htmlFor="hasShading" className="text-sm font-medium text-fg-secondary">
            El techo tiene sombras parciales
          </label>
        </div>

        {roof.shadingProfile.hasShading && (
          <div className="bg-inset p-4 rounded-lg">
            <p className="text-sm text-fg-secondary mb-3">
              Porcentaje estimado de perdida por sombra cada mes (0-100%):
            </p>
            <div className="grid grid-cols-6 gap-2">
              {['Ene', 'Feb', 'Mar', 'Abr', 'May', 'Jun', 'Jul', 'Ago', 'Sep', 'Oct', 'Nov', 'Dic'].map(
                (month, i) => (
                  <div key={month} className="text-center">
                    <label className="text-xs text-fg-tertiary">{month}</label>
                    <input
                      type="number"
                      min={0}
                      max={100}
                      value={(roof.shadingProfile.monthlyLoss[i] * 100).toFixed(0)}
                      onChange={(e) => {
                        const newLoss = [...roof.shadingProfile.monthlyLoss];
                        newLoss[i] = (parseFloat(e.target.value) || 0) / 100;
                        setFormData({
                          roof: {
                            ...roof,
                            shadingProfile: { ...roof.shadingProfile, monthlyLoss: newLoss },
                          },
                        });
                      }}
                      className="w-full bg-surface border border-[var(--color-border-default)] rounded-md text-center text-sm font-mono py-1 px-1"
                    />
                  </div>
                )
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
