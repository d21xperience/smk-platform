import { defineStore } from 'pinia'

export const useOutgoingLetterStore = defineStore('outgoingLetter', {
  state: () => ({
    // --- UI Workflow State ---
    activeTab: 'siswa',
    dialogSiswa: false,
    dialogPtk: false,
    dialogKedinasan: false,
    dialogEkspedisi: false,
    isPrintDialogOpen: false,

    // --- Filter & Loading ---
    filterSiswa: '',
    filterPtk: '',
    filterKedinasan: '',
    isLoading: false,

    // --- Data State ---
    listSiswa: [],
    listPtk: [],
    listKedinasan: [],

    // --- Form State (Single Source of Truth) ---
    formSiswa: {
      jenis: null,
      dudi: null,
      siswaIds: [],
      namaSiswa: '',
      kelas: '',
      keperluan: '',
      namaSalah: '',
      namaBenar: '',
      keterangan: '',
      nominal: null,
      rincian: '',
    },
    formPtk: {
      jenis: null,
      namaGuru: '',
      nip: '',
      jabatan: '',
      keperluan: '',
      tujuan: '',
      tanggalBerangkat: '',
      tanggalKembali: '',
      kegiatan: '',
      tanggalKegiatan: '',
    },
    formKedinasan: {
      jenis: null,
      tujuan: '',
      acara: '',
      tanggalRapat: '',
      tempat: '',
      namaMitra: '',
      bidangKerjasama: '',
      tujuanKerjasama: '',
      noSuratDituju: '',
      isiBalasan: '',
      tentang: '',
      uraian: '',
    },
    formEkspedisi: { selectedSurat: null, tanggalKirim: '', kurir: '', tandaTerima: '' },

    // --- Generated & Options ---
    generatedNomorSiswa: '',
    generatedNomorPtk: '',
    generatedNomorKedinasan: '',
    letterToPrint: null,

    jenisSuratSiswa: [
      { label: 'Surat Pengantar PKL', value: 'pkl' },
      { label: 'Surat Keterangan Aktif', value: 'aktif' },
      { label: 'Surat Koreksi Nama', value: 'kesalahan_nama' },
      { label: 'Surat Tunggakan', value: 'tunggakan' },
    ],
    jenisSuratPtk: [
      { label: 'Surat Tugas / Aktif', value: 'aktif_tugas' },
      { label: 'SPPD', value: 'sppd' },
      { label: 'Izin MGMP', value: 'mgmp' },
    ],
    jenisSuratKedinasan: [
      { label: 'Undangan', value: 'undangan' },
      { label: 'MoU', value: 'mou' },
      { label: 'Surat Balasan', value: 'balasan' },
      { label: 'SK', value: 'sk' },
    ],
    dudiOptions: [{ label: 'PT. Teknologi Maju', value: 'dudi_1' }],
    daftarSiswaOptions: [{ label: 'Ahmad (XII RPL 1)', value: 'siswa_1' }],
    semuaSuratOptions: [{ label: 'S.001/PKL/2024 - Ahmad', value: 'surat_1' }],
  }),

  actions: {
    // --- Dialog Actions ---
    openDialogSiswa() {
      this.dialogSiswa = true
      this.generatedNomorSiswa = 'PREVIEW-001/SISWA/2024'
    },
    closeDialogSiswa() {
      this.dialogSiswa = false
      this.resetFormSiswa()
    },
    resetFormSiswa() {
      this.formSiswa = {
        jenis: null,
        dudi: null,
        siswaIds: [],
        namaSiswa: '',
        kelas: '',
        keperluan: '',
        namaSalah: '',
        namaBenar: '',
        keterangan: '',
        nominal: null,
        rincian: '',
      }
      this.generatedNomorSiswa = ''
    },

    // --- Workflow Actions ---
    handleJenisChange(kategori, val) {
      if (kategori === 'siswa') {
        this.generatedNomorSiswa = `PREVIEW-${val.value.toUpperCase()}/2024`
      }
    },

    /**
     * PENTING: Store TIDAK boleh menggunakan Quasar atau UI framework apapun.
     * Store hanya menangani business logic dan state management.
     * Notification dan UI feedback adalah tanggung jawab Composable/Page.
     */
    async submitCreateSiswa() {
      this.isLoading = true
      try {
        // Simulasi delay network
        await new Promise((resolve) => setTimeout(resolve, 500))

        // TODO: CALL SERVICE DI BATCH SELANJUTNYA
        // await outgoingLetterService.createSiswa(this.formSiswa)

        // Store hanya menutup dialog dan reset form
        // Composable yang akan handle notification
        this.closeDialogSiswa()

        // Refresh list
        await this.loadList('siswa')

        // Return success indicator untuk composable
        return { success: true }
      } catch (error) {
        console.error('Gagal membuat siswa:', error)
        // Store hanya throw error, composable yang handle notification
        throw new Error('Gagal menyimpan data')
      } finally {
        this.isLoading = false
      }
    },

    openPrintDialog(row) {
      this.letterToPrint = row
      this.isPrintDialogOpen = true
    },
    closePrintDialog() {
      this.isPrintDialogOpen = false
      this.letterToPrint = null
    },
    triggerPrint() {
      window.print()
    },

    // --- Data Loading Actions (Mock) ---
    async loadList(kategori) {
      this.isLoading = true
      try {
        await new Promise((resolve) => setTimeout(resolve, 500))
        if (kategori === 'siswa') {
          this.listSiswa = [
            {
              id: 1,
              nomorSurat: '001/PKL/2024',
              jenis: 'pkl',
              tujuan: 'PT. Teknologi Maju',
              tanggal: '2024-08-20',
              status: 'Selesai',
            },
          ]
        }
      } finally {
        this.isLoading = false
      }
    },
    async loadAllLists() {
      await this.loadList('siswa')
      await this.loadList('ptk')
      await this.loadList('kedinasan')
    },
    loadOptions() {
      // Mock load options
    },
  },
})
