// FILE: src/utils/contextFormatter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

export function formatAcademicYearLabel(academicYearId) {
  if (!academicYearId) {
    return '-'
  }

  const normalized = String(academicYearId).replace('AY-', '')
  const year = Number(normalized)

  if (Number.isNaN(year)) {
    return academicYearId
  }

  return `${year}/${year + 1}`
}

export function formatSemesterLabel(semesterId) {
  if (!semesterId) {
    return '-'
  }

  const value = String(semesterId)
  const semesterCode = value.slice(-1)

  if (semesterCode === '1') {
    return 'Semester 1'
  }

  if (semesterCode === '2') {
    return 'Semester 2'
  }

  return `Semester ${semesterCode}`
}

export function formatContextLabel(context) {
  if (!context || !context.academicYearId || !context.semesterId) {
    return 'Context belum siap'
  }

  return `${formatAcademicYearLabel(context.academicYearId)} - ${formatSemesterLabel(context.semesterId)}`
}
