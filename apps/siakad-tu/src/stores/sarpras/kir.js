import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useKirStore = defineStore('sarpras-kir', () => {
  const rooms = ref([])
  const categories = ref([])
  const catalogs = ref([])
  const items = ref([])
  const isLoading = ref(false)

  // KIR Items
  const fetchRooms = async () => {
    const res = await api.get('/sarpras/kir/rooms')
    rooms.value = res.data.data || []
  }
  const fetchItemsByRoom = async (roomId) => {
    isLoading.value = true
    try {
      const res = await api.get(`/sarpras/kir/items?room_id=${roomId}`)
      items.value = res.data.data || []
    } finally {
      isLoading.value = false
    }
  }
  const fetchAllItems = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/kir/items')
      items.value = res.data.data || []
    } finally {
      isLoading.value = false
    }
  }
  const createItem = async (payload) => await api.post('/sarpras/kir/items', payload)
  const updateItem = async (id, payload) => await api.put(`/sarpras/kir/items/${id}`, payload)
  const deleteItem = async (id) => await api.delete(`/sarpras/kir/items/${id}`)

  const printQR = async (itemIds) => {
    const response = await api.post(
      '/sarpras/kir/print-qr',
      { item_ids: itemIds },
      { responseType: 'blob' },
    )
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'Stiker_QR_Aset.pdf')
    document.body.appendChild(link)
    link.click()
  }

  // Master Data
  const fetchMasterCategories = async () => {
    const res = await api.get('/sarpras/master/categories')
    categories.value = res.data.data || []
  }
  const fetchMasterCatalogs = async () => {
    const res = await api.get('/sarpras/master/catalogs')
    catalogs.value = res.data.data || []
  }
  const fetchMasterData = async () => {
    await Promise.all([fetchMasterCategories(), fetchMasterCatalogs()])
  }

  const createCategory = async (payload) => await api.post('/sarpras/master/categories', payload)
  const updateCategory = async (id, payload) =>
    await api.put(`/sarpras/master/categories/${id}`, payload)
  const deleteCategory = async (id) => await api.delete(`/sarpras/master/categories/${id}`)

  const createCatalog = async (payload) => await api.post('/sarpras/master/catalogs', payload)
  const updateCatalog = async (id, payload) =>
    await api.put(`/sarpras/master/catalogs/${id}`, payload)
  const deleteCatalog = async (id) => await api.delete(`/sarpras/master/catalogs/${id}`)

  // Upload Image
  const uploadCatalogImage = async (formData) => {
    const res = await api.post('/sarpras/master/catalogs/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return res.data.url
  }

  return {
    rooms,
    categories,
    catalogs,
    items,
    isLoading,
    fetchRooms,
    fetchItemsByRoom,
    fetchAllItems,
    createItem,
    updateItem,
    deleteItem,
    printQR,
    fetchMasterData,
    createCategory,
    updateCategory,
    deleteCategory,
    createCatalog,
    updateCatalog,
    deleteCatalog,
    uploadCatalogImage,
  }
})
