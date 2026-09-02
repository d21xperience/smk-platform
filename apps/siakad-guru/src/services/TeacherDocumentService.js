import { TeacherDocument } from '../domain/teacher-document/models/TeacherDocument.js'
import { TeacherDocumentEngine } from '../domain/teacher-document/engine/TeacherDocumentEngine.js'
import { TeacherDocumentCreatedEvent } from '../domain/teacher-document/events/TeacherDocumentCreated.js'
import { TeacherDocumentDeletedEvent } from '../domain/teacher-document/events/TeacherDocumentDeleted.js'

export class TeacherDocumentService {
  constructor({ teacherDocumentAdapter, eventDispatcher }) {
    this.teacherDocumentAdapter = teacherDocumentAdapter
    this.eventDispatcher = eventDispatcher
  }

  async listDocuments({ teacherId, documentScope, academicYearId, semesterId }) {
    const rawData = await this.teacherDocumentAdapter.listDocuments({
      teacherId,
      documentScope,
      academicYearId,
      semesterId,
    })

    return rawData.map((item) => new TeacherDocument(item))
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
    TeacherDocumentEngine.validateCreateInput({
      teacherId,
      documentType,
      documentScope,
      academicYearId,
      semesterId,
      fileName,
      fileSize,
      mimeType,
      schoolId,
    })

    const rawData = await this.teacherDocumentAdapter.saveDocument({
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
    })

    const document = new TeacherDocument(rawData)
    TeacherDocumentEngine.validateDocument(document)

    if (!id && this.eventDispatcher) {
      this.eventDispatcher.dispatch(new TeacherDocumentCreatedEvent(document))
    }

    return document
  }

  async deleteDocument({ documentId, userId, role, permissions }) {
    const rawData = await this.teacherDocumentAdapter.findById({ documentId })
    if (!rawData) {
      throw new Error(`Document not found: ${documentId}`)
    }

    const document = new TeacherDocument(rawData)
    TeacherDocumentEngine.validateOwnership(document, userId, role, permissions)

    await this.teacherDocumentAdapter.deleteDocument({ documentId })

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new TeacherDocumentDeletedEvent(document))
    }

    return { success: true, documentId }
  }

  async getDownloadUrl({ documentId, userId, role, permissions }) {
    const rawData = await this.teacherDocumentAdapter.findById({ documentId })
    if (!rawData) {
      throw new Error(`Document not found: ${documentId}`)
    }

    const document = new TeacherDocument(rawData)
    TeacherDocumentEngine.validateOwnership(document, userId, role, permissions)

    return await this.teacherDocumentAdapter.getDownloadUrl({ documentId })
  }
}
