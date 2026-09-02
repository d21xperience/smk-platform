import { IncomingLetterStatus } from './IncomingLetterStatus.js'

export function createIncomingLetter(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    academicYearId: data.academicYearId || null,
    semesterId: data.semesterId || null,
    noAgenda: data.noAgenda || '',
    nomorSuratAsal: data.nomorSuratAsal || '',
    asal: data.asal || '',
    perihal: data.perihal || '',
    tanggalSurat: data.tanggalSurat || null,
    tanggalDiterima: data.tanggalDiterima || null,
    classificationCode: data.classificationCode || '',
    fileUrl: data.fileUrl || '',
    status: data.status || IncomingLetterStatus.BARU,
    createdBy: data.createdBy || null,
    createdAt: data.createdAt || null,
    updatedAt: data.updatedAt || null,
  }
}
