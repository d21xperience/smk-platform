import { api } from '@/boot/axios';

export const FinanceApi = {
  getClassFinancialSummary: (params) => api.get('/finance/class-summary', { params }),
  getStudentFinancialDetail: (params) => api.get('/finance/student-detail', { params }),
};
