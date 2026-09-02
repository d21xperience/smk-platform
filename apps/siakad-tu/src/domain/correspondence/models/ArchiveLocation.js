export function createArchiveLocation(data = {}) {
  return {
    id: data.id || null,
    schoolId: data.schoolId || null,
    name: data.name || '',
    description: data.description || '',
    createdAt: data.createdAt || null,
    updatedAt: data.updatedAt || null,
  }
}
