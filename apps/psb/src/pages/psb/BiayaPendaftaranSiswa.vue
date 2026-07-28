<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-teal-9">
          Informasi Keuangan & Kuitansi
        </h1>
        <p class="text-caption text-grey-7 q-mb-none">
          Pantau riwayat cicilan, sisa tagihan, dan cetak bukti kuitansi pembayaran resmi PPDB
        </p>
      </div>
      <div>
        <!-- Tombol cetak hanya aktif jika siswa sudah pernah mengangsur/membayar -->
        <q-btn
          color="teal-8"
          icon="print"
          label="Cetak Struk Kuitansi"
          :disable="pembayaran.riwayat.length === 0"
          @click="cetakUlangKuitansi"
        />
      </div>
    </div>

    <div class="row q-col-gutter-md">
      <!-- PANEL KIRI: Status Tagihan & Metode Pembayaran -->
      <div class="col-12 col-md-5">
        <!-- Card Status Saldo -->
        <q-card flat class="shadow-1 q-mb-md">
          <q-card-section class="bg-teal-8 text-white q-py-sm">
            <div class="text-subtitle2 text-weight-bold flex items-center">
              <q-icon name="account_balance_wallet" class="q-mr-xs" /> Rincian Status Finansial
            </div>
          </q-card-section>
          <q-card-section class="q-gutter-xs">
            <div class="row justify-between q-py-sm border-bottom">
              <span class="text-grey-7">Tarif Resmi Pendaftaran:</span>
              <span class="text-weight-bold text-grey-9"
                >Rp {{ pembayaran.totalTarif.toLocaleString('id-ID') }}</span
              >
            </div>
            <div class="row justify-between q-py-sm border-bottom text-green-8 text-weight-medium">
              <span>Total Telah Dibayarkan:</span>
              <span>Rp {{ totalTerbayar.toLocaleString('id-ID') }}</span>
            </div>
            <div
              class="row justify-between q-py-sm text-red-9 text-weight-bold bg-red-1 q-px-sm rounded-borders q-mt-sm"
            >
              <span>Sisa Sisa Tunggakan:</span>
              <span>Rp {{ sisaTunggakan.toLocaleString('id-ID') }}</span>
            </div>
          </q-card-section>
        </q-card>

        <!-- Card Panduan Pembayaran Loket/Transfer -->
        <q-card flat class="shadow-1">
          <q-card-section class="bg-grey-3 text-grey-9 q-py-sm text-weight-bold">
            <q-icon name="info" color="teal-8" class="q-mr-xs" /> Metode Pembayaran Sah
          </q-card-section>
          <q-card-section class="text-caption text-grey-8 q-gutter-sm">
            <div>
              <strong>1. Pembayaran Tunai (Cash):</strong><br />
              Silakan datang langsung ke loket keuangan panitia PPDB di ruang Tata Usaha Sekolah
              dengan membawa nomor registrasi Anda.
            </div>
            <q-separator class="q-my-xs" />
            <div>
              <strong>2. Transfer Bank:</strong><br />
              Transfer melalui rekening Bank Mandiri virtual account resmi yayasan:<br />
              <span class="text-subtitle2 text-primary text-weight-bold">8872-0091-1223-32</span
              ><br />
              a.n <strong>PPDB SMK SWASTA INDONESIA</strong>

              <!-- FORM INPUT UPLOAD BUKTI TRANSFER KHUSUS METODE 2 -->
              <div v-if="sisaTunggakan > 0" class="q-mt-md q-pa-sm bg-grey-2 rounded-borders">
                <div class="text-weight-bold text-grey-7 q-mb-xs" style="font-size: 11px">
                  Konfirmasi Pembayaran Transfer Bank:
                </div>

                <q-file
                  v-model="berkasBuktiTransfer"
                  label="Pilih struk transfer (.jpg, .png, .pdf)"
                  outlined
                  dense
                  bg-color="white"
                  accept=".jpg, .jpeg, .png, .pdf"
                  max-file-size="3145728"
                >
                  <template v-slot:prepend>
                    <q-icon name="receipt_long" color="teal-8" />
                  </template>
                  <template v-slot:append v-if="berkasBuktiTransfer">
                    <q-icon
                      name="cancel"
                      @click.stop.prevent="berkasBuktiTransfer = null"
                      class="cursor-pointer"
                    />
                  </template>
                </q-file>

                <q-btn
                  color="teal-8"
                  label="Kirim Bukti Transfer"
                  dense
                  class="q-mt-sm full-width text-weight-bold"
                  size="sm"
                  @click="kirimBuktiTransfer"
                />
              </div>
            </div>
            <q-separator class="q-my-xs" />
            <div>
              <strong>3. Pembayaran Instan Nontunai (QRIS):</strong><br />
              Scan kode QRIS resmi sekolah menggunakan m-Banking atau e-Wallet pilihan Anda untuk
              proses pembayaran yang lebih cepat.<br />
              <q-btn
                color="teal-8"
                icon="qr_code_2"
                label="Bayar via QRIS"
                dense
                class="q-mt-sm q-px-sm"
                :disable="sisaTunggakan === 0"
                @click="openQrisPayment"
              />
            </div>
          </q-card-section>
        </q-card>
        <!-- ====== SISIPKAN METODE 3 DI BAWAH METODE 2 PADA CARD PANDUAN ====== -->
        <q-separator class="q-my-xs" />
      </div>

      <!-- PANEL KANAN: Tabel Riwayat Angsuran Masuk -->
      <div class="col-12 col-md-7">
        <q-card flat class="shadow-1 fit">
          <q-card-section class="bg-grey-1 q-py-sm flex justify-between items-center">
            <div class="text-subtitle2 text-weight-bold text-grey-9">
              Log Histori Cicilan Dana Masuk
            </div>
            <q-badge
              :color="sisaTunggakan === 0 ? 'green' : 'orange-9'"
              :label="sisaTunggakan === 0 ? 'STATUS: LUNAS' : 'STATUS: BELUM LUNAS'"
              class="q-pa-xs text-weight-bold"
            />
          </q-card-section>

          <q-card-section class="q-pa-none">
            <q-table
              flat
              :rows="pembayaran.riwayat"
              :columns="columns"
              row-key="id"
              hide-bottom
              no-data-label="Belum ada catatan setoran dana tunai masuk."
            >
              <!-- Format nominal mata uang di dalam kolom tabel -->
              <template v-slot:body-cell-jumlah="props">
                <q-td :props="props" class="text-weight-bold text-green-8">
                  Rp {{ props.row.jumlah.toLocaleString('id-ID') }}
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>

  <!-- ====== TAMBAHKAN DIALOG MODAL QRIS INI ====== -->
  <q-dialog v-model="dialogQrisOpen" persistent>
    <q-card style="min-width: 360px" class="text-center">
      <q-card-section class="bg-teal-8 text-white q-py-sm">
        <div class="text-subtitle1 text-weight-bold">Scan QRIS PPDB Sekolah</div>
        <div class="text-caption text-teal-1">
          Sisa Tagihan Anda: Rp {{ sisaTunggakan.toLocaleString('id-ID') }}
        </div>
      </q-card-section>

      <q-card-section class="q-pa-md">
        <!-- Gambar Mockup QRIS Statis Sekolah -->
        <!-- Di masa depan, gambar ini bisa berupa QRIS dinamis yang ditembak dari API Payment Gateway Golang Anda -->
        <div
          class="bg-white q-pa-sm rounded-borders shadow-1 inline-block q-mb-md"
          style="border: 1px solid #ccc; max-width: 200px; margin: 0 auto"
        >
          <q-img src="https://qrserver.com" style="width: 180px; height: 180px" />
        </div>

        <div class="text-caption text-grey-7 q-mb-md">
          Silakan scan kode QRIS di atas, masukkan nominal angsuran sesuai kemampuan Anda, kemudian
          unggah bukti transfer di bawah ini.
        </div>

        <!-- Form Upload Bukti Bayar QRIS -->
        <q-file
          v-model="berkasBuktiQris"
          label="Unggah Bukti Transfer (.jpg, .png, PDF)"
          outlined
          dense
          accept=".jpg, .jpeg, .png, .pdf"
        >
          <template v-slot:prepend><q-icon name="cloud_upload" color="teal-8" /></template>
        </q-file>
      </q-card-section>

      <q-card-actions align="right" class="q-pb-md q-pr-md">
        <q-btn flat label="Batal" v-close-popup color="grey" />
        <q-btn color="teal-8" label="Kirim Bukti Bayar" @click="uploadBuktiQris" />
      </q-card-actions>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// Struktur Struktur Kolom Histori Pembayaran q-table
const columns = [
  { name: 'tanggal', align: 'left', label: 'Tanggal Bayar', field: 'tanggal' },
  { name: 'via', align: 'center', label: 'Metode', field: 'via' },
  { name: 'jumlah', align: 'right', label: 'Jumlah Setoran', field: 'jumlah' },
  { name: 'penerima', align: 'left', label: 'Petugas TU', field: 'penerima' },
]

// Data Mock / Dummy Record Finansial Calon Siswa (Terintegrasi dengan Case Angsuran Anda)
const pembayaran = ref({
  totalTarif: 450000,
  riwayat: [
    // Contoh kondisi siswa mencicil: baru menyetor Rp 200.000 tunai ke sekolah
    {
      id: 1,
      tanggal: '10 Juni 2026',
      via: 'Tunai/Loket',
      jumlah: 200000,
      penerima: 'Siti Aminah (TU)',
    },
  ],
})

// Menghitung akumulasi dana terbayar menggunakan JavaScript reduce
const totalTerbayar = computed(() => {
  return pembayaran.value.riwayat.reduce((sum, item) => sum + item.jumlah, 0)
})

// Menghitung sisa sisa tunggakan pendaftaran pendaftar
const sisaTunggakan = computed(() => {
  return pembayaran.value.totalTarif - totalTerbayar.value
})

const cetakUlangKuitansi = () => {
  $q.notify({
    color: 'teal-9',
    message:
      'Mempersiapkan dokumen... Mengunduh duplikat salinan kuitansi pendaftaran digital resmi.',
    icon: 'print',
  })
}
// ====== TAMBAHKAN STATE & LOGIKA QRIS INI ======
const dialogQrisOpen = ref(false)
const berkasBuktiQris = ref(null)

const openQrisPayment = () => {
  if (sisaTunggakan.value === 0) {
    $q.notify({ color: 'green-8', message: 'Biaya pendaftaran Anda sudah lunas!', icon: 'check' })
    return
  }
  dialogQrisOpen.value = true
}

const uploadBuktiQris = () => {
  if (!berkasBuktiQris.value) {
    $q.notify({
      color: 'negative',
      message: 'Silakan pilih file bukti transfer terlebih dahulu!',
      icon: 'warning',
    })
    return
  }

  $q.notify({
    color: 'green-8',
    message:
      'Bukti pembayaran QRIS berhasil diunggah! Mohon tunggu verifikasi oleh tim TU Keuangan.',
    icon: 'cloud_done',
  })
  dialogQrisOpen.value = false
  berkasBuktiQris.value = null
}
// ===============================================
// ====== TAMBAHKAN STATE & LOGIKA UPLOAD BUKTI TRANSFER ======
const berkasBuktiTransfer = ref(null)

const kirimBuktiTransfer = () => {
  if (!berkasBuktiTransfer.value) {
    $q.notify({
      color: 'negative',
      message: 'Silakan pilih file bukti transfer bank terlebih dahulu!',
      icon: 'warning',
    })
    return
  }

  // Simulasi sukses upload di sisi klien
  $q.notify({
    color: 'green-8',
    message:
      'Bukti transfer bank berhasil diunggah! Status angsuran Anda akan diperbarui setelah divalidasi oleh bendahara TU.',
    icon: 'cloud_done',
  })

  // Bersihkan input setelah berhasil dikirim
  berkasBuktiTransfer.value = null
}
// ============================================================
</script>

<style scoped>
.border-bottom {
  border-bottom: 1px solid #e0e0e0;
}

.text-h5 {
  line-height: 1.2;
}
</style>
