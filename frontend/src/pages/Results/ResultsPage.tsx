import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '../../store/useProjectStore';
import { projectsApi } from '../../api/projects';
import { formatCOP, formatNumber } from '../../utils/format';
import { Project, Scenario } from '../../types';
import GenerationChart from './GenerationChart';
import IrradiationChart from './IrradiationChart';
import CashFlowChart from './CashFlowChart';
import DegradationChart from './DegradationChart';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';
import DataValue from '../../components/common/DataValue';

export default function ResultsPage() {
  const navigate = useNavigate();
  const {
    projects,
    setProjects,
    currentProject,
    setCurrentProject,
    currentScenario,
    setCurrentScenario,
    setScenarios,
    addToComparison,
  } = useProjectStore();

  const [loadingList, setLoadingList] = useState(true);
  const [loadingScenario, setLoadingScenario] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch projects on mount
  useEffect(() => {
    projectsApi
      .getAll()
      .then(setProjects)
      .catch(() => setError('Error al cargar los proyectos.'))
      .finally(() => setLoadingList(false));
  }, []);

  const handleSelectProject = async (project: Project) => {
    setCurrentProject(project);
    setCurrentScenario(null as unknown as Scenario);
    setLoadingScenario(true);
    setError(null);
    try {
      const scenarios = await projectsApi.getScenarios(project._id);
      setScenarios(scenarios);
      if (scenarios.length > 0) {
        setCurrentScenario(scenarios[0]);
      }
    } catch {
      setError('Error al cargar los escenarios del proyecto.');
    } finally {
      setLoadingScenario(false);
    }
  };

  return (
    <div className="flex gap-6 h-full">
      {/* Project list sidebar */}
      <div className="w-80 flex-shrink-0 space-y-2">
        <h3 className="text-lg font-semibold text-fg-primary mb-3">Proyectos</h3>

        {loadingList && (
          <div className="flex items-center justify-center h-32">
            <div className="w-6 h-6 border-4 border-brand border-t-transparent rounded-full animate-spin" />
          </div>
        )}

        {!loadingList && projects.length === 0 && (
          <Card className="text-center py-10">
            <p className="text-fg-muted text-sm mb-3">No hay proyectos.</p>
            <Button onClick={() => navigate('/new')} size="sm">
              Crear Proyecto
            </Button>
          </Card>
        )}

        {!loadingList &&
          projects.map((project) => (
            <Card
              key={project._id}
              hoverable
              padding="sm"
              onClick={() => handleSelectProject(project)}
              className={
                currentProject?._id === project._id
                  ? 'border-brand bg-brand-soft'
                  : ''
              }
            >
              <h4 className="font-medium text-fg-primary text-sm">{project.name}</h4>
              <p className="text-xs text-fg-secondary mt-1">
                {project.location.city}, {project.location.department}
              </p>
              <div className="flex gap-3 mt-1 text-xs text-fg-muted">
                <span>{project.systemType}</span>
                <span>{project.roof.area}m2</span>
                <span>{project.scenarios?.length || 0} esc.</span>
              </div>
            </Card>
          ))}
      </div>

      {/* Dashboard area */}
      <div className="flex-1 min-w-0">
        {!currentProject && !loadingScenario && (
          <div className="flex items-center justify-center h-64 text-fg-muted">
            Selecciona un proyecto para ver sus resultados.
          </div>
        )}

        {loadingScenario && (
          <div className="flex items-center justify-center h-64">
            <div className="w-8 h-8 border-4 border-brand border-t-transparent rounded-full animate-spin" />
          </div>
        )}

        {error && (
          <div className="p-4 bg-danger-soft text-danger border border-danger border-opacity-20 rounded-lg text-sm">
            {error}
          </div>
        )}

        {currentProject && !loadingScenario && !currentScenario && !error && (
          <Card className="text-center py-20">
            <p className="text-fg-secondary mb-4">
              Este proyecto no tiene escenarios calculados.
            </p>
            <Button onClick={() => navigate('/new')}>
              Recalcular
            </Button>
          </Card>
        )}

        {currentProject && currentScenario && !loadingScenario && (
          <ProjectDashboard
            project={currentProject}
            scenario={currentScenario}
            onAddToComparison={() => {
              addToComparison(currentScenario);
              navigate('/compare');
            }}
            onReport={() => navigate('/report')}
          />
        )}
      </div>
    </div>
  );
}

/* ─── Dashboard ─── */

function ProjectDashboard({
  project,
  scenario: s,
  onAddToComparison,
  onReport,
}: {
  project: Project;
  scenario: Scenario;
  onAddToComparison: () => void;
  onReport: () => void;
}) {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-fg-primary">{project.name}</h2>
          <p className="text-fg-secondary">
            {s.name} - Fuente: {s.irradiation.source}
          </p>
        </div>
        <div className="flex gap-3">
          <Button variant="secondary" size="sm" onClick={onAddToComparison}>
            Agregar a Comparacion
          </Button>
          <Button size="sm" onClick={onReport}>
            Generar Reporte
          </Button>
        </div>
      </div>

      {/* Summary cards */}
      <div className="grid grid-cols-4 gap-4">
        <SummaryCard
          label="Potencia Instalada"
          value={`${s.systemDesign.actualPowerKwp.toFixed(2)} kWp`}
          detail={`${s.systemDesign.numberOfPanels} paneles`}
          color="solar"
        />
        <SummaryCard
          label="Produccion Anual"
          value={`${formatNumber(s.production.annualKwh)} kWh`}
          detail={`HSP promedio: ${s.irradiation.annualAvgHSP.toFixed(1)}`}
          color="green"
        />
        <SummaryCard
          label="Costo Total"
          value={formatCOP(s.financial.installationCostCOP)}
          detail={`Payback: ${s.financial.paybackYears != null ? s.financial.paybackYears.toFixed(1) : 'N/A'} anos`}
          color="blue"
        />
        <SummaryCard
          label="TIR"
          value={`${s.financial.irrPercent.toFixed(1)}%`}
          detail={`VPN: ${formatCOP(s.financial.npvCOP)}`}
          color="purple"
        />
      </div>

      {/* Additional metrics */}
      <div className="grid grid-cols-5 gap-3">
        <MiniCard label="LCOE" value={`${formatCOP(s.financial.lcoe)}/kWh`} />
        <MiniCard label="CO2 Evitado" value={`${s.financial.co2AvoidedTonsYear.toFixed(2)} ton/ano`} />
        <MiniCard label="Uso del Techo" value={`${s.systemDesign.roofUtilization.toFixed(1)}%`} />
        <MiniCard label="Inversor" value={`${s.systemDesign.inverterCapacityKw} kW`} />
        <MiniCard label="Perdida Total" value={`${s.losses.totalSystemLoss.toFixed(1)}%`} />
      </div>

      {/* String configuration */}
      <div className="bg-surface p-4 rounded-lg border border-[var(--color-border-default)]">
        <h4 className="font-semibold text-fg-secondary mb-3">Configuracion del Array</h4>
        <div className="grid grid-cols-4 gap-4 text-sm">
          <div>
            <span className="text-fg-tertiary">Paneles por String:</span>
            <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.stringConfiguration.panelsPerString}</span>
          </div>
          <div>
            <span className="text-fg-tertiary">Numero de Strings:</span>
            <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.stringConfiguration.numberOfStrings}</span>
          </div>
          <div>
            <span className="text-fg-tertiary">Voltaje String:</span>
            <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.stringConfiguration.stringVoltage}V</span>
          </div>
          <div>
            <span className="text-fg-tertiary">Corriente String:</span>
            <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.stringConfiguration.stringCurrent}A</span>
          </div>
        </div>
      </div>

      {/* Battery (if applicable) */}
      {s.systemDesign.batteryBank && (
        <div className="bg-surface p-4 rounded-lg border border-[var(--color-border-default)]">
          <h4 className="font-semibold text-fg-secondary mb-3">Banco de Baterias</h4>
          <div className="grid grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-fg-tertiary">Capacidad:</span>
              <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.batteryBank.capacityKwh.toFixed(1)} kWh</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Autonomia:</span>
              <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.batteryBank.autonomyDays} dias</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Baterias:</span>
              <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.batteryBank.numberOfBatteries}</span>
            </div>
            <div>
              <span className="text-fg-tertiary">Voltaje Banco:</span>
              <span className="ml-2 font-mono tabular-nums text-fg-primary font-medium">{s.systemDesign.batteryBank.bankVoltage}V</span>
            </div>
          </div>
        </div>
      )}

      {/* Charts */}
      <div className="grid grid-cols-2 gap-4">
        <GenerationChart scenario={s} monthlyConsumption={project.consumption.monthly} />
        <IrradiationChart scenario={s} />
        <CashFlowChart scenario={s} />
        <DegradationChart scenario={s} />
      </div>

      {/* Losses breakdown */}
      <div className="bg-surface p-4 rounded-lg border border-[var(--color-border-default)]">
        <h4 className="font-semibold text-fg-secondary mb-3">Desglose de Perdidas</h4>
        <div className="space-y-2">
          {[
            { label: 'Sombreado', value: s.losses.shadingPercent },
            { label: 'Temperatura', value: s.losses.temperaturePercent },
            { label: 'Cableado', value: s.losses.wiringPercent },
            { label: 'Inversor', value: s.losses.inverterPercent },
            { label: 'Suciedad', value: s.losses.soilingPercent },
          ].map(({ label, value }) => (
            <div key={label} className="flex items-center gap-3">
              <span className="text-sm text-fg-tertiary w-32">{label}</span>
              <div className="flex-1 bg-[var(--color-border-default)] h-2 rounded-full">
                <div
                  className="bg-danger h-2 rounded-full"
                  style={{ width: `${Math.min(value * 3, 100)}%` }}
                />
              </div>
              <span className="text-sm font-mono tabular-nums font-medium w-16 text-right text-fg-primary">{value.toFixed(1)}%</span>
            </div>
          ))}
          <div className="flex items-center gap-3 border-t border-[var(--color-border-default)] pt-2 mt-2">
            <span className="text-sm font-semibold text-fg-primary w-32">Total</span>
            <div className="flex-1" />
            <span className="text-sm font-bold font-mono tabular-nums w-16 text-right text-fg-primary">{s.losses.totalSystemLoss.toFixed(1)}%</span>
          </div>
        </div>
      </div>
    </div>
  );
}

/* ─── Small components ─── */

function SummaryCard({
  label,
  value,
  detail,
  color,
}: {
  label: string;
  value: string;
  detail: string;
  color: string;
}) {
  const colors: Record<string, string> = {
    solar: 'border-t-brand bg-brand-soft',
    green: 'border-t-success bg-success-soft',
    blue: 'border-t-info bg-info-soft',
    purple: 'border-t-[#a855f7] bg-[rgba(168,85,247,0.1)]',
  };

  return (
    <div className={`border-t-[3px] border border-[var(--color-border-default)] rounded-lg p-4 ${colors[color] || colors.solar}`}>
      <p className="text-xs text-fg-tertiary uppercase tracking-wider">{label}</p>
      <p className="text-xl font-bold font-mono tabular-nums text-fg-primary mt-1">{value}</p>
      <p className="text-xs text-fg-tertiary mt-1">{detail}</p>
    </div>
  );
}

function MiniCard({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-surface border border-[var(--color-border-default)] rounded-lg p-3 text-center">
      <p className="text-xs text-fg-tertiary">{label}</p>
      <p className="text-sm font-semibold font-mono tabular-nums text-fg-primary mt-1">{value}</p>
    </div>
  );
}
