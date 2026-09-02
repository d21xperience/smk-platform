import {
  createIncomingListItem,
  createIncomingDetail,
} from '../../contracts/incomingLetterContract.js'
import { createIncomingLetter } from '../../domain/correspondence/models/IncomingLetter.js'
import { IncomingLetterStatus } from '../../domain/correspondence/models/IncomingLetterStatus.js'

// In-memory database simulation
let db = [
  createIncomingLetter({
    id: '1',
    schoolId: 'SCH-1',
    noAgenda: '001/MAY/2026',
    nomorSuratAsal: '123/DP/2026',
    asal: 'Dinas Pendidikan',
    perihal: 'Undangan Rapat Koordinasi',
    tanggalDiterima: '2026-05-20',
    classificationCode: '421.3',
    status: IncomingLetterStatus.BARU,
    fileUrl: 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf',
  }),
  createIncomingLetter({
    id: '2',
    schoolId: 'SCH-1',
    noAgenda: '002/MAY/2026',
    nomorSuratAsal: '456/PR/2026',
    asal: 'PT. XYZ',
    perihal: 'Permohonan Kerjasama DUDI',
    tanggalDiterima: '2026-05-21',
    classificationCode: '421.4',
    status: IncomingLetterStatus.DIDISPOSISIKAN,
    fileUrl: 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf',
  }),
]

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const IncomingLetterMockAdapter = {
  async list(context, query) {
    await delay(300)
    let filtered = db.filter((l) => l.schoolId === context.schoolId)

    if (query.search) {
      const s = query.search.toLowerCase()
      filtered = filtered.filter(
        (l) => l.perihal.toLowerCase().includes(s) || l.asal.toLowerCase().includes(s),
      )
    }
    if (query.status) filtered = filtered.filter((l) => l.status === query.status)

    return filtered.map(createIncomingListItem)
  },

  async getById(context, id) {
    await delay(200)
    const letter = db.find((l) => l.id === id && l.schoolId === context.schoolId)
    return letter ? createIncomingDetail(letter) : null
  },

  async create(context, payload) {
    await delay(400)
    const newId = String(Date.now())
    const noAgenda = `00${db.length + 1}/${new Date().toLocaleString('en', { month: 'short' }).toUpperCase()}/${new Date().getFullYear()}`

    const newLetter = createIncomingLetter({
      ...payload,
      id: newId,
      schoolId: context.schoolId,
      academicYearId: context.academicYearId,
      semesterId: context.semesterId,
      noAgenda: noAgenda,
      tanggalDiterima: new Date().toISOString().split('T')[0],
      createdBy: context.userId,
      status: IncomingLetterStatus.BARU,
    })

    db.push(newLetter)
    return createIncomingDetail(newLetter)
  },

  async updateStatus(context, id, newStatus) {
    await delay(300)
    const letter = db.find((l) => l.id === id)
    if (letter) letter.status = newStatus
    return letter ? createIncomingDetail(letter) : null
  },
}
