// apps/siakad-tu/src/adapters/mock/studentMockAdapter.js
import { mockStorage } from './mockStorage.js'

const MOCK_STUDENTS_BY_PERIOD = {
  '20251': [
    {
      peserta_didik_id: 'uuid-siswa-1',
      nis: '20251788779136',
      nisn: '0088779136',
      nm_siswa: 'Ahmad Fauzi',
      jenis_kelamin: 'L',
      tempat_lahir: 'Jakarta',
      tanggal_lahir: '2008-05-12',
      alamat_siswa: 'Jl. Merdeka No. 1, Menteng, Jakarta Pusat',
      telepon_siswa: '08123456789',
      nm_ayah: 'Budi Santoso',
      nm_ibu: 'Siti Aminah',
      kelas: 'X RPL 1',
      tingkat: 'X',
      jurusan: 'RPL',
      status: 'Aktif',
      tahun_ajaran_id: '20251',
      wali_kelas: 'Budi Santoso'
    },
    {
      peserta_didik_id: 'uuid-siswa-2',
      nis: '20251788779019',
      nisn: '0088779019',
      nm_siswa: 'Siti Nurhaliza',
      jenis_kelamin: 'P',
      tempat_lahir: 'Bandung',
      tanggal_lahir: '2008-08-20',
      alamat_siswa: 'Jl. Asia Afrika No. 2, Sumur Bandung',
      telepon_siswa: '08987654321',
      nm_ayah: 'Hendra',
      nm_ibu: 'Rina',
      kelas: 'X TKJ 2',
      tingkat: 'X',
      jurusan: 'TKJ',
      status: 'Aktif',
      tahun_ajaran_id: '20251',
      wali_kelas: 'Rina Wijaya'
    }
  ],
  '20261': [
    {
      peserta_didik_id: 'uuid-siswa-1',
      nis: '20251788779136',
      nisn: '0088779136',
      nm_siswa: 'Ahmad Fauzi',
      jenis_kelamin: 'L',
      tempat_lahir: 'Jakarta',
      tanggal_lahir: '2008-05-12',
      alamat_siswa: 'Jl. Merdeka No. 1, Menteng, Jakarta Pusat',
      telepon_siswa: '08123456789',
      nm_ayah: 'Budi Santoso',
      nm_ibu: 'Siti Aminah',
      kelas: 'XI RPL 1', 
      tingkat: 'XI',
      jurusan: 'RPL',
      status: 'Aktif',
      tahun_ajaran_id: '20261',
      wali_kelas: 'Dedi Kurniawan'
    }
  ]
}

const ensureDataSeeded = (schoolId) => {
  if (!schoolId) return
  const existing = mockStorage.read(schoolId, 'students_20251')
  if (existing.length === 0) {
    Object.keys(MOCK_STUDENTS_BY_PERIOD).forEach(periodId => {
      mockStorage.write(schoolId, `students_${periodId}`, MOCK_STUDENTS_BY_PERIOD[periodId])
    })
  }
}

const mapToFrontend = (s) => ({
  id: s.peserta_didik_id,
  nis: s.nis,
  nisn: s.nisn,
  fullName: s.nm_siswa,
  gender: s.jenis_kelamin, 
  className: s.kelas,
  tingkat: s.tingkat,
  jurusan: s.jurusan,
  status: s.status,
  waliKelas: s.wali_kelas,
  birthPlace: s.tempat_lahir,
  birthDate: s.tanggal_lahir,
  address: s.alamat_siswa,
  phone: s.telepon_siswa,
  parentName: s.nm_ayah,
  motherName: s.nm_ibu
})

const delay = (ms = 300) => new Promise(resolve => setTimeout(resolve, ms))

export const studentMockAdapter = {
  async getStudents(context, filters = {}) {
    const schoolId = context?.schoolId || 'default-school'
    const periodId = context?.semesterId || '20251' 

    ensureDataSeeded(schoolId)
    await delay()
    
    let students = mockStorage.read(schoolId, `students_${periodId}`)
    if (!students) students = []

    if (filters.search) {
      const keyword = filters.search.toLowerCase()
      students = students.filter(s => 
        s.nm_siswa.toLowerCase().includes(keyword) || 
        s.nis.includes(keyword) || 
        s.nisn.includes(keyword)
      )
    }
    if (filters.status) students = students.filter(s => s.status === filters.status)
    if (filters.tingkat) students = students.filter(s => s.tingkat === filters.tingkat)
    if (filters.jurusan) students = students.filter(s => s.jurusan === filters.jurusan)
    if (filters.kelas) students = students.filter(s => s.kelas === filters.kelas)

    return {
      success: true,
      data: {
        items: students.map(mapToFrontend),
        total: students.length,
        page: filters.page || 1,
        limit: filters.limit || 10,
        totalPages: Math.ceil(students.length / (filters.limit || 10))
      }
    }
  },

  async getStudentById(id, context) {
    const schoolId = context?.schoolId || 'default-school'
    const periodId = context?.semesterId || '20251'
    
    ensureDataSeeded(schoolId)
    await delay()
    
    const students = mockStorage.read(schoolId, `students_${periodId}`) || []
    const student = students.find(s => s.peserta_didik_id === id)
    if (!student) return { success: false, error: { code: 'STUDENT_NOT_FOUND', message: 'Siswa tidak ditemukan di periode ini' } }
    return { success: true, data: mapToFrontend(student) }
  },

  // FUNGSI BARU: Untuk menyimpan perubahan dari Form Edit
  async updateStudentProfile(command, context) {
    const schoolId = context?.schoolId || 'default-school'
    const periodId = context?.semesterId || '20251'
    await delay()
    ensureDataSeeded(schoolId)
    
    const students = mockStorage.read(schoolId, `students_${periodId}`) || []
    const index = students.findIndex(s => s.peserta_didik_id === command.studentId)
    
    if (index === -1) return { success: false, error: { message: 'Siswa tidak ditemukan' } }

    // Map dari Form (Flat) ke Database Mock (Bahasa Indonesia)
    students[index] = {
      ...students[index],
      nm_siswa: command.nama,
      tempat_lahir: command.tempatLahir,
      tanggal_lahir: command.tanggalLahir,
      jenis_kelamin: command.jenisKelamin,
      agama: command.agama,
      telepon_siswa: command.noHp,
      alamat_siswa: command.alamat,
      kelas: command.kelas,
      jurusan: command.jurusan
    }
    
    mockStorage.write(schoolId, `students_${periodId}`, students)
    return { success: true, data: mapToFrontend(students[index]) }
  },

  async registerStudent(command, context) {
    const schoolId = context?.schoolId || 'default-school'
    const periodId = context?.semesterId || '20251'
    await delay()
    ensureDataSeeded(schoolId)
    
    const students = mockStorage.read(schoolId, `students_${periodId}`) || []
    
    const newStudent = {
      peserta_didik_id: command.studentId || `uuid-${Date.now()}`,
      nis: command.nis,
      nisn: command.nisn,
      nm_siswa: command.fullName?.firstName || command.nm_siswa,
      jenis_kelamin: command.gender === 'MALE' ? 'L' : 'P',
      tahun_ajaran_id: periodId,
      status: 'Aktif',
      kelas: command.className || '-',
      tingkat: command.tingkat || 'X',
      jurusan: command.jurusan || 'UMUM'
    }
    
    students.push(newStudent)
    mockStorage.write(schoolId, `students_${periodId}`, students)
    return { success: true, data: mapToFrontend(newStudent) }
  },
  
  async enrollStudent() { return { success: true, data: null } },
  async graduateStudent() { return { success: true, data: null } },
  async checkNisnAvailability() { return { success: true, data: { available: true } } },
  
  async __resetMockData(schoolId) { 
    mockStorage.clearSchool(schoolId)
    return true 
  }
}
