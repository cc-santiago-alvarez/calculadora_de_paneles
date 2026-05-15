import { useEffect, useState } from 'react';
import { Dialog, DialogPanel, DialogTitle } from '@headlessui/react';
import { useForm, Controller } from 'react-hook-form';
import { Panel } from '../../types';
import { catalogApi } from '../../api/catalog';
import Button from '../../components/common/Button';
import Select from '../../components/common/Select';

const TYPE_OPTIONS = [
  { value: 'monocrystalline', label: 'Monocristalino' },
  { value: 'polycrystalline', label: 'Policristalino' },
  { value: 'thin-film', label: 'Capa fina' },
];

// Opciones alineadas a PANEL_FORMATS (pages/ProjectForm/panelFormats.ts).
const FORMAT_OPTIONS = [
  { value: 'Estandar', label: 'Estandar' },
  { value: 'Alta eficiencia', label: 'Alta eficiencia' },
  { value: 'Gran formato', label: 'Gran formato' },
];

interface PanelFormValues {
  manufacturer: string;
  model: string;
  type: string;
  powerWp: number;
  efficiency: number;
  area: number;
  voc: number;
  isc: number;
  vmp: number;
  imp: number;
  tempCoeffPmax: number;
  tempCoeffVoc: number;
  NOCT: number;
  weight: number;
  warranty: number;
  costCOP: number;
  format: string;
  dimensions: { length: number; width: number };
}

const emptyValues: PanelFormValues = {
  manufacturer: '',
  model: '',
  type: 'monocrystalline',
  powerWp: 0,
  efficiency: 0,
  area: 0,
  voc: 0,
  isc: 0,
  vmp: 0,
  imp: 0,
  tempCoeffPmax: 0,
  tempCoeffVoc: 0,
  NOCT: 0,
  weight: 0,
  warranty: 0,
  costCOP: 0,
  format: '',
  dimensions: { length: 0, width: 0 },
};

const NUMBER_FIELDS: { name: keyof PanelFormValues; label: string; step: string }[] = [
  { name: 'powerWp', label: 'Potencia (Wp)', step: '1' },
  { name: 'efficiency', label: 'Eficiencia (0-1)', step: '0.001' },
  { name: 'area', label: 'Area (m2)', step: '0.001' },
  { name: 'voc', label: 'Voc (V)', step: '0.01' },
  { name: 'isc', label: 'Isc (A)', step: '0.01' },
  { name: 'vmp', label: 'Vmp (V)', step: '0.01' },
  { name: 'imp', label: 'Imp (A)', step: '0.01' },
  { name: 'tempCoeffPmax', label: 'Coef. temp. Pmax (%/C)', step: '0.001' },
  { name: 'tempCoeffVoc', label: 'Coef. temp. Voc (%/C)', step: '0.001' },
  { name: 'NOCT', label: 'NOCT (C)', step: '1' },
  { name: 'weight', label: 'Peso (kg)', step: '0.1' },
  { name: 'warranty', label: 'Garantia (anios)', step: '1' },
  { name: 'costCOP', label: 'Costo (COP)', step: '1000' },
];

const inputClass =
  'w-full rounded-md bg-inset border border-[var(--color-border-default)] py-2 px-3 text-sm text-fg-primary transition-colors duration-fast focus-visible:outline-2 focus-visible:outline-brand focus-visible:outline-offset-2 hover:border-[var(--color-border-strong)]';
const labelClass = 'block text-xs font-medium text-fg-secondary mb-1';

interface Props {
  open: boolean;
  panel: Panel | null;
  onClose: () => void;
  onSaved: () => void;
}

export default function PanelFormModal({ open, panel, onClose, onSaved }: Props) {
  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<PanelFormValues>({ defaultValues: emptyValues });
  const [serverError, setServerError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    setServerError(null);
    if (panel) {
      reset({
        manufacturer: panel.manufacturer,
        model: panel.model,
        type: panel.type || 'monocrystalline',
        powerWp: panel.powerWp,
        efficiency: panel.efficiency,
        area: panel.area,
        voc: panel.voc,
        isc: panel.isc,
        vmp: panel.vmp,
        imp: panel.imp,
        tempCoeffPmax: panel.tempCoeffPmax,
        tempCoeffVoc: panel.tempCoeffVoc,
        NOCT: panel.NOCT,
        weight: panel.weight ?? 0,
        warranty: panel.warranty ?? 0,
        costCOP: panel.costCOP,
        format: panel.format ?? '',
        dimensions: {
          length: panel.dimensions?.length ?? 0,
          width: panel.dimensions?.width ?? 0,
        },
      });
    } else {
      reset(emptyValues);
    }
  }, [open, panel, reset]);

  const onSubmit = async (values: PanelFormValues) => {
    setServerError(null);
    try {
      if (panel) {
        await catalogApi.updatePanel(panel._id, values);
      } else {
        await catalogApi.createPanel(values);
      }
      onSaved();
    } catch (err: any) {
      setServerError(err?.response?.data?.error || 'No se pudo guardar el panel');
    }
  };

  return (
    <Dialog open={open} onClose={onClose} className="relative z-50">
      <div className="fixed inset-0 bg-black/40" aria-hidden="true" />
      <div className="fixed inset-0 flex items-center justify-center p-4">
        <DialogPanel className="w-full max-w-2xl max-h-[90vh] overflow-auto rounded-lg bg-surface border border-[var(--color-border-default)] p-6">
          <DialogTitle className="text-lg font-bold text-fg-primary mb-4">
            {panel ? 'Editar panel' : 'Agregar panel'}
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
              <div>
                <label className={labelClass}>Formato</label>
                <Controller
                  name="format"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onChange={field.onChange}
                      options={FORMAT_OPTIONS}
                      placeholder="Seleccionar formato..."
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

              <div>
                <label className={labelClass}>Largo (mm)</label>
                <input
                  type="number"
                  step="1"
                  className={inputClass}
                  {...register('dimensions.length', { valueAsNumber: true })}
                />
              </div>
              <div>
                <label className={labelClass}>Ancho (mm)</label>
                <input
                  type="number"
                  step="1"
                  className={inputClass}
                  {...register('dimensions.width', { valueAsNumber: true })}
                />
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
