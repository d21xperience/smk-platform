export const PROGRESS_TYPE = {
  NOTE: 'note',
  ACHIEVEMENT: 'achievement',
  VIOLATION: 'violation',
  COUNSELING: 'counseling',
}

export function isValidProgressType(type) {
  return Object.values(PROGRESS_TYPE).includes(type)
}

export function getProgressTypeLabel(type) {
  const labels = {
    [PROGRESS_TYPE.NOTE]: 'Catatan Guru',
    [PROGRESS_TYPE.ACHIEVEMENT]: 'Prestasi',
    [PROGRESS_TYPE.VIOLATION]: 'Pelanggaran',
    [PROGRESS_TYPE.COUNSELING]: 'Bimbingan Konseling',
  }
  return labels[type] || type
}
