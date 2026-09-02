// FILE: src/adapters/mock/HomeroomCommunicationMockAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Communication Mock Adapter
 * Menyediakan data log komunikasi wali kelas dengan orang tua.
 * Operasi: fetch, save, update (tidak ada delete — log komunikasi tidak dihapus).
 * Data terikat pada userId (wali kelas), schoolId, academicYearId, semesterId.
 * Menggunakan in-memory array untuk simulasi mutasi (create, update).
 */

import {
  mapHomeroomCommunications,
  mapHomeroomCommunication,
} from '@/contracts/homeroomCommunicationContract'

const MOCK_HOMEROOM_ASSIGNMENTS = {
  'USR-001': { classId: 'CLS-X-A', className: 'X-A' },
  'USR-002': { classId: 'CLS-XI-B', className: 'XI-B' },
}

let communicationIdCounter = 200

const MOCK_COMMUNICATIONS_DB = {
  'CLS-XI-B': {
    'AY-2026': {
      20261: [
        {
          communicationId: 'COMM-001',
          studentId: 'STU-001',
          studentName: 'Ahmad Fauzi',
          category: 'PHONE_CALL',
          direction: 'OUTBOUND',
          title: 'Apresiasi juara olimpiade matematika',
          summary:
            'Menghubungi orang tua Ahmad Fauzi untuk menyampaikan apresiasi atas prestasi juara 1 Olimpiade Matematika tingkat kota. Orang tua sangat senang dan berterima kasih atas dukungan sekolah. Diskusi juga mencakup rencana persiapan untuk tingkat provinsi bulan depan.',
          communicationDate: '2026-08-16',
          followUpNote:
            'Orang tua akan mendampingi latihan tambahan di rumah. Sekolah akan menyediakan buku referensi tambahan.',
          status: 'CLOSED',
          createdAt: '2026-08-16T10:00:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-002',
          studentId: 'STU-007',
          studentName: 'Gilang Ramadhan',
          category: 'WHATSAPP',
          direction: 'OUTBOUND',
          title: 'Peringatan kehadiran menurun',
          summary:
            'Mengirim pesan WhatsApp ke orang tua Gilang Ramadhan mengenai 5 kali ketidakhadiran tanpa keterangan di bulan Agustus. Meminta konfirmasi alasan dan mengharapkan kerja sama untuk memperbaiki kehadiran. Pesan sudah dibaca namun belum ada respons.',
          communicationDate: '2026-08-21',
          followUpNote: '',
          status: 'PENDING_FOLLOW_UP',
          createdAt: '2026-08-21T14:30:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-003',
          studentId: 'STU-010',
          studentName: 'Joko Susilo',
          category: 'MEETING',
          direction: 'BIDIRECTIONAL',
          title: 'Pertemuan terkait insiden perkelahian',
          summary:
            'Pertemuan tatap muka dengan kedua orang tua Joko Susilo dan siswa kelas XI-C yang terlibat perkelahian. Diskusi mencakup kronologi kejadian, dampak, dan langkah preventif. Kedua pihak sepakat untuk berdamai dan orang tua akan lebih mengawasi aktivitas anak di luar sekolah.',
          communicationDate: '2026-09-06',
          followUpNote:
            'Kedua orang tua sepakat damai. Monitoring perilaku selama 1 bulan oleh BK. Follow-up meeting dijadwalkan awal Oktober.',
          status: 'FOLLOWED_UP',
          createdAt: '2026-09-06T11:15:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-004',
          studentId: 'STU-009',
          studentName: 'Irfan Hakim',
          category: 'PHONE_CALL',
          direction: 'INBOUND',
          title: 'Orang tua menanyakan kondisi kesehatan Irfan',
          summary:
            'Orang tua Irfan Hakim menelepon untuk menanyakan kondisi kesehatan anaknya yang sering sakit. Mereka meminta saran apakah perlu surat keterangan dokter untuk keperluan administrasi sekolah. Juga menanyakan kemungkinan keringanan tugas selama masa pemulihan.',
          communicationDate: '2026-09-12',
          followUpNote: '',
          status: 'PENDING_FOLLOW_UP',
          createdAt: '2026-09-12T09:00:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-005',
          studentId: 'STU-005',
          studentName: 'Eko Prasetyo',
          category: 'HOME_VISIT',
          direction: 'OUTBOUND',
          title: 'Kunjungan rumah terkait keterlambatan berulang',
          summary:
            'Melakukan kunjungan rumah ke kediaman Eko Prasetyo untuk memahami alasan keterlambatan berulang. Diketahui bahwa Eko harus membantu orang tua di warung sebelum berangkat sekolah. Diskusi menghasilkan kesepakatan调整 jadwal bantu warung dan sekolah memberikan toleransi 15 menit.',
          communicationDate: '2026-09-20',
          followUpNote:
            'Kesepakatan: Eko membantu warung sampai jam 06.15, berangkat sekolah jam 06.20. Sekolah toleransi 15 menit. Review setelah 2 minggu.',
          status: 'FOLLOWED_UP',
          createdAt: '2026-09-20T15:00:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-006',
          studentId: 'STU-008',
          studentName: 'Hana Pertiwi',
          category: 'WHATSAPP',
          direction: 'OUTBOUND',
          title: 'Koordinasi program tutor sebaya',
          summary:
            'Mengirim pesan WhatsApp ke orang tua Hana Pertiwi mengenai inisiatif Hana yang aktif membantu teman-temannya belajar IPA. Meminta izin dan dukungan orang tua untuk menjadikan Hana sebagai tutor sebaya resmi. Orang tua sangat mendukung dan bangga.',
          communicationDate: '2026-09-19',
          followUpNote:
            'Orang tua mendukung penuh. Hana akan menjadi tutor sebaya mulai Oktober, 2x seminggu setelah jam pelajaran.',
          status: 'CLOSED',
          createdAt: '2026-09-19T13:45:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-007',
          studentId: 'STU-002',
          studentName: 'Siti Nurhaliza',
          category: 'PHONE_CALL',
          direction: 'INBOUND',
          title: 'Orang tua menanyakan persiapan class meeting',
          summary:
            'Orang tua Siti Nurhaliza menelepon untuk menanyakan persiapan class meeting di mana Siti terpilih sebagai ketua panitia. Mereka ingin tahu apakah ada yang perlu disiapkan dari rumah dan apakah ada biaya yang diperlukan. Juga menanyakan jadwal kegiatan agar bisa mendukung dari rumah.',
          communicationDate: '2026-10-11',
          followUpNote:
            'Sudah dijelaskan detail persiapan. Orang tua akan membantu transportasi untuk kegiatan luar sekolah. Tidak ada biaya tambahan.',
          status: 'CLOSED',
          createdAt: '2026-10-11T10:30:00.000Z',
          createdBy: 'USR-002',
        },
        {
          communicationId: 'COMM-008',
          studentId: 'STU-003',
          studentName: 'Budi Santoso',
          category: 'MEETING',
          direction: 'BIDIRECTIONAL',
          title: 'Diskusi motivasi belajar dengan orang tua',
          summary:
            'Pertemuan dengan orang tua Budi Santoso untuk mendiskusikan penurunan motivasi belajar. Diketahui ada masalah keluarga yang mempengaruhi konsentrasi Budi. Sekolah dan orang tua sepakat untuk memberikan dukungan psikologis dan memantau perkembangan selama sebulan ke depan.',
          communicationDate: '2026-10-18',
          followUpNote: '',
          status: 'PENDING_FOLLOW_UP',
          createdAt: '2026-10-18T11:00:00.000Z',
          createdBy: 'USR-002',
        },
      ],
    },
  },
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export class HomeroomCommunicationMockAdapter {
  /**
   * Mengambil daftar log komunikasi wali kelas berdasarkan context.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Array>} Array log komunikasi
   */
  async fetchCommunicationsByHomeroom(context) {
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
    const classData = MOCK_COMMUNICATIONS_DB[classId] || {}
    const yearData = classData[academicYearId] || {}
    const communications = yearData[semesterId] || []

    return mapHomeroomCommunications(communications)
  }

  /**
   * Menyimpan log komunikasi baru ke mock database.
   * @param {Object} communication - Objek log komunikasi baru (tanpa communicationId).
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Log komunikasi yang tersimpan dengan communicationId.
   */
  async saveCommunication(communication, context) {
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

    communicationIdCounter++
    const communicationId = `COMM-${communicationIdCounter}`

    const savedCommunication = {
      ...communication,
      communicationId: communicationId,
    }

    if (!MOCK_COMMUNICATIONS_DB[classId]) {
      MOCK_COMMUNICATIONS_DB[classId] = {}
    }
    if (!MOCK_COMMUNICATIONS_DB[classId][academicYearId]) {
      MOCK_COMMUNICATIONS_DB[classId][academicYearId] = {}
    }
    if (!MOCK_COMMUNICATIONS_DB[classId][academicYearId][semesterId]) {
      MOCK_COMMUNICATIONS_DB[classId][academicYearId][semesterId] = []
    }

    MOCK_COMMUNICATIONS_DB[classId][academicYearId][semesterId].push(savedCommunication)

    return mapHomeroomCommunication(savedCommunication)
  }

  /**
   * Mengupdate log komunikasi yang ada di mock database.
   * Digunakan untuk update status follow-up.
   * @param {Object} communication - Objek log komunikasi yang sudah diupdate.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Log komunikasi yang terupdate.
   */
  async updateCommunication(communication, context) {
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
    const classData = MOCK_COMMUNICATIONS_DB[classId] || {}
    const yearData = classData[academicYearId] || {}
    const communications = yearData[semesterId] || []

    const commIndex = communications.findIndex(
      (c) => c.communicationId === communication.communicationId,
    )

    if (commIndex === -1) {
      throw new Error(`Communication with ID ${communication.communicationId} not found.`)
    }

    communications[commIndex] = { ...communication }

    return mapHomeroomCommunication(communications[commIndex])
  }
}
