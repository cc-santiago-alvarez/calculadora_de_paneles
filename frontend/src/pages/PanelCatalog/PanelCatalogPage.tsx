import { useCallback, useEffect, useState } from 'react';
import { catalogApi } from '../../api/catalog';
import { Panel } from '../../types';
import { formatCOP, formatNumber } from '../../utils/format';
import Button from '../../components/common/Button';
import Card from '../../components/common/Card';
import PanelFormModal from './PanelFormModal';

export default function PanelCatalogPage() {
  const [panels, setPanels] = useState<Panel[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Panel | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    setError(null);
    catalogApi
      .getPanels()
      .then(setPanels)
      .catch(() => setError('No se pudieron cargar los paneles.'))
      .finally(() => setLoading(false));
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleNew = () => {
    setEditing(null);
    setModalOpen(true);
  };

  const handleEdit = (panel: Panel) => {
    setEditing(panel);
    setModalOpen(true);
  };

  const handleDelete = async (panel: Panel) => {
    if (!confirm(`Eliminar el panel "${panel.manufacturer} ${panel.model}"?`)) return;
    try {
      await catalogApi.deletePanel(panel._id);
      load();
    } catch {
      setError('No se pudo eliminar el panel.');
    }
  };

  const handleSaved = () => {
    setModalOpen(false);
    load();
  };

  return (
    <div className="max-w-5xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-fg-primary">Catalogo de paneles</h2>
        <Button onClick={handleNew} size="lg">
          + Agregar panel
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
      ) : panels.length === 0 ? (
        <Card padding="lg" className="text-center py-20">
          <p className="text-fg-muted text-lg mb-2">No hay paneles en el catalogo.</p>
          <p className="text-fg-muted text-sm mb-6">
            Agrega tu primer panel solar para usarlo en los proyectos.
          </p>
          <Button onClick={handleNew}>Agregar panel</Button>
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
                  <th className="px-4 py-3 font-medium">Formato</th>
                  <th className="px-4 py-3 font-medium text-right">Potencia</th>
                  <th className="px-4 py-3 font-medium text-right">Eficiencia</th>
                  <th className="px-4 py-3 font-medium text-right">Costo</th>
                  <th className="px-4 py-3 font-medium text-right">Acciones</th>
                </tr>
              </thead>
              <tbody>
                {panels.map((panel) => (
                  <tr
                    key={panel._id}
                    className="border-b border-[var(--color-border-subtle)] last:border-0 hover:bg-inset transition-colors duration-fast"
                  >
                    <td className="px-4 py-3 text-fg-primary">{panel.manufacturer}</td>
                    <td className="px-4 py-3 text-fg-secondary">{panel.model}</td>
                    <td className="px-4 py-3 text-fg-secondary">{panel.type}</td>
                    <td className="px-4 py-3 text-fg-secondary">{panel.format || '-'}</td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-primary">
                      {formatNumber(panel.powerWp)} Wp
                    </td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-secondary">
                      {formatNumber(panel.efficiency * 100, 1)}%
                    </td>
                    <td className="px-4 py-3 text-right font-mono tabular-nums text-fg-secondary">
                      {formatCOP(panel.costCOP)}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        <Button variant="ghost" size="sm" onClick={() => handleEdit(panel)}>
                          Editar
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="hover:text-danger hover:bg-danger-soft"
                          onClick={() => handleDelete(panel)}
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

      <PanelFormModal
        open={modalOpen}
        panel={editing}
        onClose={() => setModalOpen(false)}
        onSaved={handleSaved}
      />
    </div>
  );
}
