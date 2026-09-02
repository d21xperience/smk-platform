import { AcademicYear } from '@/models/AcademicYear'
import { Semester } from '@/models/Semester'
import { AcademicPeriod } from '@/models/AcademicPeriod'

// Data statis
const academicYears = [
  new AcademicYear({ id: 1, name: '2026/2027', isActive: true }),
  new AcademicYear({ id: 2, name: '2025/2026', isActive: false }),
]

const semesters = [
  new Semester({ id: 1, academicYearId: 1, name: 'Semester 1', type: 'ganjil' }),
  new Semester({ id: 2, academicYearId: 1, name: 'Semester 2', type: 'genap' }),
  new Semester({ id: 3, academicYearId: 2, name: 'Semester 1', type: 'ganjil' }),
  new Semester({ id: 4, academicYearId: 2, name: 'Semester 2', type: 'genap' }),
]

const academicPeriods = [
  new AcademicPeriod({
    id: 1,
    semesterId: 1,
    name: 'Tengah Semester Ganjil',
    startDate: '2026-07-15',
    endDate: '2026-09-30',
  }),
  new AcademicPeriod({
    id: 2,
    semesterId: 1,
    name: 'Akhir Semester Ganjil',
    startDate: '2026-10-01',
    endDate: '2026-12-20',
  }),
  new AcademicPeriod({
    id: 3,
    semesterId: 2,
    name: 'Tengah Semester Genap',
    startDate: '2027-01-05',
    endDate: '2027-03-10',
  }),
  new AcademicPeriod({
    id: 4,
    semesterId: 2,
    name: 'Akhir Semester Genap',
    startDate: '2027-03-11',
    endDate: '2027-06-25',
  }),
  new AcademicPeriod({
    id: 5,
    semesterId: 3,
    name: 'Tengah Semester Ganjil',
    startDate: '2025-07-10',
    endDate: '2025-10-05',
  }),
  new AcademicPeriod({
    id: 6,
    semesterId: 4,
    name: 'Akhir Semester Genap',
    startDate: '2026-03-01',
    endDate: '2026-06-20',
  }),
]

export const academicMockAdapter = {
  /**
   * Mendapatkan daftar semua tahun ajaran.
   * @returns {Promise<AcademicYear[]>}
   */
  async getAcademicYears() {
    await new Promise((resolve) => setTimeout(resolve, 300))
    return [...academicYears] // kembalikan salinan
  },

  /**
   * Mendapatkan daftar semester berdasarkan tahun ajaran.
   * @param {number} academicYearId
   * @returns {Promise<Semester[]>}
   */
  async getSemesters(academicYearId) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    return semesters.filter((s) => s.academicYearId === academicYearId?.id)
  },

  /**
   * Mendapatkan daftar periode akademik berdasarkan semester.
   * @param {number} semesterId
   * @returns {Promise<AcademicPeriod[]>}
   */
  async getAcademicPeriods(semesterId) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    return academicPeriods.filter((p) => p.semesterId === semesterId?.id)
  },
}
