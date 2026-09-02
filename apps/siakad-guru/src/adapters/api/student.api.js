// apps/siakad-guru/src/adapters/api/student.api.js
import { api } from '@/boot/axios'

export const StudentApi = {
  getStudentsByClass: (classId) => api.get(`/students/class/${classId}`),
  getStudent: (studentId) => api.get(`/students/${studentId}`),
}
