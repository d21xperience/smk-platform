import { StudentAffairsMockAdapter } from '@/adapters/mock/StudentAffairsMockAdapter'
import { StudentAffairsEngine } from '@/domain/student-affairs/engines/StudentAffairsEngine'

// TODO: Ganti dengan Adapter Registry saat Phase 3 (VITE_USE_MOCK)
const adapter = new StudentAffairsMockAdapter()

export class StudentAffairsService {
  static async getDashboardMetrics() {
    return await adapter.getDashboardMetrics()
  }

  static async getStudents(filters) {
    return await adapter.getStudents(filters)
  }

  static async getMutations(statusFilter) {
    return await adapter.getMutations(statusFilter)
  }

  static async submitMutation(data) {
    StudentAffairsEngine.validateMutationSubmission(data)
    return await adapter.createMutation(data)
  }

  static async approveMutation(id, catatan) {
    // 1. Approve mutation di adapter
    const approvedMutation = await adapter.approveMutation(id, catatan)

    // 2. Jalankan business rule dari Domain Engine
    const command = StudentAffairsEngine.processMutationApproval(approvedMutation)

    // 3. Eksekusi side-effect berdasarkan command
    if (command.action === 'CREATE_STUDENT') {
      await adapter.createStudent(command.payload)
    } else if (command.action === 'UPDATE_STUDENT_STATUS') {
      await adapter.updateStudent(command.payload.studentId, {
        status: command.payload.newStatus,
        alasanPindah: command.payload.alasanPindah,
        tujuanSekolah: command.payload.tujuanSekolah,
      })
    }

    return approvedMutation
  }

  static async rejectMutation(id, catatan) {
    return await adapter.rejectMutation(id, catatan)
  }
}
