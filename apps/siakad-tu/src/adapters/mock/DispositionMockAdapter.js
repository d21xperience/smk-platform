import { createDispositionListItem } from '../../contracts/dispositionContract.js'
import { createDisposition } from '../../domain/correspondence/models/Disposition.js'
import { DispositionStatus } from '../../domain/correspondence/models/DispositionStatus.js'

let db = [
  createDisposition({
    id: '101',
    schoolId: 'SCH-1',
    incomingLetterId: '1',
    fromUserId: 'USER-KEPSEK',
    toUserId: 'USER-WAKUR',
    instruksi: 'Hadiri rapat koordinasi',
    deadline: '2026-05-25',
    status: DispositionStatus.PENDING,
    createdAt: '2026-05-22',
  }),
  createDisposition({
    id: '102',
    schoolId: 'SCH-1',
    incomingLetterId: '2',
    fromUserId: 'USER-KEPSEK',
    toUserId: 'USER-WAKUR',
    instruksi: 'Siapkan laporan kurikulum',
    deadline: '2026-05-28',
    status: DispositionStatus.DIBACA,
    createdAt: '2026-05-23',
  }),
]

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

// Dummy mapping for UI display
const userMap = {
  'USER-KEPSEK': 'Kepala Sekolah',
  'USER-WAKUR': 'Waka Kurikulum',
  'USER-WAKES': 'Waka Kesiswaan',
  'USER-KATU': 'Kepala TU',
}

export const DispositionMockAdapter = {
  async listPendingForKepsek(context) {
    console.log(context)

    await delay(300)
    // In real app, join with IncomingLetter table
    return [
      {
        id: '1',
        perihalSurat: 'Undangan Rapat Dinas',
        asal: 'Dinas Pendidikan',
        tanggalDiterima: '2026-05-20',
      },
      {
        id: '2',
        perihalSurat: 'Permohonan Magang',
        asal: 'PT. ABC',
        tanggalDiterima: '2026-05-21',
      },
    ]
  },

  async listTasksForStaf(context, userId) {
    await delay(300)
    const tasks = db.filter((d) => d.toUserId === userId)
    return tasks.map((d) =>
      createDispositionListItem({
        ...d,
        dari: userMap[d.fromUserId] || 'Unknown',
        toUser: userMap[d.toUserId] || 'Unknown',
        perihalSurat: 'Perihal Surat (Join dari DB)', // Simplified
        tanggal: d.createdAt,
      }),
    )
  },

  async create(context, payload) {
    await delay(400)
    const newDisp = createDisposition({
      ...payload,
      id: String(Date.now()),
      schoolId: context.schoolId,
      fromUserId: context.userId,
      status: DispositionStatus.PENDING,
    })
    db.push(newDisp)
    return newDisp
  },

  async submitFollowUp(context, dispositionId, payload) {
    await delay(400)
    const disp = db.find((d) => d.id === dispositionId)
    if (disp) {
      disp.laporanTindakLanjut = payload.laporan
      disp.buktiUrl = payload.buktiUrl || 'mock-bucket://bukti.pdf'
      disp.status = DispositionStatus.SELESAI
    }
    return disp
  },

  async getStafOptions(context) {
    console.log(context)
    await delay(100)
    return ['Waka Kurikulum', 'Waka Kesiswaan', 'Kepala TU', 'Guru BK']
  },
}
