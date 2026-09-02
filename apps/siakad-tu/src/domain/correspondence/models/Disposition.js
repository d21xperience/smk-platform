import { DispositionStatus } from './DispositionStatus.js'

export function createDisposition(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    incomingLetterId: data.incomingLetterId || null,
    fromUserId: data.fromUserId || null,
    toUserId: data.toUserId || null,
    instruksi: data.instruksi || '',
    deadline: data.deadline || null,
    status: data.status || DispositionStatus.PENDING,
    laporanTindakLanjut: data.laporanTindakLanjut || '',
    buktiUrl: data.buktiUrl || '',
    createdAt: data.createdAt || null,
    updatedAt: data.updatedAt || null,
  }
}
