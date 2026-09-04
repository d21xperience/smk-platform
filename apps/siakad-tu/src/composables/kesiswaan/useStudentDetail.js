// apps/siakad-tu/src/composables/kesiswaan/useStudentDetail.js

import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'

/**
 * Composable untuk halaman detail siswa.
 */
export function useStudentDetail() {
  const router = useRouter()
  const route = useRoute()
  const studentDetailStore = useStudentDetailStore()

  // === LOCAL UI STATE ===
  const showNotification = ref(false)
  const notificationMessage = ref('')
  const notificationType = ref('success')

  // === COMPUTED ===
  const studentId = computed(() => route.params.id)

  // === HELPER FUNCTIONS ===
  function formatGender(gender) {
    return gender === 'MALE' ? 'Laki-laki' : 'Perempuan'
  }

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
      GRADUATED: 'blue',
      TRANSFERRED: 'orange',
    }
    return colors[status] || 'grey'
  }

  // === ACTIONS ===
  function editStudent() {
    router.push({ name: 'student-edit', params: { id: studentId.value } })
  }

  function goBack() {
    router.push({ name: 'student-list' })
  }

  function loadData() {
    if (studentId.value) {
      studentDetailStore.fetchStudentById(studentId.value)
    }
  }

  // === RETURN ===
  return {
    // Store access
    student: computed(() => studentDetailStore.currentStudent),
    isLoading: computed(() => studentDetailStore.isLoading),
    error: computed(() => studentDetailStore.error),

    // Local UI state
    showNotification,
    notificationMessage,
    notificationType,

    // Computed
    studentId,

    // Helper functions
    formatGender,
    formatDate,
    getStatusColor,

    // Actions
    editStudent,
    goBack,
    loadData,
  }
}
