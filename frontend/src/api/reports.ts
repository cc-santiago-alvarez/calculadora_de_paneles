import api from './client';

export const reportsApi = {
  generatePDF: (projectId: string, scenarioId: string) =>
    api.post('/reports/pdf', { projectId, scenarioId }, { responseType: 'blob' }).then((r) => r.data),

  generateExcel: (projectId: string, scenarioId: string) =>
    api.post('/reports/excel', { projectId, scenarioId }, { responseType: 'blob' }).then((r) => r.data),
};
