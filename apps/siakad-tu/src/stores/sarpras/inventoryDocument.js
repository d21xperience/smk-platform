import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useInventoryDocumentStore = defineStore('sarpras-inventory-document', () => {
  const documents = ref([])
  const isLoading = ref(false)

  const fetchDocuments = async (entityType, entityId) => {
    isLoading.value = true
    try {
      const res = await api.get(
        `/sarpras/documents?entity_type=${entityType}&entity_id=${entityId}`,
      )
      documents.value = res.data.data || []
    } catch (e) {
      console.error('Gagal mengambil dokumen:', e)
    } finally {
      isLoading.value = false
    }
  }

  const uploadDocument = async (formData) => {
    const res = await api.post('/sarpras/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    return res.data
  }

  const verifyDocument = async (id) => {
    const res = await api.post(`/sarpras/documents/${id}/verify`)
    return res.data
  }

  const rejectDocument = async (id, reason) => {
    const res = await api.post(`/sarpras/documents/${id}/reject`, { reason })
    return res.data
  }

  const getDownloadUrl = async (id) => {
    const res = await api.get(`/sarpras/documents/${id}/download`)
    return res.data.url
  }

  return {
    documents,
    isLoading,
    fetchDocuments,
    uploadDocument,
    verifyDocument,
    rejectDocument,
    getDownloadUrl,
  }
})
