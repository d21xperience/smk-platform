import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useGoodsReceiptStore = defineStore('sarpras-goods-receipt', () => {
  const receipts = ref([])
  const isLoading = ref(false)

  const fetchReceipts = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/receipts')
      receipts.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil data penerimaan:', e)
    } finally {
      isLoading.value = false
    }
  }

  const createReceipt = async (payload) => {
    const res = await api.post('/sarpras/receipts', payload)
    await fetchReceipts()
    return res.data
  }

  const verifyReceipt = async (id) => {
    const res = await api.post(`/sarpras/receipts/${id}/verify`)
    await fetchReceipts()
    return res.data
  }

  const rejectReceipt = async (id, reason) => {
    const res = await api.post(`/sarpras/receipts/${id}/reject`, { reason })
    await fetchReceipts()
    return res.data
  }

  const uploadDocument = async (formData) => {
    const res = await api.post('/sarpras/receipts/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return res.data.url
  }

  const downloadPDF = async (id) => {
    const res = await api.get(`/sarpras/receipts/${id}/pdf`, { responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `BAST_${id}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  }

  return {
    receipts,
    isLoading,
    fetchReceipts,
    createReceipt,
    verifyReceipt,
    rejectReceipt,
    uploadDocument,
    downloadPDF,
  }
})
