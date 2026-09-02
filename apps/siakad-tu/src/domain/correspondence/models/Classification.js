export function createClassification(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    parentId: data.parentId || null,
    code: data.code || '',
    label: data.label || '',
    createdAt: data.createdAt || null,
    updatedAt: data.updatedAt || null,
  }
}
