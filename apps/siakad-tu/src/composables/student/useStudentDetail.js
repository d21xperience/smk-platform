// apps/siakad-tu/src/composables/student/useStudentDetail.js

import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'

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
  }

  // === RETURN ===
  return {
    student,
    isLoading,
    error,
    showNotification,
    notificationMessage,
    notificationType,
    studentId,
    formatDate,
    getStatusColor,
    editStudent,
    goBack,
    loadData,
  }
}
