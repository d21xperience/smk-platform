<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header Halaman -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Arsip Kelulusan & Legalisir</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pelayanan verifikasi nomor ijazah, pengarsipan data kelulusan, dan
          log permohonan legalisir alumni SMK</p>
      </div>
      <div>
        <q-btn color="primary" icon="note_add" label="Permohonan Legalisir" @click="openLegalisirDialog" />
      </div>
    </div>

    <!-- Ringkasan Layanan Dokumen Alumni -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card class="bg-indigo-1 text-indigo-10 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Antrean Legalisir Hari Ini</div>
              <div class="text-h5 text-weight-bold">14 Berkas</div>
            </div>
            <q-icon name="history_edu" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-amber-1 text-amber-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Proses TTD Kepala Sekolah</div>
              <div class="text-h5 text-weight-bold">5 Pemohon</div>
            </div>
            <q-icon name="draw" size="md" />
          </q-card-section>
        </q-card>
      </div>
      <div class="col-12 col-sm-4">
        <q-card class="bg-green-1 text-green-9 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Siap Diambil / Selesai</div>
              <div class="text-h5 text-weight-bold">9 Pemohon</div>
            </div>
            <q-icon name="inventory" size="md" />
          </q-card-section>
        </q-card>
      </div>
    </div>

    <!-- Panel Filter Pencarian & Tabel Utama -->
    <q-card class="shadow-1">
      <q-card-section class="q-py-sm">
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input v-model="search" dense outlined placeholder="Cari Nama Alumni, No. Ijazah, atau Tahun Lulus..."
              clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterLulusan"
              :options="['Semua Angkatan', 'Lulusan 2025', 'Lulusan 2024', 'Lulusan 2023']" dense outlined
              label="Filter Tahun Kelulusan" />
          </div>
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn flat color="secondary" icon="folder_zip" label="Buka Digital Arsip" @click="openArsipNasional" />
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Log Pengajuan Legalisir -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredLegalisir" :columns="columns" row-key="id" :pagination="{ rowsPerPage: 10 }">
          <!-- Custom Tampilan Kolom Detail Ijazah -->
          <template v-slot:body-cell-ijazah="props">
            <q-td :props="props">
              <div class="text-weight-bold text-grey-9">{{ props.row.noIjazah }}</div>
              <div class="text-caption text-primary">Angkatan {{ props.row.tahunLulus }} • {{ props.row.jurusan }}</div>
            </q-td>
          </template>

          <!-- Custom Tampilan Status Pemrosesan Berkas -->
          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-chip clickable @click="advanceStatus(props.row)" size="sm"
                :color="props.row.status === 'Siap Diambil' ? 'green-1' : 'orange-1'"
                :text-color="props.row.status === 'Siap Diambil' ? 'green-9' : 'orange-9'"
                :icon="props.row.status === 'Siap Diambil' ? 'check_circle' : 'hourglass_empty'">
                {{ props.row.status }}
              </q-chip>
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Cetak Bukti Permohonan -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="indigo" icon="confirmation_number" round flat @click="printResi(props.row)">
                <q-tooltip>Cetak Resi Pengambilan</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="green-8" icon="done_all" round flat :disable="props.row.status !== 'Siap Diambil'"
                @click="finalizePermohonan(props.row)">
                <q-tooltip>Selesai Diambil Alumni</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- DIALOG POPUP: Form Pencatatan Permohonan Legalisir Baru -->
    <q-dialog v-model="dialogLegalisirOpen" persistent>
      <q-card style="min-width: 400px;">
        <q-card-section class="bg-primary text-white">
          <div class="text-h6">Registrasi Pengajuan Legalisir</div>
          <div class="text-caption text-blue-2">Pencatatan tanda terima dokumen alumni sekolah</div>
        </q-card-section>

        <q-card-section class="q-gutter-md q-pt-md">
          <q-input v-model="formLegalisir.nama" label="Nama Lengkap Alumni" outlined dense
            placeholder="Ketik nama lengkap sesuai ijazah..." />
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select v-model="formLegalisir.tahunLulus" :options="['2025', '2024', '2023', '2022']"
                label="Tahun Kelulusan" outlined dense />
            </div>
            <div class="col-6">
              <q-select v-model="formLegalisir.jurusan" :options="['RPL', 'TKRO', 'AKL']" label="Jurusan" outlined
                dense />
            </div>
          </div>
          <q-input v-model="formLegalisir.noIjazah" label="Nomor Seri Ijazah (Nasional)" outlined dense
            placeholder="Contoh: DN-01/M-SMK/0012345..." />
          <q-input v-model="formLegalisir.jumlah" type="number" label="Jumlah Lembar yang Dilegalisir" outlined dense
            placeholder="Contoh: 5" />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Simpan Permohonan" @click="saveLegalisirBaru" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useQuasar } from 'quasar'

const $q = useQuasar()

// State Komponen
const search = ref('')
const filterLulusan = ref('Semua Angkatan')
const dialogLegalisirOpen = ref(false)

const formLegalisir = ref({ nama: '', tahunLulus: '2025', jurusan: 'RPL', noIjazah: '', jumlah: 5 })

// Struktur Kolom Tabel q-table Log Legalisir
const columns = [
  { name: 'nama', align: 'left', label: 'Nama Alumni', field: 'nama', sortable: true },
  { name: 'ijazah', align: 'left', label: 'Arsip Dokumen Kelulusan', field: 'noIjazah' },
  { name: 'jumlah', align: 'center', label: 'Jumlah (Lbr)', field: 'jumlah' },
  { name: 'status', align: 'center', label: 'Status Proses', field: 'status', sortable: true },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

// Data Mock / Dummy Log Pengajuan Layanan Dokumen Alumni
const dataLegalisir = ref([
  { id: 1, nama: 'Rian Hidayat', tahunLulus: '2025', jurusan: 'RPL', noIjazah: 'M-SMK/23/0012391', jumlah: 5, status: 'Siap Diambil' },
  { id: 2, nama: 'Siti Sarah', tahunLulus: '2025', jurusan: 'AKL', noIjazah: 'M-SMK/23/0052312', jumlah: 8, status: 'Proses TTD Kepala Sekolah' },
  { id: 3, nama: 'Denny Wahyudi', tahunLulus: '2024', jurusan: 'TKRO', noIjazah: 'M-SMK/22/0088123', jumlah: 5, status: 'Siap Diambil' },
  { id: 4, nama: 'Gita Permata', tahunLulus: '2023', jurusan: 'RPL', noIjazah: 'M-SMK/21/0091211', jumlah: 10, status: 'Proses TTD Kepala Sekolah' }
])

// Logika Compute Filter Pencarian Alumni di Frontend
const filteredLegalisir = computed(() => {
  return dataLegalisir.value.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) ||
      item.noIjazah.includes(search.value)
      : true

    const matchAngkatan = filterLulusan.value !== 'Semua Angkatan'
      ? item.tahunLulus === filterLulusan.value.split(' ')[1] // Mengambil angka tahun saja dari string "Lulusan 2025"
      : true

    return matchSearch && matchAngkatan
  })
})

// Fungsi Penanganan Event Aksi Kesiswaan & Alumni
const advanceStatus = (row) => {
  if (row.status === 'Proses TTD Kepala Sekolah') {
    row.status = 'Siap Diambil'
    $q.notify({ color: 'green-8', message: `Dokumen ${row.nama} selesai diproses dan siap diambil!`, icon: 'check_circle' })
  } else {
    row.status = 'Proses TTD Kepala Sekolah'
    $q.notify({ color: 'orange-9', message: `Dokumen ${row.nama} dikembalikan ke status antrean penandatanganan.`, icon: 'draw' })
  }
}

const openLegalisirDialog = () => {
  formLegalisir.value = { nama: '', tahunLulus: '2025', jurusan: 'RPL', noIjazah: '', jumlah: 5 }
  dialogLegalisirOpen.value = true
}

const saveLegalisirBaru = () => {
  if (!formLegalisir.value.nama || !formLegalisir.value.noIjazah) {
    $q.notify({ color: 'negative', message: 'Nama alumni dan Nomor seri ijazah wajib diisi!', icon: 'warning' })
    return
  }

  dataLegalisir.value.unshift({
    id: Date.now(),
    nama: formLegalisir.value.nama,
    tahunLulus: formLegalisir.value.tahunLulus,
    jurusan: formLegalisir.value.jurusan,
    noIjazah: formLegalisir.value.noIjazah,
    jumlah: formLegalisir.value.jumlah, status: 'Proses TTD Kepala Sekolah'
  })
  $q.notify({ color: 'primary', message: `Permohonan legalisir ijazah atas nama ${formLegalisir.value.nama} berhasil didaftarkan ke sistem.`, icon: 'done' })
  dialogLegalisirOpen.value = false
}
const printResi = (row) => { $q.notify({ color: 'indigo-8', message: `Mencetak lembar nomor tanda terima resi pengambilan legalisir untuk ${row.nama}.`, icon: 'print' }) }
const finalizePermohonan = (row) => {
  // Menghapus data dari antrean aktif karena telah selesai serah terima berkas fisik
  dataLegalisir.value = dataLegalisir.value.filter(item => item.id !== row.id)
  $q.notify({ color: 'green-9', message: `Penyerahan dokumen sukses. Berkas legalisir ${row.nama} resmi diambil dari arsip aktif.`, icon: 'archive' })
}
const openArsipNasional = () => {
  $q.notify({ color: 'grey-8', message: 'Membuka folder kompilasi file digital ijazah/transkrip nilai alumni (Sisi Klien).', icon: 'folder' })
}


</script>
