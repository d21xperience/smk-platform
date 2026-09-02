// apps/siakad-tu/src/composables/student/useStudentActions.js

import { ref } from 'vue'
import { useStudentStore } from '../../stores/student/studentStore.js'
import { useOperationalContext } from '../context/useOperationalContext.js'

/**
 * useStudentActions — Composable untuk actions siswa (graduate, transfer, dll).
 *
 * Menyediakan:
 * - Graduate action dengan confirmation
 * - Transfer action dengan confirmation & form
 * - Loading & error state per action
 */
export function useStudentActions() {
  const studentStore = useStudentStore()
  const { current: context } = useOperationalContext()

  const loading = ref(false)
  const error = ref(null)

  /**
   * Graduate student dengan confirmation
   * @param {string} studentId
   * @param {string} studentName - untuk confirmation message
   * @returns {Promise<{ success: boolean, data?: Object }>}
   */
  async function graduateStudent(studentId, studentName) {
    const confirmed = window.confirm(
      `Luluskan siswa ${studentName}?\n\n` +
        `Status siswa akan berubah menjadi "Lulus" dan semua enrollment aktif akan diselesaikan.\n` +
        `Tindakan ini tidak dapat dibatalkan.`,
    )

    if (!confirmed) {
      return { success: false, cancelled: true }
    }

    loading.value = true
    error.value = null

    try {
      const result = await studentStore.graduateStudent(studentId, context.value)
      if (!result.success) {
        error.value = result.error
      }
      return result
    } finally {
      loading.value = false
    }
  }

  /**
   * Transfer student dengan confirmation
   * @param {Object} data - { studentId, studentName, targetSchool, reason }
   * @returns {Promise<{ success: boolean, data?: Object }>}
   */
  async function transferStudent(data) {
    const confirmed = window.confirm(
      `Pindahkan siswa ${data.studentName} ke ${data.targetSchool}?\n\n` +
        `Status siswa akan berubah menjadi "Pindah" dan semua enrollment aktif akan diselesaikan.\n` +
        `Tindakan ini tidak dapat dibatalkan.`,
    )

    if (!confirmed) {
      return { success: false, cancelled: true }
    }

    loading.value = true
    error.value = null

    try {
      const result = await studentStore.transferStudent(data, context.value)
      if (!result.success) {
        error.value = result.error
      }
      return result
    } finally {
      loading.value = false
    }
  }

  /**
   * Enroll student ke kelas
   * @param {Object} data - { studentId, classId }
   * @returns {Promise<{ success: boolean, data?: Object }>}
   */
  async function enrollStudent(data) {
    loading.value = true
    error.value = null

    try {
      const result = await studentStore.enrollStudent(data, context.value)
      if (!result.success) {
        error.value = result.error
      }
      return result
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    loading,
    error,
    graduateStudent,
    transferStudent,
    enrollStudent,
    clearError,
  }
}
