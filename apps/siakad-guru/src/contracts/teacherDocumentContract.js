export const TeacherDocumentContract = {
  listDocuments: {
    request: {
      teacherId: 'string',
      documentScope: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: [
      {
        id: 'string',
        teacherId: 'string',
        teacherName: 'string',
        documentType: 'string',
        documentScope: 'string',
        academicYearId: 'string',
        academicYearName: 'string',
        semesterId: 'string',
        semesterName: 'string',
        fileName: 'string',
        fileSize: 'number',
        mimeType: 'string',
        storageKey: 'string',
        status: 'string',
        createdAt: 'string',
        updatedAt: 'string',
        schoolId: 'string',
      },
    ],
  },

  saveDocument: {
    request: {
      id: 'string',
      teacherId: 'string',
      teacherName: 'string',
      documentType: 'string',
      documentScope: 'string',
      academicYearId: 'string',
      academicYearName: 'string',
      semesterId: 'string',
      semesterName: 'string',
      fileName: 'string',
      fileSize: 'number',
      mimeType: 'string',
      schoolId: 'string',
    },
    response: {
      id: 'string',
      teacherId: 'string',
      teacherName: 'string',
      documentType: 'string',
      documentScope: 'string',
      academicYearId: 'string',
      academicYearName: 'string',
      semesterId: 'string',
      semesterName: 'string',
      fileName: 'string',
      fileSize: 'number',
      mimeType: 'string',
      storageKey: 'string',
      status: 'string',
      createdAt: 'string',
      updatedAt: 'string',
      schoolId: 'string',
    },
  },

  deleteDocument: {
    request: {
      documentId: 'string',
    },
    response: {
      success: 'boolean',
      documentId: 'string',
    },
  },

  getDownloadUrl: {
    request: {
      documentId: 'string',
    },
    response: {
      url: 'string',
      storageKey: 'string',
      fileName: 'string',
      mimeType: 'string',
      expiresAt: 'string',
    },
  },
}
