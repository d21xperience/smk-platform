const MOCK_STORAGE_BASE_URL = 'https://storage.mock.local/teacher-documents'

let mockDocuments = []
let documentCounter = 1

function generateId() {
  return `TDOC-${String(documentCounter++).padStart(3, '0')}`
}

function generateStorageKey(documentId, fileName) {
  const timestamp = Date.now()
  return `teacher-documents/${documentId}/${timestamp}_${fileName}`
}

function generateMockSignedUrl(storageKey) {
  const expiresInMs = 60 * 60 * 1000
  return `${MOCK_STORAGE_BASE_URL}/${storageKey}?signature=mock-signature-${Date.now()}&expires=${Date.now() + expiresInMs}`
}

export class TeacherDocumentMockAdapter {
  async listDocuments({ teacherId, documentScope, academicYearId, semesterId }) {
    await new Promise((resolve) => setTimeout(resolve, 500))

    return mockDocuments.filter((doc) => {
      if (doc.status !== 'active') return false
      if (teacherId && doc.teacherId !== teacherId) return false
      if (documentScope && doc.documentScope !== documentScope) return false
      if (academicYearId && doc.academicYearId !== academicYearId) return false
      if (semesterId && doc.semesterId !== semesterId) return false
      return true
    })
  }

  async saveDocument({
    id,
    teacherId,
    teacherName,
    documentType,
    documentScope,
    academicYearId,
    academicYearName,
    semesterId,
    semesterName,
    fileName,
    fileSize,
    mimeType,
    schoolId,
  }) {
    await new Promise((resolve) => setTimeout(resolve, 600))

    const now = new Date().toISOString()

    if (id) {
      const existingIndex = mockDocuments.findIndex((d) => d.id === id)
      if (existingIndex === -1) {
        throw new Error(`Document not found: ${id}`)
      }

      const existing = mockDocuments[existingIndex]
      const updated = {
        ...existing,
        documentType,
        documentScope,
        academicYearId: academicYearId || null,
        academicYearName: academicYearName || null,
        semesterId: semesterId || null,
        semesterName: semesterName || null,
        fileName,
        fileSize,
        mimeType,
        updatedAt: now,
      }

      mockDocuments[existingIndex] = updated
      return updated
    }

    const newId = generateId()
    const storageKey = generateStorageKey(newId, fileName)

    const newDocument = {
      id: newId,
      teacherId,
      teacherName,
      documentType,
      documentScope,
      academicYearId: academicYearId || null,
      academicYearName: academicYearName || null,
      semesterId: semesterId || null,
      semesterName: semesterName || null,
      fileName,
      fileSize,
      mimeType,
      storageKey,
      status: 'active',
      createdAt: now,
      updatedAt: now,
      schoolId,
    }

    mockDocuments.push(newDocument)
    return newDocument
  }

  async deleteDocument({ documentId }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    const index = mockDocuments.findIndex((d) => d.id === documentId)
    if (index === -1) {
      throw new Error(`Document not found: ${documentId}`)
    }

    mockDocuments[index] = {
      ...mockDocuments[index],
      status: 'deleted',
      updatedAt: new Date().toISOString(),
    }

    return {
      success: true,
      documentId,
    }
  }

  async getDownloadUrl({ documentId }) {
    await new Promise((resolve) => setTimeout(resolve, 300))

    const document = mockDocuments.find((d) => d.id === documentId && d.status === 'active')
    if (!document) {
      throw new Error(`Document not found: ${documentId}`)
    }

    const expiresAt = new Date(Date.now() + 60 * 60 * 1000).toISOString()

    return {
      url: generateMockSignedUrl(document.storageKey),
      storageKey: document.storageKey,
      fileName: document.fileName,
      mimeType: document.mimeType,
      expiresAt,
    }
  }

  async findById({ documentId }) {
    await new Promise((resolve) => setTimeout(resolve, 200))
    const document = mockDocuments.find((d) => d.id === documentId && d.status === 'active')
    return document || null
  }
}

if (mockDocuments.length === 0) {
  const now = new Date().toISOString()
  mockDocuments = [
    {
      id: 'TDOC-001',
      teacherId: 'USR-001',
      teacherName: 'Budi Santoso, S.Pd.',
      documentType: 'degree',
      documentScope: 'global',
      academicYearId: null,
      academicYearName: null,
      semesterId: null,
      semesterName: null,
      fileName: 'Ijazah_S1_Matematika.pdf',
      fileSize: 2048576,
      mimeType: 'application/pdf',
      storageKey: 'teacher-documents/TDOC-001/ijazah_s1.pdf',
      status: 'active',
      createdAt: now,
      updatedAt: now,
      schoolId: 'SCH-001',
    },
    {
      id: 'TDOC-002',
      teacherId: 'USR-001',
      teacherName: 'Budi Santoso, S.Pd.',
      documentType: 'teaching_decree',
      documentScope: 'academic_year',
      academicYearId: 'AY-2026',
      academicYearName: '2026/2027',
      semesterId: null,
      semesterName: null,
      fileName: 'SK_Mengajar_2026_2027.pdf',
      fileSize: 1024000,
      mimeType: 'application/pdf',
      storageKey: 'teacher-documents/TDOC-002/sk_mengajar.pdf',
      status: 'active',
      createdAt: now,
      updatedAt: now,
      schoolId: 'SCH-001',
    },
    {
      id: 'TDOC-003',
      teacherId: 'USR-001',
      teacherName: 'Budi Santoso, S.Pd.',
      documentType: 'assignment_decree',
      documentScope: 'semester',
      academicYearId: 'AY-2026',
      academicYearName: '2026/2027',
      semesterId: '20261',
      semesterName: 'Semester 1',
      fileName: 'SK_Pembagian_Tugas_Sem1_2026.pdf',
      fileSize: 819200,
      mimeType: 'application/pdf',
      storageKey: 'teacher-documents/TDOC-003/sk_pembagian_tugas.pdf',
      status: 'active',
      createdAt: now,
      updatedAt: now,
      schoolId: 'SCH-001',
    },
    {
      id: 'TDOC-004',
      teacherId: 'USR-001',
      teacherName: 'Budi Santoso, S.Pd.',
      documentType: 'certificate',
      documentScope: 'global',
      academicYearId: null,
      academicYearName: null,
      semesterId: null,
      semesterName: null,
      fileName: 'Sertifikat_Pelatihan_Kurikulum_Merdeka.pdf',
      fileSize: 1536000,
      mimeType: 'application/pdf',
      storageKey: 'teacher-documents/TDOC-004/sertifikat_pelatihan.pdf',
      status: 'active',
      createdAt: now,
      updatedAt: now,
      schoolId: 'SCH-001',
    },
  ]
  documentCounter = 5
}
