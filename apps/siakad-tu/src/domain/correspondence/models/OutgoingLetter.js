import { OutgoingLetterStatus } from './OutgoingLetterStatus.js'

export function createOutgoingLetter(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    academicYearId: data.academicYearId || null,
    semesterId: data.semesterId || null,
    nomorSurat: data.nomorSurat || '',
    kategori: data.kategori || null, // SISWA, PTK, KEDINASAN
    jenis: data.jenis || '', // pkl, aktif, sppd, dll
    tujuan: data.tujuan || '',
    perihal: data.perihal || '',
    tanggal: data.tanggal || null,
    status: data.status || OutgoingLetterStatus.DRAF,
    // templateData menyimpan JSON untuk field spesifik berdasarkan 'jenis' surat
    // Contoh: { dudi: 'PT Maju Jaya', siswaIds: [1, 2] } untuk surat PKL
    templateData: data.templateData || {},
    fileUrl: data.fileUrl || '',
    createdBy: data.createdBy || null,
    createdAt: data.createdAt || null,
    updatedAt: data.updatedAt || null,
  }
}
