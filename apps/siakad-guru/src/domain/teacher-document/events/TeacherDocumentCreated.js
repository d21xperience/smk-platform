import { EVENT_TEACHER_DOCUMENT_CREATED } from '../../../events/Type.js'

export class TeacherDocumentCreatedEvent {
  constructor(document) {
    this.type = EVENT_TEACHER_DOCUMENT_CREATED
    this.payload = {
      documentId: document.id,
      teacherId: document.teacherId,
      documentType: document.documentType,
      documentScope: document.documentScope,
      academicYearId: document.academicYearId,
      semesterId: document.semesterId,
      fileName: document.fileName,
      fileSize: document.fileSize,
      mimeType: document.mimeType,
      storageKey: document.storageKey,
      schoolId: document.schoolId,
      timestamp: new Date().toISOString(),
    }
  }
}
