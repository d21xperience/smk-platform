// apps/siakad-tu/src/composables/student/useStudentForm.js

import { ref, computed } from 'vue'
import { studentApi } from '@/adapters/api/studentApi.js'
import { useContextStore } from '@/stores/contextStore.js'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'

export function useStudentForm(studentId = null) {
  const contextStore = useContextStore()
  const studentDetailStore = useStudentDetailStore()

  const loading = ref(false)
  const isLocked = ref(false)
  const isEditMode = computed(() => !!studentId)

  // Flat structure (Bahasa Indonesia)
  const form = ref({
    nisn: '',
    nis: '',
    nama: '',
    tempatLahir: '',
    tanggalLahir: '',
    jenisKelamin: 'L',
    agama: 'Islam',
    noHp: '',
    email: '',
    asalSekolah: '',
    jurusan: '',
    kelas: '',
    alamat: '',
  })

  async function loadStudent() {
    if (!studentId) return
    loading.value = true

    try {
      const context = contextStore.currentContext
      const res = await studentApi.getStudentById(studentId, context)

      if (res.success && res.data) {
        const d = res.data
        form.value = {
          nisn: d.nisn || '',
          nis: d.nis || '',
          nama: d.fullName || '',
          tempatLahir: d.birthPlace || '',
          tanggalLahir: d.birthDate ? d.birthDate.split('T')[0] : '',
          jenisKelamin: d.gender || 'L',
          agama: 'Islam', // Default
          noHp: d.phone || '',
          email: '',
          asalSekolah: '',
          jurusan: d.jurusan || '',
          kelas: d.className || '',
          alamat: d.address || '',
        }
      }
    } catch (err) {
      console.error('Error loading student:', err)
    } finally {
      loading.value = false
    }
  }

  async function saveStudent() {
    loading.value = true
    try {
      const context = contextStore.currentContext

      const payload = {
        studentId: studentId, // undefined jika create
        nis: form.value.nis,
        nisn: form.value.nisn,
        nama: form.value.nama,
        tempatLahir: form.value.tempatLahir,
        tanggalLahir: form.value.tanggalLahir,
        jenisKelamin: form.value.jenisKelamin,
        alamat: form.value.alamat,
        noHp: form.value.noHp,
        kelas: form.value.kelas,
        jurusan: form.value.jurusan,
      }

      let res
      if (isEditMode.value) {
        res = await studentApi.updateStudentProfile(payload)
        if (res.success) studentDetailStore.invalidate(studentId)
      } else {
        res = await studentApi.registerStudent(payload, context)
      }

      if (res.success) {
        return true
      }
      throw new Error(res.error?.message || 'Gagal menyimpan data')
    } catch (err) {
      console.error('Error saving student:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  return { form, loading, isLocked, isEditMode, loadStudent, saveStudent }
}
