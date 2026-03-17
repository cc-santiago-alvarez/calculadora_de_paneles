import { useProjectStore } from '../../store/useProjectStore';
import { SystemType } from '../../types';

const METER_INFO: Record<SystemType, { meters: { name: string; detail: string }[]; note: string }> = {
  'on-grid': {
    meters: [
      { name: '1 medidor bidireccional', detail: 'Lo instala la empresa distribuidora' },
    ],
    note: 'Registra lo que consumes de la red y lo que exportas. Tramite gestionado por la distribuidora (Codensa, EPM, Celsia, etc.).',
  },
  'off-grid': {
    meters: [
      { name: '1 medidor interno', detail: 'Consumo propio y monitoreo' },
    ],
    note: 'No requiere medidores oficiales. Solo un sensor/medidor interno para control del propietario.',
  },
  'hybrid': {
    meters: [
      { name: '1 medidor bidireccional', detail: 'Lo instala la empresa distribuidora' },
      { name: '1 medidor de produccion solar', detail: 'Normalmente integrado en el inversor hibrido' },
    ],
    note: 'Minimo 2 medidores. Algunos proyectos agregan un 3er medidor de consumo para auditoria energetica completa.',
  },
};

const SYSTEM_TYPES: { value: SystemType; label: string; description: string }[] = [
  {
    value: 'on-grid',
    label: 'Conectado a Red (On-Grid)',
    description:
      'El sistema se conecta a la red electrica. Los excedentes se inyectan a la red y se descuenta de la factura. No requiere baterias.',
  },
  {
    value: 'off-grid',
    label: 'Aislado (Off-Grid)',
    description:
      'Sistema independiente de la red electrica. Requiere banco de baterias para almacenamiento. Ideal para zonas sin acceso a red.',
  },
  {
    value: 'hybrid',
    label: 'Hibrido',
    description:
      'Combinacion de conexion a red con baterias de respaldo. Permite almacenar energia y tener respaldo ante cortes.',
  },
];

export default function SystemTypeStep() {
  const { formData, setFormData } = useProjectStore();

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold text-fg-primary mb-2">Tipo de Sistema</h3>
        <p className="text-sm text-fg-secondary">
          Selecciona el tipo de instalacion fotovoltaica.
        </p>
      </div>

      <div className="grid gap-4">
        {SYSTEM_TYPES.map((type) => {
          const isSelected = formData.systemType === type.value;
          return (
            <button
              key={type.value}
              onClick={() => setFormData({ systemType: type.value })}
              className={`p-6 rounded-lg border-2 text-left transition-colors duration-fast ease-decel ${
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
                  {isSelected && (
                    <div className="w-3 h-3 rounded-full bg-brand" />
                  )}
                </div>
                <h4 className="font-semibold text-fg-primary">{type.label}</h4>
              </div>
              <p className="text-sm text-fg-secondary ml-8">{type.description}</p>
            </button>
          );
        })}
      </div>

      {formData.systemType && (
        <div className="p-4 bg-info-soft border border-[var(--color-border-default)] rounded-lg">
          <h4 className="text-sm font-semibold text-fg-primary mb-3">
            Medidores y complementos necesarios
          </h4>
          <ul className="space-y-2">
            {METER_INFO[formData.systemType].meters.map((meter, i) => (
              <li key={i} className="text-sm text-fg-secondary flex items-start gap-2">
                <span className="bg-brand-soft text-brand px-2 py-0.5 rounded-full text-xs font-medium shrink-0 mt-0.5">
                  Requerido
                </span>
                <div>
                  <span className="font-medium text-fg-primary">{meter.name}</span>
                  <span className="text-fg-tertiary ml-1">— {meter.detail}</span>
                </div>
              </li>
            ))}
          </ul>
          <p className="mt-3 text-xs text-fg-tertiary">
            {METER_INFO[formData.systemType].note}
          </p>
        </div>
      )}
    </div>
  );
}
