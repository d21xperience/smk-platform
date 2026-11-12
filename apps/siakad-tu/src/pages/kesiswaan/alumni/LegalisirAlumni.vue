<template>
  <q-page class="q-pa-md bg-grey-2">
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Arsip Kelulusan & Legalisir</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pelayanan verifikasi nomor ijazah, pengarsipan data kelulusan, dan log permohonan legalisir alumni SMK</p>
      </div>
      <div>
        <q-btn color="primary" icon="note_add" label="Permohonan Legalisir" @click="openLegalisirDialog" />
      </div>
    </div>

    <q-banner v-if="store.error" class="bg-red-1 text-red-9 q-mb-md" rounded>
      <q-icon name="error" class="q-mr-sm" />{{ store.error }}
    </q-banner>

    <!-- Ringkasan Layanan Dokumen Alumni -->
    <div class="row q-col-gutter-sm q-mb-md">
      <div class="col-12 col-sm-4">
        <q-card class="bg-indigo-1 text-indigo-10 flat shadow-1">
          <q-card-section class="q-py-sm flex justify-between items-center">
            <div>
              <div class="text-caption text-weight-medium">Total Antrean</div>
              <div class="text-h5 text-weight-bold">{{ stats.antrean }} Berkas</div>
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
              <div class="text-h5 text-weight-bold">{{ stats.ttd }} Pemohon</div>
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
              <div class="text-h5 text-weight-bold">{{ stats.siap }} Pemohon</div>
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
            <q-input v-model="search" dense outlined placeholder="Cari Nama Alumni, No. Ijazah, atau Tahun Lulus..." clearable>
              <template v-slot:prepend><q-icon name="search" /></template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select v-model="filterLulusan" :options="['Semua Angkatan', '2025', '2024', '2023']" dense outlined label="Filter Tahun Kelulusan" />
          </div>
          <div class="col-12 col-md-3 flex justify-end">
            <q-btn icon="refresh" flat color="primary" :loading="store.loading" @click="store.fetchList">
              <q-tooltip>Muat Ulang</q-tooltip>
            </q-btn>
          </div>
        </div>
      </q-card-section>

      <!-- Tabel Log Pengajuan Legalisir -->
      <q-card-section class="q-pa-none">
        <q-table flat :rows="filteredList" :columns="columns" row-key="id" :loading="store.loading" :pagination="{ rowsPerPage: 10 }">
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
                :color="props.row.status === 'SIAP DIAMBIL' ? 'green-1' : 'orange-1'"
                :text-color="props.row.status === 'SIAP DIAMBIL' ? 'green-9' : 'orange-9'"
                :icon="props.row.status === 'SIAP DIAMBIL' ? 'check_circle' : 'hourglass_empty'">
                {{ props.row.status === 'SIAP DIAMBIL' ? 'Siap Diambil' : 'Proses TTD' }}
              </q-chip>
            </q-td>
          </template>

          <!-- Custom Kolom Aksi Cetak Bukti Permohonan -->
          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <q-btn size="sm" color="indigo" icon="confirmation_number" round flat @click="printResi(props.row)">
                <q-tooltip>Cetak Resi Pengambilan</q-tooltip>
              </q-btn>
              <q-btn size="sm" color="green-8" icon="done_all" round flat :disable="props.row.status !== 'SIAP DIAMBIL'"
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
          <q-input v-model="formLegalisir.nama" label="Nama Lengkap Alumni" outlined dense placeholder="Ketik nama lengkap sesuai ijazah..." />
          <div class="row q-col-gutter-sm">
            <div class="col-6">
              <q-select v-model="formLegalisir.tahunLulus" :options="['2025', '2024', '2023', '2022']" label="Tahun Kelulusan" outlined dense />
            </div>
            <div class="col-6">
              <q-select v-model="formLegalisir.jurusan" :options="['RPL', 'TKRO', 'AKL', 'TKJ']" label="Jurusan" outlined dense />
            </div>
          </div>
          <q-input v-model="formLegalisir.noIjazah" label="Nomor Seri Ijazah (Nasional)" outlined dense placeholder="Contoh: DN-01/M-SMK/0012345..." />
          <q-input v-model="formLegalisir.jumlahLembar" type="number" label="Jumlah Lembar yang Dilegalisir" outlined dense />
        </q-card-section>

        <q-card-actions align="right" class="q-pb-md q-pr-md">
          <q-btn flat label="Batal" v-close-popup color="grey" />
          <q-btn color="primary" label="Simpan Permohonan" :loading="saving" @click="saveLegalisirBaru" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useQuasar } from 'quasar'
import { useAlumniLegalisirStore } from '@/stores/kesiswaan/alumniLegalisirStore.js'

const $q = useQuasar()
const store = useAlumniLegalisirStore()

const search = ref('')
const filterLulusan = ref('Semua Angkatan')
const dialogLegalisirOpen = ref(false)
const saving = ref(false)

const formLegalisir = ref({ nama: '', tahunLulus: '2025', jurusan: 'RPL', noIjazah: '', jumlahLembar: 5 })

const columns = [
  { name: 'nama', align: 'left', label: 'Nama Alumni', field: 'nama', sortable: true },
  { name: 'ijazah', align: 'left', label: 'Arsip Dokumen Kelulusan', field: 'noIjazah' },
  { name: 'jumlahLembar', align: 'center', label: 'Jumlah (Lbr)', field: 'jumlahLembar' },
  { name: 'status', align: 'center', label: 'Status Proses', field: 'status', sortable: true },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

onMounted(() => {
  store.fetchList()
})

const filteredList = computed(() => {
  return store.list.filter(item => {
    const matchSearch = search.value
      ? item.nama.toLowerCase().includes(search.value.toLowerCase()) || item.noIjazah.includes(search.value)
      : true
    const matchAngkatan = filterLulusan.value !== 'Semua Angkatan' ? item.tahunLulus === filterLulusan.value : true
    return matchSearch && matchAngkatan
  })
})

const stats = computed(() => {
  let antrean = store.list.length
  let ttd = 0
  let siap = 0
  for (const item of store.list) {
    if (item.status === 'PROSES TTD') ttd++
    if (item.status === 'SIAP DIAMBIL') siap++
  }
  return { antrean, ttd, siap }
})

function openLegalisirDialog() {
  formLegalisir.value = { nama: '', tahunLulus: '2025', jurusan: 'RPL', noIjazah: '', jumlahLembar: 5 }
  dialogLegalisirOpen.value = true
}

async function saveLegalisirBaru() {
  if (!formLegalisir.value.nama || !formLegalisir.value.noIjazah) {
    $q.notify({ color: 'negative', message: 'Nama alumni dan Nomor seri ijazah wajib diisi!', icon: 'warning' })
    return
  }
  saving.value = true
  try {
    const result = await store.create(formLegalisir.value)
    if (result.success) {
      $q.notify({ color: 'primary', message: 'Permohonan legalisir berhasil didaftarkan.', icon: 'done' })
      dialogLegalisirOpen.value = false
      store.fetchList()
    } else {
      $q.notify({ color: 'negative', message: result.error, icon: 'error' })
    }
  } finally {
    saving.value = false
  }
}

async function advanceStatus(row) {
  const newStatus = row.status === 'PROSES TTD' ? 'SIAP DIAMBIL' : 'PROSES TTD'
  const result = await store.updateStatus(row.id, newStatus)
  if (result.success) {
    $q.notify({ color: 'green-8', message: `Status berkas ${row.nama} diubah ke ${newStatus}.`, icon: 'check_circle' })
    store.fetchList()
  }
}

async function finalizePermohonan(row) {
  $q.dialog({
    title: 'Konfirmasi',
    message: `Apakah berkas legalisir atas nama ${row.nama} sudah diserahkan ke alumni?`,
    cancel: true,
    persistent: true
  }).onOk(async () => {
    const result = await store.remove(row.id)
    if (result.success) {
      $q.notify({ color: 'green-9', message: `Berkas ${row.nama} resmi diambil.`, icon: 'archive' })
      store.fetchList()
    }
  })
}

function printResi(row) {
  $q.notify({ color: 'indigo-8', message: `Mencetak resi untuk ${row.nama}.`, icon: 'print' })
}
</script>
