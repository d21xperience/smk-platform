import {
  createOutgoingListItem,
  createOutgoingDetail,
  createNextNumberResponse,
} from '../../contracts/outgoingLetterContract.js'
import { createOutgoingLetter } from '../../domain/correspondence/models/OutgoingLetter.js'
import { OutgoingLetterStatus } from '../../domain/correspondence/models/OutgoingLetterStatus.js'
import { LetterNumberEngine } from '../../domain/correspondence/engine/LetterNumberEngine.js'

let db = [
  createOutgoingLetter({
    id: '10',
    schoolId: 'SCH-1',
    nomorSurat: '422.1/001/SMK-01/2026',
    kategori: 'SISWA',
    jenis: 'pkl',
    tujuan: 'PT. Maju Jaya',
    perihal: 'Pengantar PKL',
    tanggal: '2026-05-15',
    status: OutgoingLetterStatus.TERKIRIM,
  }),
  createOutgoingLetter({
    id: '11',
    schoolId: 'SCH-1',
    nomorSurat: '800/001/SMK-01/2026',
    kategori: 'PTK',
    jenis: 'aktif_tugas',
    tujuan: 'Kepala Dinas',
    perihal: 'Aktif Tugas',
    tanggal: '2026-05-14',
    status: OutgoingLetterStatus.TERKIRIM,
  }),
  createOutgoingLetter({
    id: '12',
    schoolId: 'SCH-1',
    nomorSurat: '001/SMK-01/2026',
    kategori: 'KEDINASAN',
    jenis: 'undangan',
    tujuan: 'Orang Tua/Wali',
    perihal: 'Rapat Komite',
    tanggal: '2026-05-20',
    status: OutgoingLetterStatus.DRAF,
  }),
]

// Simulasi sequence counter per kategori & tahun ajaran
let sequenceCounter = { SISWA: 1, PTK: 1, KEDINASAN: 1 }

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const OutgoingLetterMockAdapter = {
  async listByCategory(context, kategori, query) {
    await delay(300)
    let filtered = db.filter((l) => l.schoolId === context.schoolId && l.kategori === kategori)
    if (query.search) {
      const s = query.search.toLowerCase()
      filtered = filtered.filter(
        (l) => l.perihal.toLowerCase().includes(s) || l.tujuan.toLowerCase().includes(s),
      )
    }
    return filtered.map(createOutgoingListItem)
  },

  async getById(context, id) {
    await delay(200)
    const letter = db.find((l) => l.id === id)
    return letter ? createOutgoingDetail(letter) : null
  },

  async create(context, payload) {
    await delay(400)
    const newId = String(Date.now())
    const newLetter = createOutgoingLetter({
      ...payload,
      id: newId,
      schoolId: context.schoolId,
      academicYearId: context.academicYearId,
      createdBy: context.userId,
      status: OutgoingLetterStatus.DRAF,
    })
    db.push(newLetter)
    return createOutgoingDetail(newLetter)
  },

  async getNextLetterNumber(context, kategori, jenis) {
    await delay(200)
    // Simulasi pengambilan nomor dari backend (mencegah race condition)
    const currentSeq = sequenceCounter[kategori] || 1
    sequenceCounter[kategori] = currentSeq + 1

    const kodeMap = {
      SISWA: { pkl: '422.1', aktif: '422.2', kesalahan_nama: '422.3', tunggakan: '422.4' },
      PTK: { aktif_tugas: '800', sppd: '900', mgmp: '800' },
      KEDINASAN: { undangan: '005', mou: '006', balasan: '007', sk: '008' }, // Simplified
    }

    const categoryCode = kodeMap[kategori]?.[jenis] || '000'
    const schoolCode = 'SMK-01' // Should come from context/school profile
    const year = new Date().getFullYear()

    const nomorSurat = LetterNumberEngine.format(categoryCode, currentSeq, schoolCode, year)
    return createNextNumberResponse({ nomorSurat })
  },

  // Dummy options for UI (In real app, this comes from Student/DUDI modules)
  async getStudentOptions(context) {
    console.log(context)
    await delay(100)
    return ['Ahmad Fauzi', 'Siti Nurhaliza', 'Budi Hartono']
  },
  async getDudiOptions(context) {
    console.log(context)

    await delay(100)
    return ['PT. Maju Jaya', 'CV. Karya Mandiri', 'Hotel Santika']
  },
}
