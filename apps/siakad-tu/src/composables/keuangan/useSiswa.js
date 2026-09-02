// composables/useSiswa.js
import { ref } from 'vue'
import { siswaService } from '@/services/siswaService'

export function useSiswa() {
  const items = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchAll = async () => {
    loading.value = true
    try {
      items.value = await siswaService.getAll()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const create = async (data) => {
    loading.value = true
    try {
      const newItem = await siswaService.create(data)
      items.value.push(newItem)
      return newItem
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const update = async (id, data) => {
    loading.value = true
    try {
      const updated = await siswaService.update(id, data)
      const index = items.value.findIndex((i) => i.id === id)
      if (index !== -1) items.value[index] = updated
      return updated
    } finally {
      loading.value = false
    }
  }

  const remove = async (id) => {
    loading.value = true
    try {
      await siswaService.delete(id)
      items.value = items.value.filter((i) => i.id !== id)
    } finally {
      loading.value = false
    }
  }

  return { items, loading, error, fetchAll, create, update, remove }
}
