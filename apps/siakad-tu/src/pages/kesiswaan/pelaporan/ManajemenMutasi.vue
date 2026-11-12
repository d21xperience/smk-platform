<template>
  <q-page class="q-pa-md bg-grey-2">
    <!-- Header -->
    <div class="row items-center justify-between q-mb-md">
      <div>
        <h1 class="text-h5 text-weight-bold q-my-none text-primary">Manajemen Mutasi Siswa</h1>
        <p class="text-caption text-grey-7 q-mb-none">Pengelolaan administrasi keluar-masuk peserta didik SMK</p>
      </div>
      <div>
        <q-btn :color="activeTab === 'masuk' ? 'green-8' : 'red-9'" :icon="activeTab === 'masuk' ? 'login' : 'logout'"
          :label="activeTab === 'masuk' ? 'Catat Mutasi Masuk' : 'Catat Mutasi Keluar'" @click="openMutasiDialog" />
      </div>
    </div>

    <!-- Tabs -->
    <q-card class="shadow-1 q-mb-md">
      <q-tabs v-model="activeTab" dense class="text-grey" active-color="primary" indicator-color="primary" align="left"
        narrow-indicator>
        <q-tab name="masuk" icon="call_received" label="Mutasi Masuk (Pindahan)" />
        <q-tab name="keluar" icon="call_made" label="Mutasi Keluar (Pindah Sekolah)" />
      </q-tabs>
    </q-card>

    <!-- Tabel -->
    <q-card class="shadow-1">
      <q-card-section class="q-pa-none">
        <q-table flat :rows="activeTab === 'masuk' ? riwayatMasuk : riwayatKeluar"
          :columns="activeTab === 'masuk' ? kolomMasuk : kolomKeluar" row-key="id" :pagination="{ rowsPerPage: 10 }"
          :loading="isLoading">

          <template v-slot:body-cell-status="props">
            <q-td :props="props" class="text-center">
              <q-badge :color="getStatusColor(props.row.status)" :label="props.row.status" />
            </q-td>
          </template>

          <template v-slot:body-cell-aksi="props">
            <q-td :props="props" class="text-center q-gutter-xs">
              <!-- Tombol Approve/Reject hanya untuk PENDING -->
              <template v-if="props.row.status === 'PENDING'">
                <q-btn size="sm" color="green-8" icon="check" round flat @click="konfirmasiApprove(props.row)">
                  <q-tooltip>Setujui Mutasi</q-tooltip>
                </q-btn>
                <q-btn size="sm" color="red-8" icon="close" round flat @click="konfirmasiReject(props.row)">
                  <q-tooltip>Tolak Mutasi</q-tooltip>
                </q-btn>
              </template>
              <q-btn size="sm" color="primary" icon="print" round flat @click="cetakSuratMutasi(props.row)">
                <q-tooltip>Cetak Surat</q-tooltip>
              </q-btn>
            </q-td>
          </template>
        </q-table>
      </q-card-section>
    </q-card>

    <!-- Dialog Form -->
    <q-dialog v-model="dialogOpen" persistent>
      <q-card style="min-width: 500px; max-width: 600px;">
        <q-card-section class="row items-center q-pb-none">
          <div class="text-h6 text-weight-bold">
            {{ activeTab === 'masuk' ? 'Formulir Siswa Masuk Pindahan' : 'Formulir Siswa Pindah Keluar' }}
          </div>
          <q-space />
          <q-btn icon="close" flat round dense v-close-popup />
        </q-card-section>

        <q-separator class="q-mt-sm" />

        <q-card-section class="q-scroll q-py-md" style="max-height: 70vh;">
          <q-form @submit.prevent="simpanMutasi" class="row q-col-gutter-sm">

            <template v-if="activeTab === 'masuk'">
              <div class="col-12"><q-input v-model="form.nama" label="Nama Lengkap Siswa" outlined dense required />
              </div>
              <div class="col-6"><q-input v-model="form.nisn" label="NISN" outlined dense required type="number" />
              </div>
              <div class="col-6">
                <q-select v-model="form.jurusan" :options="opsiJurusan" label="Pilihan Jurusan SMK" outlined dense
                  required />
              </div>
              <div class="col-12"><q-input v-model="form.asalSekolah" label="Sekolah Asal (SMK/SMA Asal)" outlined dense
                  required /></div>
              <div class="col-12"><q-input v-model="form.alasan" label="Alasan Mutasi" outlined dense /></div>
            </template>

            <template v-else>
              <div class="col-12">
                <q-input v-model="form.nama" label="Nama Siswa Aktif" placeholder="Ketik nama siswa..." outlined dense
                  required />
                <p class="text-caption text-grey-6 q-mt-xs q-mb-none">*Sistem otomatis menonaktifkan status siswa
                  setelah
                  disetujui.</p>
              </div>
              <div class="col-12"><q-input v-model="form.tujuanSekolah" label="Sekolah Tujuan (Tempat Pindah)" outlined
                  dense required /></div>
              <div class="col-12"><q-input v-model="form.alasan" label="Alasan Pindah / Mutasi Keluar" outlined dense
                  type="textarea" rows="2" required /></div>
            </template>

            <div class="col-12 flex justify-end q-mt-md">
              <q-btn label="Batal" color="grey" flat v-close-popup class="q-mr-sm" />
              <q-btn :label="activeTab === 'masuk' ? 'Simpan & Daftarkan' : 'Proses Keluar'"
                :color="activeTab === 'masuk' ? 'green-8' : 'red-9'" type="submit" :loading="isLoading" />
            </div>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useQuasar } from 'quasar'
import { useMutasiStore } from '@/stores/kesiswaan/mutasiStore'

const $q = useQuasar()
const mutasiStore = useMutasiStore()

const activeTab = ref('masuk')
const dialogOpen = ref(false)

const form = ref({
  nama: '', nisn: '', jurusan: '', asalSekolah: '', tujuanSekolah: '',
  alasan: '',
})

const opsiJurusan = ['RPL', 'TKRO', 'AKL']

const riwayatMasuk = computed(() => mutasiStore.riwayatMasuk)
const riwayatKeluar = computed(() => mutasiStore.riwayatKeluar)
const isLoading = computed(() => mutasiStore.isLoading)

const kolomMasuk = [
  { name: 'tgl', align: 'left', label: 'Tanggal Masuk', field: (r) => formatDate(r.createdAt), sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama' },
  { name: 'nisn', align: 'left', label: 'NISN', field: 'nisn' },
  { name: 'asalSekolah', align: 'left', label: 'Sekolah Asal', field: 'asalSekolah' },
  { name: 'status', align: 'center', label: 'Status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

const kolomKeluar = [
  { name: 'tgl', align: 'left', label: 'Tanggal Keluar', field: (r) => formatDate(r.createdAt), sortable: true },
  { name: 'nama', align: 'left', label: 'Nama Siswa', field: 'nama' },
  { name: 'tujuanSekolah', align: 'left', label: 'Sekolah Tujuan', field: 'tujuanSekolah' },
  { name: 'alasan', align: 'left', label: 'Alasan Mutasi', field: 'alasan' },
  { name: 'status', align: 'center', label: 'Status' },
  { name: 'aksi', align: 'center', label: 'Aksi' }
]

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}

function getStatusColor(status) {
  const colors = { PENDING: 'orange', APPROVED: 'green', REJECTED: 'red' }
  return colors[status] || 'grey'
}

function openMutasiDialog() {
  form.value = { nama: '', nisn: '', jurusan: '', asalSekolah: '', tujuanSekolah: '', alasan: '' }
  dialogOpen.value = true
}

async function simpanMutasi() {
  const payload = {
    mutationType: activeTab.value === 'masuk' ? 'MASUK' : 'KELUAR',
    nama: form.value.nama,
    nisn: form.value.nisn,
    asalSekolah: form.value.asalSekolah,
    tujuanSekolah: form.value.tujuanSekolah,
    alasan: form.value.alasan,
  }

  const success = await mutasiStore.createMutation(payload)
  if (success) {
    $q.notify({ color: 'positive', message: 'Data mutasi berhasil disimpan!', icon: 'check' })
    dialogOpen.value = false
  } else {
    $q.notify({ color: 'negative', message: mutasiStore.error?.message || 'Gagal menyimpan mutasi', icon: 'error' })
  }
}

function konfirmasiApprove(row) {
  $q.dialog({
    title: 'Konfirmasi Setujui Mutasi',
    message: `Setujui mutasi untuk ${row.nama}? Siswa akan ${row.mutationType === 'KELUAR' ? 'dinonaktifkan dari sekolah' : 'diaktifkan'}.`,
    cancel: true,
    persistent: true,
  }).onOk(async () => {
    const success = await mutasiStore.approveMutation(row.id, 'Disetujui admin')
    if (success) {
      $q.notify({ color: 'green-9', message: 'Mutasi disetujui.', icon: 'check' })
    } else {
      $q.notify({ color: 'negative', message: mutasiStore.error?.message || 'Gagal', icon: 'error' })
    }
  })
}

function konfirmasiReject(row) {
  $q.dialog({
    title: 'Tolak Mutasi',
    message: 'Masukkan alasan penolakan:',
    prompt: { model: '', type: 'text' },
    cancel: true,
    persistent: true,
  }).onOk(async (reason) => {
    const success = await mutasiStore.rejectMutation(row.id, reason || 'Ditolak admin')
    if (success) {
      $q.notify({ color: 'orange-9', message: 'Mutasi ditolak.', icon: 'close' })
    }
  })
}

function cetakSuratMutasi(row) {
  $q.notify({ color: 'primary', message: `Mencetak surat mutasi untuk: ${row.nama}`, icon: 'print' })
}

onMounted(() => {
  mutasiStore.fetchMutations()
})

watch(activeTab, () => {
  // Refresh dengan filter type yang sesuai
  const type = activeTab.value === 'masuk' ? 'MASUK' : 'KELUAR'
  mutasiStore.fetchMutations({ type })
})
</script>
