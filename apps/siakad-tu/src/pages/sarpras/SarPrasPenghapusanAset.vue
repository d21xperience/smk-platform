<template>
  <q-page class="q-pa-md">
    <!-- PANEL MONITORING REKAP DISPOSAL -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-blue-grey-1 text-blue-grey-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Menunggu Persetujuan Yayasan</div>
              <div class="text-h6 text-weight-bold">{{ totalDraft }} Usulan</div>
            </div>
            <q-space />
            <q-icon name="gavel" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-green-1 text-green-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Disetujui & Dihapus</div>
              <div class="text-h6 text-weight-bold">{{ totalDisetujui }} Aset</div>
            </div>
            <q-space />
            <q-icon name="delete_forever" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-4 text-right">
        <q-btn color="negative" icon="delete_sweep" label="Buat Usulan Hapus Aset" class="full-width-xs"
          @click="dialogHapusTerbuka = true" />
      </div>
    </div>

    <!-- TABEL UTAMA USULAN PENGHAPUSAN -->
    <q-table title="Daftar Pengajuan Penghapusan Barang Milik Sekolah" :rows="daftarPenghapusan" :columns="columns"
      row-key="id" :filter="filterText" flat bordered>
      <!-- Slot Cari Data -->
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari kode/nama aset...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>

      <!-- Kustomisasi Tampilan Kode Unik Aset -->
      <template v-slot:body-cell-kodeAset="props">
        <q-td :props="props">
          <span class="text-monospace text-weight-bold text-red-7">{{ props.value }}</span>
        </q-td>
      </template>

      <!-- Kustomisasi Status Dokumen Usulan -->
      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-chip :color="props.value === 'Disetujui Yayasan' ? 'green-2' : 'amber-2'"
            :text-color="props.value === 'Disetujui Yayasan' ? 'green-9' : 'amber-9'" dense class="text-weight-bold">
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <!-- Kolom Unduh Berita Acara (PDF) / Aksi Persetujuan -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <q-btn v-if="props.row.status === 'Draft Usulan'" dense size="sm" color="positive" icon="verified_user"
            label="Setujui (Yayasan)" @click="eksekusiSetujuHapus(props.row)" />
          <q-btn v-else dense flat size="sm" color="primary" icon="picture_as_pdf" label="Unduh BAP"
            @click="unduhBAP(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- POP-UP DIALOG: FORM USULAN PENGHAPUSAN BARU -->
    <q-dialog v-model="dialogHapusTerbuka">
      <q-card style="width: 450px;">
        <q-card-section class="bg-negative text-white row items-center">
          <div class="text-h6"><q-icon name="warning" /> Form Pengajuan Penghapusan Aset</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <!-- Pilih Aset dari Buku Induk -->
          <q-select outlined dense v-model="formHapus.aset" :options="opsiAset" option-label="label"
            label="Pilih Aset yang Akan Dihapus *" />

          <!-- Alasan Utama Penghapusan -->
          <q-select outlined dense v-model="formHapus.alasan"
            :options="['Rusak Berat (Tidak Bisa Diperbaiki)', 'Hilang / Kecurian', 'Dijual / Dilelang', 'Dihibahkan ke Pihak Lain']"
            label="Alasan Penghapusan *" />

          <!-- Tindakan Akhir Fisik Barang -->
          <q-input outlined dense v-model="formHapus.tindakanFisik" label="Rencana Tindakan Fisik Barang *"
            placeholder="Contoh: Dimusnahkan/Dijual sebagai loakan besi tua" />

          <q-input outlined dense type="number" v-model.number="formHapus.nilaiSisa"
            label="Taksiran Nilai Sisa Jual (Rp)" prefix="Rp" hint="Isi 0 jika barang dimusnahkan/tidak bernilai" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="negative" label="Ajukan ke Yayasan" @click="kirimUsulanHapus" v-close-popup />
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
const dialogHapusTerbuka = ref(false)

// 1. Data Mock Riwayat Pengajuan Penghapusan Aset
const daftarPenghapusan = ref([
  { id: 'BAP-2026-001', tanggal: '2026/06/02', kodeAset: 'YYS/SAR-ACUN/2023/004', namaBarang: 'AC Split Sharp 1 PK', alasan: 'Rusak Berat (Kompresor Terbakar)', nilaiSisa: 150000, status: 'Draft Usulan' },
  { id: 'BAP-2026-002', tanggal: '2026/05/20', kodeAset: 'YYS/SAR-COMP/2021/088', namaBarang: 'Monitor Tabung CRT LG', alasan: 'Dijual / Dilelang (Peremajaan Lab)', nilaiSisa: 50000, status: 'Disetujui Yayasan' }
])

// 2. State untuk Input Form Usulan Baru
const formHapus = ref({
  aset: null,
  alasan: '',
  tindakanFisik: '',
  nilaiSisa: 0
})

// 3. Dropdown Data Master Inventaris untuk Dihapus
const opsiAset = [
  { kode: 'YYS/SAR-PROJ/2024/012', label: 'Proyektor BenQ Wall Mount (Rusak Total)' },
  { kode: 'YYS/SAR-MEBL/2020/011', name: 'Lemari Kayu Arsip Guru (Lapuk Dimakan Rayap)', label: 'Lemari Kayu Arsip Guru (Lapuk)' }
]

// 4. Struktur Kolom Tabel Penghapusan Aset
const columns = [
  { name: 'id', label: 'No. Berita Acara', align: 'left', field: 'id', sortable: true },
  { name: 'tanggal', label: 'Tanggal Ajuan', align: 'center', field: 'tanggal' },
  { name: 'kodeAset', label: 'Kode Unik Aset', align: 'left', field: 'kodeAset' },
  { name: 'namaBarang', label: 'Nama Inventaris', align: 'left', field: 'namaBarang', sortable: true },
  { name: 'alasan', label: 'Alasan Utama', align: 'left', field: 'alasan' },
  { name: 'nilaiSisa', label: 'Nilai Jual Sisa', align: 'right', field: 'nilaiSisa', format: val => `Rp ${val.toLocaleString('id-ID')}` },
  { name: 'status', label: 'Status Legalitas', align: 'center', field: 'status' },
  { name: 'aksi', label: 'Opsi Dokumen', align: 'center' }
]

// 5. Perhitungan Metrik Statistik Panel Atas
const totalDraft = computed(() => {
  return daftarPenghapusan.value.filter(item => item.status === 'Draft Usulan').length
})

const totalDisetujui = computed(() => {
  return daftarPenghapusan.value.filter(item => item.status === 'Disetujui Yayasan').length
})

// 6. Alur Aksi Eksekusi Sistem (Persetujuan & Simulasi Cetak PDF)
function eksekusiSetujuHapus(row) {
  $q.dialog({
    title: 'Persetujuan Final Yayasan',
    message: `Dengan menyetujui, aset "${row.namaBarang}" dengan kode ${row.kodeAset} akan dihapus secara permanen dari Buku Induk Aktiva Tetap Sekolah Swasta. Lanjutkan?`,
    cancel: true,
    persistent: true
  }).onOk(() => {
    row.status = 'Disetujui Yayasan'
    $q.notify({
      color: 'green',
      icon: 'gavel',
      message: 'Otorisasi berhasil. Nilai aset resmi dikeluarkan dari neraca yayasan.'
    })
  })
}

function kirimUsulanHapus() {
  if (!formHapus.value.aset || !formHapus.value.alasan || !formHapus.value.tindakanFisik) {
    $q.notify({ color: 'warning', message: 'Mohon lengkapi parameter bertanda bintang (*).' })
    return
  }

  daftarPenghapusan.value.unshift({
    id: `BAP-2026-00${daftarPenghapusan.value.length + 1}`,
    tanggal: '2026/06/06',
    kodeAset: formHapus.value.aset.kode,
    namaBarang: formHapus.value.aset.label.split(' (')[0],
    alasan: formHapus.value.alasan,
    nilaiSisa: formHapus.value.nilaiSisa,
    status: 'Draft Usulan'
  })

  $q.notify({
    color: 'orange-9',
    icon: 'send',
    message: 'Draf berita acuan pemusnahan berhasil dikirim ke antrean approval yayasan.'
  })

  // Reset form isian
  formHapus.value = { aset: null, alasan: '', tindakanFisik: '', nilaiSisa: 0 }
}

function unduhBAP(row) {
  $q.notify({
    color: 'indigo',
    icon: 'picture_as_pdf',
    message: `Mengunduh Dokumen Berita Acara Penghapusan (BAP) #${row.id} format PDF...`
  })
}
</script>

<style scoped>
.text-monospace {
  font-family: 'Courier New', Courier, monospace;
}

@media (max-width: 599px) {
  .full-width-xs {
    width: 100%;
  }
}
</style>
