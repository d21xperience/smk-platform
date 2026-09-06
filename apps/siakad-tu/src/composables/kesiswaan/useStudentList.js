// apps/siakad-tu/src/composables/kesiswaan/useStudentList.js

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStudentListStore } from '@/stores/student/studentListStore'

/**
 * Composable untuk halaman daftar siswa.
 *
 * Tanggung jawab:
 * - Mengelola state lokal UI (filter, dialog, notifikasi)
 * - Menyediakan helper functions (format, color mapping)
 * - Mengkoordinasi interaksi antara Page dan Store
 */
export function useStudentList() {
  const router = useRouter()
  const studentListStore = useStudentListStore()

  // === LOCAL UI STATE ===
  const localFilters = ref({
    status: null,
    search: '',
  })

  const showDeleteDialog = ref(false)
  const showNotification = ref(false)
  const notificationMessage = ref('')
  const notificationType = ref('success')
  const selectedStudent = ref(null)

  // === OPTIONS ===
  const statusOptions = [
    { label: 'Aktif', value: 'ACTIVE' },
    { label: 'Lulus', value: 'GRADUATED' },
    { label: 'Mutasi', value: 'TRANSFERRED' },
  ]

  // === COLUMNS ===
  // const columns = [
  //   { name: 'nisn', label: 'NISN', field: 'nisn', align: 'left', sortable: true },
  //   { name: 'nis', label: 'NIS', field: 'nis', align: 'left', sortable: true },
  //   { name: 'fullName', label: 'Nama Lengkap', field: row => `${row.fullName?.firstName || ''} ${row.fullName?.lastName || ''}`.trim(), align: 'left', sortable: true },
  //   { name: 'gender', label: 'L/P', field: row => row.gender === 'MALE' ? 'L' : 'P', align: 'center' },
  //   { name: 'status', label: 'Status', field: 'status', align: 'center' },
  //   { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
  // ]

  // === COMPUTED ===
  const tablePagination = computed(() => ({
    page: studentListStore.pagination.page,
    rowsPerPage: studentListStore.pagination.limit,
    rowsNumber: studentListStore.pagination.total,
  }))

  // === HELPER FUNCTIONS ===
  function getStatusColor(status) {
    const colors = {
      ACTIVE: 'positive',
      GRADUATED: 'blue',
      TRANSFERRED: 'orange',
    }
    return colors[status] || 'grey'
  }

  function getStatusLabel(status) {
    const labels = {
      ACTIVE: 'Aktif',
      GRADUATED: 'Lulus',
      TRANSFERRED: 'Mutasi',
    }
    return labels[status] || status
  }

  // === ACTIONS ===
  function applyFilters() {
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
    localFilters.value = { status: null, search: '' }
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
    const result = await studentListStore.deleteStudent(selectedStudent.value.studentId)

    if (result) {
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

  // === RETURN ===
  return {
    // Store access (readonly)
    students: computed(() => studentListStore.students),
    isLoading: computed(() => studentListStore.isLoading),

    // Local UI state
    localFilters,
    showDeleteDialog,
    showNotification,
    notificationMessage,
    notificationType,
    selectedStudent,

    // Options & columns
    statusOptions,
    // columns,
    tablePagination,

    // Helper functions
    getStatusColor,
    getStatusLabel,

    // Actions
    applyFilters,
    debouncedApplyFilters,
    resetFilters,
    onRequest,
    viewDetail,
    openDeleteDialog,
    confirmDelete,
    loadData,
  }
}
