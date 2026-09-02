export function createIncomingListItem(data = {}) {
  return {
    id: data.id,
    noAgenda: data.noAgenda,
    nomorSuratAsal: data.nomorSuratAsal,
    asal: data.asal,
    perihal: data.perihal,
    tanggalDiterima: data.tanggalDiterima,
    classificationCode: data.classificationCode,
    status: data.status,
  }
}

export function createIncomingDetail(data = {}) {
  return {
    ...createIncomingListItem(data),
    tanggalSurat: data.tanggalSurat,
    fileUrl: data.fileUrl,
    createdBy: data.createdBy,
    createdAt: data.createdAt,
    updatedAt: data.updatedAt,
  }
}

export function createIncomingCreatePayload(data = {}) {
  return {
    nomorSuratAsal: data.nomorSuratAsal,
    asal: data.asal,
    perihal: data.perihal,
    tanggalSurat: data.tanggalSurat,
    tanggalDiterima: data.tanggalDiterima,
    classificationCode: data.classificationCode,
    fileUrl: data.fileUrl,
  }
}

export function createIncomingListQuery(data = {}) {
  return {
    search: data.search || '',
    status: data.status || null,
    classificationCode: data.classificationCode || null,
    page: data.page || 1,
    limit: data.limit || 10,
  }
}
