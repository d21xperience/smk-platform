// FILE: src/composables/useHomeroomStudent.js
// STATUS: MODIFY (koreksi pola arsitektur)
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomStudentStore } from '../stores/homeroomStudentStore.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

export function useHomeroomStudent() {
  const store = useHomeroomStudentStore()
  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const students = computed(() => store.students)
  const classId = computed(() => store.classId)
  const className = computed(() => store.className)
  const loading = computed(() => store.loading)
  const error = computed(() => store.error)
  const activeStudents = computed(() => store.activeStudents)
  const totalActiveStudents = computed(() => store.totalActiveStudents)
  const totalStudents = computed(() => store.totalStudents)

  async function loadStudents() {
    if (!currentContext.value || !currentUser.value) return
    await store.loadStudents({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
    })
  }

  function clearData() {
    store.clearData()
  }

  return {
    students,
    classId,
    className,
    loading,
    error,
    activeStudents,
    totalActiveStudents,
    totalStudents,
    loadStudents,
    clearData,
  }
}
