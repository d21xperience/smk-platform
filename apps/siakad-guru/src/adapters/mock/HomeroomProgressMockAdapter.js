// FILE: src/adapters/mock/HomeroomProgressMockAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Progress Mock Adapter
 * Menyediakan data mentah nilai siswa per kelas untuk diproses Domain Engine.
 * Data terikat pada userId (wali kelas), schoolId, dan academicYearId.
 */

const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' },
}

const MOCK_SUBJECTS = [
  { subjectId: 'SUB-001', subjectName: 'Matematika', kkm: 75 },
  { subjectId: 'SUB-002', subjectName: 'Bahasa Indonesia', kkm: 78 },
  { subjectId: 'SUB-003', subjectName: 'Bahasa Inggris', kkm: 75 },
  { subjectId: 'SUB-004', subjectName: 'IPA', kkm: 75 },
  { subjectId: 'SUB-005', subjectName: 'IPS', kkm: 78 },
  { subjectId: 'SUB-006', subjectName: 'Pendidikan Agama', kkm: 80 },
]

const MOCK_PROGRESS_DB = {
  'CLS-XI-B': {
    'AY-2026': [
      {
        studentId: 'STU-001',
        studentName: 'Ahmad Fauzi',
        seatNumber: 1,
        scores: {
          'SUB-001': { formative: 88, summative: 90, previousFinalScore: 82 },
          'SUB-002': { formative: 85, summative: 87, previousFinalScore: 80 },
          'SUB-003': { formative: 92, summative: 94, previousFinalScore: 88 },
          'SUB-004': { formative: 90, summative: 91, previousFinalScore: 85 },
          'SUB-005': { formative: 84, summative: 86, previousFinalScore: 79 },
          'SUB-006': { formative: 88, summative: 90, previousFinalScore: 84 },
        },
      },
      {
        studentId: 'STU-002',
        studentName: 'Siti Nurhaliza',
        seatNumber: 2,
        scores: {
          'SUB-001': { formative: 82, summative: 84, previousFinalScore: 78 },
          'SUB-002': { formative: 90, summative: 92, previousFinalScore: 85 },
          'SUB-003': { formative: 85, summative: 88, previousFinalScore: 80 },
          'SUB-004': { formative: 78, summative: 80, previousFinalScore: 75 },
          'SUB-005': { formative: 88, summative: 90, previousFinalScore: 82 },
          'SUB-006': { formative: 92, summative: 94, previousFinalScore: 88 },
        },
      },
      {
        studentId: 'STU-003',
        studentName: 'Budi Santoso',
        seatNumber: 3,
        scores: {
          'SUB-001': { formative: 72, summative: 74, previousFinalScore: 70 },
          'SUB-002': { formative: 78, summative: 80, previousFinalScore: 76 },
          'SUB-003': { formative: 70, summative: 72, previousFinalScore: 68 },
          'SUB-004': { formative: 75, summative: 76, previousFinalScore: 74 },
          'SUB-005': { formative: 80, summative: 82, previousFinalScore: 78 },
          'SUB-006': { formative: 82, summative: 84, previousFinalScore: 80 },
        },
      },
      {
        studentId: 'STU-004',
        studentName: 'Dewi Lestari',
        seatNumber: 4,
        scores: {
          'SUB-001': { formative: 80, summative: 82, previousFinalScore: 76 },
          'SUB-002': { formative: 84, summative: 86, previousFinalScore: 80 },
          'SUB-003': { formative: 82, summative: 85, previousFinalScore: 78 },
          'SUB-004': { formative: 86, summative: 88, previousFinalScore: 82 },
          'SUB-005': { formative: 78, summative: 80, previousFinalScore: 74 },
          'SUB-006': { formative: 85, summative: 88, previousFinalScore: 82 },
        },
      },
      {
        studentId: 'STU-005',
        studentName: 'Eko Prasetyo',
        seatNumber: 5,
        scores: {
          'SUB-001': { formative: 68, summative: 70, previousFinalScore: 72 },
          'SUB-002': { formative: 74, summative: 76, previousFinalScore: 78 },
          'SUB-003': { formative: 65, summative: 68, previousFinalScore: 70 },
          'SUB-004': { formative: 70, summative: 72, previousFinalScore: 74 },
          'SUB-005': { formative: 76, summative: 78, previousFinalScore: 80 },
          'SUB-006': { formative: 78, summative: 80, previousFinalScore: 82 },
        },
      },
      {
        studentId: 'STU-006',
        studentName: 'Fitri Handayani',
        seatNumber: 6,
        scores: {
          'SUB-001': { formative: 76, summative: 78, previousFinalScore: 74 },
          'SUB-002': { formative: 82, summative: 84, previousFinalScore: 78 },
          'SUB-003': { formative: 74, summative: 76, previousFinalScore: 72 },
          'SUB-004': { formative: 80, summative: 82, previousFinalScore: 76 },
          'SUB-005': { formative: 84, summative: 86, previousFinalScore: 80 },
          'SUB-006': { formative: 86, summative: 88, previousFinalScore: 82 },
        },
      },
      {
        studentId: 'STU-007',
        studentName: 'Gilang Ramadhan',
        seatNumber: 7,
        scores: {
          'SUB-001': { formative: 65, summative: 68, previousFinalScore: 70 },
          'SUB-002': { formative: 72, summative: 74, previousFinalScore: 76 },
          'SUB-003': { formative: 60, summative: 64, previousFinalScore: 68 },
          'SUB-004': { formative: 68, summative: 70, previousFinalScore: 72 },
          'SUB-005': { formative: 74, summative: 76, previousFinalScore: 78 },
          'SUB-006': { formative: 76, summative: 78, previousFinalScore: 80 },
        },
      },
      {
        studentId: 'STU-008',
        studentName: 'Hana Pertiwi',
        seatNumber: 8,
        scores: {
          'SUB-001': { formative: 92, summative: 94, previousFinalScore: 88 },
          'SUB-002': { formative: 88, summative: 90, previousFinalScore: 84 },
          'SUB-003': { formative: 90, summative: 92, previousFinalScore: 86 },
          'SUB-004': { formative: 94, summative: 95, previousFinalScore: 90 },
          'SUB-005': { formative: 86, summative: 88, previousFinalScore: 82 },
          'SUB-006': { formative: 90, summative: 92, previousFinalScore: 86 },
        },
      },
      {
        studentId: 'STU-009',
        studentName: 'Irfan Hakim',
        seatNumber: 9,
        scores: {
          'SUB-001': { formative: 74, summative: 76, previousFinalScore: 72 },
          'SUB-002': { formative: 76, summative: 78, previousFinalScore: 74 },
          'SUB-003': { formative: 72, summative: 74, previousFinalScore: 70 },
          'SUB-004': { formative: 78, summative: 80, previousFinalScore: 76 },
          'SUB-005': { formative: 80, summative: 82, previousFinalScore: 78 },
          'SUB-006': { formative: 80, summative: 82, previousFinalScore: 78 },
        },
      },
      {
        studentId: 'STU-010',
        studentName: 'Joko Susilo',
        seatNumber: 10,
        scores: {
          'SUB-001': { formative: 60, summative: 62, previousFinalScore: 65 },
          'SUB-002': { formative: 68, summative: 70, previousFinalScore: 72 },
          'SUB-003': { formative: 58, summative: 60, previousFinalScore: 64 },
          'SUB-004': { formative: 64, summative: 66, previousFinalScore: 68 },
          'SUB-005': { formative: 70, summative: 72, previousFinalScore: 74 },
          'SUB-006': { formative: 74, summative: 76, previousFinalScore: 78 },
        },
      },
    ],
    'AY-2025': [
      {
        studentId: 'STU-101',
        studentName: 'Zainal Abidin',
        seatNumber: 1,
        scores: {
          'SUB-001': { formative: 80, summative: 82, previousFinalScore: null },
          'SUB-002': { formative: 84, summative: 86, previousFinalScore: null },
          'SUB-003': { formative: 78, summative: 80, previousFinalScore: null },
          'SUB-004': { formative: 82, summative: 84, previousFinalScore: null },
          'SUB-005': { formative: 86, summative: 88, previousFinalScore: null },
          'SUB-006': { formative: 88, summative: 90, previousFinalScore: null },
        },
      },
    ],
  },
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export class HomeroomProgressMockAdapter {
  /**
   * Mengambil data mentah nilai siswa per kelas berdasarkan wali kelas.
   * Data ini akan diproses oleh Domain Engine di Service layer.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} Data mentah nilai siswa per kelas
   */
  async fetchProgressByHomeroom(context) {
    await delay(500)

    const { userId, academicYearId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    if (!academicYearId) {
      throw new Error('Academic Year ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      return {
        classId: null,
        className: null,
        subjects: [],
        students: [],
      }
    }

    const { classId, className } = assignment
    const classData = MOCK_PROGRESS_DB[classId] || {}
    const rawStudents = classData[academicYearId] || []

    return {
      classId: classId,
      className: className,
      subjects: MOCK_SUBJECTS,
      students: rawStudents,
    }
  }
}
