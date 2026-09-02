import { defineStore } from 'pinia'
import { teacherDocumentService } from '../boot/services.js'

export const useTeacherDocumentStore = defineStore('teacherDocument', {
  state: () => ({
    documents: [],
    selectedDocument: null,
    loading: false,
    saving: false,
    deleting: false,
    error: null,
    formDialogOpen: false,
    editingDocument: null,
  }),

  getters: {
    hasDocuments: (state) => state.documents.length > 0,
    activeDocuments: (state) => state.documents.filter((d) => d.status === 'active'),
    isFormDialogOpen: (state) => state.formDialogOpen,
    isEditing: (state) => !!state.editingDocument,
  },

  actions: {
    async listDocuments({ teacherId, documentScope, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        this.documents = await teacherDocumentService.listDocuments({
          teacherId,
          documentScope,
          academicYearId,
          semesterId,
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async saveDocument(payload) {
      this.saving = true
      this.error = null
      try {
        const saved = await teacherDocumentService.saveDocument(payload)
        const index = this.documents.findIndex((d) => d.id === saved.id)
        if (index !== -1) {
          this.documents[index] = saved
        } else {
          this.documents.unshift(saved)
        }
        return saved
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.saving = false
      }
    },

    async deleteDocument({ documentId, userId, role, permissions }) {
      this.deleting = true
      this.error = null
      try {
        await teacherDocumentService.deleteDocument({ documentId, userId, role, permissions })
        this.documents = this.documents.filter((d) => d.id !== documentId)
        if (this.selectedDocument && this.selectedDocument.id === documentId) {
          this.selectedDocument = null
        }
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.deleting = false
      }
    },

    async getDownloadUrl({ documentId, userId, role, permissions }) {
      this.loading = true
      this.error = null
      try {
        return await teacherDocumentService.getDownloadUrl({
          documentId,
          userId,
          role,
          permissions,
        })
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    setSelectedDocument(document) {
      this.selectedDocument = document
    },

    openFormDialog(document) {
      this.editingDocument = document || null
      this.formDialogOpen = true
    },

    closeFormDialog() {
      this.formDialogOpen = false
      this.editingDocument = null
    },
  },
})
