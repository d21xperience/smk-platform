<template>
  <q-page class="q-pa-md">
    <!-- PANEL SELEKTOR RUANGAN -->
    <q-card flat bordered class="q-mb-md bg-blue-grey-1">
      <q-card-section class="row items-center q-col-gutter-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-select outlined dense bg-color="white" v-model="ruanganTerpilih" :options="opsiRuangan" option-label="nama"
            label="Pilih Lokasi Ruangan / Kelas" @update:model-value="onRuanganChange">
            <template v-slot:prepend>
              <q-icon name="room" color="primary" />
            </template>
          </q-select>
        </div>

        <q-space />

        <!-- Tombol Cetak Dokumen KIR Pintu -->
        <div class="col-12 col-sm-auto">
          <q-btn color="secondary" icon="print" label="Cetak Lembar KIR Pintu" :disabled="!ruanganTerpilih"
            @click="cetakKIR" />
        </div>
      </q-card-section>
    </q-card>

    <!-- KONTEN DATA KIR (Hanya muncul jika ruangan sudah dipilih) -->
    <div v-if="ruanganTerpilih">
      <!-- INFO DETAIL PENANGGUNG JAWAB RUANG -->
      <q-card flat bordered class="q-mb-md border-primary shadow-1">
        <q-card-section class="row items-center justify-between q-py-md">
          <div class="row items-center">
            <q-avatar icon="person" color="primary" text-color="white" size="md" class="q-mr-sm" />
            <div>
              <div class="text-caption text-grey-7">Penanggung Jawab Ruangan (PJR)</div>
              <div class="text-subtitle1 text-weight-bold text-grey-9">
                {{ infoRuangan.penanggungJawab }}
              </div>
            </div>
          </div>

          <q-separator vertical class="gt-xs q-mx-md" />

          <div>
            <div class="text-caption text-grey-7 text-right-gt-xs">Total Inventaris Ruang</div>
            <div class="text-subtitle1 text-weight-bold text-primary text-right-gt-xs">
              {{ asetDitemukan.length }} Unit Barang
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- TABEL DATA BARANG DI RUANGAN -->
      <q-table :title="`Daftar Inventaris: ${ruanganTerpilih.nama}`" :rows="asetDitemukan" :columns="columns"
        row-key="id" flat bordered :rows-per-page-options="[0]">
        <!-- Kustomisasi Tampilan Kolom Status Kondisi -->
        <template v-slot:body-cell-status="props">
          <q-td :props="props">
            <q-chip :color="props.value === 'Bagus' ? 'green-2' : props.value === 'Rusak Ringan' ? 'orange-2' : 'red-2'"
              :text-color="props.value === 'Bagus' ? 'green-9' : props.value === 'Rusak Ringan' ? 'orange-9' : 'red-9'"
              dense class="text-weight-medium">
              {{ props.value }}
            </q-chip>
          </q-td>
        </template>

        <!-- Kustomisasi Kolom Barcode/QR teks -->
        <template v-slot:body-cell-id="props">
          <q-td :props="props">
            <span class="text-monospace text-weight-bold text-primary">{{ props.value }}</span>
          </q-td>
        </template>
      </q-table>
    </div>

    <!-- TAMPILAN KOSONG BILA BELUM MEMILIH RUANG -->
    <div v-else class="text-center q-pa-xl text-grey-6">
      <q-icon name="maps_home_work" size="100px" class="q-mb-md text-grey-4" />
      <div class="text-h6 font-weight-light">Silakan pilih ruangan terlebih dahulu untuk melihat Kartu Inventaris Ruang
        (KIR).</div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// 1. Data Pilihan Ruangan di Sekolah Swasta
const opsiRuangan = ref([
  { id: 'R01', nama: 'Ruang Laboratorium Komputer' },
  { id: 'R02', nama: 'Ruang Kelas 10-A (Gedung Utama)' },
  { id: 'R03', nama: 'Ruang Teori Guru & Kepala Sekolah' }
])

const ruanganTerpilih = ref(null)

// 2. Master Data Mock Seluruh Aset Terdaftar di Sekolah
const masterAset = ref([
  { id: 'YYS/SAR-COMP/2025/001', name: 'PC Komputer Rakitan Core i5', category: 'Elektronik', merk: 'Asus', status: 'Bagus', lokasiId: 'R01' },
  { id: 'YYS/SAR-COMP/2025/002', name: 'PC Komputer Rakitan Core i5', category: 'Elektronik', merk: 'Asus', status: 'Bagus', lokasiId: 'R01' },
  { id: 'YYS/SAR-PROJ/2024/012', name: 'Proyektor BenQ Wall Mount', category: 'Elektronik', merk: 'BenQ', status: 'Rusak Ringan', lokasiId: 'R01' },
  { id: 'YYS/SAR-MEBL/2023/105', name: 'Meja Komputer Kayu Panjang', category: 'Mebel', merk: 'Custom', status: 'Bagus', lokasiId: 'R01' },

  { id: 'YYS/SAR-PROJ/2024/015', name: 'Proyektor BenQ Wall Mount', category: 'Elektronik', merk: 'BenQ', status: 'Bagus', lokasiId: 'R02' },
  { id: 'YYS/SAR-MEBL/2026/040', name: 'Meja Siswa Chitose', category: 'Mebel', merk: 'Chitose', status: 'Bagus', lokasiId: 'R02' },
  { id: 'YYS/SAR-ACUN/2023/004', name: 'AC Split Sharp 1 PK', category: 'Elektronik', merk: 'Sharp', status: 'Rusak Berat', lokasiId: 'R02' },

  { id: 'YYS/SAR-MEBL/2022/001', name: 'Sofa Tamu Eksekutif Yayasan', category: 'Mebel', merk: 'Informa', status: 'Bagus', lokasiId: 'R03' }
])

// 3. Mock Data Info Tambahan untuk Ruangan (Penanggung Jawab)
const masterPJR = {
  'R01': { penanggungJawab: 'Pak Angga Saputra, S.Kom (Kepala Lab)' },
  'R02': { penanggungJawab: 'Ibu Nenden Herlina, S.Pd (Wali Kelas 10-A)' },
  'R03': { penanggungJawab: 'Drs. H. Mulyadi (Kepala Sekolah)' }
}

// 4. Struktur Kolom Tabel KIR
const columns = [
  { name: 'id', label: 'Kode Unik Aset (Kodifikasi)', align: 'left', field: 'id', sortable: true },
  { name: 'nama', label: 'Nama Barang', align: 'left', field: 'name', sortable: true },
  { name: 'merk', label: 'Merk/Brand', align: 'left', field: 'merk' },
  { name: 'category', label: 'Jenis Kategori', align: 'center', field: 'category' },
  { name: 'status', label: 'Kondisi Fisik', align: 'center', field: 'status', sortable: true }
]

// 5. Computed Properties untuk memfilter barang berdasarkan ruangan aktif
const asetDitemukan = computed(() => {
  if (!ruanganTerpially.value) return []
  return masterAset.value.filter(aset => aset.lokasiId === ruanganTerpilih.value.id)
})

const infoRuangan = computed(() => {
  if (!ruanganTerpilih.value) return { penanggungJawab: '-' }
  return masterPJR[ruanganTerpilih.value.id] || { penanggungJawab: 'Belum Ditunjuk' }
})

// Alias untuk reactivity drop-down
const ruanganTerpially = ruanganTerpilih

function onRuanganChange(val) {
  if (val) {
    $q.notify({
      message: `Menampilkan data aset untuk ${val.nama}`,
      color: 'primary',
      icon: 'info',
      timeout: 1000
    })
  }
}

// 6. Simulasi Ekspor / Cetak Dokumen Fisik Pintu
function cetakKIR() {
  $q.notify({
    message: `Menyiapkan cetakan lembar KIR resmi untuk ${ruanganTerpilih.value.nama}...`,
    color: 'secondary',
    icon: 'local_printshop'
  })
  // Di masa mendatang, Anda bisa mengintegrasikan window.print() dengan CSS @media print
}
</script>

<style scoped>
.border-primary {
  border-left: 5px solid var(--q-primary);
}

.text-monospace {
  font-family: 'Courier New', Courier, monospace;
}

@media gt-xs {
  .text-right-gt-xs {
    text-align: right;
  }
}
</style>
