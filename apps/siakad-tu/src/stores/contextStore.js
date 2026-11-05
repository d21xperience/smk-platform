import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ContextService } from '@sdp/context'
import { Notify } from 'quasar'
import { useAuthStore } from './authStore.js'

export const useContextStore = defineStore('context', () => {
  const authStore = useAuthStore()

  // === STATE ===
  const availablePeriods = ref([])
  
  // FIX: Inisialisasi sinkron dari AuthStore / LocalStorage agar tidak kosong saat di-refresh
  let initialSemesterId = ''
  if (authStore.isAuthenticated && authStore.user) {
    const cacheKey = `ctx:${authStore.user.id}:${authStore.user.school_id}`
    const cachedPref = localStorage.getItem(cacheKey)
    if (cachedPref) {
      initialSemesterId = cachedPref
    } else if (authStore.user.academic_period_id) {
      initialSemesterId = String(authStore.user.academic_period_id)
    }
  }
  
  const selectedSemesterId = ref(initialSemesterId)
  const isLoading = ref(false)
  const error = ref(null)

  // === COMPUTED: identity dari authStore ===
  const userId = computed(() => authStore.user?.id || '')
  const schoolId = computed(() => authStore.user?.school_id || '')

  const hasContext = computed(() => !!userId.value && !!schoolId.value)

  // === COMPUTED: context object ===
  const currentContext = computed(() => ({
    userId: userId.value,
    schoolId: schoolId.value,
    semesterId: selectedSemesterId.value,
    toHeaders: () => ({
      'x-school-id': schoolId.value,
      'x-user-id': userId.value,
      'x-academic-period-id': selectedSemesterId.value,
    }),
  }))

  const current = currentContext

  // === COMPUTED: derivasi academicYear & semester dari ID ===
  const academicYear = computed(() => {
    const sid = String(selectedSemesterId.value || '')
    if (sid.length === 5) {
      const year = parseInt(sid.substring(0, 4), 10)
      if (!isNaN(year)) return `${year}/${year + 1}`
    }
    return authStore.user?.academic_year || ''
  })

  const semester = computed(() => {
    const sid = String(selectedSemesterId.value || '')
    if (sid.length === 5) {
      const s = parseInt(sid.substring(4, 5), 10)
      if (s === 1 || s === 2) return s
    }
    return authStore.user?.semester || 1
  })

  // === COMPUTED: label yang lebih bersih ===
  const selectedYearLabel = computed(() => {
    const selected = availablePeriods.value.find((p) => String(p.semester_id) === String(selectedSemesterId.value))
    
    if (!selected) return '-'
    
    const yearText = `${selected.tahun_ajaran_id}/${parseInt(selected.tahun_ajaran_id) + 1}`
    let semText = selected.nama_semester || ''
    if (semText.includes('/')) {
      semText = semText.split(' ')[0]
    }
    
    return `${yearText} - ${semText}`
  })

  // === ACTIONS ===
  async function initialize() {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    try {
      let fetchCtx = currentContext.value
      if (!authStore.isAuthenticated || !authStore.user) {
        fetchCtx = {
          userId: 'public',
          schoolId: '11111111-1111-1111-1111-111111111111',
          semesterId: ''
        }
      }

      availablePeriods.value = await ContextService.fetchPeriods(fetchCtx)

      if (!availablePeriods.value || availablePeriods.value.length === 0) {
        console.warn("[contextStore] Data periode kosong dari server.")
        selectedSemesterId.value = ''
        return
      }

      // Jika selectedSemesterId sudah ada di array, jangan timpa
      const currentIsValid = availablePeriods.value.some(p => String(p.semester_id) === String(selectedSemesterId.value))
      
      if (!currentIsValid) {
        if (authStore.isAuthenticated) {
          const cacheKey = `ctx:${fetchCtx.userId}:${fetchCtx.schoolId}`
          const cachedPref = localStorage.getItem(cacheKey)
          const validCached = cachedPref && availablePeriods.value.some((p) => String(p.semester_id) === String(cachedPref))

          if (validCached) {
            selectedSemesterId.value = cachedPref
          } else if (
            authStore.user?.academic_period_id &&
            availablePeriods.value.some((p) => String(p.semester_id) === String(authStore.user.academic_period_id))
          ) {
            selectedSemesterId.value = authStore.user.academic_period_id
            localStorage.setItem(cacheKey, authStore.user.academic_period_id)
          } else {
            const pref = await ContextService.fetchPreference(fetchCtx).catch(() => null)
            if (pref && pref.semester_id) {
              selectedSemesterId.value = pref.semester_id
            } else {
              const active = availablePeriods.value.find((p) => p.periode_aktif === true || p.periode_aktif === '1' || p.is_active === true)
              selectedSemesterId.value = active?.semester_id || availablePeriods.value[0]?.semester_id || ''
            }
          }
        } else {
          const active = availablePeriods.value.find((p) => p.periode_aktif === true || p.periode_aktif === '1' || p.is_active === true)
          selectedSemesterId.value = active?.semester_id || availablePeriods.value[0]?.semester_id || ''
        }
      }
    } catch (err) {
      error.value = err.message
      console.error('[contextStore] Initialize error:', err)
      Notify.create({
        type: 'negative',
        message: 'Gagal memuat data tahun ajaran.',
        icon: 'error',
      })
    } finally {
      isLoading.value = false
    }
  }

  async function setSemester(semesterId) {
    isLoading.value = true
    try {
      const ctx = currentContext.value
      const sidStr = String(semesterId)
      
      if (authStore.isAuthenticated) {
        await ContextService.updatePreference(ctx, sidStr)
        localStorage.setItem(`ctx:${ctx.userId}:${ctx.schoolId}`, sidStr)
      }
      
      selectedSemesterId.value = sidStr

      const selected = availablePeriods.value.find((p) => String(p.semester_id) === sidStr)
      const label = selected ? `${selected.tahun_ajaran_id} - ${selected.nama_semester}` : sidStr

      Notify.create({
        type: 'positive',
        message: `Tahun ajaran berhasil diubah ke ${label}`,
        icon: 'check_circle',
        position: 'top-right',
      })
    } catch (err) {
      error.value = err.message
      Notify.create({
        type: 'negative',
        message: 'Gagal menyimpan tahun ajaran: ' + err.message,
        icon: 'error',
        position: 'top-right',
      })
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    availablePeriods,
    selectedSemesterId,
    isLoading,
    error,
    userId,
    schoolId,
    hasContext,
    current,
    currentContext,
    academicYear,
    semester,
    selectedYearLabel,
    initialize,
    setSemester,
  }
})
