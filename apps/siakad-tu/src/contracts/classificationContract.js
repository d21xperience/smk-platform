export function createClassificationNode(data = {}) {
  return {
    id: data.id,
    parentId: data.parentId,
    code: data.code,
    label: data.label,
    children: data.children || [],
  }
}

export function createArchiveLocationItem(data = {}) {
  return {
    id: data.id,
    name: data.name,
    description: data.description,
  }
}

export function createClassificationCreatePayload(data = {}) {
  return {
    parentId: data.parentId || null,
    code: data.code,
    label: data.label,
  }
}

export function createLocationCreatePayload(data = {}) {
  return {
    name: data.name,
    description: data.description,
  }
}

// BARU: Payload untuk update lokasi
export function createLocationUpdatePayload(data = {}) {
  return {
    name: data.name,
    description: data.description,
  }
}
