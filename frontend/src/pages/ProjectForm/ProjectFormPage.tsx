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
import EquipmentStep from './EquipmentStep';
import ErrorBanner from '../../components/common/ErrorBanner';
import LoadingOverlay from '../../components/common/LoadingOverlay';
import StepIndicator from '../../components/common/StepIndicator';
import Button from '../../components/common/Button';
import StepSidebar from './StepSidebar';

const STEPS = [
  { label: 'Consumo', component: ConsumptionStep },
  { label: 'Ubicacion', component: LocationStep },
  { label: 'Techo', component: RoofStep },
  { label: 'Sistema', component: SystemTypeStep },
  { label: 'Cobertura', component: CoverageStep },
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
      case 2: // Techo
        return formData.roof.area > 0;
      case 3: // Sistema
        return true;
      case 4: // Cobertura
        return formData.coveragePercentage > 0 && formData.coveragePercentage <= 100 && !!formData.panelFormat;
      case 5: // Equipos
        return formData.equipment.panelId && formData.equipment.inverterId;
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

  const hasSidebar = [0, 1, 4, 5].includes(currentStep);

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
