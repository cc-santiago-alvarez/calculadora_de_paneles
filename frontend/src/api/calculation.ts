import api from './client';
import { IrradiationResult, Scenario } from '../types';

export const calculationApi = {
  fetchIrradiation: (data: { latitude: number; longitude: number; tilt: number; azimuth: number }) =>
    api.post<IrradiationResult>('/irradiation/fetch', data).then((r) => r.data),

  fullCalculation: (data: { projectId: string; scenarioName?: string }) =>
    api.post<{ scenario: Scenario; warnings: string[] }>('/calculation/full', data).then((r) => r.data),

  analyzeFinancial: (data: Record<string, unknown>) =>
    api.post('/financial/analyze', data).then((r) => r.data),
};
