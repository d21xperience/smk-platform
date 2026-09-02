import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useDispositionStore } from '../stores/DispositionStore.js'
import { useOperationalContext } from './useOperationalContext.js'
import { createDispositionCreatePayload, createFollowUpPayload } from '../contracts/dispositionContract.js'

export function useDisposition() {
  const $q = useQuasar()
  const store = useDispositionStore()
  const context = useOperationalContext()

  // UI States
  const dialogDisposisi = ref(false)
  const dialogTindakLanjut = ref(false)

  const selectedSurat = ref(null)
  const selectedTugas = ref(null)

  const formDisposisi = ref({
    tujuan: null,
    instruksi: ''
  })

  const formTindakLanjut = ref({
    laporan: '',
    bukti: null
  })

  // Computed
  const pendingIncomingLetters = computed(() => store.pendingIncomingLetters)
  const myTasks = computed(() => store.myTasks)
  const stafOptions = computed(() => store.stafOptions)
  const isLoading = computed(() => store.isLoading)

  // Actions
  async function loadKepsekInbox() {
    await store.loadPendingIncomingLetters(context)
  }

  async function loadStafTasks() {
    await store.loadMyTasks(context, context.userId)
  }

  async function loadStafOptions() {
    await store.loadStafOptions(context)
  }

  function openDialogDisposisi(surat) {
    selectedSurat.value = surat
    formDisposisi.value = { tujuan: null, instruksi: '' }
    dialogDisposisi.value = true
  }

  async function submitDisposisi() {
    try {
      const payload = createDispositionCreatePayload({
        incomingLetterId: selectedSurat.value.id,
        toUserId: formDisposisi.value.tujuan, // Di mock ini masih string nama, di real API akan jadi userId
        instruksi: formDisposisi.value.instruksi,
        deadline: null
      })

      await store.createDisposition(context, payload)
      $q.notify({ type: 'positive', message: 'Disposisi berhasil dikirim' })
      dialogDisposisi.value = false
    } catch (e) {
      $q.notify({ type: 'negative', message: e.message || 'Gagal mengirim disposisi' })
    }
  }

  function openDialogTindakLanjut(tugas) {
    selectedTugas.value = tugas
    formTindakLanjut.value = { laporan: '', bukti: null }
    dialogTindakLanjut.value = true
  }

  async function submitTindakLanjut() {
    try {
      const payload = createFollowUpPayload({
        dispositionId: selectedTugas.value.id,
        laporan: formTindakLanjut.value.laporan,
        buktiUrl: formTindakLanjut.value.bukti || 'mock-bucket://bukti.pdf'
      })

      await store.submitFollowUp(context, selectedTugas.value.id, payload)
      $q.notify({ type: 'positive', message: 'Tindak lanjut berhasil disimpan' })
      dialogTindakLanjut.value = false
    } catch (e) {
      $q.notify({ type: 'negative', message: e.message || 'Gagal menyimpan tindak lanjut' })
    }
  }

  // Mock upload for bukti (Nanti diganti real upload service)
  function uploadBukti(files) {
    return new Promise((resolve) => {
      resolve({ files: [files[0]], url: 'https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf' })
    })
  }

  return {
    dialogDisposisi,
    dialogTindakLanjut,
    formDisposisi,
    formTindakLanjut,
    pendingIncomingLetters,
    myTasks,
    stafOptions,
    isLoading,
    loadKepsekInbox,
    loadStafTasks,
    loadStafOptions,
    openDialogDisposisi,
    submitDisposisi,
    openDialogTindakLanjut,
    submitTindakLanjut,
    uploadBukti
  }
}
