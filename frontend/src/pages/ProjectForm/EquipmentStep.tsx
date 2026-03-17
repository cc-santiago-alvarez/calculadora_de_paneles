import { useEffect, useCallback, useState } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { catalogApi } from '../../api/catalog';
import { formatCOP } from '../../utils/format';
import { PANEL_FORMATS } from './panelFormats';
import Select from '../../components/common/Select';

export default function EquipmentStep() {
  const { formData, setFormData, panels, inverters, setPanels, setInverters } = useProjectStore();
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const selectedFormat = PANEL_FORMATS.find((f) => f.key === formData.panelFormat);

  const fetchCatalogs = useCallback(() => {
    setLoading(true);
    setFetchError(null);

    const panelParams: Record<string, string> = {};
    if (selectedFormat) {
      panelParams.minPower = String(selectedFormat.minPower);
      panelParams.maxPower = String(selectedFormat.maxPower);
    }
    const panelPromise = catalogApi.getPanels(panelParams).then(setPanels);

    const inverterParams: Record<string, string> = {};
    if (formData.recommendedInverterKw > 0) {
      inverterParams.minPower = String(Math.floor(formData.recommendedInverterKw * 0.8));
    }
    const systemToInverterType: Record<string, string> = {
      'on-grid': 'string',
      'off-grid': 'off-grid',
      'hybrid': 'hybrid',
    };
    if (formData.systemType && systemToInverterType[formData.systemType]) {
      inverterParams.type = systemToInverterType[formData.systemType];
    }
    const inverterPromise = catalogApi.getInverters(inverterParams).then(setInverters);

    Promise.all([panelPromise, inverterPromise])
      .catch(() => {
        setFetchError('Error al cargar el catalogo de equipos. Verifica tu conexion e intenta de nuevo.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [formData.panelFormat, formData.recommendedInverterKw, formData.systemType]);

  // Always fetch on mount and when filter criteria change
  useEffect(() => {
    fetchCatalogs();
  }, [fetchCatalogs]);

  const selectedPanel = panels.find((p) => p._id === formData.equipment.panelId);
  const selectedInverter = inverters.find((i) => i._id === formData.equipment.inverterId);

  const panelOptions = panels.map((panel) => ({
    value: panel._id,
    label: `${panel.manufacturer} ${panel.model} - ${panel.powerWp}W - ${formatCOP(panel.costCOP)}`,
  }));

  const inverterOptions = inverters.map((inv) => ({
    value: inv._id,
    label: `${inv.manufacturer} ${inv.model} - ${inv.ratedPowerKw}kW ${inv.type} - ${formatCOP(inv.costCOP)}`,
  }));

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Seleccion de Equipos</h3>
        <p className="text-sm text-fg-secondary">
          Selecciona el panel solar y el inversor del catalogo.
        </p>
        {selectedFormat && (
          <div className="mt-2 flex gap-4 text-sm">
            <span className="bg-brand-soft text-brand px-3 py-1 rounded-full">
              Formato: {selectedFormat.label} ({selectedFormat.minPower}-{selectedFormat.maxPower}W)
            </span>
            {formData.recommendedInverterKw > 0 && (
              <span className="bg-info-soft text-info px-3 py-1 rounded-full">
                Inversor recomendado: &ge;{formData.recommendedInverterKw} kW
              </span>
            )}
          </div>
        )}
      </div>

      {loading && (
        <div className="p-4 bg-info-soft border border-[var(--color-border-default)] rounded-lg text-sm text-info flex items-center gap-2">
          <svg className="animate-spin h-4 w-4 text-brand" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          Cargando catalogo de equipos...
        </div>
      )}

      {fetchError && !loading && (
        <div className="p-4 bg-danger-soft border border-[var(--color-border-default)] rounded-lg text-sm text-danger">
          {fetchError}
        </div>
      )}

      {/* Panel Selection */}
      {!loading && !fetchError && (
      <>
      <div>
        <label className="block text-sm font-medium text-fg-secondary mb-2">Panel Solar</label>
        {panels.length === 0 ? (
          <div className="p-4 bg-warning-soft border border-[var(--color-border-default)] rounded-lg text-sm text-warning">
            No hay paneles disponibles para el formato seleccionado. Vuelve al paso anterior para
            cambiar la seleccion.
          </div>
        ) : (
          <Select
            value={formData.equipment.panelId}
            onChange={(value) =>
              setFormData({ equipment: { ...formData.equipment, panelId: value } })
            }
            options={panelOptions}
            placeholder="Seleccionar panel..."
          />
        )}

        {selectedPanel && (
          <div className="mt-3 bg-inset rounded-md p-4 grid grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-fg-tertiary">Potencia:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.powerWp}W</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Eficiencia:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{(selectedPanel.efficiency * 100).toFixed(1)}%</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Area:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.area} m2</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Voc:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.voc}V</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Isc:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.isc}A</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Vmp:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.vmp}V</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Imp:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.imp}A</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Coef. Temp:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedPanel.tempCoeffPmax}%/C</span>
            </div>
          </div>
        )}
      </div>

      {/* Inverter Selection */}
      <div>
        <label className="block text-sm font-medium text-fg-secondary mb-2">Inversor</label>
        {inverters.length === 0 ? (
          <div className="p-4 bg-warning-soft border border-[var(--color-border-default)] rounded-lg text-sm text-warning">
            No hay inversores disponibles para los criterios seleccionados. Vuelve al paso anterior
            para ajustar la configuracion.
          </div>
        ) : (
          <Select
            value={formData.equipment.inverterId}
            onChange={(value) =>
              setFormData({ equipment: { ...formData.equipment, inverterId: value } })
            }
            options={inverterOptions}
            placeholder="Seleccionar inversor..."
          />
        )}

        {selectedInverter && (
          <div className="mt-3 bg-inset rounded-md p-4 grid grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-fg-tertiary">Potencia:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedInverter.ratedPowerKw}kW</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Eficiencia:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{(selectedInverter.efficiency * 100).toFixed(1)}%</span>
            </div>
            <div>
              <span className="text-fg-tertiary">MPPT:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedInverter.mpptCount} canales</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Tipo:</span>
              <span className="ml-1 font-medium text-fg-primary">{selectedInverter.type}</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Rango MPPT:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedInverter.mpptVoltageMin}-{selectedInverter.mpptVoltageMax}V</span>
            </div>
            <div>
              <span className="text-fg-tertiary">V max entrada:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedInverter.maxInputVoltage}V</span>
            </div>
            <div>
              <span className="text-fg-tertiary">I max:</span>
              <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedInverter.maxInputCurrent}A</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Bateria:</span>
              <span className="ml-1 font-medium text-fg-primary">{selectedInverter.hasBatteryPort ? 'Si' : 'No'}</span>
            </div>
          </div>
        )}
      </div>
      </>
      )}
    </div>
  );
}
