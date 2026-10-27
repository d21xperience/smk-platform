import { mockSemesters } from './data/context.mock.js';

export const ContextMockAdapter = {
  async getAvailablePeriods(context) {
    if (!context.schoolId) throw new Error('schoolId required');
    await new Promise(resolve => setTimeout(resolve, 300));
    return mockSemesters;
  },
  async getUserPreference(context) {
    if (!context.userId || !context.schoolId) throw new Error('userId and schoolId required');
    await new Promise(resolve => setTimeout(resolve, 200));
    return { semester_id: '20261' }; 
  },
  async saveUserPreference(context, semesterId) {
    if (!context.userId || !context.schoolId) throw new Error('userId and schoolId required');
    await new Promise(resolve => setTimeout(resolve, 200));
    return { success: true, semester_id: semesterId };
  }
};
