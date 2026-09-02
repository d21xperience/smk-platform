<template>
  <q-page class="q-pa-md">
    <!-- PANEL RINGKASAN REKAP PERMINTAAN -->
    <div class="row q-col-gutter-sm q-mb-md items-center justify-between">
      <div class="col-12 col-sm-4">
        <q-card flat bordered class="bg-amber-1 text-amber-9">
          <q-card-section class="row items-center no-wrap q-py-sm">
            <div>
              <div class="text-caption text-grey-7">Menunggu Persetujuan</div>
              <div class="text-h6 text-weight-bold">{{ totalPending }} Pengajuan</div>
            </div>
            <q-space />
            <q-icon name="pending_actions" size="sm" />
          </q-card-section>
        </q-card>
      </div>

      <div class="col-12 col-sm-auto">
        <q-btn color="primary" icon="rate_review" label="Ajukan ATK (Mode Guru)" @click="dialogMintaTerbuka = true" />
      </div>
    </div>

    <!-- TABEL LOG PERMINTAAN ATK -->
    <q-table title="Daftar Permintaan Logistik ATK Guru" :rows="daftarPermintaan" :columns="columns" row-key="id"
      :filter="filterText" flat bordered>
      <!-- Slot Cari Data -->
      <template v-slot:top-right>
        <q-input borderless dense debounce="300" v-model="filterText" placeholder="Cari nama guru/item...">
          <template v-slot:append>
            <q-icon name="search" />
          </template>
        </q-input>
      </template>

      <!-- Kustomisasi Tampilan Status Persetujuan -->
      <template v-slot:body-cell-status="props">
        <q-td :props="props" class="text-center">
          <q-chip :color="props.value === 'Disetujui' ? 'green-2' : props.value === 'Ditolak' ? 'red-2' : 'amber-2'"
            :text-color="props.value === 'Disetujui' ? 'green-9' : props.value === 'Ditolak' ? 'red-9' : 'amber-9'"
            dense class="text-weight-bold">
            {{ props.value }}
          </q-chip>
        </q-td>
      </template>

      <!-- Kolom Aksi Persetujuan untuk Staf Tata Usaha -->
      <template v-slot:body-cell-aksi="props">
        <q-td :props="props" class="text-center">
          <div v-if="props.row.status === 'Pending'" class="row justify-center q-gutter-xs">
            <q-btn dense round size="sm" color="positive" icon="check"
              @click="updateStatusPermintaan(props.row, 'Disetujui')">
              <q-tooltip>Setujui & Kurangi Stok</q-tooltip>
            </q-btn>
            <q-btn dense round size="sm" color="negative" icon="close"
              @click="updateStatusPermintaan(props.row, 'Ditolak')">
              <q-tooltip>Tolak Pengajuan</q-tooltip>
            </q-btn>
          </div>
          <span v-else class="text-caption text-grey-5">Sudah Diproses</span>
        </q-td>
      </template>
    </q-table>

    <!-- POP-UP DIALOG: SIMULASI FORM PENGAJUAN OLEH GURU -->
    <q-dialog v-model="dialogMintaTerbuka">
      <q-card style="width: 400px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6"><q-icon name="post_add" /> Formulir Permintaan ATK</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-input outlined dense v-model="formMinta.namaGuru" label="Nama Guru Pemohon *" />
          <q-select outlined dense v-model="formMinta.peruntukan"
            :options="['Kelas 10-A', 'Kelas 11-B', 'Ruang Guru Utama', 'Ujian Sekolah']"
            label="Peruntukan / Lokasi Kelas *" />
          <q-select outlined dense v-model="formMinta.itemAtk"
            :options="['Spidol Whiteboard Hitam (Snowman)', 'Kertas HVS A4 80gr', 'Penghapus Papan Tulis', 'Isi Staples']"
            label="Pilih Jenis ATK *" />
          <q-input outlined dense type="number" v-model.number="formMinta.jumlah" label="Jumlah Kebutuhan *" />
        </q-card-section>

        <q-card-actions align="right" class="q-pa-md">
          <q-btn flat label="Batal" v-close-popup />
          <q-btn color="primary" label="Kirim Permintaan" @click="kirimPermintaanBaru" v-close-popup />
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
const dialogMintaTerbuka = ref(false)

// 1. Mock Data Riwayat Pengajuan ATK Guru
const daftarPermintaan = ref([
  { id: 'REQ-001', tanggal: '2026/06/06', namaGuru: 'Pak Hermawan, S.Pd', peruntukan: 'Kelas 10-A', itemAtk: 'Spidol Whiteboard Hitam (Snowman)', jumlah: 3, satuan: 'Pcs', status: 'Pending' },
  { id: 'REQ-002', tanggal: '2026/06/05', namaGuru: 'Bu Linda, M.Si', peruntukan: 'Ujian Sekolah', itemAtk: 'Kertas HVS A4 80gr', jumlah: 2, satuan: 'Rim', status: 'Disetujui' },
  { id: 'REQ-003', tanggal: '2026/06/02', namaGuru: 'Pak Danang, S.T', peruntukan: 'Kelas 11-B', itemAtk: 'Isi Staples', jumlah: 10, satuan: 'Kotak', status: 'Ditolak' }
])

// 2. State untuk Form Input Pengajuan Baru (Sisi Guru)
const formMinta = ref({
  namaGuru: '',
  peruntukan: '',
  itemAtk: '',
  jumlah: 1
})

// 3. Definisi Struktur Kolom Tabel Permintaan
const columns = [
  { name: 'id', label: 'Kode Nota', align: 'left', field: 'id', sortable: true },
  { name: 'tanggal', label: 'Tgl Pengajuan', align: 'center', field: 'tanggal' },
  { name: 'namaGuru', label: 'Guru Pemohon', align: 'left', field: 'namaGuru', sortable: true },
  { name: 'peruntukan', label: 'Untuk Keperluan', align: 'left', field: 'peruntukan' },
  { name: 'itemAtk', label: 'Nama ATK', align: 'left', field: 'itemAtk' },
  { name: 'jumlah', label: 'Vol', align: 'center', field: 'jumlah', format: (val, row) => `${val} ${row.satuan || 'Unit'}` },
  { name: 'status', label: 'Status Verifikasi', align: 'center', field: 'status', sortable: true },
  { name: 'aksi', label: 'Tindakan TU', align: 'center' }
]

// 4. Perhitungan Statistik Pengajuan yang Masih Menggantung
const totalPending = computed(() => {
  return daftarPermintaan.value.filter(req => req.status === 'Pending').length
})

// 5. Fungsi Aksi Pembaruan Status Dokumen oleh Staf TU
function updateStatusPermintaan(row, statusBaru) {
  row.status = statusBaru

  const warnaNotif = statusBaru === 'Disetujui' ? 'green' : 'red'
  const pesan = statusBaru === 'Disetujui'
    ? `Permintaan ${row.namaGuru} disetujui. Stok otomatis terpotong.`
    : `Permintaan ${row.namaGuru} telah ditolak.`

  $q.notify({
    color: warnaNotif,
    icon: statusBaru === 'Disetujui' ? 'done' : 'block',
    message: pesan
  })
}

// 6. Fungsi Kirim Pengajuan Baru dari Sisi Interface Guru
function kirimPermintaanBaru() {
  if (!formMinta.value.namaGuru || !formMinta.value.itemAtk || formMinta.value.jumlah <= 0) {
    $q.notify({ color: 'negative', message: 'Gagal mengirim. Mohon lengkapi seluruh kolom.' })
    return
  }

  // Tentukan satuan mock otomatis berdasarkan barang yang dipilih
  const satuanMock = formMinta.value.itemAtk.includes('Kertas') ? 'Rim' : 'Pcs'

  daftarPermintaan.value.unshift({
    id: `REQ-00${daftarPermintaan.value.length + 1}`,
    tanggal: '2026/06/06', // Sesuai tanggal hari ini
    namaGuru: formMinta.value.namaGuru,
    peruntukan: formMinta.value.peruntukan,
    itemAtk: formMinta.value.itemAtk,
    jumlah: formMinta.value.jumlah,
    satuan: satuanMock,
    status: 'Pending'
  })

  $q.notify({
    color: 'info',
    icon: 'send',
    message: 'Sukses! Permintaan Anda telah dikirim ke bagian Sarpras Tata Utara.'
  })

  // Reset form
  formMinta.value.namaGuru = ''
  formMinta.value.jumlah = 1
}
</script>
