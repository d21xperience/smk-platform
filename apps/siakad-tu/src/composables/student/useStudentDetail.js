// apps/siakad-tu/src/composables/student/useStudentDetail.js

import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'
import { studentApi } from '@/adapters/api/studentApi.js' // Tambahkan ini
import { useQuasar } from 'quasar' // Untuk notifikasi

export function useStudentDetail() {
  const router = useRouter()
  const route = useRoute()
  const $q = useQuasar()
  const studentDetailStore = useStudentDetailStore()

  // === LOCAL UI STATE ===
  const showNotification = ref(false)
  const notificationMessage = ref('')
  const notificationType = ref('success')
  const documents = ref([])
  const isUploading = ref(false)

  // === COMPUTED ===
  const studentId = computed(() => route.params.id)
  const student = computed(() => studentDetailStore.getById(studentId.value))
  const isLoading = computed(() => studentDetailStore.isLoading(studentId.value))
  const error = computed(() => studentDetailStore.getError(studentId.value))

  // === HELPER FUNCTIONS ===
  function formatDate(dateString) {
    if (!dateString) return '-'
    return new Date(dateString).toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  }

  function getStatusColor(status) {
    const colors = {
      ACTIVE: 'positive',
      Aktif: 'positive',
      GRADUATED: 'blue',
      TRANSFERRED: 'orange',
    }
    return colors[status] || 'grey'
  }

  // === ACTIONS ===
  function editStudent() {
    router.push(`/kesiswaan/siswa/edit/${studentId.value}`)
  }

  function goBack() {
    router.push('/kesiswaan/siswa')
  }

  function loadData() {
    if (studentId.value) {
      studentDetailStore.fetchStudentById(studentId.value)
    }
    fetchDocuments()
  }

  async function fetchDocuments() {
    if (!studentId.value) return
    try {
      const res = await studentApi.getStudentDocuments(studentId.value)
      if (res.success) {
        documents.value = res.data
      }
    } catch (err) {
      console.error('Gagal mengambil dokumen:', err)
    }
  }

  async function uploadDocument(file, jenisBerkas) {
    isUploading.value = true
    try {
      const formData = new FormData()
      formData.append('file', file)
      formData.append('jenis_berkas', jenisBerkas)
      
      const res = await studentApi.uploadStudentDocument(studentId.value, formData)
      if (res.success) {
        $q.notify({ type: 'positive', message: 'Berkas berhasil diupload' })
        fetchDocuments()
      } else {
        $q.notify({ type: 'negative', message: res.error.message })
      }
    } catch (err) {
      console.error(err)
      $q.notify({ type: 'negative', message: 'Gagal upload' })
    } finally {
      isUploading.value = false
    }
  }

  async function deleteDocument(docId) {
    try {
      const res = await studentApi.deleteStudentDocument(docId)
      if (res.success) {
        $q.notify({ type: 'positive', message: 'Berkas dihapus' })
        fetchDocuments()
      }
    } catch (err) {
      console.error('Error delete:', err)
    }
  }

  // === RETURN ===
  return {
    student,
    isLoading,
    error,
    documents,
    isUploading,
    showNotification,
    notificationMessage,
    notificationType,
    studentId,
    formatDate,
    getStatusColor,
    editStudent,
    goBack,
    loadData,
    uploadDocument,
    deleteDocument
  }
}
