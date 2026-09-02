export function createOutgoingListItem(data = {}) {
  return {
    id: data.id,
    nomorSurat: data.nomorSurat,
    kategori: data.kategori,
    jenis: data.jenis,
    tujuan: data.tujuan,
    perihal: data.perihal,
    tanggal: data.tanggal,
    status: data.status,
  }
}

export function createOutgoingDetail(data = {}) {
  return {
    ...createOutgoingListItem(data),
    templateData: data.templateData,
    fileUrl: data.fileUrl,
    createdBy: data.createdBy,
    createdAt: data.createdAt,
    updatedAt: data.updatedAt,
  }
}

export function createOutgoingCreatePayload(data = {}) {
  return {
    kategori: data.kategori,
    jenis: data.jenis,
    tujuan: data.tujuan,
    perihal: data.perihal,
    tanggal: data.tanggal,
    templateData: data.templateData || {},
    fileUrl: data.fileUrl || '',
  }
}

export function createOutgoingListQuery(data = {}) {
  return {
    kategori: data.kategori || null,
    search: data.search || '',
    status: data.status || null,
    page: data.page || 1,
    limit: data.limit || 10,
  }
}

export function createNextNumberRequest(data = {}) {
  return {
    kategori: data.kategori,
    jenis: data.jenis,
    academicYearId: data.academicYearId,
  }
}

export function createNextNumberResponse(data = {}) {
  return {
    nomorSurat: data.nomorSurat,
  }
}
