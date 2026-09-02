<template>
  <q-page class="q-pa-md">
    <!-- PANEL REKAP & AKSI MUTASI -->
    <div class="row q-col-gutter-sm q-mb-md items-center justify-between">
      <div class="col-12 col-sm-6 col-md-4">
        <q-input outlined dense v-model="filterText" placeholder="Cari nama barang, kode, atau ruangan...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>

      <div class="col-12 col-sm-auto">
        <q-btn color="primary" icon="move_down" label="Catat Mutasi Lokasi" @click="dialogMutasiTerbuka = true" />
      </div>
    </div>

    <!-- TABEL LOG RIWAYAT MUTASI ASET -->
    <q-table title="Riwayat Perpindahan (Mutasi) Aset Internal" :rows="daftarMutasi" :columns="columns" row-key="id"
      :filter="filterText" flat bordered>
      <!-- Kustomisasi Visual Alur Perpindahan Ruang -->
      <template v-slot:body-cell-alurPerpindahan="props">
        <q-td :props="props">
          <div class="row items-center no-wrap text-weight-medium">
            <q-chip dense color="amber-2" text-color="amber-10" icon="door_sliding">
              {{ props.row.asal }}
            </q-chip>
            <q-icon name="trending_flat" size="sm" class="q-mx-xs text-primary text-weight-bolder" />
            <q-chip dense color="green-2" text-color="green-10" icon="meeting_room">
              {{ props.row.tujuan }}
            </q-chip>
          </div>
        </q-td>
      </template>

      <!-- Kustomisasi Tampilan Kode Unik Aset -->
      <template v-slot:body-cell-kodeAset="props">
        <q-td :props="props">
          <span class="text-monospace text-weight-bold text-indigo-7">{{ props.value }}</span>
        </q-td>
      </template>
    </q-table>

    <!-- POP-UP DIALOG: FORM CATAT MUTASI INTERNAL BARU -->
    <q-dialog v-model="dialogMutasiTerbuka">
      <q-card style="width: 450px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6"><q-icon name="swap_horiz" /> Form Pemindahan Lokasi Barang</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <!-- Pilih Aset Fisik -->
          <q-select outlined dense v-model="formMutasi.aset" :options="opsiAset" option-label="label"
            label="Pilih Barang / Kode Unik *" />

          <div class="row q-col-gutter-sm">
            <!-- Ruang Asal -->
            <div class="col-6">
              <q-select outlined dense v-model="formMutasi.asal" :options="opsiRuangan" label="Ruang Asal *" />
            </div>
            <!-- Ruang Tujuan -->
            <div class="col-6">
              <q-select outlined dense v-model="formMutasi.tujuan" :options="opsiRuangan"
                label="Ruang Tujuan Pemindahan *" />
            </div>
          </div>

          <!-- Penanggung Jawab Eksekutor TU -->
          <q-input outlined dense v-model="formMutasi.pj" label="Penanggung Jawab Pemindahan (Staf TU/Kurir) *" />

          <!-- Catatan / Alasan Mutasi -->
          <q-input outlined dense type="textarea" v-model="formMutasi.alasan"
            label="Alasan Pemindahan / Catatan Referensi" rows="2"
            placeholder="Contoh: Lab Komputer sedang dicat, kursi dialihkan sementara." />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Simpan Perubahan Lokasi" @click="simpanMutasiBaru" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()
const filterText = ref('')
const dialogMutasiTerbuka = ref(false)

// 1. Data Mock Riwayat Mutasi Internal Sekolah
const daftarMutasi = ref([
  { id: 'MTS-001', tanggal: '2026/06/05', kodeAset: 'YYS/SAR-MEBL/2026/040', namaBarang: 'Kursi Belajar Siswa Chitose', asal: 'Gudang Utama', tujuan: 'Ruang Kelas 10-A', pj: 'Pak Joko (Staf TU)', alasan: 'Penambahan kuantitas siswa baru di kelas 10-A.' },
  { id: 'MTS-002', tanggal: '2026/06/02', kodeAset: 'YYS/SAR-PROJ/2025/012', namaBarang: 'Proyektor BenQ Portable', asal: 'Ruang Lab Fisika', tujuan: 'Ruang Rapat Yayasan', pj: 'Pak Angga Saputra', alasan: 'Dipindahkan permanen untuk keperluan presentasi pengurus yayasan.' }
])

// 2. State untuk Form Pengisian Mutasi Baru
const formMutasi = ref({
  aset: null,
  asal: '',
  tujuan: '',
  pj: '',
  alasan: ''
})

// 3. Dropdown Opsi Pilihan Mock Data
const opsiRuangan = ['Gudang Utama', 'Ruang Kelas 10-A', 'Ruang Kelas 11-B', 'Lab Fisika', 'Ruang Rapat Yayasan', 'Ruang Guru']
const opsiAset = [
  { kode: 'YYS/SAR-COMP/2026/001', label: 'PC Komputer Rakitan (YYS/SAR-COMP/2026/001)' },
  { kode: 'YYS/SAR-MEBL/2026/099', label: 'Meja Lipat Kayu (YYS/SAR-MEBL/2026/099)' }
]

// 4. Struktur Arsitektur Kolom Tabel Riwayat Mutasi
const columns = [
  { name: 'id', label: 'ID Log', align: 'left', field: 'id', sortable: true },
  { name: 'tanggal', label: 'Tanggal Pindah', align: 'center', field: 'tanggal', sortable: true },
  { name: 'kodeAset', label: 'Kode Unik Aset', align: 'left', field: 'kodeAset' },
  { name: 'namaBarang', label: 'Nama Inventaris', align: 'left', field: 'namaBarang', sortable: true },
  { name: 'alurPerpindahan', label: 'Sirkulasi Ruangan (Asal -> Tujuan)', align: 'left' },
  { name: 'pj', label: 'Eksekutor / PJ', align: 'left', field: 'pj' },
  { name: 'alasan', label: 'Keterangan Mutasi', align: 'left', field: 'alasan', format: val => val || '-' }
]

// 5. Fungsi Menyimpan Data Mutasi Lokasi Baru ke Tabel
function simpanMutasiBaru() {
  if (!formMutasi.value.aset || !formMutasi.value.asal || !formMutasi.value.tujuan || !formMutasi.value.pj) {
    $q.notify({ color: 'negative', message: 'Gagal memproses. Mohon isi semua kolom bertanda bintang (*).' })
    return
  }

  if (formMutasi.value.asal === formMutasi.value.tujuan) {
    $q.notify({ color: 'warning', message: 'Ruang asal dan tujuan tidak boleh sama.' })
    return
  }

  // Ekstrak nama murni dari label objek dropdown aset
  const namaMurniAset = formMutasi.value.aset.label.split(' (')[0]

  daftarMutasi.value.unshift({
    id: `MTS-00${daftarMutasi.value.length + 1}`,
    tanggal: '2026/06/06', // Mengikuti tanggal sistem terkini
    kodeAset: formMutasi.value.aset.kode,
    namaBarang: namaMurniAset,
    asal: formMutasi.value.asal,
    tujuan: formMutasi.value.tujuan,
    pj: formMutasi.value.pj,
    alasan: formMutasi.value.alasan
  })

  $q.notify({
    color: 'positive',
    icon: 'local_shipping',
    message: 'Lokasi penempatan aset berhasil dimutasi dalam Kartu Inventaris Ruang.'
  })

  // Reset form isian
  formMutasi.value = { aset: null, asal: '', tujuan: '', pj: '', alasan: '' }
}
</script>

<style scoped>
.text-monospace {
  font-family: 'Courier New', Courier, monospace;
}
</style>
