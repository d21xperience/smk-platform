// apps/siakad-tu/src/adapters/mock/mockDataSeeder.js

import { mockStorage } from './mockStorage.js'
// import { Gender } from '../../../domain/student/value-objects/Gender.js'

/**
 * MockDataSeeder — Mengisi data awal untuk demo/development.
 *
 * Data ini akan muncul saat pertama kali app dibuka,
 * memberikan pengalaman realistis tanpa backend.
 */
export const mockDataSeeder = {
  /**
   * Seed data siswa awal untuk sekolah tertentu
   * @param {string} schoolId
   * @param {string} periodId
   */
  seedStudents(schoolId, periodId = '20261') {
    const existing = mockStorage.read(schoolId, 'students')
    if (existing.length > 0) {
      console.log(`[MockSeeder] Data siswa sudah ada untuk ${schoolId}, skip seeding.`)
      return
    }

    const students = [
      {
        studentId: 'seed-student-001',
        schoolId,
        nisn: '0012345678',
        nis: '2024001',
        fullName: { firstName: 'Andi', middleName: '', lastName: 'Pratama' },
        birthDate: '2009-03-15',
        gender: 'MALE',
        address: {
          street: 'Jl. Raya Jatinangor No. 10',
          rtRw: '01/02',
          village: 'Cibeusi',
          district: 'Jatinangor',
          city: 'Sumedang',
          postalCode: '45363',
        },
        contactInfo: { phone: '081234567801', email: 'andi@example.com' },
        guardianInfo: {
          name: 'Budi Pratama',
          relation: 'father',
          phone: '081234567800',
          occupation: 'Wiraswasta',
        },
        status: 'ACTIVE',
        enrollments: [
          {
            enrollmentId: 'seed-enr-001',
            schoolId,
            periodId,
            classId: 'class-xi-1',
            enrollmentDate: '2025-07-15T00:00:00.000Z',
            status: 'active',
          },
        ],
      },
      {
        studentId: 'seed-student-002',
        schoolId,
        nisn: '0012345679',
        nis: '2024002',
        fullName: { firstName: 'Siti', middleName: 'Nur', lastName: 'Aisyah' },
        birthDate: '2009-07-22',
        gender: 'FEMALE',
        address: {
          street: 'Jl. Cipacing No. 5',
          rtRw: '02/03',
          village: 'Cipacing',
          district: 'Jatinangor',
          city: 'Sumedang',
          postalCode: '45363',
        },
        contactInfo: { phone: '081234567802', email: '' },
        guardianInfo: {
          name: 'Hajah Aminah',
          relation: 'mother',
          phone: '081234567803',
          occupation: 'Guru',
        },
        status: 'ACTIVE',
        enrollments: [
          {
            enrollmentId: 'seed-enr-002',
            schoolId,
            periodId,
            classId: 'class-xi-2',
            enrollmentDate: '2025-07-15T00:00:00.000Z',
            status: 'active',
          },
        ],
      },
      {
        studentId: 'seed-student-003',
        schoolId,
        nisn: '0012345680',
        nis: '2023001',
        fullName: { firstName: 'Rizky', middleName: '', lastName: 'Firmansyah' },
        birthDate: '2008-11-08',
        gender: 'MALE',
        address: {
          street: 'Jl. Hegarmanah No. 20',
          rtRw: '03/01',
          village: 'Hegarmanah',
          district: 'Jatinangor',
          city: 'Sumedang',
          postalCode: '45364',
        },
        contactInfo: { phone: '081234567804', email: 'rizky@example.com' },
        guardianInfo: {
          name: 'H. Firmansyah',
          relation: 'father',
          phone: '081234567805',
          occupation: 'PNS',
        },
        status: 'ACTIVE',
        enrollments: [
          {
            enrollmentId: 'seed-enr-003',
            schoolId,
            periodId,
            classId: 'class-xii-1',
            enrollmentDate: '2024-07-15T00:00:00.000Z',
            status: 'active',
          },
        ],
      },
    ]

    mockStorage.write(schoolId, 'students', students)
    console.log(`[MockSeeder] ${students.length} siswa di-seed untuk ${schoolId}.`)
  },

  /**
   * Seed semua data yang diperlukan
   * @param {string} schoolId
   * @param {string} periodId
   */
  seedAll(schoolId, periodId = '20261') {
    this.seedStudents(schoolId, periodId)
  },
}
