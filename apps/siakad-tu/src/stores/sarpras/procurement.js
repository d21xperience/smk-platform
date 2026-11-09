import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useProcurementStore = defineStore('sarpras-procurement', () => {
  const requests = ref([])
  const ptkList = ref([])
  const catalogList = ref([])
  const isLoading = ref(false)

  // Fetch List Pengajuan
  const fetchRequests = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/procurement/requests')
      requests.value = res.data.data || []
    } catch (e) {
      console.error(e)
    } finally {
      isLoading.value = false
    }
  }

  // Fetch Data Dropdown (PTK & Katalog)
  const fetchDropdownData = async () => {
    try {
      // PERBAIKAN: Endpoint PTK
      const resPtk = await api.get('/teachers')
      ptkList.value = resPtk.data.data || []

      // PERBAIKAN: Endpoint Master Catalogs
      const resCatalog = await api.get('/sarpras/master/catalogs')
      catalogList.value = resCatalog.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data dropdown:', e)
    }
  }

  // CRUD
  const createRequest = async (payload) => {
    const res = await api.post('/procurement/requests', payload)
    await fetchRequests()
    return res.data
  }

  const updateRequest = async (id, payload) => {
    const res = await api.put(`/procurement/requests/${id}`, payload)
    await fetchRequests()
    return res.data
  }

  const deleteRequest = async (id) => {
    const res = await api.delete(`/procurement/requests/${id}`)
    await fetchRequests()
    return res.data
  }

  const updateStatus = async (id, status, reason = '', notes = '') => {
    const res = await api.put(`/procurement/requests/${id}/status`, {
      status,
      rejection_reason: reason,
      notes,
    })
    await fetchRequests()
    return res.data
  }

  return {
    requests,
    ptkList,
    catalogList,
    isLoading,
    fetchRequests,
    fetchDropdownData,
    createRequest,
    updateRequest,
    deleteRequest,
    updateStatus,
  }
})
