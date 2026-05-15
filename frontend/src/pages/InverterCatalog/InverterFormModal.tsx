import { useEffect, useState } from 'react';
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/react';
import { useForm, Controller } from 'react-hook-form';
import { Inverter } from '../../types';
import { catalogApi } from '../../api/catalog';
import Button from '../../components/common/Button';
import Select from '../../components/common/Select';

const TYPE_OPTIONS = [
  { value: 'string', label: 'String' },
  { value: 'hybrid', label: 'Hibrido' },
  { value: 'off-grid', label: 'Off-grid' },
];

interface InverterFormValues {
  manufacturer: string;
  model: string;
  type: string;
  ratedPowerKw: number;
  maxDCPowerKw: number;
  efficiency: number;
  mpptCount: number;
  mpptVoltageMin: number;
  mpptVoltageMax: number;
  maxInputVoltage: number;
  maxInputCurrent: number;
  outputVoltage: number;
  outputPhases: number;
  hasBatteryPort: boolean;
  weight: number;
  warranty: number;
  costCOP: number;
}

const emptyValues: InverterFormValues = {
  manufacturer: '',
  model: '',
  type: 'string',
  ratedPowerKw: 0,
  maxDCPowerKw: 0,
  efficiency: 0,
  mpptCount: 0,
  mpptVoltageMin: 0,
  mpptVoltageMax: 0,
  maxInputVoltage: 0,
  maxInputCurrent: 0,
  outputVoltage: 0,
  outputPhases: 0,
  hasBatteryPort: false,
  weight: 0,
  warranty: 0,
  costCOP: 0,
};

const NUMBER_FIELDS: { name: keyof InverterFormValues; label: string; step: string }[] = [
  { name: 'ratedPowerKw', label: 'Potencia nominal (kW)', step: '0.1' },
  { name: 'maxDCPowerKw', label: 'Potencia DC max (kW)', step: '0.1' },
  { name: 'efficiency', label: 'Eficiencia (0-1)', step: '0.001' },
  { name: 'mpptCount', label: 'Cantidad de MPPT', step: '1' },
  { name: 'mpptVoltageMin', label: 'Voltaje MPPT min (V)', step: '1' },
  { name: 'mpptVoltageMax', label: 'Voltaje MPPT max (V)', step: '1' },
  { name: 'maxInputVoltage', label: 'Voltaje entrada max (V)', step: '1' },
  { name: 'maxInputCurrent', label: 'Corriente entrada max (A)', step: '0.1' },
  { name: 'outputVoltage', label: 'Voltaje de salida (V)', step: '1' },
  { name: 'outputPhases', label: 'Fases de salida', step: '1' },
  { name: 'weight', label: 'Peso (kg)', step: '0.1' },
  { name: 'warranty', label: 'Garantia (anios)', step: '1' },
  { name: 'costCOP', label: 'Costo (COP)', step: '1000' },
];

const inputClass =
  'w-full rounded-md bg-inset border border-[var(--color-border-default)] py-2 px-3 text-sm text-fg-primary transition-colors duration-fast focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 hover:border-[var(--color-border-strong)]';
const labelClass = 'block text-xs font-medium text-fg-secondary mb-1';

interface Props {
  open: boolean;
  inverter: Inverter | null;
  onClose: () => void;
  onSaved: () => void;
}

export default function InverterFormModal({ open, inverter, onClose, onSaved }: Props) {
  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<InverterFormValues>({ defaultValues: emptyValues });
  const [serverError, setServerError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    setServerError(null);
    if (inverter) {
      reset({
        manufacturer: inverter.manufacturer,
        model: inverter.model,
        type: inverter.type || 'string',
        ratedPowerKw: inverter.ratedPowerKw,
        maxDCPowerKw: inverter.maxDCPowerKw,
        efficiency: inverter.efficiency,
        mpptCount: inverter.mpptCount,
        mpptVoltageMin: inverter.mpptVoltageMin,
        mpptVoltageMax: inverter.mpptVoltageMax,
        maxInputVoltage: inverter.maxInputVoltage,
        maxInputCurrent: inverter.maxInputCurrent,
        outputVoltage: inverter.outputVoltage,
        outputPhases: inverter.outputPhases,
        hasBatteryPort: inverter.hasBatteryPort,
        weight: inverter.weight ?? 0,
        warranty: inverter.warranty ?? 0,
        costCOP: inverter.costCOP,
      });
    } else {
      reset(emptyValues);
    }
  }, [open, inverter, reset]);

  const onSubmit = async (values: InverterFormValues) => {
    setServerError(null);
    try {
      if (inverter) {
        await catalogApi.updateInverter(inverter._id, values);
      } else {
        await catalogApi.createInverter(values);
      }
      onSaved();
    } catch (err: any) {
      setServerError(err?.response?.data?.error || 'No se pudo guardar el inversor');
    }
  };

  return (
    <Dialog open={open} onClose={onClose} className="relative z-50">
      <div className="fixed inset-0 bg-black/40" aria-hidden="true" />
      <div className="fixed inset-0 flex items-center justify-center p-4">
        <DialogPanel className="w-full max-w-2xl max-h-[90vh] overflow-auto rounded-lg bg-surface border border-[var(--color-border-default)] p-6">
          <DialogTitle className="text-lg font-bold text-fg-primary mb-4">
            {inverter ? 'Editar inversor' : 'Agregar inversor'}
          </DialogTitle>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className={labelClass}>Fabricante *</label>
                <input
                  className={inputClass}
                  {...register('manufacturer', { required: 'Obligatorio' })}
                />
                {errors.manufacturer && (
                  <p className="text-xs text-danger mt-1">{errors.manufacturer.message}</p>
                )}
              </div>
              <div>
                <label className={labelClass}>Modelo *</label>
                <input
                  className={inputClass}
                  {...register('model', { required: 'Obligatorio' })}
                />
                {errors.model && (
                  <p className="text-xs text-danger mt-1">{errors.model.message}</p>
                )}
              </div>

              <div>
                <label className={labelClass}>Tipo</label>
                <Controller
                  name="type"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onChange={field.onChange}
                      options={TYPE_OPTIONS}
                    />
                  )}
                />
              </div>

              {NUMBER_FIELDS.map((f) => (
                <div key={f.name}>
                  <label className={labelClass}>{f.label}</label>
                  <input
                    type="number"
                    step={f.step}
                    className={inputClass}
                    {...register(f.name, { valueAsNumber: true })}
                  />
                </div>
              ))}

              <div className="flex items-center gap-2 sm:col-span-2">
                <input
                  id="hasBatteryPort"
                  type="checkbox"
                  className="h-4 w-4 rounded border-[var(--color-border-default)] accent-brand"
                  {...register('hasBatteryPort')}
                />
                <label htmlFor="hasBatteryPort" className="text-sm text-fg-secondary">
                  Tiene puerto de bateria
                </label>
              </div>
            </div>

            {serverError && (
              <div className="rounded-md bg-danger-soft text-danger text-sm px-3 py-2">
                {serverError}
              </div>
            )}

            <div className="flex justify-end gap-2 pt-2">
              <Button type="button" variant="secondary" onClick={onClose} disabled={isSubmitting}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? 'Guardando...' : 'Guardar'}
              </Button>
            </div>
          </form>
        </DialogPanel>
      </div>
    </Dialog>
  );
}
