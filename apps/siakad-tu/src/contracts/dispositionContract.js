export function createDispositionListItem(data = {}) {
  return {
    id: data.id,
    incomingLetterId: data.incomingLetterId,
    perihalSurat: data.perihalSurat, // Denormalized from IncomingLetter for UI
    instruksi: data.instruksi,
    dari: data.dari, // Denormalized user name
    toUser: data.toUser, // Denormalized user name
    tanggal: data.tanggal,
    deadline: data.deadline,
    status: data.status,
  }
}

export function createDispositionCreatePayload(data = {}) {
  return {
    incomingLetterId: data.incomingLetterId,
    toUserId: data.toUserId,
    instruksi: data.instruksi,
    deadline: data.deadline,
  }
}

export function createFollowUpPayload(data = {}) {
  return {
    dispositionId: data.dispositionId,
    laporan: data.laporan,
    buktiUrl: data.buktiUrl,
  }
}
