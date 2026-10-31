import { ContextEngine } from '../engine/ContextEngine.js';
import { ContextApiAdapter } from '../adapters/api/ContextApiAdapter.js';

const adapter = ContextApiAdapter;

export const ContextService = {
  async fetchPeriods(context) {
    ContextEngine.validateContext(context);
    return await adapter.getAvailablePeriods(context);
  },
  async fetchPreference(context) {
    ContextEngine.validateContext(context);
    return await adapter.getUserPreference(context);
  },
  async updatePreference(context, semesterId) {
    ContextEngine.validateContext(context);
    if (!semesterId) throw new Error('semesterId is required');
    return await adapter.saveUserPreference(context, semesterId);
  }
};
