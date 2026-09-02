import { EVENT_TEACHER_DOCUMENT_DELETED } from '../../../events/Type.js'

export class TeacherDocumentDeletedEvent {
  constructor(document) {
    this.type = EVENT_TEACHER_DOCUMENT_DELETED
    this.payload = {
      documentId: document.id,
      teacherId: document.teacherId,
      documentType: document.documentType,
      storageKey: document.storageKey,
      schoolId: document.schoolId,
      timestamp: new Date().toISOString(),
    }
  }
}
