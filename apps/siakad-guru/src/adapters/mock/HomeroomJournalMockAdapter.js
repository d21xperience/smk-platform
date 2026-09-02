// FILE: src/adapters/mock/HomeroomJournalMockAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Journal Mock Adapter
 * Menyediakan data catatan wali kelas dengan operasi CRUD.
 * Data terikat pada userId (wali kelas), schoolId, academicYearId, semesterId.
 * Menggunakan in-memory array untuk simulasi mutasi (create, update, delete).
 */

import { mapHomeroomNotes, mapHomeroomNote } from '@/contracts/homeroomJournalContract'

const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' },
}

let noteIdCounter = 100

const MOCK_NOTES_DB = {
  'CLS-XI-B': {
    'AY-2026': {
      20261: [
        {
          noteId: 'NOTE-001',
          studentId: 'STU-001',
          studentName: 'Ahmad Fauzi',
          category: 'POSITIVE',
          priority: 'LOW',
          title: 'Juara 1 Olimpiade Matematika Tingkat Kota',
          description:
            'Ahmad Fauzi berhasil meraih juara 1 dalam Olimpiade Matematika tingkat kota. Prestasi ini sangat membanggakan kelas dan sekolah. Perlu diberikan apresiasi di depan kelas untuk memotivasi siswa lainnya.',
          noteDate: '2026-08-15',
          createdAt: '2026-08-15T10:30:00.000Z',
          createdBy: 'USR-002',
          status: 'SUBMITTED',
        },
        {
          noteId: 'NOTE-002',
          studentId: 'STU-007',
          studentName: 'Gilang Ramadhan',
          category: 'NEGATIVE',
          priority: 'HIGH',
          title: 'Sering bolos pelajaran tanpa keterangan',
          description:
            'Gilang tercatat sudah 5 kali tidak masuk sekolah tanpa keterangan dalam bulan Agustus. Sudah dilakukan pemanggilan orang tua namun belum ada perubahan signifikan. Perlu tindak lanjut lebih lanjut dengan BK.',
          noteDate: '2026-08-20',
          createdAt: '2026-08-20T14:15:00.000Z',
          createdBy: 'USR-002',
          status: 'SUBMITTED',
        },
        {
          noteId: 'NOTE-003',
          studentId: 'STU-010',
          studentName: 'Joko Susilo',
          category: 'INCIDENT',
          priority: 'URGENT',
          title: 'Perkelahian dengan siswa kelas lain',
          description:
            'Joko terlibat perkelahian dengan siswa kelas XI-C saat istirahat. Kejadian dipicu oleh kesalahpahaman saat bermain bola. Kedua siswa sudah dipanggil ke ruang BK dan orang tua sudah dihubungi. Perlu monitoring perilaku selanjutnya.',
          noteDate: '2026-09-05',
          createdAt: '2026-09-05T11:00:00.000Z',
          createdBy: 'USR-002',
          status: 'SUBMITTED',
        },
        {
          noteId: 'NOTE-004',
          studentId: 'STU-003',
          studentName: 'Budi Santoso',
          category: 'NEUTRAL',
          priority: 'MEDIUM',
          title: 'Konseling mengenai motivasi belajar',
          description:
            'Budi terlihat kurang termotivasi dalam beberapa minggu terakhir. Sudah dilakukan sesi konseling ringan dan diketahui ada masalah keluarga. Akan dilakukan follow-up minggu depan untuk melihat perkembangan.',
          noteDate: '2026-09-10',
          createdAt: '2026-09-10T09:45:00.000Z',
          createdBy: 'USR-002',
          status: 'SUBMITTED',
        },
        {
          noteId: 'NOTE-005',
          studentId: 'STU-008',
          studentName: 'Hana Pertiwi',
          category: 'POSITIVE',
          priority: 'LOW',
          title: 'Aktif membantu teman yang kesulitan belajar',
          description:
            'Hana secara sukarela membantu teman-temannya yang kesulitan memahami materi IPA setelah jam pelajaran. Inisiatif ini sangat positif dan perlu didukung. Pertimbangkan untuk menjadikan Hana sebagai tutor sebaya.',
          noteDate: '2026-09-18',
          createdAt: '2026-09-18T13:20:00.000Z',
          createdBy: 'USR-002',
          status: 'SUBMITTED',
        },
        {
          noteId: 'NOTE-006',
          studentId: 'STU-005',
          studentName: 'Eko Prasetyo',
          category: 'NEGATIVE',
          priority: 'MEDIUM',
          title: 'Sering terlambat masuk kelas',
          description:
            'Eko tercatat 4 kali terlambat masuk kelas dalam minggu ini. Alasan yang diberikan beragam. Sudah diingatkan secara personal. Jika berlanjut, akan dilakukan pemanggilan orang tua.',
          noteDate: '2026-10-02',
          createdAt: '2026-10-02T08:30:00.000Z',
          createdBy: 'USR-002',
          status: 'DRAFT',
        },
        {
          noteId: 'NOTE-007',
          studentId: 'STU-002',
          studentName: 'Siti Nurhaliza',
          category: 'POSITIVE',
          priority: 'LOW',
          title: 'Terpilih sebagai ketua panitia class meeting',
          description:
            'Siti terpilih secara aklamasi sebagai ketua panitia class meeting semester ini. Menunjukkan jiwa kepemimpinan yang baik. Perlu didukung dan dibimbing dalam menjalankan tugasnya.',
          noteDate: '2026-10-10',
          createdAt: '2026-10-10T10:00:00.000Z',
          createdBy: 'USR-002',
          status: 'DRAFT',
        },
        {
          noteId: 'NOTE-008',
          studentId: 'STU-009',
          studentName: 'Irfan Hakim',
          category: 'NEUTRAL',
          priority: 'HIGH',
          title: 'Sering sakit, perlu koordinasi dengan orang tua',
          description:
            'Irfan tercatat 6 kali tidak masuk karena sakit dalam 2 bulan terakhir. Perlu koordinasi dengan orang tua untuk mengetahui kondisi kesehatan sebenarnya dan kemungkinan surat keterangan dokter untuk keperluan administrasi.',
          noteDate: '2026-10-15',
          createdAt: '2026-10-15T11:30:00.000Z',
          createdBy: 'USR-002',
          status: 'DRAFT',
        },
      ],
    },
  },
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export class HomeroomJournalMockAdapter {
  /**
   * Mengambil daftar catatan wali kelas berdasarkan context.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Array>} Array catatan
   */
  async fetchNotesByHomeroom(context) {
    await delay(400)

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
      return []
    }

    const { classId } = assignment
    const classData = MOCK_NOTES_DB[classId] || {}
    const yearData = classData[academicYearId] || {}
    const notes = yearData[semesterId] || []

    return mapHomeroomNotes(notes)
  }

  /**
   * Menyimpan catatan baru ke mock database.
   * @param {Object} note - Objek catatan baru (tanpa noteId).
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Catatan yang tersimpan dengan noteId.
   */
  async saveNote(note, context) {
    await delay(300)

    const { userId, academicYearId, semesterId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      throw new Error('User is not assigned as homeroom teacher.')
    }

    const { classId } = assignment

    noteIdCounter++
    const noteId = `NOTE-${noteIdCounter}`

    const savedNote = {
      ...note,
      noteId: noteId,
    }

    if (!MOCK_NOTES_DB[classId]) {
      MOCK_NOTES_DB[classId] = {}
    }
    if (!MOCK_NOTES_DB[classId][academicYearId]) {
      MOCK_NOTES_DB[classId][academicYearId] = {}
    }
    if (!MOCK_NOTES_DB[classId][academicYearId][semesterId]) {
      MOCK_NOTES_DB[classId][academicYearId][semesterId] = []
    }

    MOCK_NOTES_DB[classId][academicYearId][semesterId].push(savedNote)

    return mapHomeroomNote(savedNote)
  }

  /**
   * Mengupdate catatan yang ada di mock database.
   * @param {Object} note - Objek catatan yang sudah diupdate.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Catatan yang terupdate.
   */
  async updateNote(note, context) {
    await delay(300)

    const { userId, academicYearId, semesterId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      throw new Error('User is not assigned as homeroom teacher.')
    }

    const { classId } = assignment
    const classData = MOCK_NOTES_DB[classId] || {}
    const yearData = classData[academicYearId] || {}
    const notes = yearData[semesterId] || []

    const noteIndex = notes.findIndex((n) => n.noteId === note.noteId)

    if (noteIndex === -1) {
      throw new Error(`Note with ID ${note.noteId} not found.`)
    }

    notes[noteIndex] = { ...note }

    return mapHomeroomNote(notes[noteIndex])
  }

  /**
   * Menghapus catatan dari mock database.
   * @param {string} noteId - ID catatan yang akan dihapus.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Hasil penghapusan.
   */
  async deleteNote(noteId, context) {
    await delay(300)

    const { userId, academicYearId, semesterId } = context

    if (!userId) {
      throw new Error('User ID is required in operational context.')
    }

    const assignment = MOCK_HOMEROOM_ASSIGNMENTS[userId]

    if (!assignment) {
      throw new Error('User is not assigned as homeroom teacher.')
    }

    const { classId } = assignment
    const classData = MOCK_NOTES_DB[classId] || {}
    const yearData = classData[academicYearId] || {}
    const notes = yearData[semesterId] || []

    const noteIndex = notes.findIndex((n) => n.noteId === noteId)

    if (noteIndex === -1) {
      throw new Error(`Note with ID ${noteId} not found.`)
    }

    notes.splice(noteIndex, 1)

    return { success: true, deletedNoteId: noteId }
  }
}
