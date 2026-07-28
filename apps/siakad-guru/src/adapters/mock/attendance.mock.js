import { AttendanceRecord } from '@/models/Attendance'
import { AttendanceStatus } from '@/engine/AttendanceEngine'

// --- Data siswa per kelas (dummy) ---
const studentsByClass = {
  'XII RPL 1': [
    { id: 1, name: 'Ahmad Fauzi', nis: '20241001' },
    { id: 2, name: 'Bunga Lestari', nis: '20241002' },
    { id: 3, name: 'Cahya Pratama', nis: '20241003' },
    { id: 4, name: 'Dewi Sartika', nis: '20241004' },
    { id: 5, name: 'Eko Saputro', nis: '20241005' },
    { id: 6, name: 'Fitriani', nis: '20241006' },
    { id: 7, name: 'Galang Ramadhan', nis: '20241007' },
    { id: 8, name: 'Hani Nuraini', nis: '20241008' },
    { id: 9, name: 'Irfan Maulana', nis: '20241009' },
    { id: 10, name: 'Jasmine Azahra', nis: '20241010' },
  ],
  'XII RPL 2': [
    { id: 11, name: 'Kartika Sari', nis: '20241011' },
    { id: 12, name: 'Lutfi Hakim', nis: '20241012' },
    { id: 13, name: 'Mega Permata', nis: '20241013' },
    { id: 14, name: 'Naufal Arif', nis: '20241014' },
    { id: 15, name: 'Olivia Putri', nis: '20241015' },
    { id: 16, name: 'Pandu Winata', nis: '20241016' },
    { id: 17, name: 'Qori Aini', nis: '20241017' },
    { id: 18, name: 'Rizki Akbar', nis: '20241018' },
  ],
}

// --- State in-memory untuk attendance records ---
let attendanceRecords = []

// Helper: buat attendance record kosong untuk suatu sesi
// eslint-disable-next-line no-unused-vars
function createEmptyRecord(session) {
  const students = (studentsByClass[session.className] || []).map((s) => ({
    studentId: s.id,
    studentName: s.name,
    nis: s.nis,
    status: AttendanceStatus.UNMARKED,
    note: '',
  }))

  return new AttendanceRecord({
    id: 0, // akan di-set saat disimpan
    sessionId: session.id,
    date: session.date,
    className: session.className,
    subject: session.subject,
    students,
    status: 'draft',
  })
}

export const attendanceMockAdapter = {
  /**
   * Mendapatkan daftar siswa untuk sesi tertentu.
   * @param {number} sessionId
   * @param {string} className
   */
  async fetchStudentsForSession(sessionId, className) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    
    const list = studentsByClass[className] || []
    return list.map((s) => ({
      studentId: s.id,
      studentName: s.name,
      nis: s.nis,
      status: AttendanceStatus.UNMARKED,
      note: '',
    }))
  },

  /**
   * Mendapatkan attendance record yang sudah ada untuk sesi tertentu.
   * Mengembalikan null jika belum ada.
   */
  async fetchAttendanceRecord(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const existing = attendanceRecords.find((r) => r.sessionId === sessionId)
    return existing ? new AttendanceRecord({ ...existing }) : null
  },

  /**
   * Menyimpan record sebagai draft (create or update).
   * @param {Object} session - data sesi dari TeachingSession
   * @param {Array} students - array student attendance items
   */
  async saveDraft(session, students) {
    await new Promise((resolve) => setTimeout(resolve, 250))
    const index = attendanceRecords.findIndex((r) => r.sessionId === session.id)
    const now = new Date().toISOString()

    if (index >= 0) {
      // Update existing draft
      attendanceRecords[index].students = students
      attendanceRecords[index].lastSavedAt = now
      return new AttendanceRecord({ ...attendanceRecords[index] })
    } else {
      // Buat baru
      const newRecord = new AttendanceRecord({
        id: attendanceRecords.length + 1,
        sessionId: session.id,
        date: session.date,
        className: session.className,
        subject: session.subject,
        students,
        status: 'draft',
        lastSavedAt: now,
      })
      attendanceRecords.push(newRecord)
      return new AttendanceRecord({ ...newRecord })
    }
  },

  /**
   * Submit final absensi (ubah status jadi 'submitted').
   * Hanya bisa jika semua siswa ditandai.
   */
  async submitAttendance(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 300))
    const record = attendanceRecords.find((r) => r.sessionId === sessionId)
    if (!record) throw new Error('Record tidak ditemukan')
    if (record.status === 'submitted') throw new Error('Absensi sudah disubmit')

    record.status = 'submitted'
    return new AttendanceRecord({ ...record })
  },

  /**
   * Mendapatkan ringkasan absensi untuk sesi tertentu.
   */
  async getAttendanceSummary(sessionId) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const record = attendanceRecords.find((r) => r.sessionId === sessionId)
    if (!record) return null
    // Ringkasan akan dihitung oleh engine
    return {
      sessionId: record.sessionId,
      date: record.date,
      className: record.className,
      subject: record.subject,
      students: record.students,
      status: record.status,
      lastSavedAt: record.lastSavedAt,
    }
  },
}
