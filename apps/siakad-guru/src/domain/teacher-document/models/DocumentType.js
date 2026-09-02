export const DOCUMENT_TYPE = {
  CERTIFICATE: 'certificate',
  TEACHING_DECREE: 'teaching_decree',
  ASSIGNMENT_DECREE: 'assignment_decree',
  DEGREE: 'degree',
  TRAINING: 'training',
  OTHER: 'other',
}

const DOCUMENT_TYPE_LABELS = {
  [DOCUMENT_TYPE.CERTIFICATE]: 'Sertifikat',
  [DOCUMENT_TYPE.TEACHING_DECREE]: 'SK Mengajar',
  [DOCUMENT_TYPE.ASSIGNMENT_DECREE]: 'SK Pembagian Tugas',
  [DOCUMENT_TYPE.DEGREE]: 'Ijazah',
  [DOCUMENT_TYPE.TRAINING]: 'Dokumen Pelatihan',
  [DOCUMENT_TYPE.OTHER]: 'Dokumen Lainnya',
}

const ALLOWED_MIME_TYPES = {
  [DOCUMENT_TYPE.CERTIFICATE]: ['application/pdf', 'image/jpeg', 'image/png'],
  [DOCUMENT_TYPE.TEACHING_DECREE]: ['application/pdf'],
  [DOCUMENT_TYPE.ASSIGNMENT_DECREE]: ['application/pdf'],
  [DOCUMENT_TYPE.DEGREE]: ['application/pdf', 'image/jpeg', 'image/png'],
  [DOCUMENT_TYPE.TRAINING]: ['application/pdf', 'image/jpeg', 'image/png'],
  [DOCUMENT_TYPE.OTHER]: [
    'application/pdf',
    'image/jpeg',
    'image/png',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
  ],
}

export function isValidDocumentType(type) {
  return Object.values(DOCUMENT_TYPE).includes(type)
}

export function getDocumentTypeLabel(type) {
  return DOCUMENT_TYPE_LABELS[type] || type
}

export function getAllDocumentTypes() {
  return Object.values(DOCUMENT_TYPE).map((type) => ({
    value: type,
    label: DOCUMENT_TYPE_LABELS[type],
  }))
}

export function isMimeTypeAllowed(documentType, mimeType) {
  const allowed = ALLOWED_MIME_TYPES[documentType]
  if (!allowed) return false
  return allowed.includes(mimeType)
}

export function getMaxFileSize() {
  return 10 * 1024 * 1024
}
