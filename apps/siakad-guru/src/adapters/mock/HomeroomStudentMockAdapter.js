// FILE: src/adapters/mock/HomeroomStudentMockAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { mapHomeroomStudents } from '@/contracts/homeroomStudentContract'

/**
 * Mock Homeroom Assignments
 * Memetakan userId guru/wali kelas ke kelas yang diwalinya.
 */
const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' }, // User 'wali'
}

/**
 * Mock Students Database
 * Dikelompokkan berdasarkan classId dan academicYearId.
 */
const MOCK_STUDENTS_DB = {
  'CLS-XI-B': {
    'AY-2026': [
      {
        id: 'STU-001',
        nisn: '0012345678',
        nis: '12345',
        fullName: 'Ahmad Fauzi',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=ahmad',
        status: 'ACTIVE',
        parentContact: '081234567890',
        seatNumber: 1,
      },
      {
        id: 'STU-002',
        nisn: '0012345679',
        nis: '12346',
        fullName: 'Siti Nurhaliza',
        gender: 'P',
        photoUrl: 'https://i.pravatar.cc/150?u=siti',
        status: 'ACTIVE',
        parentContact: '081234567891',
        seatNumber: 2,
      },
      {
        id: 'STU-003',
        nisn: '0012345680',
        nis: '12347',
        fullName: 'Budi Santoso',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=budi',
        status: 'ACTIVE',
        parentContact: '081234567892',
        seatNumber: 3,
      },
      {
        id: 'STU-004',
        nisn: '0012345681',
        nis: '12348',
        fullName: 'Dewi Lestari',
        gender: 'P',
        photoUrl: 'https://i.pravatar.cc/150?u=dewi',
        status: 'ACTIVE',
        parentContact: '081234567893',
        seatNumber: 4,
      },
      {
        id: 'STU-005',
        nisn: '0012345682',
        nis: '12349',
        fullName: 'Eko Prasetyo',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=eko',
        status: 'TRANSFERRED',
        parentContact: '081234567894',
        seatNumber: 5,
      },
      {
        id: 'STU-006',
        nisn: '0012345683',
        nis: '12350',
        fullName: 'Fitri Handayani',
        gender: 'P',
        photoUrl: 'https://i.pravatar.cc/150?u=fitri',
        status: 'ACTIVE',
        parentContact: '081234567895',
        seatNumber: 6,
      },
      {
        id: 'STU-007',
        nisn: '0012345684',
        nis: '12351',
        fullName: 'Gilang Ramadhan',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=gilang',
        status: 'ACTIVE',
        parentContact: '081234567896',
        seatNumber: 7,
      },
      {
        id: 'STU-008',
        nisn: '0012345685',
        nis: '12352',
        fullName: 'Hana Pertiwi',
        gender: 'P',
        photoUrl: 'https://i.pravatar.cc/150?u=hana',
        status: 'ACTIVE',
        parentContact: '081234567897',
        seatNumber: 8,
      },
      {
        id: 'STU-009',
        nisn: '0012345686',
        nis: '12353',
        fullName: 'Irfan Hakim',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=irfan',
        status: 'ACTIVE',
        parentContact: '081234567898',
        seatNumber: 9,
      },
      {
        id: 'STU-010',
        nisn: '0012345687',
        nis: '12354',
        fullName: 'Joko Susilo',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=joko',
        status: 'ACTIVE',
        parentContact: '081234567899',
        seatNumber: 10,
      },
    ],
    'AY-2025': [
      {
        id: 'STU-101',
        nisn: '0012345600',
        nis: '11345',
        fullName: 'Zainal Abidin',
        gender: 'L',
        photoUrl: 'https://i.pravatar.cc/150?u=zainal',
        status: 'GRADUATED',
        parentContact: '081234567800',
        seatNumber: 1,
      },
    ],
  },
}

/**
 * Simulasi latency jaringan
 */
function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/**
 * Homeroom Student Mock Adapter
 * Mengimplementasikan kontrak untuk mengambil data siswa wali kelas.
 */
export class HomeroomStudentMockAdapter {
  /**
   * Mengambil daftar siswa yang diwalikan berdasarkan context.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID dari wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @returns {Promise<Object>} Hasil berisi info kelas dan daftar siswa
   */
  async fetchStudentsByHomeroom(context) {
    await delay(400) // Simulasi latency jaringan

    const { userId, academicYearId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    if (!academicYearId) {
      throw new Error('Academic Year ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      // User bukan wali kelas atau belum memiliki penugasan kelas
      return {
        classId: null,
        className: null,
        students: [],
      }
    }

    const { classId, className } = assignment
    const classData = MOCK_STUDENTS_DB[classId] || {}
    const rawStudents = classData[academicYearId] || []

    // Petakan ke kontrak frontend
    const students = mapHomeroomStudents(rawStudents)

    // Urutkan berdasarkan nomor absen
    students.sort((a, b) => a.seatNumber - b.seatNumber)

    return {
      classId,
      className,
      students,
    }
  }
}
