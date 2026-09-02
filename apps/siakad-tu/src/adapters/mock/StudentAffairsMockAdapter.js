import { StudentAffairsAdapterContract } from '@/contracts/student-affairs/StudentAffairsAdapterContract'
import {
  mockStudents,
  mockMutations,
  mockDashboard,
} from '@/adapters/mock/data/student-affairs.mock'

export class StudentAffairsMockAdapter extends StudentAffairsAdapterContract {
  constructor() {
    super()
    this.students = [...mockStudents]
    this.mutations = [...mockMutations]
    this.dashboard = { ...mockDashboard }
  }

  async _delay(ms = 500) {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }

  async getDashboardMetrics() {
    await this._delay()
    return this.dashboard
  }

  async getStudents(filters = {}) {
    await this._delay()
    let result = [...this.students]

    if (filters.search) {
      const q = filters.search.toLowerCase()
      result = result.filter(
        (s) => s.nama.toLowerCase().includes(q) || s.nis.includes(q) || s.nisn.includes(q),
      )
    }
    if (filters.status) {
      result = result.filter((s) => s.status === filters.status)
    }
    if (filters.jurusan && filters.jurusan !== 'Semua') {
      result = result.filter((s) => s.jurusan === filters.jurusan)
    }
    if (filters.tingkat && filters.tingkat !== 'Semua') {
      result = result.filter((s) => s.kelas.startsWith(filters.tingkat))
    }

    return {
      data: result,
      total: result.length,
      page: 1,
      totalPages: 1,
    }
  }

  async getMutations(statusFilter = 'semua') {
    await this._delay()
    let result = [...this.mutations]
    if (statusFilter !== 'semua') {
      result = result.filter((m) => m.status.toLowerCase() === statusFilter.toLowerCase())
    }
    return result
  }

  async createMutation(data) {
    await this._delay()
    const newMutation = {
      id: `M${Date.now()}`,
      schoolId: 'SCH001',
      academicYearId: 'AY2023',
      tanggalPengajuan: new Date().toISOString(),
      status: 'PENDING',
      ...data,
    }
    this.mutations.unshift(newMutation)
    return newMutation
  }

  async approveMutation(id, catatan) {
    await this._delay()
    const mutation = this.mutations.find((m) => m.id === id)
    if (!mutation) throw new Error('Mutasi tidak ditemukan')

    mutation.status = 'DISETUJUI'
    mutation.catatanVerifikasi = catatan
    return mutation
  }

  async rejectMutation(id, catatan) {
    await this._delay()
    const mutation = this.mutations.find((m) => m.id === id)
    if (!mutation) throw new Error('Mutasi tidak ditemukan')

    mutation.status = 'DITOLAK'
    mutation.catatanVerifikasi = catatan
    return mutation
  }

  async createStudent(data) {
    await this._delay()
    const newStudent = {
      id: `S${Date.now()}`,
      schoolId: 'SCH001',
      academicYearId: 'AY2023',
      isLocked: false,
      createdAt: new Date().toISOString(),
      ...data,
    }
    this.students.push(newStudent)
    return newStudent
  }

  async updateStudent(id, data) {
    await this._delay()
    const index = this.students.findIndex((s) => s.id === id)
    if (index === -1) throw new Error('Siswa tidak ditemukan')

    this.students[index] = { ...this.students[index], ...data, updatedAt: new Date().toISOString() }
    return this.students[index]
  }
}
