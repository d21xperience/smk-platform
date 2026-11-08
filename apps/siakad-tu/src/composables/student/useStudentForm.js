// apps/siakad-tu/src/composables/student/useStudentForm.js

import { ref, computed, onMounted } from 'vue'
import { studentApi } from '@/adapters/api/studentApi.js'
import { useContextStore } from '@/stores/contextStore.js'
import { useStudentDetailStore } from '@/stores/student/studentDetailStore'
import { api } from '@/boot/axios.js'

export function useStudentForm(studentId = null) {
  const contextStore = useContextStore()
  const studentDetailStore = useStudentDetailStore()

  const loading = ref(false)
  const isLocked = ref(false)
  const isEditMode = computed(() => !!studentId)
  const loadingClasses = ref(false)
  
  const classOptions = ref([])

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
    kelas: '', // Menyimpan rombongan_belajar_id (UUID)
    alamat: '',
    parentName: '',
    parentPhone: '',
  })

  async function fetchClasses() {
    loadingClasses.value = true
    try {
      const response = await api.get('/classes', { params: { limit: 100 } })
      const responseData = response.data
      
      // PERBAIKAN (Hard Rule #5): Deteksi semua kemungkinan bentuk wrapper
      let items = []
      if (Array.isArray(responseData)) {
        // Format: [ ... ]
        items = responseData
      } else if (responseData && Array.isArray(responseData.items)) {
        // Format: { items: [ ... ] }
        items = responseData.items
      } else if (responseData && Array.isArray(responseData.data)) {
        // Format: { data: [ ... ] }
        items = responseData.data
      } else if (responseData && responseData.data && Array.isArray(responseData.data.data)) {
        // Format: { data: { data: [ ... ] } } (double wrapped, seperti yang terjadi saat ini)
        items = responseData.data.data
      } else {
        console.warn('[useStudentForm] Format respons kelas tidak dikenali:', responseData)
      }

      classOptions.value = items.map(c => ({
        label: c.name || c.nm_kelas || 'Tanpa Nama',
        value: c.id || c.rombongan_belajar_id,
        major: c.major || c.nama_jurusan_sp || '' // Mengambil dari nama_jurusan_sp
      }))
    } catch (err) {
      console.error('Gagal mengambil data kelas:', err)
    } finally {
      loadingClasses.value = false
    }
  }

  function onClassChange(selectedId) {
    const selected = classOptions.value.find(c => c.value === selectedId)
    if (selected && selected.major) {
      form.value.jurusan = selected.major
    } else {
      form.value.jurusan = '' // Reset jika tidak ada jurusan
    }
  }

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
          agama: 'Islam',
          noHp: d.phone || '',
          email: '',
          asalSekolah: '',
          jurusan: '',
          kelas: d.classId || d.className || '', // Gunakan classId (UUID) jika ada
          alamat: d.address || '',
          parentName: d.parentName || '',
          parentPhone: d.parentPhone || '',
        }
        
        // Isi jurusan otomatis jika kelas ditemukan di dropdown
        const matchedClass = classOptions.value.find(c => c.value === form.value.kelas || c.label === d.className)
        if (matchedClass) {
          form.value.kelas = matchedClass.value
          form.value.jurusan = matchedClass.major
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
        studentId: studentId,
        name: form.value.nama,
        nis: form.value.nis,
        nisn: form.value.nisn,
        birthPlace: form.value.tempatLahir,
        birthDate: form.value.tanggalLahir,
        gender: form.value.jenisKelamin === 'L' ? 'MALE' : 'FEMALE',
        address: form.value.alamat,
        phone: form.value.noHp,
        classId: form.value.kelas, // Mengirim UUID
        parentName: form.value.parentName,
        parentPhone: form.value.parentPhone,
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

  onMounted(async () => {
    await fetchClasses()
  })

  return { 
    form, 
    loading, 
    isLocked, 
    isEditMode, 
    classOptions, 
    loadingClasses,
    onClassChange,
    loadStudent, 
    saveStudent 
  }
}
