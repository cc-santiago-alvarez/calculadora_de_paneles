import { useEffect, useCallback, useState } from 'react';
import { useProjectStore } from '../../store/useProjectStore';
import { catalogApi } from '../../api/catalog';
import { formatCOP } from '../../utils/format';
import { PANEL_FORMATS } from './panelFormats';
import Select from '../../components/common/Select';

export default function EquipmentStep() {
  const { formData, setFormData, inverters, chargeControllers, setInverters, setChargeControllers } = useProjectStore();
  const [loading, setLoading] = useState(true);
  const [fetchError, setFetchError] = useState<string | null>(null);

  const selectedFormat = PANEL_FORMATS.find((f) => f.key === formData.panelFormat);
  const isOffGrid = formData.systemType === 'off-grid';

  const fetchCatalogs = useCallback(() => {
    setLoading(true);
    setFetchError(null);

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

    // Fetch charge controllers only for off-grid
    let chargeControllerPromise = Promise.resolve();
    if (isOffGrid) {
      chargeControllerPromise = catalogApi.getChargeControllers({ type: 'mppt' }).then(setChargeControllers);
    } else {
      setChargeControllers([]);
      if (formData.equipment.chargeControllerId) {
        setFormData({ equipment: { ...formData.equipment, chargeControllerId: '' } });
      }
    }

    Promise.all([inverterPromise, chargeControllerPromise])
      .catch(() => {
        setFetchError('Error al cargar el catalogo de equipos. Verifica tu conexion e intenta de nuevo.');
      })
      .finally(() => {
        setLoading(false);
      });
  }, [formData.recommendedInverterKw, formData.systemType]);

  // Always fetch on mount and when filter criteria change
  useEffect(() => {
    fetchCatalogs();
  }, [fetchCatalogs]);

  const selectedInverter = inverters.find((i) => i._id === formData.equipment.inverterId);
  const selectedChargeController = chargeControllers.find((cc) => cc._id === formData.equipment.chargeControllerId);

  const inverterHasBuiltInMPPT = selectedInverter && selectedInverter.mpptCount > 0;
  const chargeControllerRequired = isOffGrid && selectedInverter && selectedInverter.mpptCount === 0;
  const showChargeController = isOffGrid;

  const inverterOptions = inverters.map((inv) => ({
    value: inv._id,
    label: `${inv.manufacturer} ${inv.model} - ${inv.ratedPowerKw}kW ${inv.type} - ${formatCOP(inv.costCOP)}`,
  }));

  const chargeControllerOptions = chargeControllers.map((cc) => ({
    value: cc._id,
    label: `${cc.manufacturer} ${cc.model} - ${cc.maxChargeCurrentA}A - ${formatCOP(cc.costCOP)}`,
  }));

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Seleccion de Inversor</h3>
        <p className="text-sm text-fg-secondary">
          Selecciona el inversor y, si aplica, el regulador de carga del catalogo.
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

      {/* Inverter Selection */}
      {!loading && !fetchError && (
      <>
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
              setFormData({ equipment: { ...formData.equipment, inverterId: value, chargeControllerId: '' } })
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

      {/* Charge Controller Selection (off-grid only) */}
      {showChargeController && (
        <div>
          <label className="block text-sm font-medium text-fg-secondary mb-2">
            Regulador de Carga
            {chargeControllerRequired ? (
              <span className="ml-2 bg-danger-soft text-danger px-2 py-0.5 rounded-full text-xs font-medium">
                Requerido
              </span>
            ) : (
              <span className="ml-2 bg-info-soft text-info px-2 py-0.5 rounded-full text-xs font-medium">
                Opcional
              </span>
            )}
          </label>

          {selectedInverter && !inverterHasBuiltInMPPT && (
            <div className="mb-3 p-3 bg-warning-soft border border-[var(--color-border-default)] rounded-lg text-sm text-warning">
              El inversor seleccionado no tiene MPPT integrado. Se requiere un regulador de carga externo MPPT
              entre los paneles y el banco de baterias.
            </div>
          )}

          {selectedInverter && inverterHasBuiltInMPPT && (
            <div className="mb-3 p-3 bg-info-soft border border-[var(--color-border-default)] rounded-lg text-sm text-info">
              El inversor seleccionado tiene MPPT integrado ({selectedInverter.mpptCount} canal{selectedInverter.mpptCount > 1 ? 'es' : ''}).
              Un regulador externo es opcional, pero puede ser util para arreglos solares mas grandes.
            </div>
          )}

          {chargeControllers.length === 0 ? (
            <div className="p-4 bg-warning-soft border border-[var(--color-border-default)] rounded-lg text-sm text-warning">
              No hay reguladores de carga disponibles.
            </div>
          ) : (
            <Select
              value={formData.equipment.chargeControllerId}
              onChange={(value) =>
                setFormData({ equipment: { ...formData.equipment, chargeControllerId: value } })
              }
              options={chargeControllerOptions}
              placeholder="Seleccionar regulador de carga..."
            />
          )}

          {selectedChargeController && (
            <div className="mt-3 bg-inset rounded-md p-4 grid grid-cols-3 gap-4 text-sm">
              <div>
                <span className="text-fg-tertiary">Tipo:</span>
                <span className="ml-1 font-medium text-fg-primary">{selectedChargeController.type.toUpperCase()}</span>
              </div>
              <div>
                <span className="text-fg-tertiary">Corriente max:</span>
                <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedChargeController.maxChargeCurrentA}A</span>
              </div>
              <div>
                <span className="text-fg-tertiary">V PV max:</span>
                <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedChargeController.maxPVVoltage}V</span>
              </div>
              <div>
                <span className="text-fg-tertiary">Eficiencia:</span>
                <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{(selectedChargeController.efficiency * 100).toFixed(1)}%</span>
              </div>
              <div>
                <span className="text-fg-tertiary">Baterias:</span>
                <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedChargeController.batteryVoltages.join('/')}V</span>
              </div>
              <div>
                <span className="text-fg-tertiary">Potencia:</span>
                <span className="ml-1 font-medium font-mono tabular-nums text-fg-primary">{selectedChargeController.ratedPowerW}W</span>
              </div>
            </div>
          )}
        </div>
      )}
      </>
      )}
    </div>
  );
}
