import api from './client';
import { Project, Scenario } from '../types';

export const projectsApi = {
  getAll: () => api.get<Project[]>('/projects').then((r) => r.data),
  getById: (id: string) => api.get<Project>(`/projects/${id}`).then((r) => r.data),
  create: (data: Omit<Project, '_id' | 'createdAt' | 'updatedAt' | 'scenarios'>) =>
    api.post<Project>('/projects', data).then((r) => r.data),
  update: (id: string, data: Partial<Project>) =>
    api.put<Project>(`/projects/${id}`, data).then((r) => r.data),
  delete: (id: string) => api.delete(`/projects/${id}`).then((r) => r.data),
  getScenarios: (id: string) =>
    api.get<Scenario[]>(`/projects/${id}/scenarios`).then((r) => r.data),
};
