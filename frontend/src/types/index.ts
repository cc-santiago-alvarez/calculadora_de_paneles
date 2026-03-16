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

export interface Roof {
  area: number;
  azimuth: number;
  tilt: number;
  usablePercentage: number;
  shadingProfile: ShadingProfile;
}

export interface Equipment {
  panelId: string;
  inverterId: string;
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
  costCOP: number;
  dimensions: { length: number; width: number };
}

export interface Inverter {
  _id: string;
  manufacturer: string;
  model: string;
  type: string;
  ratedPowerKw: number;
  efficiency: number;
  mpptCount: number;
  mpptVoltageMin: number;
  mpptVoltageMax: number;
  maxInputVoltage: number;
  maxInputCurrent: number;
  hasBatteryPort: boolean;
  costCOP: number;
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
