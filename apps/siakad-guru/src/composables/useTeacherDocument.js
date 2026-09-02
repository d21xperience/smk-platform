import { computed } from 'vue'
import { useTeacherDocumentStore } from '../stores/teacherDocumentStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'

export function useTeacherDocument() {
  const teacherDocumentStore = useTeacherDocumentStore()
  const contextStore = useContextStore()
  const authStore = useAuthStore()

  const currentContext = computed(() => contextStore.currentContext)
  const currentUser = computed(() => authStore.currentUser)
  const documents = computed(() => teacherDocumentStore.documents)
  const selectedDocument = computed(() => teacherDocumentStore.selectedDocument)
  const isLoading = computed(() => teacherDocumentStore.loading)
  const isSaving = computed(() => teacherDocumentStore.saving)
  const isDeleting = computed(() => teacherDocumentStore.deleting)
  const error = computed(() => teacherDocumentStore.error)
  const formDialogOpen = computed(() => teacherDocumentStore.formDialogOpen)
  const editingDocument = computed(() => teacherDocumentStore.editingDocument)
  const isEditing = computed(() => teacherDocumentStore.isEditing)

  const listDocuments = async (filters = {}) => {
    if (!currentContext.value || !currentUser.value) return
    await teacherDocumentStore.listDocuments({
      teacherId: currentUser.value.id,
      documentScope: filters.documentScope,
      academicYearId: filters.academicYearId,
      semesterId: filters.semesterId,
    })
  }

  const saveDocument = async (payload) => {
    if (!currentContext.value || !currentUser.value) return
    return await teacherDocumentStore.saveDocument({
      ...payload,
      teacherId: currentUser.value.id,
      teacherName: currentUser.value.name,
      schoolId: currentContext.value.schoolId,
    })
  }

  const deleteDocument = async (documentId) => {
    if (!currentContext.value || !currentUser.value) return
    await teacherDocumentStore.deleteDocument({
      documentId,
      userId: currentUser.value.id,
      role: currentUser.value.role,
      permissions: currentUser.value.permissions,
    })
  }

  const downloadDocument = async (documentId) => {
    if (!currentContext.value || !currentUser.value) return
    const downloadInfo = await teacherDocumentStore.getDownloadUrl({
      documentId,
      userId: currentUser.value.id,
      role: currentUser.value.role,
      permissions: currentUser.value.permissions,
    })
    return downloadInfo
  }

  const openFormDialog = (document) => {
    teacherDocumentStore.openFormDialog(document)
  }

  const closeFormDialog = () => {
    teacherDocumentStore.closeFormDialog()
  }

  const setSelectedDocument = (document) => {
    teacherDocumentStore.setSelectedDocument(document)
  }

  return {
    documents,
    selectedDocument,
    isLoading,
    isSaving,
    isDeleting,
    error,
    formDialogOpen,
    editingDocument,
    isEditing,
    listDocuments,
    saveDocument,
    deleteDocument,
    downloadDocument,
    openFormDialog,
    closeFormDialog,
    setSelectedDocument,
  }
}
