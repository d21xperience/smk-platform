import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useLogisticStore = defineStore('sarpras-logistic', () => {
  const items = ref([])
  const isLoading = ref(false)

  const fetchItems = async () => {
    isLoading.value = true
    try {
      // KOREKSI: Hapus prefix /api/v1 karena baseURL sudah menanganinya
      const response = await api.get('/logistics')
      items.value = response.data.data || []
    } catch (error) {
      console.error('Gagal mengambil data logistik:', error)
    } finally {
      isLoading.value = false
    }
  }

  const createItem = async (payload) => {
    const response = await api.post('/logistics', payload)
    await fetchItems()
    return response.data
  }

  const restockItem = async (id, payload) => {
    const response = await api.post(`/logistics/${id}/restock`, payload)
    await fetchItems()
    return response.data
  }

  const useItem = async (id, payload) => {
    const response = await api.post(`/logistics/${id}/use`, payload)
    await fetchItems()
    return response.data
  }

  return {
    items,
    isLoading,
    fetchItems,
    createItem,
    restockItem,
    useItem,
  }
})
