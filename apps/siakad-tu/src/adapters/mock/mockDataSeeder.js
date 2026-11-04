// apps/siakad-tu/src/adapters/mock/mockDataSeeder.js
import { mockStorage } from './mockStorage.js'

export const mockDataSeeder = {
  seedStudents(schoolId, periodId = '20261') {
    const existing = mockStorage.getItem('students') || []
    if (existing.length > 0) {
      console.log(`[MockSeeder] Data siswa sudah ada, skip seeding.`)
      return
    }
    
    const students = [
      {
        peserta_didik_id: 'seed-student-001',
        nisn: '0012345678',
        nis: '2024001',
        nm_siswa: 'Andi Pratama',
        tanggal_lahir: '2009-03-15',
        jenis_kelamin: 'L',
        alamat_siswa: 'Jl. Raya Jatinangor No. 10',
        telepon_siswa: '081234567801',
        nm_ayah: 'Budi Pratama',
        tahun_ajaran_id: periodId,
        kelas: 'XI IPA 1',
        status: 'Aktif',
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
        peserta_didik_id: 'seed-student-002',
        nisn: '0012345679',
        nis: '2024002',
        nm_siswa: 'Siti Nur Aisyah',
        tanggal_lahir: '2009-07-22',
        jenis_kelamin: 'P',
        alamat_siswa: 'Jl. Cipacing No. 5',
        telepon_siswa: '081234567802',
        nm_ayah: 'Hajah Aminah',
        tahun_ajaran_id: periodId,
        kelas: 'XI IPA 2',
        status: 'Aktif',
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
        peserta_didik_id: 'seed-student-003',
        nisn: '0012345680',
        nis: '2023001',
        nm_siswa: 'Rizky Firmansyah',
        tanggal_lahir: '2008-11-08',
        jenis_kelamin: 'L',
        alamat_siswa: 'Jl. Hegarmanah No. 20',
        telepon_siswa: '081234567804',
        nm_ayah: 'H. Firmansyah',
        tahun_ajaran_id: periodId,
        kelas: 'XII IPA 1',
        status: 'Aktif',
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
    
    mockStorage.setItem('students', students)
    console.log(`[MockSeeder] ${students.length} siswa di-seed.`)
  },
  
  seedAll(schoolId, periodId = '20261') {
    this.seedStudents(schoolId, periodId)
  },
}
