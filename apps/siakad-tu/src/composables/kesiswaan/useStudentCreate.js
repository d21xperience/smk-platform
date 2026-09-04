// apps/siakad-tu/src/composables/kesiswaan/useStudentCreate.js

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStudentFormStore } from '@/stores/student/studentFormStore'

/**
 * Composable untuk halaman create/edit siswa.
 */
export function useStudentCreate() {
  const router = useRouter()
  const studentFormStore = useStudentFormStore()

  // === LOCAL UI STATE ===
  const showNotification = ref(false)
  const notificationMessage = ref('')
  const notificationType = ref('success')

  const formData = ref({
    nisn: '',
    nis: '',
    fullName: {
      firstName: '',
      middleName: '',
      lastName: '',
    },
    gender: null,
    birthDate: '',
    address: {
      street: '',
      rtRw: '',
      village: '',
      district: '',
      city: '',
      postalCode: '',
    },
    contactInfo: {
      phone: '',
      email: '',
    },
    guardianInfo: {
      name: '',
      relation: '',
      phone: '',
      occupation: '',
    },
  })

  // === OPTIONS ===
  const genderOptions = [
    { label: 'Laki-laki', value: 'MALE' },
    { label: 'Perempuan', value: 'FEMALE' },
  ]

  // === COMPUTED ===
  const isEditMode = computed(() => !!studentFormStore.currentStudent)

  // === ACTIONS ===
  function resetForm() {
    formData.value = {
      nisn: '',
      nis: '',
      fullName: { firstName: '', middleName: '', lastName: '' },
      gender: null,
      birthDate: '',
      address: { street: '', rtRw: '', village: '', district: '', city: '', postalCode: '' },
      contactInfo: { phone: '', email: '' },
      guardianInfo: { name: '', relation: '', phone: '', occupation: '' },
    }
  }

  async function onSubmit() {
    const result = await studentFormStore.submitStudent(formData.value)

    if (result) {
      notificationMessage.value = isEditMode.value
        ? 'Data siswa berhasil diupdate'
        : 'Siswa berhasil didaftarkan'
      notificationType.value = 'success'
      showNotification.value = true

      setTimeout(() => {
        router.push({ name: 'student-list' })
      }, 1500)
    } else {
      notificationMessage.value = studentFormStore.error?.message || 'Gagal menyimpan data siswa'
      notificationType.value = 'error'
      showNotification.value = true
    }
  }

  function onCancel() {
    router.push({ name: 'student-list' })
  }

  function loadData(studentId) {
    if (studentId) {
      studentFormStore.fetchStudentById(studentId)
    }
  }

  // === RETURN ===
  return {
    // Store access
    isLoading: computed(() => studentFormStore.isLoading),
    currentStudent: computed(() => studentFormStore.currentStudent),

    // Local UI state
    showNotification,
    notificationMessage,
    notificationType,
    formData,

    // Options
    genderOptions,

    // Computed
    isEditMode,

    // Actions
    resetForm,
    onSubmit,
    onCancel,
    loadData,
  }
}
