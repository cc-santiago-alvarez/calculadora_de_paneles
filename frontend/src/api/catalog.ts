import api from './client';
import { Panel, Inverter } from '../types';

export const catalogApi = {
  getPanels: (params?: Record<string, string>) =>
    api.get<Panel[]>('/catalog/panels', { params }).then((r) => r.data),
  getInverters: (params?: Record<string, string>) =>
    api.get<Inverter[]>('/catalog/inverters', { params }).then((r) => r.data),
};
