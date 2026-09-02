import { api } from '@/boot/axios';

export const StudentProgressApi = {
  getClassProgress: (params) => api.get('/student-progress/class', { params }),
  getStudentProgress: (params) => api.get('/student-progress/student', { params }),
  addTeacherNote: (data) => api.post('/student-progress/note', data),
  addAchievement: (data) => api.post('/student-progress/achievement', data),
  addViolation: (data) => api.post('/student-progress/violation', data),
  syncFromAssessment: (data) => api.post('/student-progress/sync', data),
};
