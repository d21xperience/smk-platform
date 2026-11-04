// apps/siakad-tu/src/services/StudentAffairsService.js
import { studentMockAdapter } from '@/adapters/mock/studentMockAdapter.js'

// Fallback kompatibilitas mundur
const adapter = studentMockAdapter 

export const StudentAffairsService = {
  getStudents: (context, filters) => adapter.getStudents(context, filters),
  createStudent: (data) => adapter.registerStudent(data, {}),
  // Tambahkan method lain sesuai kebutuhan UI lama jika ada
}
