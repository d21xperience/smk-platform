// apps/siakad-guru/src/adapters/api/context.api.js
import { api } from '@/boot/axios.js';

export const ContextApi = {
  getContext: () => api.get('/context'),
  switchContext: (payload) => api.post('/context/switch', payload),
  getAvailableContexts: (userId) => api.get(`/context/available?userId=${userId}`),
  selectContext: (payload) => api.post('/context/select', payload),
};
