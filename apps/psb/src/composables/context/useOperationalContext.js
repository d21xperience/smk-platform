// src/composables/context/useOperationalContext.js

import { ref } from 'vue'
// import { useAuthStore } from 'src/stores/auth'
// import { useSchoolStore } from 'src/stores/school'

export function useOperationalContext() {
  const context = ref({
    user: null,
    school: null,
    academicYear: null,
    semester: null,
  })
  const isContextReady = ref(false)

  async function loadContext() {
    // Di sini kita ambil dari store atau API
    // Contoh: ambil dari Pinia
    // const authStore = useAuthStore()
    // const schoolStore = useSchoolStore()
    // context.value.user = authStore.user
    // context.value.school = schoolStore.currentSchool
    // context.value.academicYear = schoolStore.currentAcademicYear
    // context.value.semester = schoolStore.currentSemester

    // Simulasi async
    await new Promise((resolve) => setTimeout(resolve, 300))
    context.value = {
      user: { id: 1, name: 'Guru A', role: 'teacher' },
      school: { id: 1, name: 'SMA Negeri 1' },
      academicYear: '2025/2026',
      semester: 'Ganjil',
    }
    isContextReady.value = true
  }

  return {
    context,
    isContextReady: () => isContextReady.value,
    loadContext,
  }
}
