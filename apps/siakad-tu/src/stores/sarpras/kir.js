import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useKirStore = defineStore('sarpras-kir', () => {
  const rooms = ref([])
  const categories = ref([])
  const catalogs = ref([])
  const items = ref([])
  const history = ref([])
  const isLoading = ref(false)

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

  const fetchMasterData = async () => {
    const resCat = await api.get('/sarpras/master/categories')
    categories.value = resCat.data.data || []
    const resCata = await api.get('/sarpras/master/catalogs')
    catalogs.value = resCata.data.data || []
  }

  const createItem = async (payload) => await api.post('/sarpras/kir/items', payload)
  const updateItem = async (id, payload) => await api.put(`/sarpras/kir/items/${id}`, payload)
  const deleteItem = async (id) => await api.delete(`/sarpras/kir/items/${id}`)

  const mutateItem = async (id, payload) => {
    const res = await api.post(`/sarpras/inventory/${id}/mutate`, payload)
    return res.data
  }

  const fetchHistory = async (id) => {
    const res = await api.get(`/sarpras/inventory/${id}/history`)
    history.value = res.data.data || []
  }

  const downloadReport = async () => {
    const res = await api.get('/sarpras/inventory/report', { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', 'Buku_Induk_Inventaris.pdf')
    document.body.appendChild(link)
    link.click()
    link.remove()
  }

  const downloadKIRReport = async (roomId) => {
    const res = await api.get(`/sarpras/kir/report/${roomId}`, { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `KIR_${roomId}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  }

  return {
    rooms,
    categories,
    catalogs,
    items,
    history,
    isLoading,
    fetchRooms,
    fetchItemsByRoom,
    fetchAllItems,
    fetchMasterData,
    createItem,
    updateItem,
    deleteItem,
    mutateItem,
    fetchHistory,
    downloadReport,
    downloadKIRReport,
  }
})
