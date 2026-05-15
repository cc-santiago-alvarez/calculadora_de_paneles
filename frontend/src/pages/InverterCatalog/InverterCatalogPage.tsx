import { useCallback, useEffect, useState } from 'react';
import { catalogApi } from '../../api/catalog';
import { Inverter } from '../../types';
import { formatCOP, formatNumber } from '../../utils/format';
import Button from '../../components/common/Button';
import Card from '../../components/common/Card';
import InverterFormModal from './InverterFormModal';

export default function InverterCatalogPage() {
  const [inverters, setInverters] = useState<Inverter[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Inverter | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    catalogApi
      .getInverters()
      .then(setInverters)
      .catch(() => setError('No se pudieron cargar los inversores.'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleNew = () => {
    setEditing(null);
    setModalOpen(true);
  };

  const handleEdit = (inverter: Inverter) => {
    setEditing(inverter);
    setModalOpen(true);
  };

  const handleDelete = async (inverter: Inverter) => {
    if (!confirm(`Eliminar el inversor "${inverter.manufacturer} ${inverter.model}"?`)) return;
    try {
      await catalogApi.deleteInverter(inverter._id);
      load();
    } catch {
      setError('No se pudo eliminar el inversor.');
    }
  };

  const handleSaved = () => {
    setModalOpen(false);
    load();
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-fg-primary">Catalogo de inversores</h2>
        <Button onClick={handleNew} size="lg">
          + Agregar inversor
        </Button>
      </div>

      {error && (
        <div className="rounded-md bg-danger-soft text-danger text-sm px-3 py-2 mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="flex items-center justify-center h-64">
          <div className="w-8 h-8 border-4 border-brand border-t-transparent rounded-full animate-spin" />
        </div>
      ) : inverters.length === 0 ? (
        <Card padding="lg" className="text-center py-20">
          <p className="text-fg-muted text-lg mb-2">No hay inversores en el catalogo.</p>
          <p className="text-fg-muted text-sm mb-6">
            Agrega tu primer inversor para usarlo en los proyectos.
          </p>
          <Button onClick={handleNew}>Agregar inversor</Button>
        </Card>
      ) : (
        <Card padding="none" className="overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[var(--color-border-default)] text-left text-xs uppercase tracking-wide text-fg-tertiary">
                  <th className="px-4 py-3 font-medium">Fabricante</th>
                  <th className="px-4 py-3 font-medium">Modelo</th>
                  <th className="px-4 py-3 font-medium">Tipo</th>
                  <th className="px-4 py-3 font-medium text-right">Potencia</th>
                  <th className="px-4 py-3 font-medium text-right">Eficiencia</th>
                  <th className="px-4 py-3 font-medium text-right">MPPT</th>
                  <th className="px-4 py-3 font-medium">Bateria</th>
                  <th className="px-4 py-3 font-medium text-right">Costo</th>
                  <th className="px-4 py-3 font-medium text-right">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {inverters.map((inverter) => (
                  <tr
                    key={inverter._id}
                    className="border-b border-[var(--color-border-subtle)] last:border-0 hover:bg-inset transition-colors duration-fast"
                  >
                    <td className="px-4 py-3 text-fg-primary">{inverter.manufacturer}</td>
                    <td className="px-4 py-3 text-fg-secondary">{inverter.model}</td>
                    <td className="px-4 py-3 text-fg-secondary">{inverter.type}</td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-primary">
                      {formatNumber(inverter.ratedPowerKw, 1)} kW
                    </td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-secondary">
                      {formatNumber(inverter.efficiency * 100, 1)}%
                    </td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-secondary">
                      {inverter.mpptCount}
                    </td>
                    <td className="px-4 py-3 text-fg-secondary">
                      {inverter.hasBatteryPort ? 'Si' : 'No'}
                    </td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-secondary">
                      {formatCOP(inverter.costCOP)}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        <Button variant="ghost" size="sm" onClick={() => handleEdit(inverter)}>
                          Editar
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="hover:text-danger hover:bg-danger-soft"
                          onClick={() => handleDelete(inverter)}
                        >
                          Eliminar
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      <InverterFormModal
        open={modalOpen}
        inverter={editing}
        onClose={() => setModalOpen(false)}
        onSaved={handleSaved}
      />
    </div>
  );
}
