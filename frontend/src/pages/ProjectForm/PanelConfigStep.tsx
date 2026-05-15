import { useEffect, useCallback, useState, useMemo } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { catalogApi } from '../../api/catalog';
import { formatCOP } from '../../utils/format';
import { PANEL_FORMATS } from './panelFormats';
import { PanelConnectionType } from '../../types';
import Select from '../../components/common/Select';

const CONNECTION_TYPES: { value: PanelConnectionType; label: string; description: string }[] = [
  {
    value: 'serie',
    label: 'Serie',
    description:
      'Se conecta el positivo de un panel con el negativo del siguiente. Los voltajes se suman y la corriente se mantiene constante. Ideal para sistemas de alto voltaje.',
  },
  {
    value: 'paralelo',
    label: 'Paralelo',
    description:
      'Todos los positivos se conectan entre si y todos los negativos igualmente. Las corrientes se suman y el voltaje se mantiene constante. Ideal para sistemas de bajo voltaje.',
  },
  {
    value: 'mixto',
    label: 'Mixto serie-paralelo',
    description:
      'Se arman varias cadenas en serie y luego se conectan en paralelo. Permite alcanzar el voltaje y la corriente adecuados para el sistema.',
  },
];

const SYSTEM_EFFICIENCY = (1 - 0.03) * (1 - 0.02) * (1 - 0.02) * 0.96;

function ValueItem({ label, value, unit, highlight }: { label: string; value: string; unit?: string; highlight?: boolean }) {
  return (
    <div>
      <span className="text-fg-tertiary text-xs block">{label}</span>
      <span className={`font-medium font-mono tabular-nums ${highlight ? 'text-brand' : 'text-fg-primary'}`}>
        {value}
        {unit && <span className="text-fg-tertiary ml-0.5 text-xs">{unit}</span>}
      </span>
    </div>
  );
}

export default function PanelConfigStep() {
  const { formData, setFormData, panels, setPanels, irradiationPreview } = useProjectStore();
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const selectedFormat = PANEL_FORMATS.find((f) => f.key === formData.panelFormat);
  const cfg = formData.equipment.panelConfiguration;

  const avgHSP = irradiationPreview?.annualAvgHSP ?? 4.5;
  const requiredPowerKwp = useMemo(() => {
    const annual = formData.consumption.monthly.reduce((sum, v) => sum + v, 0);
    const adjustedDaily = (annual / 365) * (formData.coveragePercentage / 100);
    return adjustedDaily / avgHSP / SYSTEM_EFFICIENCY;
  }, [formData.consumption.monthly, formData.coveragePercentage, avgHSP]);

  const handleFormatSelect = (formatKey: typeof PANEL_FORMATS[number]['key']) => {
    if (formData.panelFormat === formatKey) return;
    setFormData({
      panelFormat: formatKey,
      equipment: {
        panelId: '',
        inverterId: '',
        chargeControllerId: '',
        panelConfiguration: { connectionType: 'serie', panelsPerString: 1, numberOfStrings: 1 },
      },
    });
  };

  const fetchPanels = useCallback(() => {
    setLoading(true);
    setFetchError(null);

    const panelParams: Record<string, string> = {};
    if (selectedFormat) {
      panelParams.minPower = String(selectedFormat.minPower);
      panelParams.maxPower = String(selectedFormat.maxPower);
    }
    catalogApi
      .getPanels(panelParams)
      .then(setPanels)
      .catch(() => {
        setFetchError('Error al cargar el catalogo de paneles. Verifica tu conexion e intenta de nuevo.');
      })
      .finally(() => setLoading(false));
  }, [formData.panelFormat]);

  useEffect(() => {
    fetchPanels();
  }, [fetchPanels]);

  // Auto-select the panel when only one is available
  useEffect(() => {
    if (panels.length === 1 && !formData.equipment.panelId) {
      setFormData({ equipment: { ...formData.equipment, panelId: panels[0]._id } });
    }
  }, [panels]);

  const selectedPanel = panels.find((p) => p._id === formData.equipment.panelId);

  const panelOptions = panels.map((panel) => ({
    value: panel._id,
    label: `${panel.manufacturer} ${panel.model} - ${panel.powerWp}W - ${formatCOP(panel.costCOP)}`,
  }));

  const updateConfig = (changes: Partial<typeof cfg>) => {
    setFormData({
      equipment: {
        ...formData.equipment,
        panelConfiguration: { ...cfg, ...changes },
      },
    });
  };

  const selectConnectionType = (value: PanelConnectionType) => {
    if (value === 'serie') {
      updateConfig({ connectionType: 'serie', numberOfStrings: 1 });
    } else if (value === 'paralelo') {
      updateConfig({ connectionType: 'paralelo', panelsPerString: 1 });
    } else {
      updateConfig({ connectionType: 'mixto' });
    }
  };

  const parseCount = (raw: string) => {
    const n = parseInt(raw, 10);
    return isNaN(n) ? 1 : Math.max(1, n);
  };

  const aggregate = useMemo(() => {
    if (!selectedPanel) return null;
    const s = cfg.panelsPerString;
    const p = cfg.numberOfStrings;
    return {
      voc: selectedPanel.voc * s,
      vmp: selectedPanel.vmp * s,
      isc: selectedPanel.isc * p,
      imp: selectedPanel.imp * p,
      power: selectedPanel.powerWp * s * p,
      totalPanels: s * p,
    };
  }, [selectedPanel, cfg.panelsPerString, cfg.numberOfStrings]);

  // Inverter sizing depends on the actual array power (coverage % + panel configuration)
  const recommendedInverterKw = useMemo(() => {
    if (!selectedPanel) return 0;
    const arrayKwp = (selectedPanel.powerWp * cfg.panelsPerString * cfg.numberOfStrings) / 1000;
    return Math.ceil(arrayKwp * 1.1);
  }, [selectedPanel, cfg.panelsPerString, cfg.numberOfStrings]);

  useEffect(() => {
    if (recommendedInverterKw !== formData.recommendedInverterKw) {
      setFormData({ recommendedInverterKw });
    }
  }, [recommendedInverterKw]);

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Configuracion de Paneles</h3>
        <p className="text-sm text-fg-secondary">
          Elige el formato y el panel a usar, y define como se conectan entre si. La configuracion
          determina los valores electricos del arreglo y se usa en el calculo del sistema.
        </p>
      </div>

      {/* Panel format selection */}
      <div>
        <h4 className="font-semibold text-fg-primary text-base mb-3">Formato de Panel</h4>
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

      {loading && (
        <div className="p-4 bg-info-soft border border-[var(--color-border-default)] rounded-lg text-sm text-info flex items-center gap-2">
          <svg className="animate-spin h-4 w-4 text-brand" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          Cargando catalogo de paneles...
        </div>
      )}

      {fetchError && !loading && (
        <div className="p-4 bg-danger-soft border border-[var(--color-border-default)] rounded-lg text-sm text-danger">
          {fetchError}
        </div>
      )}

      {!loading && !fetchError && (
        <>
          {/* Panel selection */}
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
                <ValueItem label="Potencia" value={String(selectedPanel.powerWp)} unit="W" />
                <ValueItem label="Voc" value={String(selectedPanel.voc)} unit="V" />
                <ValueItem label="Isc" value={String(selectedPanel.isc)} unit="A" />
                <ValueItem label="Vmp" value={String(selectedPanel.vmp)} unit="V" />
                <ValueItem label="Imp" value={String(selectedPanel.imp)} unit="A" />
                <ValueItem label="Area" value={String(selectedPanel.area)} unit="m2" />
                <ValueItem label="NOCT" value={String(selectedPanel.NOCT)} unit="C" />
                <ValueItem label="Eficiencia" value={(selectedPanel.efficiency * 100).toFixed(1)} unit="%" />
              </div>
            )}
          </div>

          {/* Connection type selector */}
          <div>
            <label className="block text-sm font-medium text-fg-secondary mb-2">Tipo de Conexion</label>
            <div className="grid gap-4">
              {CONNECTION_TYPES.map((type) => {
                const isSelected = cfg.connectionType === type.value;
                return (
                  <button
                    key={type.value}
                    onClick={() => selectConnectionType(type.value)}
                    className={`p-5 rounded-lg border-2 text-left transition-colors duration-fast ease-decel ${
                      isSelected
                        ? 'border-brand bg-brand-soft'
                        : 'border-[var(--color-border-default)] hover:border-[var(--color-border-strong)] bg-surface'
                    }`}
                  >
                    <div className="flex items-center gap-3 mb-2">
                      <div
                        className={`w-5 h-5 rounded-full border-2 flex items-center justify-center ${
                          isSelected ? 'border-brand' : 'border-[var(--color-border-default)]'
                        }`}
                      >
                        {isSelected && <div className="w-3 h-3 rounded-full bg-brand" />}
                      </div>
                      <h4 className="font-semibold text-fg-primary">{type.label}</h4>
                    </div>
                    <p className="text-sm text-fg-secondary ml-8">{type.description}</p>
                  </button>
                );
              })}
            </div>
          </div>

          {/* Manual inputs */}
          <div>
            <label className="block text-sm font-medium text-fg-secondary mb-2">Cantidad de Paneles</label>
            <div className="grid grid-cols-2 gap-4">
              {cfg.connectionType !== 'paralelo' && (
                <div>
                  <span className="text-xs text-fg-tertiary block mb-1">
                    {cfg.connectionType === 'mixto' ? 'Paneles en serie por cadena' : 'Paneles en serie'}
                  </span>
                  <input
                    type="number"
                    min={1}
                    value={cfg.panelsPerString}
                    onChange={(e) => updateConfig({ panelsPerString: parseCount(e.target.value) })}
                    className="w-full px-3 py-2 rounded-md border border-[var(--color-border-default)] bg-surface text-fg-primary font-mono tabular-nums focus:outline-none focus:border-brand"
                  />
                </div>
              )}
              {cfg.connectionType !== 'serie' && (
                <div>
                  <span className="text-xs text-fg-tertiary block mb-1">
                    {cfg.connectionType === 'mixto' ? 'Numero de cadenas en paralelo' : 'Paneles en paralelo'}
                  </span>
                  <input
                    type="number"
                    min={1}
                    value={cfg.numberOfStrings}
                    onChange={(e) => updateConfig({ numberOfStrings: parseCount(e.target.value) })}
                    className="w-full px-3 py-2 rounded-md border border-[var(--color-border-default)] bg-surface text-fg-primary font-mono tabular-nums focus:outline-none focus:border-brand"
                  />
                </div>
              )}
            </div>
          </div>

          {/* Aggregate results */}
          {aggregate && (
            <div>
              <label className="block text-sm font-medium text-fg-secondary mb-2">
                Valores del Arreglo ({cfg.panelsPerString} en serie x {cfg.numberOfStrings} en paralelo)
              </label>
              <div className="bg-inset rounded-md p-4 grid grid-cols-3 gap-4 text-sm">
                <ValueItem label="Voc total" value={aggregate.voc.toFixed(2)} unit="V" highlight />
                <ValueItem label="Vmp total" value={aggregate.vmp.toFixed(2)} unit="V" highlight />
                <ValueItem label="Isc total" value={aggregate.isc.toFixed(2)} unit="A" highlight />
                <ValueItem label="Imp total" value={aggregate.imp.toFixed(2)} unit="A" highlight />
                <ValueItem label="Potencia total" value={(aggregate.power / 1000).toFixed(2)} unit="kWp" highlight />
                <ValueItem label="Total de paneles" value={String(aggregate.totalPanels)} unit="uds" highlight />
              </div>
            </div>
          )}

          {/* Invariant properties */}
          {selectedPanel && (
            <div className="p-4 bg-info-soft border border-[var(--color-border-default)] rounded-lg">
              <h4 className="text-sm font-semibold text-fg-primary mb-2">
                Propiedades del modulo (no cambian con la conexion)
              </h4>
              <div className="grid grid-cols-4 gap-4 text-sm">
                <ValueItem label="NOCT" value={String(selectedPanel.NOCT)} unit="C" />
                <ValueItem label="Eficiencia" value={(selectedPanel.efficiency * 100).toFixed(1)} unit="%" />
                <ValueItem label="Coef. Pmax" value={String(selectedPanel.tempCoeffPmax)} unit="%/C" />
                <ValueItem label="Coef. Voc" value={String(selectedPanel.tempCoeffVoc)} unit="%/C" />
              </div>
              <p className="mt-3 text-xs text-fg-tertiary">
                Estos valores son propiedades fisicas de cada modulo individual. Conectar paneles en
                serie o paralelo no los altera; solo cambian el voltaje y la corriente del arreglo.
              </p>
            </div>
          )}
        </>
      )}
    </div>
  );
}
