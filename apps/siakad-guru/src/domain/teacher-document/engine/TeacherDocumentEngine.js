import { TeacherDocument } from '../models/TeacherDocument.js'
import { isValidDocumentType, isMimeTypeAllowed, getMaxFileSize } from '../models/DocumentType.js'
import {
  isValidDocumentScope,
  scopeRequiresAcademicYear,
  scopeRequiresSemester,
} from '../models/DocumentScope.js'

export class TeacherDocumentEngine {
  static validateDocument(document) {
    if (!document || !(document instanceof TeacherDocument)) {
      throw new Error('Invalid teacher document object')
    }
    if (!document.isValid()) {
      throw new Error('Teacher document is missing required fields')
    }
    return true
  }

  static validateDocumentType(documentType) {
    if (!isValidDocumentType(documentType)) {
      throw new Error(`Invalid document type: ${documentType}`)
    }
    return true
  }

  static validateScopeRules({ documentScope, academicYearId, semesterId }) {
    if (!isValidDocumentScope(documentScope)) {
      throw new Error(`Invalid document scope: ${documentScope}`)
    }

    if (scopeRequiresAcademicYear(documentScope) && !academicYearId) {
      throw new Error(`Scope "${documentScope}" requires academicYearId`)
    }

    if (scopeRequiresSemester(documentScope)) {
      if (!academicYearId) {
        throw new Error(`Scope "${documentScope}" requires academicYearId`)
      }
      if (!semesterId) {
        throw new Error(`Scope "${documentScope}" requires semesterId`)
      }
    }

    if (!scopeRequiresAcademicYear(documentScope)) {
      if (academicYearId || semesterId) {
        throw new Error(`Scope "${documentScope}" must not have academicYearId or semesterId`)
      }
    }

    return true
  }

  static validateFileMetadata({ documentType, fileName, fileSize, mimeType }) {
    if (!fileName || fileName.trim().length === 0) {
      throw new Error('File name is required')
    }

    if (!mimeType || mimeType.trim().length === 0) {
      throw new Error('MIME type is required')
    }

    if (!isMimeTypeAllowed(documentType, mimeType)) {
      throw new Error(`MIME type "${mimeType}" is not allowed for document type "${documentType}"`)
    }

    if (typeof fileSize !== 'number' || fileSize <= 0) {
      throw new Error('File size must be a positive number')
    }

    const maxSize = getMaxFileSize()
    if (fileSize > maxSize) {
      throw new Error(`File size exceeds maximum allowed (${maxSize} bytes)`)
    }

    return true
  }

  static validateOwnership(document, userId, role, permissions) {
    TeacherDocumentEngine.validateDocument(document)

    if (!userId) {
      throw new Error('User ID is required for ownership validation')
    }

    const isAdmin = role === 'admin' || role === 'superadmin'
    const hasManageAllPermission =
      permissions && permissions.includes('teacher_document.manage_all')

    if (isAdmin || hasManageAllPermission) {
      return true
    }

    if (!document.isOwnedBy(userId)) {
      throw new Error('You do not have permission to modify this document')
    }

    return true
  }

  static validateCreateInput(input) {
    if (!input.teacherId) {
      throw new Error('Teacher ID is required')
    }
    if (!input.schoolId) {
      throw new Error('School ID is required')
    }

    TeacherDocumentEngine.validateDocumentType(input.documentType)
    TeacherDocumentEngine.validateScopeRules({
      documentScope: input.documentScope,
      academicYearId: input.academicYearId,
      semesterId: input.semesterId,
    })
    TeacherDocumentEngine.validateFileMetadata({
      documentType: input.documentType,
      fileName: input.fileName,
      fileSize: input.fileSize,
      mimeType: input.mimeType,
    })

    return true
  }
}
