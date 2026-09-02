import { computed } from 'vue'
import { useHomeroomBillingStore } from '../stores/homeroomBillingStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'

export function useHomeroomBilling() {
  const homeroomBillingStore = useHomeroomBillingStore()
  const contextStore = useContextStore()
  const authStore = useAuthStore()

  const currentContext = computed(() => contextStore.currentContext)
  const currentUser = computed(() => authStore.currentUser)
  const assignment = computed(() => homeroomBillingStore.assignment)
  const bills = computed(() => homeroomBillingStore.bills)
  const summary = computed(() => homeroomBillingStore.summary)
  const isLoading = computed(() => homeroomBillingStore.loading)
  const error = computed(() => homeroomBillingStore.error)
  const hasAssignment = computed(() => homeroomBillingStore.hasAssignment)
  const hasData = computed(() => homeroomBillingStore.hasData)

  const loadBillingData = async () => {
    if (!currentContext.value || !currentUser.value) return

    const assignmentData = await homeroomBillingStore.loadHomeroomAssignment({
      teacherId: currentUser.value.id
    })

    if (!assignmentData) {
      return
    }

    await homeroomBillingStore.loadBillingData({
      classId: assignmentData.classId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      teacherId: currentUser.value.id
    })
  }

  const clearData = () => {
    homeroomBillingStore.clearData()
  }

  return {
    assignment,
    bills,
    summary,
    isLoading,
    error,
    hasAssignment,
    hasData,
    loadBillingData,
    clearData
  }
}
