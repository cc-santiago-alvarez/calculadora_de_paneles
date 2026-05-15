import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProjectStore } from '../../store/useProjectStore';
import { projectsApi } from '../../api/projects';
import { calculationApi } from '../../api/calculation';
import LocationStep from './LocationStep';
import ConsumptionStep from './ConsumptionStep';
import RoofStep from './RoofStep';
import SystemTypeStep from './SystemTypeStep';
import CoverageStep from './CoverageStep';
import PanelConfigStep from './PanelConfigStep';
import EquipmentStep from './EquipmentStep';
import ErrorBanner from '../../components/common/ErrorBanner';
import LoadingOverlay from '../../components/common/LoadingOverlay';
import StepIndicator from '../../components/common/StepIndicator';
import Button from '../../components/common/Button';
import StepSidebar from './StepSidebar';

const STEPS = [
  { label: 'Consumo', component: ConsumptionStep },
  { label: 'Ubicacion', component: LocationStep },
  { label: 'Cobertura', component: CoverageStep },
  { label: 'Sistema', component: SystemTypeStep },
  { label: 'Techo', component: RoofStep },
  { label: 'Configuracion', component: PanelConfigStep },
  { label: 'Equipos', component: EquipmentStep },
];

export default function ProjectFormPage() {
  const navigate = useNavigate();
  const {
    formData,
    currentStep,
    setCurrentStep,
    setCurrentProject,
    setCurrentScenario,
    setIsLoading,
    isLoading,
    error,
    setError,
    resetForm,
  } = useProjectStore();

  // Reset form on mount so "Nuevo Proyecto" always starts fresh
  useEffect(() => {
    resetForm();
  }, []);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const StepComponent = STEPS[currentStep].component;

  const canGoNext = () => {
    switch (currentStep) {
      case 0: // Consumo - requiere nombre y al menos un mes con consumo
        return formData.name.trim().length > 0 && formData.consumption.monthly.some((v) => v > 0);
      case 1: // Ubicacion
        return true;
      case 2: // Cobertura
        return formData.coveragePercentage > 0 && formData.coveragePercentage <= 100;
      case 3: // Sistema
        return true;
      case 4: // Techo
        if (formData.roof.roofType === 'plana') {
          return formData.roof.area > 0;
        }
        return formData.roof.slopes.length > 0 && formData.roof.slopes.every((s) => s.area > 0);
      case 5: { // Configuracion de paneles
        const cfg = formData.equipment.panelConfiguration;
        return !!formData.panelFormat && !!formData.equipment.panelId && cfg.panelsPerString >= 1 && cfg.numberOfStrings >= 1;
      }
      case 6: { // Equipos
        if (!formData.equipment.inverterId) return false;
        // For off-grid with inverter without built-in MPPT, charge controller is required
        if (formData.systemType === 'off-grid') {
          const { inverters } = useProjectStore.getState();
          const selectedInv = inverters.find((i) => i._id === formData.equipment.inverterId);
          if (selectedInv && selectedInv.mpptCount === 0 && !formData.equipment.chargeControllerId) {
            return false;
          }
        }
        return true;
      }
      default:
        return true;
    }
  };

  const handleSubmit = async () => {
    try {
      setIsSubmitting(true);
      setIsLoading(true);
      setError(null);

      // Create project
      const project = await projectsApi.create(formData as any);
      setCurrentProject(project);

      // Run full calculation
      const result = await calculationApi.fullCalculation({
        projectId: project._id,
        scenarioName: 'Escenario Principal',
      });

      setCurrentScenario(result.scenario);
      navigate('/results');
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Error al crear el proyecto');
    } finally {
      setIsSubmitting(false);
      setIsLoading(false);
    }
  };

  const hasSidebar = [0, 1, 2, 5, 6].includes(currentStep);

  return (
    <div className="max-w-6xl mx-auto">
      {/* Step indicator */}
      <div className="mb-8">
        <StepIndicator steps={STEPS} currentStep={currentStep} />
      </div>

      {error && <ErrorBanner message={error} onDismiss={() => setError(null)} />}

      <div className={`flex gap-6 ${hasSidebar ? '' : 'max-w-4xl mx-auto'}`}>
        {/* Step content */}
        <div className="flex-1 min-w-0">
          <div className="bg-surface rounded-lg border border-[var(--color-border-default)] p-6">
            <StepComponent />
          </div>

          {/* Navigation */}
          <div className="flex justify-between mt-6">
            <Button
              variant="secondary"
              size="lg"
              onClick={() => setCurrentStep(currentStep - 1)}
              disabled={currentStep === 0}
            >
              Anterior
            </Button>

            {currentStep < STEPS.length - 1 ? (
              <Button
                size="lg"
                onClick={() => setCurrentStep(currentStep + 1)}
                disabled={!canGoNext()}
              >
                Siguiente
              </Button>
            ) : (
              <Button
                size="lg"
                className="bg-brand hover:bg-brand-hover text-white"
                onClick={handleSubmit}
                disabled={!canGoNext() || isSubmitting}
              >
                {isSubmitting ? 'Calculando...' : 'Calcular Sistema'}
              </Button>
            )}
          </div>
        </div>

        {/* Sidebar */}
        <StepSidebar step={currentStep} />
      </div>

      {isSubmitting && <LoadingOverlay message="Calculando dimensionamiento del sistema..." />}
    </div>
  );
}
