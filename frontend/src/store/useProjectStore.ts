import { create } from 'zustand';
import { Project, Scenario, Panel, Inverter, IrradiationResult, SystemType, PanelFormat } from '../types';

interface ProjectFormData {
  name: string;
  location: {
    latitude: number;
    longitude: number;
    altitude: number;
    climateZone: string;
    department: string;
    city: string;
  };
  consumption: {
    monthly: number[];
    tariffPerKwh: number;
    estrato: number;
    connectionType: 'monofasica' | 'bifasica' | 'trifasica';
  };
  roof: {
    area: number;
    azimuth: number;
    tilt: number;
    usablePercentage: number;
    shadingProfile: {
      hasShading: boolean;
      monthlyLoss: number[];
    };
  };
  systemType: SystemType;
  coveragePercentage: number;
  panelFormat: PanelFormat;
  recommendedInverterKw: number;
  equipment: {
    panelId: string;
    inverterId: string;
  };
}

interface ProjectStore {
  // Projects list
  projects: Project[];
  setProjects: (projects: Project[]) => void;

  // Current project being edited
  currentProject: Project | null;
  setCurrentProject: (project: Project | null) => void;

  // Form wizard state
  formData: ProjectFormData;
  currentStep: number;
  setFormData: (data: Partial<ProjectFormData>) => void;
  setCurrentStep: (step: number) => void;
  resetForm: () => void;

  // Catalogs
  panels: Panel[];
  inverters: Inverter[];
  setPanels: (panels: Panel[]) => void;
  setInverters: (inverters: Inverter[]) => void;

  // Irradiation preview
  irradiationPreview: IrradiationResult | null;
  setIrradiationPreview: (data: IrradiationResult | null) => void;

  // Scenarios & results
  currentScenario: Scenario | null;
  scenarios: Scenario[];
  setCurrentScenario: (scenario: Scenario | null) => void;
  setScenarios: (scenarios: Scenario[]) => void;

  // Comparison
  comparisonScenarios: Scenario[];
  addToComparison: (scenario: Scenario) => void;
  removeFromComparison: (id: string) => void;
  clearComparison: () => void;

  // Loading states
  isLoading: boolean;
  setIsLoading: (loading: boolean) => void;
  error: string | null;
  setError: (error: string | null) => void;
}

const defaultFormData: ProjectFormData = {
  name: '',
  location: {
    latitude: 4.61,
    longitude: -74.08,
    altitude: 2640,
    climateZone: '',
    department: 'Cundinamarca',
    city: 'Bogotá',
  },
  consumption: {
    monthly: new Array(12).fill(200),
    tariffPerKwh: 800,
    estrato: 4,
    connectionType: 'monofasica',
  },
  roof: {
    area: 50,
    azimuth: 0,
    tilt: 10,
    usablePercentage: 80,
    shadingProfile: {
      hasShading: false,
      monthlyLoss: new Array(12).fill(0),
    },
  },
  systemType: 'on-grid',
  coveragePercentage: 100,
  panelFormat: 'standard',
  recommendedInverterKw: 0,
  equipment: {
    panelId: '',
    inverterId: '',
  },
};

export const useProjectStore = create<ProjectStore>((set) => ({
  projects: [],
  setProjects: (projects) => set({ projects }),

  currentProject: null,
  setCurrentProject: (project) => set({ currentProject: project }),

  formData: { ...defaultFormData },
  currentStep: 0,
  setFormData: (data) =>
    set((state) => ({
      formData: { ...state.formData, ...data },
    })),
  setCurrentStep: (step) => set({ currentStep: step }),
  resetForm: () => set({ formData: { ...defaultFormData }, currentStep: 0 }),

  panels: [],
  inverters: [],
  setPanels: (panels) => set({ panels }),
  setInverters: (inverters) => set({ inverters }),

  irradiationPreview: null,
  setIrradiationPreview: (data) => set({ irradiationPreview: data }),

  currentScenario: null,
  scenarios: [],
  setCurrentScenario: (scenario) => set({ currentScenario: scenario }),
  setScenarios: (scenarios) => set({ scenarios }),

  comparisonScenarios: [],
  addToComparison: (scenario) =>
    set((state) => ({
      comparisonScenarios: state.comparisonScenarios.find((s) => s._id === scenario._id)
        ? state.comparisonScenarios
        : [...state.comparisonScenarios, scenario],
    })),
  removeFromComparison: (id) =>
    set((state) => ({
      comparisonScenarios: state.comparisonScenarios.filter((s) => s._id !== id),
    })),
  clearComparison: () => set({ comparisonScenarios: [] }),

  isLoading: false,
  setIsLoading: (loading) => set({ isLoading: loading }),
  error: null,
  setError: (error) => set({ error }),
}));
