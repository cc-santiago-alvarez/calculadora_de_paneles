import { useProjectStore } from '../../store/useProjectStore';

export default function RoofStep() {
  const { formData, setFormData } = useProjectStore();

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Caracteristicas del Techo</h3>
        <p className="text-sm text-fg-secondary">
          Define el area disponible y orientacion de tu techo.
        </p>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">
            Area Total del Techo (m2)
          </label>
          <input
            type="number"
            value={formData.roof.area}
            onChange={(e) =>
              setFormData({ roof: { ...formData.roof, area: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
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
            value={formData.roof.usablePercentage}
            onChange={(e) =>
              setFormData({ roof: { ...formData.roof, usablePercentage: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
          />
          <p className="text-xs text-fg-muted mt-1">Area efectiva: <span className="font-mono tabular-nums">{((formData.roof.area * formData.roof.usablePercentage) / 100).toFixed(1)}</span> m2</p>
        </div>

        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-1">
            Inclinacion (grados)
          </label>
          <input
            type="number"
            min={0}
            max={90}
            value={formData.roof.tilt}
            onChange={(e) =>
              setFormData({ roof: { ...formData.roof, tilt: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
          />
          <p className="text-xs text-fg-muted mt-1">
            0° = horizontal, 90° = vertical. Recomendado: ~10° para Colombia
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
            value={formData.roof.azimuth}
            onChange={(e) =>
              setFormData({ roof: { ...formData.roof, azimuth: parseFloat(e.target.value) || 0 } })
            }
            className="w-full bg-inset border border-[var(--color-border-default)] rounded-md py-2 px-3 text-sm text-fg-primary focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2"
          />
          <p className="text-xs text-fg-muted mt-1">
            0° = Sur, 90° = Oeste, -90° = Este, 180° = Norte
          </p>
        </div>
      </div>

      {/* Shading */}
      <div className="border-t border-[var(--color-border-subtle)] pt-4">
        <div className="flex items-center gap-3 mb-4">
          <input
            type="checkbox"
            id="hasShading"
            checked={formData.roof.shadingProfile.hasShading}
            onChange={(e) =>
              setFormData({
                roof: {
                  ...formData.roof,
                  shadingProfile: {
                    ...formData.roof.shadingProfile,
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

        {formData.roof.shadingProfile.hasShading && (
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
                      value={(formData.roof.shadingProfile.monthlyLoss[i] * 100).toFixed(0)}
                      onChange={(e) => {
                        const newLoss = [...formData.roof.shadingProfile.monthlyLoss];
                        newLoss[i] = (parseFloat(e.target.value) || 0) / 100;
                        setFormData({
                          roof: {
                            ...formData.roof,
                            shadingProfile: { ...formData.roof.shadingProfile, monthlyLoss: newLoss },
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
