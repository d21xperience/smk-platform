import { api } from '@/boot/axios';

export const InventoryApi = {
  getAssets: (params) => api.get('/inventory/assets', { params }),
  getAsset: (assetId) => api.get(`/inventory/assets/${assetId}`),
  getIssues: (assetId) => api.get('/inventory/issues', { params: { assetId } }),
  getIssue: (issueId) => api.get(`/inventory/issues/${issueId}`),
  borrowAsset: (assetId, data) => api.post(`/inventory/assets/${assetId}/borrow`, data),
  returnAsset: (issueId) => api.post(`/inventory/issues/${issueId}/return`),
  reportDamage: (assetId, data) => api.post(`/inventory/assets/${assetId}/damage`, data),
};
