// FILE: src/adapters/mock/HomeroomAttendanceMockAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Attendance Mock Adapter
 * Menyediakan data mentah kehadiran siswa per kelas untuk diproses Domain Engine.
 * Data terikat pada userId (wali kelas), schoolId, dan academicYearId.
 */

const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' },
}

const MOCK_SEMESTER_MONTHS = {
  20261: [
    { monthId: '2026-07', monthName: 'Juli 2026', totalDays: 22 },
    { monthId: '2026-08', monthName: 'Agustus 2026', totalDays: 24 },
    { monthId: '2026-09', monthName: 'September 2026', totalDays: 23 },
    { monthId: '2026-10', monthName: 'Oktober 2026', totalDays: 25 },
    { monthId: '2026-11', monthName: 'November 2026', totalDays: 22 },
    { monthId: '2026-12', monthName: 'Desember 2026', totalDays: 18 },
  ],
  20262: [
    { monthId: '2027-01', monthName: 'Januari 2027', totalDays: 22 },
    { monthId: '2027-02', monthName: 'Februari 2027', totalDays: 23 },
    { monthId: '2027-03', monthName: 'Maret 2027', totalDays: 24 },
    { monthId: '2027-04', monthName: 'April 2027', totalDays: 22 },
    { monthId: '2027-05', monthName: 'Mei 2027', totalDays: 23 },
    { monthId: '2027-06', monthName: 'Juni 2027', totalDays: 18 },
  ],
}

const MOCK_ATTENDANCE_DB = {
  'CLS-XI-B': {
    'AY-2026': [
      {
        studentId: 'STU-001',
        studentName: 'Ahmad Fauzi',
        seatNumber: 1,
        attendance: {
          '2026-07': { present: 21, sick: 0, permitted: 1, absent: 0 },
          '2026-08': { present: 23, sick: 1, permitted: 0, absent: 0 },
          '2026-09': { present: 22, sick: 0, permitted: 1, absent: 0 },
          '2026-10': { present: 24, sick: 0, permitted: 1, absent: 0 },
          '2026-11': { present: 21, sick: 1, permitted: 0, absent: 0 },
          '2026-12': { present: 17, sick: 0, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-002',
        studentName: 'Siti Nurhaliza',
        seatNumber: 2,
        attendance: {
          '2026-07': { present: 20, sick: 1, permitted: 1, absent: 0 },
          '2026-08': { present: 22, sick: 1, permitted: 1, absent: 0 },
          '2026-09': { present: 21, sick: 1, permitted: 1, absent: 0 },
          '2026-10': { present: 23, sick: 1, permitted: 1, absent: 0 },
          '2026-11': { present: 20, sick: 1, permitted: 1, absent: 0 },
          '2026-12': { present: 16, sick: 1, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-003',
        studentName: 'Budi Santoso',
        seatNumber: 3,
        attendance: {
          '2026-07': { present: 21, sick: 0, permitted: 1, absent: 0 },
          '2026-08': { present: 22, sick: 1, permitted: 1, absent: 0 },
          '2026-09': { present: 21, sick: 0, permitted: 1, absent: 1 },
          '2026-10': { present: 23, sick: 1, permitted: 0, absent: 1 },
          '2026-11': { present: 20, sick: 1, permitted: 1, absent: 0 },
          '2026-12': { present: 16, sick: 1, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-004',
        studentName: 'Dewi Lestari',
        seatNumber: 4,
        attendance: {
          '2026-07': { present: 21, sick: 1, permitted: 0, absent: 0 },
          '2026-08': { present: 23, sick: 0, permitted: 1, absent: 0 },
          '2026-09': { present: 22, sick: 0, permitted: 1, absent: 0 },
          '2026-10': { present: 24, sick: 0, permitted: 1, absent: 0 },
          '2026-11': { present: 21, sick: 0, permitted: 1, absent: 0 },
          '2026-12': { present: 17, sick: 0, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-005',
        studentName: 'Eko Prasetyo',
        seatNumber: 5,
        attendance: {
          '2026-07': { present: 18, sick: 1, permitted: 2, absent: 1 },
          '2026-08': { present: 19, sick: 1, permitted: 3, absent: 1 },
          '2026-09': { present: 18, sick: 1, permitted: 2, absent: 2 },
          '2026-10': { present: 20, sick: 1, permitted: 3, absent: 1 },
          '2026-11': { present: 17, sick: 1, permitted: 3, absent: 1 },
          '2026-12': { present: 14, sick: 1, permitted: 2, absent: 1 },
        },
      },
      {
        studentId: 'STU-006',
        studentName: 'Fitri Handayani',
        seatNumber: 6,
        attendance: {
          '2026-07': { present: 20, sick: 1, permitted: 1, absent: 0 },
          '2026-08': { present: 22, sick: 0, permitted: 1, absent: 1 },
          '2026-09': { present: 21, sick: 1, permitted: 1, absent: 0 },
          '2026-10': { present: 23, sick: 0, permitted: 1, absent: 1 },
          '2026-11': { present: 20, sick: 1, permitted: 1, absent: 0 },
          '2026-12': { present: 16, sick: 1, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-007',
        studentName: 'Gilang Ramadhan',
        seatNumber: 7,
        attendance: {
          '2026-07': { present: 14, sick: 1, permitted: 2, absent: 5 },
          '2026-08': { present: 15, sick: 1, permitted: 2, absent: 6 },
          '2026-09': { present: 14, sick: 1, permitted: 2, absent: 6 },
          '2026-10': { present: 16, sick: 1, permitted: 2, absent: 6 },
          '2026-11': { present: 13, sick: 1, permitted: 2, absent: 6 },
          '2026-12': { present: 11, sick: 1, permitted: 2, absent: 4 },
        },
      },
      {
        studentId: 'STU-008',
        studentName: 'Hana Pertiwi',
        seatNumber: 8,
        attendance: {
          '2026-07': { present: 22, sick: 0, permitted: 0, absent: 0 },
          '2026-08': { present: 23, sick: 0, permitted: 1, absent: 0 },
          '2026-09': { present: 23, sick: 0, permitted: 0, absent: 0 },
          '2026-10': { present: 24, sick: 0, permitted: 1, absent: 0 },
          '2026-11': { present: 22, sick: 0, permitted: 0, absent: 0 },
          '2026-12': { present: 17, sick: 0, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-009',
        studentName: 'Irfan Hakim',
        seatNumber: 9,
        attendance: {
          '2026-07': { present: 18, sick: 3, permitted: 1, absent: 0 },
          '2026-08': { present: 19, sick: 3, permitted: 1, absent: 1 },
          '2026-09': { present: 18, sick: 3, permitted: 1, absent: 1 },
          '2026-10': { present: 20, sick: 3, permitted: 1, absent: 1 },
          '2026-11': { present: 17, sick: 3, permitted: 1, absent: 1 },
          '2026-12': { present: 14, sick: 3, permitted: 1, absent: 0 },
        },
      },
      {
        studentId: 'STU-010',
        studentName: 'Joko Susilo',
        seatNumber: 10,
        attendance: {
          '2026-07': { present: 13, sick: 1, permitted: 1, absent: 7 },
          '2026-08': { present: 14, sick: 1, permitted: 1, absent: 8 },
          '2026-09': { present: 13, sick: 1, permitted: 1, absent: 8 },
          '2026-10': { present: 15, sick: 1, permitted: 1, absent: 8 },
          '2026-11': { present: 12, sick: 1, permitted: 1, absent: 8 },
          '2026-12': { present: 10, sick: 1, permitted: 1, absent: 6 },
        },
      },
    ],
  },
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export class HomeroomAttendanceMockAdapter {
  /**
   * Mengambil data mentah kehadiran siswa per kelas berdasarkan wali kelas.
   * Data ini akan diproses oleh Domain Engine di Service layer.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} Data mentah kehadiran siswa per kelas
   */
  async fetchAttendanceByHomeroom(context) {
    await delay(500)

    const { userId, academicYearId, semesterId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    if (!academicYearId) {
      throw new Error('Academic Year ID is required in operational context.')
    }

    if (!semesterId) {
      throw new Error('Semester ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      return {
        classId: null,
        className: null,
        months: [],
        students: [],
      }
    }

    const { classId, className } = assignment
    const months = MOCK_SEMESTER_MONTHS[semesterId] || []
    const classData = MOCK_ATTENDANCE_DB[classId] || {}
    const rawStudents = classData[academicYearId] || []

    return {
      classId: classId,
      className: className,
      months: months,
      students: rawStudents,
    }
  }
}
