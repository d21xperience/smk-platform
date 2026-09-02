export const DOCUMENT_STATUS = {
  ACTIVE: 'active',
  ARCHIVED: 'archived',
}

export class TeacherDocument {
  constructor({
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
    storageKey,
    status,
    createdAt,
    updatedAt,
    schoolId,
  }) {
    this.id = id
    this.teacherId = teacherId
    this.teacherName = teacherName
    this.documentType = documentType
    this.documentScope = documentScope
    this.academicYearId = academicYearId || null
    this.academicYearName = academicYearName || null
    this.semesterId = semesterId || null
    this.semesterName = semesterName || null
    this.fileName = fileName
    this.fileSize = fileSize
    this.mimeType = mimeType
    this.storageKey = storageKey
    this.status = status || DOCUMENT_STATUS.ACTIVE
    this.createdAt = createdAt
    this.updatedAt = updatedAt
    this.schoolId = schoolId
  }

  isOwnedBy(userId) {
    return this.teacherId === userId
  }

  isActive() {
    return this.status === DOCUMENT_STATUS.ACTIVE
  }

  isValid() {
    return !!(
      this.teacherId &&
      this.documentType &&
      this.documentScope &&
      this.fileName &&
      this.fileSize &&
      this.mimeType &&
      this.schoolId
    )
  }
}
