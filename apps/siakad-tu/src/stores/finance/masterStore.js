import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { masterAdapter } from '@/adapters/api/keuangan/masterAdapter'

export const useMasterStore = defineStore('keuangan-master', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  function makeCrud(adapter) {
    const items = ref([])
    const detail = ref(null)
    const loading = ref(false)
    const error = ref(null)
    const pagination = ref({ total: 0, limit: 100, offset: 0 })

    async function run(fn) {
      loading.value = true
      error.value = null
      try { return await fn() }
      catch (e) {
        error.value = e.response?.data?.error || e.message
        throw e
      } finally {
        loading.value = false
      }
    }

    return {
      items, detail, loading, error, pagination,
      async list(params = {}) {
        return run(async () => {
          const res = await adapter.list(baseParams(params))
          items.value = res.data || []
          if (res.pagination) pagination.value = res.pagination
          return res
        })
      },
      async get(id) {
        return run(async () => {
          const res = await adapter.get(id, baseParams())
          detail.value = res.data
          return res
        })
      },
      async create(payload) {
        return run(async () => {
          const res = await adapter.create(payload, baseParams())
          items.value = [res.data, ...items.value]
          return res
        })
      },
      async update(id, payload) {
        return run(async () => {
          const res = await adapter.update(id, payload, baseParams())
          const idx = items.value.findIndex(x => x.id === id)
          if (idx >= 0) items.value[idx] = res.data
          detail.value = res.data
          return res
        })
      },
      async remove(id) {
        return run(async () => {
          await adapter.remove(id, baseParams())
          items.value = items.value.filter(x => x.id !== id)
        })
      },
      reset() {
        items.value = []
        detail.value = null
        pagination.value = { total: 0, limit: 100, offset: 0 }
        error.value = null
      },
    }
  }

  return {
    schoolID, periodID,
    jenisIuran: makeCrud({
      list: (p) => masterAdapter.listJenisIuran(p),
      get: (id, p) => masterAdapter.getJenisIuran(id, p),
      create: (payload, p) => masterAdapter.createJenisIuran(payload, p),
      update: (id, payload, p) => masterAdapter.updateJenisIuran(id, payload, p),
      remove: (id, p) => masterAdapter.deleteJenisIuran(id, p),
    }),
    tarifIuran: makeCrud({
      list: (p) => masterAdapter.listTarifIuran(p),
      get: (id, p) => masterAdapter.getTarifIuran(id, p),
      create: (payload, p) => masterAdapter.createTarifIuran(payload, p),
      update: (id, payload, p) => masterAdapter.updateTarifIuran(id, payload, p),
      remove: (id, p) => masterAdapter.deleteTarifIuran(id, p),
    }),
    kodeRekening: makeCrud({
      list: (p) => masterAdapter.listKodeRekening(p),
      get: (id, p) => masterAdapter.getKodeRekening(id, p),
      create: (payload, p) => masterAdapter.createKodeRekening(payload, p),
      update: (id, payload, p) => masterAdapter.updateKodeRekening(id, payload, p),
      remove: (id, p) => masterAdapter.deleteKodeRekening(id, p),
    }),
    rekeningBank: makeCrud({
      list: (p) => masterAdapter.listRekeningBank(p),
      get: (id, p) => masterAdapter.getRekeningBank(id, p),
      create: (payload, p) => masterAdapter.createRekeningBank(payload, p),
      update: (id, payload, p) => masterAdapter.updateRekeningBank(id, payload, p),
      remove: (id, p) => masterAdapter.deleteRekeningBank(id, p),
    }),
  }
})
