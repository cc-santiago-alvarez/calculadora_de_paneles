import api from './client';
import { Panel, Inverter, ChargeController } from '../types';

export const catalogApi = {
  getPanels: (params?: Record<string, string>) =>
    api.get<Panel[]>('/catalog/panels', { params }).then((r) => r.data),
  getInverters: (params?: Record<string, string>) =>
    api.get<Inverter[]>('/catalog/inverters', { params }).then((r) => r.data),
  getChargeControllers: (params?: Record<string, string>) =>
    api.get<ChargeController[]>('/catalog/charge-controllers', { params }).then((r) => r.data),
  createPanel: (data: Partial<Panel>) =>
    api.post<Panel>('/catalog/panels', data).then((r) => r.data),
  updatePanel: (id: string, data: Partial<Panel>) =>
    api.put<Panel>(`/catalog/panels/${id}`, data).then((r) => r.data),
  deletePanel: (id: string) =>
    api.delete(`/catalog/panels/${id}`).then((r) => r.data),
  createInverter: (data: Partial<Inverter>) =>
    api.post<Inverter>('/catalog/inverters', data).then((r) => r.data),
  updateInverter: (id: string, data: Partial<Inverter>) =>
    api.put<Inverter>(`/catalog/inverters/${id}`, data).then((r) => r.data),
  deleteInverter: (id: string) =>
    api.delete(`/catalog/inverters/${id}`).then((r) => r.data),
};
