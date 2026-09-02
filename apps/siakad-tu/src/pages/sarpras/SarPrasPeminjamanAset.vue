<template>
  <q-page class="q-pa-md">
    <!-- PANEL MONITORING RINGKAS -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-blue-1 text-blue-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Sedang Dipinjam</div>
              <div class="text-h6 text-weight-bold">{{ totalDipinjam }} Transaksi</div>
            </div>
            <q-space />
            <q-icon name="outbound" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-red-1 text-red-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Terlambat Kembali</div>
              <div class="text-h6 text-weight-bold">{{ totalTerlambat }} Transaksi</div>
            </div>
            <q-space />
            <q-icon name="running_with_errors" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4 text-right">
        <q-btn color="primary" icon="add" label="Catat Pinjaman Baru" class="full-width-xs height-full"
          @click="dialogPinjamTerbuka = true" />
      </div>
    </div>

    <!-- TABEL UTAMA SIRKULASI PEMINJAMAN -->
    <q-table title="Log Transaksi Peminjaman Aset" :rows="daftarPeminjaman" :columns="columns" row-key="id"
      :filter="filterText" flat bordered>
      <!-- Slot Cari Data -->
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari peminjam/barang...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>

      <!-- Kustomisasi Tampilan Status Transaksi -->
      <template v-slot:body-cell-status="props">
        <q-td :props="props">
          <q-chip
            :color="props.value === 'Kembali' ? 'green-2' : apakahTerlambat(props.row.tglKembali) ? 'red-2' : 'orange-2'"
            :text-color="props.value === 'Kembali' ? 'green-9' : apakahTerlambat(props.row.tglKembali) ? 'red-9' : 'orange-9'"
            dense class="text-weight-bold">
            {{ props.value === 'Dipinjam' && apakahTerlambat(props.row.tglKembali) ? 'Terlambat' : props.value }}
          </q-chip>
        </q-td>
      </template>

      <!-- Kolom Tombol Pengembalian Cepat -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <q-btn v-if="props.row.status === 'Dipinjam'" dense size="sm" color="secondary" icon="assignment_returned"
            label="Proses Kembali" @click="prosesPengembalian(props.row)" />
          <span v-else class="text-grey-5 text-caption">Selesai</span>
        </q-td>
      </template>
    </q-table>

    <!-- POP-UP DIALOG: FORM CATAT PINJAMAN BARU -->
    <q-dialog v-model="dialogPinjamTerbuka">
      <q-card style="width: 450px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6"><q-icon name="edit" /> Form Peminjaman Baru</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-input outlined dense v-model="formBaru.peminjam" label="Nama Peminjam (Guru/Siswa) *" />
          <q-select outlined dense v-model="formBaru.peran" :options="['Guru', 'Staf TU', 'Siswa / OSIS']"
            label="Status Peminjam" />
          <q-select outlined dense v-model="formBaru.barang"
            :options="['Kamera DSLR Canon 02', 'Laptop ASUS TUF 01', 'Speaker Wireless Portable']"
            label="Pilih Barang Milik Sekolah *" />

          <q-input outlined dense v-model="formBaru.tglKembali" label="Tanggal Batas Pengembalian" mask="date">
            <template v-slot:append>
              <q-icon name="event" class="cursor-pointer">
                <q-popup-proxy cover transition-show="scale" transition-hide="scale">
                  <q-date v-model="formBaru.tglKembali">
                    <div class="row items-center justify-end">
                      <q-btn v-close-popup label="Pilih" color="primary" flat />
                    </div>
                  </q-date>
                </q-popup-proxy>
              </q-icon>
            </template>
          </q-input>
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Transaksi" @click="simpanPinjamanBaru" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')
const dialogPinjamTerbuka = ref(false)

// Konstanta Waktu Simulasi Sistem Saat Ini
const HARI_INI = new Date('2026/06/06')

// 1. Data Mock Sirkulasi Peminjaman Aset Sekolah
const daftarPeminjaman = ref([
  { id: 'TRX-001', peminjam: 'Pak Sandi, S.Pd', peran: 'Guru', barang: 'Laptop ASUS TUF 01', tglPinjam: '2026/06/01', tglKembali: '2026/06/05', status: 'Dipinjam' },
  { id: 'TRX-002', peminjam: 'Rian Hidayat (Ketua OSIS)', peran: 'Siswa / OSIS', barang: 'Kamera DSLR Canon 02', tglPinjam: '2026/06/05', tglKembali: '2026/06/07', status: 'Dipinjam' },
  { id: 'TRX-003', peminjam: 'Bu Anita, M.Pd', peran: 'Guru', barang: 'Speaker Wireless Portable', tglPinjam: '2026/05/28', tglKembali: '2026/05/30', status: 'Kembali' }
])

// 2. State untuk Form Input Peminjaman Baru
const formBaru = ref({
  peminjam: '',
  peran: 'Guru',
  barang: '',
  tglKembali: '2026/06/08'
})

// 3. Definisi Struktur Kolom Tabel Transaksi
const columns = [
  { name: 'id', label: 'ID Transaksi', align: 'left', field: 'id', sortable: true },
  { name: 'peminjam', label: 'Nama Peminjam', align: 'left', field: 'peminjam', sortable: true },
  { name: 'barang', label: 'Aset yang Dibawa', align: 'left', field: 'barang' },
  { name: 'tglPinjam', label: 'Tgl Pinjam', align: 'center', field: 'tglPinjam' },
  { name: 'tglKembali', label: 'Batas Kembali', align: 'center', field: 'tglKembali' },
  { name: 'status', label: 'Status Barang', align: 'center', field: 'status', sortable: true },
  { name: 'aksi', label: 'Aksi Pengembalian', align: 'center' }
]

// 4. Logika Validasi Waktu Jatuh Tempo Keterlambatan
function apakahTerlambat(tglKembaliTeks) {
  const tglTarget = new Date(tglKembaliTeks)
  return HARI_INI > tglTarget
}

// 5. Perhitungan Computed Statistik Dashboard Atas
const totalDipinjam = computed(() => {
  return daftarPeminjaman.value.filter(t => t.status === 'Dipinjam').length
})

const totalTerlambat = computed(() => {
  return daftarPeminjaman.value.filter(t => t.status === 'Dipinjam' && apakahTerlambat(t.tglKembali)).length
})

// 6. Fungsi Eksekusi Manajemen Alur Peminjaman
function prosesPengembalian(row) {
  $q.dialog({
    title: 'Konfirmasi Pengembalian',
    message: `Apakah Anda menyatakan barang "${row.barang}" yang dibawa oleh ${row.peminjam} telah kembali dalam kondisi baik?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    row.status = 'Kembali'
    $q.notify({
      color: 'green',
      icon: 'check_circle',
      message: 'Pengembalian aset sukses dicatat dalam database.'
    })
  })
}

function simpanPinjamanBaru() {
  if (!formBaru.value.peminjam || !formBaru.value.barang) {
    $q.notify({ color: 'negative', message: 'Mohon lengkapi kolom nama dan barang.' })
    return
  }

  daftarPeminjaman.value.unshift({
    id: `TRX-00${daftarPeminjaman.value.length + 1}`,
    peminjam: formBaru.value.peminjam,
    peran: formBaru.value.peran,
    barang: formBaru.value.barang,
    tglPinjam: '2026/06/06',
    tglKembali: formBaru.value.tglKembali,
    status: 'Dipinjam'
  })

  $q.notify({ color: 'positive', icon: 'done', message: 'Transaksi peminjam baru berhasil didaftarkan.' })

  // Reset form
  formBaru.value.peminjam = ''
  formBaru.value.barang = ''
}
</script>

<style scoped>
@media (max-width: 599px) {
  .full-width-xs {
    width: 100%;
  }
}
</style>
