import { storeToRefs } from 'pinia'
import { useQuasar } from 'quasar'
import { useOutgoingLetterStore } from '@/stores/OutgoingLetterStore.js'

export function useOutgoingLetter() {
  const store = useOutgoingLetterStore()
  const $q = useQuasar()

  // PENTING: storeToRefs menjaga reaktivitas state saat di-destructure.
  const {
    activeTab,
    dialogSiswa,
    dialogPtk,
    dialogKedinasan,
    dialogEkspedisi,
    isPrintDialogOpen,
    filterSiswa,
    filterPtk,
    filterKedinasan,
    isLoading,
    listSiswa,
    listPtk,
    listKedinasan,
    formSiswa,
    formPtk,
    formKedinasan,
    formEkspedisi,
    generatedNomorSiswa,
    generatedNomorPtk,
    generatedNomorKedinasan,
    letterToPrint,
    jenisSuratSiswa,
    jenisSuratPtk,
    jenisSuratKedinasan,
    dudiOptions,
    daftarSiswaOptions,
    semuaSuratOptions
  } = storeToRefs(store)

  /**
   * Composable bertanggung jawab untuk:
   * 1. Memanggil store action
   * 2. Menangani notifikasi UI (success/error)
   * 3. Memformat data untuk UI jika diperlukan
   */
  async function submitCreateSiswaWithNotification() {
    try {
      await store.submitCreateSiswa()
      $q.notify({
        type: 'positive',
        message: 'Surat siswa berhasil dibuat'
      })
    } catch (error) {
      $q.notify({
        type: 'negative',
        message: error.message || 'Gagal membuat surat siswa'
      })
    }
  }

  return {
    // 1. State (Reaktif via storeToRefs)
    activeTab,
    dialogSiswa,
    dialogPtk,
    dialogKedinasan,
    dialogEkspedisi,
    isPrintDialogOpen,
    filterSiswa,
    filterPtk,
    filterKedinasan,
    isLoading,
    listSiswa,
    listPtk,
    listKedinasan,
    formSiswa,
    formPtk,
    formKedinasan,
    formEkspedisi,
    generatedNomorSiswa,
    generatedNomorPtk,
    generatedNomorKedinasan,
    letterToPrint,
    jenisSuratSiswa,
    jenisSuratPtk,
    jenisSuratKedinasan,
    dudiOptions,
    daftarSiswaOptions,
    semuaSuratOptions,

    // 2. Actions (Fungsi langsung dari store)
    loadList: store.loadList,
    loadAllLists: store.loadAllLists,
    loadOptions: store.loadOptions,
    handleJenisChange: store.handleJenisChange,
    openDialogSiswa: store.openDialogSiswa,
    closeDialogSiswa: store.closeDialogSiswa,
    openPrintDialog: store.openPrintDialog,
    closePrintDialog: store.closePrintDialog,
    triggerPrint: store.triggerPrint,

    // 3. Composable-specific actions (dengan notification wrapper)
    submitCreateSiswa: submitCreateSiswaWithNotification
  }
}
