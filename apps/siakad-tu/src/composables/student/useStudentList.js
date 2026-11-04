// apps/siakad-tu/src/composables/student/useStudentList.js

import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useStudentListStore } from '@/stores/student/studentListStore'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'
import { useContextStore } from '@/stores/contextStore.js'

export function useStudentList() {
  const router = useRouter()
  const studentListStore = useStudentListStore()
  const studentDetailStore = useStudentDetailStore()
  const contextStore = useContextStore()

  const localFilters = ref({
    status: null,
    search: '',
    tingkat: null,
    jurusan: null,
    kelas: null,
  })

  const showDeleteDialog = ref(false)
  const showNotification = ref(false)
  const notificationMessage = ref('')
  const notificationType = ref('success')
  const selectedStudent = ref(null)

  const statusOptions = [
    { label: 'Aktif', value: 'ACTIVE' },
    { label: 'Lulus', value: 'GRADUATED' },
    { label: 'Mutasi', value: 'TRANSFERRED' },
  ]

  const tablePagination = computed(() => ({
    page: studentListStore.pagination.page,
    rowsPerPage: studentListStore.pagination.limit,
    rowsNumber: studentListStore.pagination.total,
  }))

  watch(
    () => contextStore.selectedSemesterId,
    (newSemesterId, oldSemesterId) => {
      if (newSemesterId && newSemesterId !== oldSemesterId) {
        console.log(`[useStudentList] Context berubah ke ${newSemesterId}, reset cache & memuat ulang data...`)
        studentDetailStore.reset() 
        studentListStore.fetchStudents({ page: 1 })
      }
    }
  )

  function getStatusColor(status) {
    const colors = { ACTIVE: 'positive', Aktif: 'positive', GRADUATED: 'blue', TRANSFERRED: 'orange' }
    return colors[status] || 'grey'
  }

  function getStatusLabel(status) {
    const labels = { ACTIVE: 'Aktif', Aktif: 'Aktif', GRADUATED: 'Lulus', TRANSFERRED: 'Mutasi' }
    return labels[status] || status
  }

  function applyFilters() {
    // PERBAIKAN: Kirim seluruh state localFilters, biarkan Adapter yang mengabaikan nilai null
    studentListStore.updateFilters(localFilters.value)
    studentListStore.fetchStudents({ page: 1 })
  }

  let searchTimeout
  function debouncedApplyFilters() {
    clearTimeout(searchTimeout)
    searchTimeout = setTimeout(() => {
      applyFilters()
    }, 500)
  }

  function resetFilters() {
    // Set semua ke default, lalu apply
    localFilters.value = { status: null, search: '', tingkat: null, jurusan: null, kelas: null }
    applyFilters()
  }

  function onRequest(props) {
    const { page, rowsPerPage } = props.pagination
    studentListStore.fetchStudents({ page, limit: rowsPerPage })
  }

  function viewDetail(studentId) {
    router.push({ name: 'student-detail', params: { id: studentId } })
  }

  function openDeleteDialog(student) {
    selectedStudent.value = student
    showDeleteDialog.value = true
  }

  async function confirmDelete() {
    const result = await studentListStore.deleteStudent(selectedStudent.value.id)
    if (result) {
      studentDetailStore.invalidate(selectedStudent.value.id) 
      showDeleteDialog.value = false
      notificationMessage.value = 'Siswa berhasil dihapus'
      notificationType.value = 'success'
      showNotification.value = true
      studentListStore.fetchStudents()
    } else {
      notificationMessage.value = studentListStore.error?.message || 'Gagal menghapus siswa'
      notificationType.value = 'error'
      showNotification.value = true
    }
  }

  function loadData() {
    studentListStore.fetchStudents()
  }

  function applySearch(searchText) {
    localFilters.value.search = searchText
    applyFilters()
  }

  function resetSearch() {
    localFilters.value.search = ''
    applyFilters()
  }

  return {
    students: computed(() => studentListStore.students),
    isLoading: computed(() => studentListStore.isLoading),
    localFilters, 
    showDeleteDialog,
    showNotification,
    notificationMessage,
    notificationType,
    selectedStudent,
    statusOptions,
    tablePagination,
    getStatusColor,
    getStatusLabel,
    applyFilters,
    debouncedApplyFilters,
    resetFilters,
    onRequest,
    viewDetail,
    openDeleteDialog,
    confirmDelete,
    loadData,
    applySearch,
    resetSearch
  }
}
