import { ref } from 'vue'
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore'

export function useStudentDetail(studentId) {
  const store = useStudentAffairsStore()
  const loading = ref(false)
  const error = ref(null)
  const student = ref(null)

  async function loadDetail() {
    loading.value = true
    try {
      student.value = await store.loadStudentById(studentId)
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  return { student, loading, error, loadDetail }
}
