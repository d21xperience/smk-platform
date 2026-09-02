import { computed } from 'vue'
import { useInventoryStore } from '../stores/inventoryStore.js'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'
import { useRouter } from 'vue-router'

export function useInventory() {
  const inventoryStore = useInventoryStore()
  const contextStore = useContextStore()
  const authStore = useAuthStore()
  const router = useRouter()

  const currentContext = computed(() => contextStore.currentContext)
  const currentUser = computed(() => authStore.currentUser)
  const items = computed(() => inventoryStore.items)
  const selectedItem = computed(() => inventoryStore.selectedItem)
  const isLoading = computed(() => inventoryStore.loading)
  const error = computed(() => inventoryStore.error)

  const loadItems = async ({ classId }) => {
    await inventoryStore.loadItems({ classId })
  }

  const loadItemById = async ({ itemId }) => {
    await inventoryStore.loadItemById({ itemId })
  }

  const submitDamageReport = async ({
    itemId,
    itemName,
    condition,
    description,
    reportMethod,
    date,
    classId,
    className,
  }) => {
    if (!currentContext.value || !currentUser.value) return

    const report = await inventoryStore.submitDamageReport({
      itemId,
      itemName,
      condition,
      description,
      reportMethod,
      date,
      teacherId: currentUser.value.id,
      teacherName: currentUser.value.name,
      classId,
      className,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
    return report
  }

  const navigateToReport = (itemId) => {
    router.push({
      name: 'inventory-report',
      params: { itemId },
    })
  }

  return {
    items,
    selectedItem,
    isLoading,
    error,
    loadItems,
    loadItemById,
    submitDamageReport,
    navigateToReport,
  }
}
