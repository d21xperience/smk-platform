import { computed } from 'vue'
import { useProgressStore } from '../stores/progressStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'
import { useRouter } from 'vue-router'

export function useProgress() {
  const progressStore = useProgressStore()
  const contextStore = useContextStore()
  const authStore = useAuthStore()
  const router = useRouter()

  const currentContext = computed(() => contextStore.currentContext)
  const currentUser = computed(() => authStore.currentUser)
  const students = computed(() => progressStore.students)
  const progressRecords = computed(() => progressStore.progressRecords)
  const isLoading = computed(() => progressStore.loading)
  const error = computed(() => progressStore.error)

  const loadStudents = async ({ classId }) => {
    await progressStore.loadStudents({ classId })
  }

  const loadProgress = async ({ studentId }) => {
    if (!currentContext.value) return
    await progressStore.loadProgress({
      studentId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  const createProgress = async ({
    studentId,
    studentName,
    classId,
    className,
    type,
    title,
    description,
    date,
  }) => {
    if (!currentContext.value || !currentUser.value) return

    const progress = await progressStore.createProgress({
      studentId,
      studentName,
      classId,
      className,
      type,
      title,
      description,
      date,
      teacherId: currentUser.value.id,
      teacherName: currentUser.value.name,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
    return progress
  }

  const navigateToStudent = (student) => {
    router.push({
      name: 'progress-student',
      params: { studentId: student.studentId },
      query: {
        studentName: student.studentName,
        classId: student.classId,
        className: student.className,
      },
    })
  }

  return {
    students,
    progressRecords,
    isLoading,
    error,
    loadStudents,
    loadProgress,
    createProgress,
    navigateToStudent,
  }
}
