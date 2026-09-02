import { ref } from 'vue'
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore'

export function useStudentForm(studentId = null) {
  const store = useStudentAffairsStore()
  const loading = ref(false)
  const error = ref(null)
  const isLocked = ref(false)

  const form = ref({
    nama: '',
    nisn: '',
    tempatLahir: '',
    tanggalLahir: '',
    jenisKelamin: '',
    agama: '',
    alamat: '',
    noHp: '',
    email: '',
    asalSekolah: '',
    jurusan: '',
    kelas: '',
    status: 'AKTIF',
  })

  async function loadStudent() {
    if (!studentId) return
    loading.value = true
    try {
      const student = await store.loadStudentById(studentId)
      form.value = { ...student }
      isLocked.value = student.isLocked
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  async function saveStudent() {
    loading.value = true
    error.value = null
    try {
      if (studentId) {
        await store.updateStudentAction(studentId, form.value)
      } else {
        await store.createStudentAction(form.value)
      }
      return true
    } catch (err) {
      error.value = err
      return false
    } finally {
      loading.value = false
    }
  }

  return { form, loading, error, isLocked, loadStudent, saveStudent }
}
