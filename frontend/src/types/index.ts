export interface Location {
  latitude: number;
  longitude: number;
  altitude: number;
  climateZone: string;
  department: string;
  city: string;
}

export interface Consumption {
  monthly: number[];
  tariffPerKwh: number;
  estrato: number;
  connectionType: 'monofasica' | 'bifasica' | 'trifasica';
}

export interface ShadingProfile {
  hasShading: boolean;
  monthlyLoss: number[];
}

export type RoofType = 'plana' | 'una_agua' | 'dos_aguas' | 'cuatro_aguas';

export interface Slope {
  area: number;
  tilt: number;
  azimuth: number;
}

export interface Roof {
  roofType: RoofType;
  area: number;
  azimuth: number;
  tilt: number;
  slopes: Slope[];
  usablePercentage: number;
  shadingProfile: ShadingProfile;
}

export type PanelConnectionType = 'serie' | 'paralelo' | 'mixto';

export interface PanelConfiguration {
  connectionType: PanelConnectionType;
  panelsPerString: number; // S - paneles en serie por cadena
  numberOfStrings: number; // P - cadenas en paralelo
}

export interface Equipment {
  panelId: string;
  inverterId: string;
  chargeControllerId?: string;
  panelConfiguration: PanelConfiguration;
  panelOverride?: {
    watts: number;
    area: number;
  };
}

export type SystemType = 'on-grid' | 'off-grid' | 'hybrid';

export type PanelFormat = 'standard' | 'high-efficiency' | 'large-format';

export interface Project {
  _id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  location: Location;
  consumption: Consumption;
  roof: Roof;
  systemType: SystemType;
  coveragePercentage: number;
  panelFormat: PanelFormat;
  equipment: Equipment;
  scenarios: string[];
}

export interface Panel {
  _id: string;
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
  weight?: number;
  warranty?: number;
  costCOP: number;
  format?: string;
  isActive?: boolean;
  dimensions: { length: number; width: number };
}

export interface Inverter {
  _id: string;
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
  weight?: number;
  warranty?: number;
  costCOP: number;
  isActive?: boolean;
}

export interface ChargeController {
  _id: string;
  manufacturer: string;
  model: string;
  type: string;
  ratedPowerW: number;
  maxPVVoltage: number;
  maxPVCurrent: number;
  batteryVoltages: number[];
  maxChargeCurrentA: number;
  efficiency: number;
  costCOP: number;
  warranty: number;
  weight: number;
}

export interface Scenario {
  _id: string;
  projectId: string;
  name: string;
  createdAt: string;
  irradiation: {
    source: string;
    monthlyGHI: number[];
    monthlyPOA: number[];
    annualAvgHSP: number;
  };
  systemDesign: {
    requiredPowerKwp: number;
    numberOfPanels: number;
    actualPowerKwp: number;
    roofUtilization: number;
    inverterCapacityKw: number;
    stringConfiguration: {
      panelsPerString: number;
      numberOfStrings: number;
      stringVoltage: number;
      stringCurrent: number;
    };
    batteryBank?: {
      capacityKwh: number;
      autonomyDays: number;
      numberOfBatteries: number;
      bankVoltage: number;
    };
  };
  production: {
    monthlyKwh: number[];
    annualKwh: number;
    degradationRate: number;
    yearly25: number[];
  };
  financial: {
    installationCostCOP: number;
    monthlySavingsCOP: number[];
    annualSavingsCOP: number;
    paybackYears: number | null;
    irrPercent: number;
    npvCOP: number;
    co2AvoidedTonsYear: number;
    cumulativeSavings25: number[];
    lcoe: number;
  };
  losses: {
    shadingPercent: number;
    temperaturePercent: number;
    wiringPercent: number;
    inverterPercent: number;
    soilingPercent: number;
    totalSystemLoss: number;
  };
}

export interface IrradiationResult {
  source: string;
  monthlyGHI: number[];
  monthlyPOA: number[];
  annualAvgHSP: number;
  elevation: number;
}
