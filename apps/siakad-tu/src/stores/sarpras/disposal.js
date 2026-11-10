import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useDisposalStore = defineStore('sarpras-disposal', () => {
  const requests = ref([])
  const activeAssets = ref([])
  const isLoading = ref(false)

  const fetchRequests = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/disposal')
      requests.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data penghapusan:', e)
    } finally {
      isLoading.value = false
    }
  }

  // Mengambil aset yang AKTIF (untuk dipilih di form usulan)
  const fetchActiveAssets = async () => {
    try {
      const res = await api.get('/sarpras/kir/items')
      // Filter hanya aset yang berstatus AKTIF
      activeAssets.value = (res.data.data || []).filter((a) => a.status_aset === 'AKTIF')
    } catch (e) {
      console.error('Gagal mengambil data aset:', e)
    }
  }

  const createRequest = async (payload) => {
    const res = await api.post('/sarpras/disposal', payload)
    await fetchRequests()
    return res.data
  }

  const verifyRequest = async (id, payload) => {
    const res = await api.post(`/sarpras/disposal/${id}/verify`, payload)
    await fetchRequests()
    return res.data
  }

  const approveRequest = async (id, payload) => {
    const res = await api.post(`/sarpras/disposal/${id}/approve`, payload)
    await fetchRequests()
    return res.data
  }

  const executeRequest = async (id, payload) => {
    const res = await api.post(`/sarpras/disposal/${id}/execute`, payload)
    await fetchRequests()
    return res.data
  }

  return {
    requests,
    activeAssets,
    isLoading,
    fetchRequests,
    fetchActiveAssets,
    createRequest,
    verifyRequest,
    approveRequest,
    executeRequest,
  }
})
