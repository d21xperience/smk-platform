import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useIncomingLetterStore } from '../stores/IncomingLetterStore.js'
import { useOperationalContext } from './useOperationalContext.js'
import { createIncomingCreatePayload, createIncomingListQuery } from '../contracts/incomingLetterContract.js'
import { ClassificationMockAdapter } from '../adapters/mock/ClassificationMockAdapter.js' // Import mock untuk opsi

export function useIncomingLetter() {
  const $q = useQuasar()
  const store = useIncomingLetterStore()
  const context = useOperationalContext()

  // UI States
  const isDialogOpen = ref(false)
  const isDetailDialogOpen = ref(false)
  const filter = ref('')

  const form = ref({
    asal: '',
    nomorSurat: '',
    tanggalSurat: '',
    perihal: '',
    kodeKlasifikasi: null,
    fileSurat: null
  })

  // Klasifikasi Options (untuk q-select)
  const klasifikasiOptions = ref([])

  // Computed
  const list = computed(() => store.list)
  const currentItem = computed(() => store.currentItem)
  const isLoading = computed(() => store.isLoading)

  // Actions
  async function loadList(query = {}) {
    const q = createIncomingListQuery({ ...query, search: filter.value })
    await store.loadList(context, q)
  }

  async function loadKlasifikasiOptions() {
    klasifikasiOptions.value = await ClassificationMockAdapter.getKlasifikasiOptions(context)
  }

  async function submitCreate() {
    try {
      const payload = createIncomingCreatePayload({
        asal: form.value.asal,
        nomorSuratAsal: form.value.nomorSurat,
        tanggalSurat: form.value.tanggalSurat,
        tanggalDiterima: new Date().toISOString().split('T')[0],
        perihal: form.value.perihal,
        classificationCode: form.value.kodeKlasifikasi?.split(' - ')[0] || form.value.kodeKlasifikasi,
        fileUrl: form.value.fileSurat || 'mock-url.pdf'
      })

      await store.createItem(context, payload)
      $q.notify({ type: 'positive', message: 'Surat masuk berhasil disimpan' })
      closeDialog()
      await loadList()
    } catch (e) {
      $q.notify({ type: 'negative', message: e.message || 'Gagal menyimpan surat' })
    }
  }

  async function viewDetail(id) {
    await store.loadDetail(context, id)
    isDetailDialogOpen.value = true
  }

  function openDialog() {
    resetForm()
    isDialogOpen.value = true
  }

  function closeDialog() {
    isDialogOpen.value = false
    resetForm()
  }

  function resetForm() {
    form.value = {
      asal: '',
      nomorSurat: '',
      tanggalSurat: '',
      perihal: '',
      kodeKlasifikasi: null,
      fileSurat: null
    }
  }

  function statusColor(status) {
    if (status === 'BARU') return 'blue'
    if (status === 'DIDISPOSISIKAN') return 'orange'
    return 'green'
  }

  return {
    list,
    currentItem,
    isLoading,
    filter,
    form,
    klasifikasiOptions,
    isDialogOpen,
    isDetailDialogOpen,
    loadList,
    loadKlasifikasiOptions,
    submitCreate,
    viewDetail,
    openDialog,
    closeDialog,
    statusColor
  }
}
