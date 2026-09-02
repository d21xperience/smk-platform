export const DOCUMENT_SCOPE = {
  GLOBAL: 'global',
  ACADEMIC_YEAR: 'academic_year',
  SEMESTER: 'semester',
}

const DOCUMENT_SCOPE_LABELS = {
  [DOCUMENT_SCOPE.GLOBAL]: 'Global (Seluruh Periode)',
  [DOCUMENT_SCOPE.ACADEMIC_YEAR]: 'Tahun Pelajaran',
  [DOCUMENT_SCOPE.SEMESTER]: 'Semester',
}

export function isValidDocumentScope(scope) {
  return Object.values(DOCUMENT_SCOPE).includes(scope)
}

export function getDocumentScopeLabel(scope) {
  return DOCUMENT_SCOPE_LABELS[scope] || scope
}

export function getAllDocumentScopes() {
  return Object.values(DOCUMENT_SCOPE).map((scope) => ({
    value: scope,
    label: DOCUMENT_SCOPE_LABELS[scope],
  }))
}

export function scopeRequiresAcademicYear(scope) {
  return scope === DOCUMENT_SCOPE.ACADEMIC_YEAR || scope === DOCUMENT_SCOPE.SEMESTER
}

export function scopeRequiresSemester(scope) {
  return scope === DOCUMENT_SCOPE.SEMESTER
}
